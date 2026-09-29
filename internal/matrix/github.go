// Copyright © 2026 es-3581100. ALL RIGHTS RESERVED. See LICENSE and LEGAL.md.
package matrix

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"strings"
)

type GitHubResolver struct {
	BaseURL string
	Token   string
	Client  *http.Client
}

type githubRepo struct {
	DefaultBranch string `json:"default_branch"`
}
type githubTree struct {
	SHA       string `json:"sha"`
	Truncated bool   `json:"truncated"`
	Tree      []struct {
		Path string `json:"path"`
		Type string `json:"type"`
		Size *int64 `json:"size"`
		SHA  string `json:"sha"`
	} `json:"tree"`
}

func (g GitHubResolver) Pin(ctx context.Context, refs NormalizedRefGlob, projectID, buildID string) (Project, error) {
	if g.BaseURL == "" {
		g.BaseURL = "https://api.github.com"
	}
	if g.Client == nil {
		g.Client = http.DefaultClient
	}
	p := Project{
		SchemaVersion: ProjectSchemaVersion,
		CoreVersion:   "0.2.0",
		ProjectID:     projectID,
		Build:         BuildInfo{ID: buildID, Source: "GitHub repository metadata + recursive Git Trees API"},
		RefPack:       map[string]any{"raw_ref_glob": refs.Raw, "external_docs": refs.ExternalDocs},
	}
	for i, note := range refs.Annotations {
		p.Evidence = append(p.Evidence, Evidence{ID: fmt.Sprintf("annotation-%03d", i+1), Class: "annotation", Claim: note, Source: "ref_glob", Authority: "none"})
	}
	for _, nr := range refs.Repositories {
		var meta githubRepo
		if err := g.getJSON(ctx, "/repos/"+nr.Repo, &meta); err != nil {
			return Project{}, fmt.Errorf("resolve %s metadata: %w", nr.Repo, err)
		}
		ref := meta.DefaultBranch
		if ref == "" {
			return Project{}, fmt.Errorf("repository %s has empty default branch", nr.Repo)
		}
		explicit := map[string]bool{}
		for _, sp := range refs.Subpaths {
			if sp.Repo == nr.Repo && sp.Ref != "" {
				explicit[sp.Ref] = true
			}
		}
		if len(explicit) > 1 {
			keys := []string{}
			for k := range explicit {
				keys = append(keys, k)
			}
			sort.Strings(keys)
			return Project{}, fmt.Errorf("repository %s has conflicting explicit refs: %s", nr.Repo, strings.Join(keys, ","))
		}
		for k := range explicit {
			ref = k
		}

		var tree githubTree
		path := "/repos/" + nr.Repo + "/git/trees/" + url.PathEscape(ref) + "?recursive=1"
		if err := g.getJSON(ctx, path, &tree); err != nil {
			return Project{}, fmt.Errorf("resolve %s tree at %s: %w", nr.Repo, ref, err)
		}
		if tree.Truncated {
			return Project{}, fmt.Errorf("repository %s recursive tree is truncated", nr.Repo)
		}
		entries := make([]TreeEntry, 0, len(tree.Tree))
		exists := map[string]bool{}
		for _, e := range tree.Tree {
			if e.Type != "blob" && e.Type != "tree" {
				continue
			}
			entries = append(entries, TreeEntry{Path: e.Path, Type: e.Type, Size: e.Size, SHA: e.SHA})
			exists[e.Path] = true
		}
		sort.Slice(entries, func(i, j int) bool { return entries[i].Path < entries[j].Path })
		r := Repository{Repo: nr.Repo, Branch: ref, TreeSHA: tree.SHA, Deployment: "reference", SourceClass: "unknown", Description: "Pinned repository imported from ref glob.", Boundary: "Structural/reference evidence only; no execution authority.", Tree: entries, Truncated: false}
		p.Repositories = append(p.Repositories, r)
		p.Evidence = append(p.Evidence,
			Evidence{ID: "repo-" + safeID(nr.Repo), Class: "repo_metadata", Claim: "Resolved repository default branch metadata; selected ref=" + ref, Source: nr.URL, Repo: nr.Repo, Value: map[string]any{"default_branch": meta.DefaultBranch, "selected_ref": ref}, Authority: "none"},
			Evidence{ID: "tree-" + safeID(nr.Repo), Class: "git_tree", Claim: fmt.Sprintf("Pinned complete recursive Git tree with %d entries", len(entries)), Source: g.BaseURL, Repo: nr.Repo, Value: map[string]any{"tree_sha": tree.SHA, "entries": len(entries)}, Authority: "none"},
		)
		for _, sp := range refs.Subpaths {
			if sp.Repo != nr.Repo || sp.Path == "" {
				continue
			}
			if !exists[sp.Path] {
				return Project{}, fmt.Errorf("verified subpath missing: %s:%s", nr.Repo, sp.Path)
			}
			p.Evidence = append(p.Evidence, Evidence{ID: "path-" + safeID(nr.Repo+"-"+sp.Path), Class: "verified_subpath", Claim: "Named subpath exists in pinned recursive tree", Source: sp.URL, Repo: nr.Repo, Path: sp.Path, Value: true, Authority: "none"})
		}
	}
	sort.Slice(p.Repositories, func(i, j int) bool { return p.Repositories[i].Repo < p.Repositories[j].Repo })
	if err := p.Validate(); err != nil {
		return Project{}, err
	}
	return p, nil
}

func (g GitHubResolver) getJSON(ctx context.Context, path string, out any) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, strings.TrimRight(g.BaseURL, "/")+path, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")
	if g.Token != "" {
		req.Header.Set("Authorization", "Bearer "+g.Token)
	}
	res, err := g.Client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return fmt.Errorf("HTTP %d %s", res.StatusCode, res.Status)
	}
	if err := json.NewDecoder(res.Body).Decode(out); err != nil {
		return err
	}
	return nil
}

func safeID(s string) string {
	r := strings.NewReplacer("/", "-", "\", "-", " ", "-", ".", "-")
	return strings.ToLower(r.Replace(s))
}
