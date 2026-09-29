// Copyright © 2026 es-3581100. ALL RIGHTS RESERVED. See LICENSE and LEGAL.md.
package matrix

import "sort"

type DriftReport struct {
	Equal       bool        `json:"equal"`
	RepoChanges []RepoDrift `json:"repo_changes"`
}

type RepoDrift struct {
	Repo      string   `json:"repo"`
	PinBefore string   `json:"pin_before,omitempty"`
	PinAfter  string   `json:"pin_after,omitempty"`
	Added     []string `json:"added,omitempty"`
	Removed   []string `json:"removed,omitempty"`
	Changed   []string `json:"changed,omitempty"`
}

func CompareProjects(a, b Project) DriftReport {
	am, bm := map[string]Repository{}, map[string]Repository{}
	for _, r := range a.Repositories {
		am[r.Repo] = r
	}
	for _, r := range b.Repositories {
		bm[r.Repo] = r
	}
	keys := map[string]bool{}
	for k := range am {
		keys[k] = true
	}
	for k := range bm {
		keys[k] = true
	}
	names := []string{}
	for k := range keys {
		names = append(names, k)
	}
	sort.Strings(names)
	out := DriftReport{Equal: true}
	for _, name := range names {
		ra, oka := am[name]
		rb, okb := bm[name]
		d := RepoDrift{Repo: name}
		if oka {
			d.PinBefore = ra.TreeSHA
		}
		if okb {
			d.PinAfter = rb.TreeSHA
		}
		if !oka {
			for _, e := range rb.Tree {
				d.Added = append(d.Added, e.Path)
			}
		} else if !okb {
			for _, e := range ra.Tree {
				d.Removed = append(d.Removed, e.Path)
			}
		} else {
			pa, pb := map[string]TreeEntry{}, map[string]TreeEntry{}
			for _, e := range ra.Tree {
				pa[e.Path] = e
			}
			for _, e := range rb.Tree {
				pb[e.Path] = e
			}
			pk := map[string]bool{}
			for k := range pa {
				pk[k] = true
			}
			for k := range pb {
				pk[k] = true
			}
			paths := []string{}
			for k := range pk {
				paths = append(paths, k)
			}
			sort.Strings(paths)
			for _, p := range paths {
				ea, ao := pa[p]
				eb, bo := pb[p]
				if !ao {
					d.Added = append(d.Added, p)
				} else if !bo {
					d.Removed = append(d.Removed, p)
				} else if ea.SHA != eb.SHA || ea.Type != eb.Type {
					d.Changed = append(d.Changed, p)
				}
			}
		}
		if d.PinBefore != d.PinAfter || len(d.Added)+len(d.Removed)+len(d.Changed) > 0 {
			out.Equal = false
			out.RepoChanges = append(out.RepoChanges, d)
		}
	}
	return out
}
