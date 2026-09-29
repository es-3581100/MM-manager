// Copyright © 2026 es-3581100. ALL RIGHTS RESERVED. See LICENSE and LEGAL.md.
package matrix

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

type VerifyReport struct {
	OK              bool     `json:"ok"`
	ProjectID       string   `json:"project_id"`
	SchemaVersion   string   `json:"schema_version"`
	RepositoryCount int      `json:"repository_count"`
	EntryCount      int      `json:"entry_count"`
	ArtifactSHA256  string   `json:"artifact_sha256"`
	Checks          []string `json:"checks"`
}

func VerifyArtifact(path string) (VerifyReport, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return VerifyReport{}, err
	}
	text := string(b)
	pj, err := extractScriptJSON(text, "embeddedProjectDocument")
	if err != nil {
		return VerifyReport{}, err
	}
	var project Project
	if err := json.Unmarshal(pj, &project); err != nil {
		return VerifyReport{}, fmt.Errorf("decode embedded project: %w", err)
	}
	if err := project.Validate(); err != nil {
		return VerifyReport{}, fmt.Errorf("embedded project invalid: %w", err)
	}

	sj, err := extractScriptJSON(text, "embeddedSnapshot")
	if err != nil {
		return VerifyReport{}, err
	}
	var snap map[string]snapshotPayload
	if err := json.Unmarshal(sj, &snap); err != nil {
		return VerifyReport{}, fmt.Errorf("decode embedded snapshot: %w", err)
	}

	report := VerifyReport{OK: true, ProjectID: project.ProjectID, SchemaVersion: project.SchemaVersion, RepositoryCount: len(project.Repositories), ArtifactSHA256: SHA256Bytes(b)}
	report.Checks = append(report.Checks, "embedded project document validates", "embedded snapshot parses")
	for _, r := range project.Repositories {
		p, ok := snap[r.Repo]
		if !ok {
			return VerifyReport{}, fmt.Errorf("snapshot missing repository %s", r.Repo)
		}
		if p.Truncated {
			return VerifyReport{}, fmt.Errorf("snapshot %s is truncated", r.Repo)
		}
		if p.SHA != r.TreeSHA {
			return VerifyReport{}, fmt.Errorf("snapshot pin mismatch for %s: %s != %s", r.Repo, p.SHA, r.TreeSHA)
		}
		if len(p.Tree) != len(r.Tree) {
			return VerifyReport{}, fmt.Errorf("entry count mismatch for %s", r.Repo)
		}
		want := map[string]TreeEntry{}
		for _, e := range r.Tree {
			want[e.Path] = e
		}
		for _, e := range p.Tree {
			w, ok := want[e.Path]
			if !ok || w.Type != e.Type || w.SHA != e.SHA {
				return VerifyReport{}, fmt.Errorf("tree mismatch %s:%s", r.Repo, e.Path)
			}
			report.EntryCount++
		}
	}
	report.Checks = append(report.Checks, "repository pins match snapshot", "tree entries match project document", "evidence classes preserve authority=none")
	return report, nil
}

func VerifySHAFile(artifactPath, shaPath string) error {
	b, err := os.ReadFile(shaPath)
	if err != nil {
		return err
	}
	fields := strings.Fields(string(b))
	if len(fields) == 0 {
		return fmt.Errorf("empty sha256 file")
	}
	got, err := SHA256File(artifactPath)
	if err != nil {
		return err
	}
	if fields[0] != got {
		return fmt.Errorf("sha256 mismatch: sidecar=%s actual=%s", fields[0], got)
	}
	return nil
}

func extractScriptJSON(text, id string) ([]byte, error) {
	marker := `<script id="` + id + `"`
	start := strings.Index(text, marker)
	if start < 0 {
		return nil, fmt.Errorf("missing script %s", id)
	}
	gt := strings.Index(text[start:], ">")
	if gt < 0 {
		return nil, fmt.Errorf("malformed script %s", id)
	}
	contentStart := start + gt + 1
	end := strings.Index(text[contentStart:], "</script>")
	if end < 0 {
		return nil, fmt.Errorf("unterminated script %s", id)
	}
	b := []byte(text[contentStart : contentStart+end])
	b = []byte(strings.ReplaceAll(string(b), `<\/`, `</`))
	return b, nil
}
