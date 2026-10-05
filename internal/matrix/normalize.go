// Copyright © 2026 es-3581100. ALL RIGHTS RESERVED. See LICENSE and LEGAL.md.
package matrix

import (
	"net/url"
	"regexp"
	"sort"
	"strings"
)

type NormalizedRefGlob struct {
	Raw          string           `json:"raw"`
	Repositories []NormalizedRepo `json:"repositories"`
	Subpaths     []NormalizedPath `json:"subpaths"`
	ExternalDocs []ExternalDoc    `json:"external_docs"`
	Annotations  []string         `json:"annotations"`
}

type NormalizedRepo struct {
	Repo string `json:"repo"`
	URL  string `json:"url"`
}

type NormalizedPath struct {
	Repo string `json:"repo"`
	Mode string `json:"mode"`
	Ref  string `json:"ref,omitempty"`
	Path string `json:"path,omitempty"`
	URL  string `json:"url"`
}

type ExternalDoc struct {
	URL  string `json:"url"`
	Host string `json:"host"`
}

var mdURL = regexp.MustCompile(`\[[^\]]*\]\((https?://[^)\s]+)\)`)
var plainURL = regexp.MustCompile(`https?://[^\s<>"'` + "`" + `]+`)

func NormalizeRefGlob(input string) NormalizedRefGlob {
	raw := strings.ReplaceAll(input, `\://`, `://`)
	raw = strings.ReplaceAll(raw, `\/`, `/`)
	urls := []string{}
	for _, m := range mdURL.FindAllStringSubmatch(raw, -1) {
		urls = append(urls, m[1])
	}
	for _, m := range plainURL.FindAllString(raw, -1) {
		urls = append(urls, strings.TrimRight(m, `),.;\`))
	}
	seen := map[string]bool{}
	unique := []string{}
	for _, u := range urls {
		if !seen[u] {
			seen[u] = true
			unique = append(unique, u)
		}
	}

	repoMap := map[string]NormalizedRepo{}
	paths := []NormalizedPath{}
	external := []ExternalDoc{}
	for _, rawURL := range unique {
		u, err := url.Parse(rawURL)
		if err != nil || u.Scheme == "" || u.Host == "" {
			continue
		}
		parts := splitPath(u.Path)
		if (u.Host == "github.com" || u.Host == "www.github.com") && len(parts) >= 2 {
			repo := parts[0] + "/" + strings.TrimSuffix(parts[1], ".git")
			repoMap[repo] = NormalizedRepo{Repo: repo, URL: "https://github.com/" + repo}
			if len(parts) > 2 {
				mode := parts[2]
				ref, path := "", ""
				if mode == "tree" || mode == "blob" {
					if len(parts) > 3 {
						ref = parts[3]
					}
					if len(parts) > 4 {
						path = strings.Join(parts[4:], "/")
					}
				} else {
					mode = "path"
					path = strings.Join(parts[2:], "/")
				}
				paths = append(paths, NormalizedPath{Repo: repo, Mode: mode, Ref: ref, Path: path, URL: rawURL})
			}
		} else {
			external = append(external, ExternalDoc{URL: rawURL, Host: u.Host})
		}
	}
	repos := make([]NormalizedRepo, 0, len(repoMap))
	for _, r := range repoMap {
		repos = append(repos, r)
	}
	sort.Slice(repos, func(i, j int) bool { return repos[i].Repo < repos[j].Repo })
	sort.Slice(paths, func(i, j int) bool { return paths[i].URL < paths[j].URL })
	sort.Slice(external, func(i, j int) bool { return external[i].URL < external[j].URL })
	annotations := []string{}
	for _, line := range strings.Split(raw, "\n") {
		s := strings.TrimSpace(line)
		if s == "" || strings.HasPrefix(s, "<ref-block") || strings.HasPrefix(s, "</ref-block") {
			continue
		}
		without := mdURL.ReplaceAllString(s, "")
		without = plainURL.ReplaceAllString(without, "")
		without = strings.TrimSpace(strings.Trim(without, "-•* "))
		if without != "" {
			annotations = append(annotations, without)
		}
	}
	return NormalizedRefGlob{Raw: raw, Repositories: repos, Subpaths: paths, ExternalDocs: external, Annotations: annotations}
}

func splitPath(p string) []string {
	out := []string{}
	for _, s := range strings.Split(p, "/") {
		if s != "" {
			out = append(out, s)
		}
	}
	return out
}
