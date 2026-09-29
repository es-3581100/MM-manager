// Copyright © 2026 es-3581100. ALL RIGHTS RESERVED. See LICENSE and LEGAL.md.
package matrix

import (
	"fmt"
	"strings"
)

const ProjectSchemaVersion = "app-dir-matrix.project/v1"

type Project struct {
	SchemaVersion string       `json:"schema_version"`
	CoreVersion   string       `json:"core_version"`
	ProjectID     string       `json:"project_id"`
	Build         BuildInfo    `json:"build"`
	Repositories  []Repository `json:"repositories"`
	Evidence      []Evidence   `json:"evidence"`
	RefPack       any          `json:"ref_pack,omitempty"`
}

type BuildInfo struct {
	ID     string `json:"id"`
	Source string `json:"source"`
}

type Repository struct {
	Repo        string      `json:"repo"`
	Branch      string      `json:"branch"`
	TreeSHA     string      `json:"tree_sha"`
	Deployment  string      `json:"deployment"`
	SourceClass string      `json:"source_class"`
	Description string      `json:"description"`
	Boundary    string      `json:"boundary"`
	Site        string      `json:"site,omitempty"`
	Tree        []TreeEntry `json:"tree"`
	Truncated   bool        `json:"truncated"`
}

type TreeEntry struct {
	Path string `json:"path"`
	Type string `json:"type"`
	Size *int64 `json:"size"`
	SHA  string `json:"sha"`
}

type Evidence struct {
	ID        string `json:"id"`
	Class     string `json:"class"`
	Claim     string `json:"claim"`
	Source    string `json:"source,omitempty"`
	Repo      string `json:"repo,omitempty"`
	Path      string `json:"path,omitempty"`
	Value     any    `json:"value,omitempty"`
	Authority string `json:"authority"`
}

var allowedEvidenceClasses = map[string]bool{
	"annotation": true, "repo_metadata": true, "git_tree": true,
	"verified_subpath": true, "semantic_graph": true, "derived_level": true,
}

func (p Project) Validate() error {
	if p.SchemaVersion != ProjectSchemaVersion {
		return fmt.Errorf("schema_version %q != %q", p.SchemaVersion, ProjectSchemaVersion)
	}
	if strings.TrimSpace(p.CoreVersion) == "" || strings.TrimSpace(p.ProjectID) == "" {
		return fmt.Errorf("core_version and project_id are required")
	}
	if strings.TrimSpace(p.Build.ID) == "" || strings.TrimSpace(p.Build.Source) == "" {
		return fmt.Errorf("build.id and build.source are required")
	}
	seenRepo := map[string]bool{}
	for _, r := range p.Repositories {
		if strings.Count(r.Repo, "/") != 1 || r.Branch == "" || r.TreeSHA == "" {
			return fmt.Errorf("invalid repository identity/pin: %q", r.Repo)
		}
		if r.Truncated {
			return fmt.Errorf("repository %s is truncated and cannot be packaged as complete", r.Repo)
		}
		if seenRepo[r.Repo] {
			return fmt.Errorf("duplicate repository %s", r.Repo)
		}
		seenRepo[r.Repo] = true
		seenPath := map[string]bool{}
		for _, e := range r.Tree {
			if e.Path == "" || (e.Type != "tree" && e.Type != "blob") || e.SHA == "" {
				return fmt.Errorf("invalid tree entry in %s: %+v", r.Repo, e)
			}
			if seenPath[e.Path] {
				return fmt.Errorf("duplicate path %s:%s", r.Repo, e.Path)
			}
			seenPath[e.Path] = true
		}
	}
	seenEvidence := map[string]bool{}
	for _, e := range p.Evidence {
		if e.ID == "" || e.Claim == "" || !allowedEvidenceClasses[e.Class] {
			return fmt.Errorf("invalid evidence record %q class=%q", e.ID, e.Class)
		}
		if e.Authority != "none" {
			return fmt.Errorf("evidence %s attempts authority inflation: %q", e.ID, e.Authority)
		}
		if seenEvidence[e.ID] {
			return fmt.Errorf("duplicate evidence id %s", e.ID)
		}
		seenEvidence[e.ID] = true
		if e.Class == "derived_level" {
			want := LevelForPath(e.Path)
			got, ok := e.Value.(string)
			if !ok || got != want {
				return fmt.Errorf("derived_level %s=%v, want %s for %s", e.ID, e.Value, want, e.Path)
			}
		}
	}
	return nil
}

func LevelForPath(path string) string {
	path = strings.Trim(path, "/")
	if path == "" {
		return "application"
	}
	depth := len(strings.Split(path, "/"))
	switch depth {
	case 1:
		return "root"
	case 2:
		return "core"
	case 3:
		return "trunk"
	case 4:
		return "branch"
	default:
		return "leaves"
	}
}
