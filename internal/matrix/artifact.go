// Copyright © 2026 es-3581100. ALL RIGHTS RESERVED. See LICENSE and LEGAL.md.
package matrix

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"
)

type uiRepo struct {
	Repo        string `json:"repo"`
	Branch      string `json:"branch"`
	TreeSHA     string `json:"treeSha"`
	Entries     int    `json:"entries"`
	Files       int    `json:"files"`
	Dirs        int    `json:"dirs"`
	Deployment  string `json:"deployment"`
	SourceClass string `json:"sourceClass"`
	Description string `json:"description"`
	Boundary    string `json:"boundary"`
	Site        string `json:"site,omitempty"`
}

type snapshotPayload struct {
	SHA       string      `json:"sha"`
	Truncated bool        `json:"truncated"`
	Tree      []TreeEntry `json:"tree"`
}

func BuildArtifact(templatePath string, project Project) ([]byte, error) {
	if err := project.Validate(); err != nil {
		return nil, err
	}
	tpl, err := os.ReadFile(templatePath)
	if err != nil {
		return nil, err
	}
	text := string(tpl)

	repos := make([]uiRepo, 0, len(project.Repositories))
	snapshot := map[string]snapshotPayload{}
	verifiedEntries, verifiedFiles, verifiedDirs := 0, 0, 0
	for _, r := range project.Repositories {
		files, dirs := 0, 0
		entries := append([]TreeEntry(nil), r.Tree...)
		sort.Slice(entries, func(i, j int) bool { return entries[i].Path < entries[j].Path })
		for _, e := range entries {
			if e.Type == "tree" {
				dirs++
			} else {
				files++
			}
		}
		verifiedEntries += len(entries)
		verifiedFiles += files
		verifiedDirs += dirs
		repos = append(repos, uiRepo{r.Repo, r.Branch, r.TreeSHA, len(entries), files, dirs, r.Deployment, r.SourceClass, r.Description, r.Boundary, r.Site})
		snapshot[r.Repo] = snapshotPayload{SHA: r.TreeSHA, Truncated: false, Tree: entries}
	}
	sort.Slice(repos, func(i, j int) bool { return repos[i].Repo < repos[j].Repo })

	buildMeta := map[string]any{
		"version":         project.CoreVersion,
		"schemaVersion":   project.SchemaVersion,
		"buildId":         project.Build.ID,
		"source":          project.Build.Source,
		"verifiedEntries": verifiedEntries,
		"verifiedFiles":   verifiedFiles,
		"verifiedDirs":    verifiedDirs,
		"repositoryCount": len(repos),
		"reproducible":    true,
		"refGlobContract": "app_dir_matrix.ingest_ref_glob@1.0.0",
	}

	reposJSON, _ := json.Marshal(repos)
	metaJSON, _ := json.Marshal(buildMeta)
	refPackJSON, _ := json.Marshal(project.RefPack)
	snapJSON, _ := json.Marshal(snapshot)
	projectJSON, _ := json.Marshal(project)
	for _, p := range []*[]byte{&reposJSON, &metaJSON, &refPackJSON, &snapJSON, &projectJSON} {
		*p = bytes.ReplaceAll(*p, []byte("</"), []byte("<\\/"))
	}

	var ok bool
	text, ok = replaceJSConstLine(text, "REPOS", string(reposJSON))
	if !ok {
		return nil, fmt.Errorf("template missing const REPOS")
	}
	text, ok = replaceJSConstLine(text, "BUILD_META", string(metaJSON))
	if !ok {
		return nil, fmt.Errorf("template missing const BUILD_META")
	}
	text, ok = replaceJSConstLine(text, "REF_PACK", string(refPackJSON))
	if !ok {
		return nil, fmt.Errorf("template missing const REF_PACK")
	}

	oldSnapshot := `<script id="embeddedSnapshot" type="application/json">{}</script>`
	if !strings.Contains(text, oldSnapshot) {
		return nil, fmt.Errorf("template missing empty embeddedSnapshot marker")
	}
	newSnapshot := `<script id="embeddedSnapshot" type="application/json">` + string(snapJSON) + `</script>` +
		"\n" + `<script id="embeddedProjectDocument" type="application/json">` + string(projectJSON) + `</script>`
	text = strings.Replace(text, oldSnapshot, newSnapshot, 1)

	header := fmt.Sprintf("<!-- APP-DIR-MATRIX GENERATED ARTIFACT | schema=%s | project=%s | build=%s | reproducible=true | Copyright © 2026 es-3581100 | ALL RIGHTS RESERVED | See repository LICENSE and LEGAL.md -->\n", project.SchemaVersion, project.ProjectID, project.Build.ID)
	if strings.HasPrefix(text, "<!doctype html>\n") {
		text = strings.Replace(text, "<!doctype html>\n", "<!doctype html>\n"+header, 1)
	} else {
		text = header + text
	}
	return []byte(text), nil
}

func replaceJSConstLine(text, name, value string) (string, bool) {
	marker := "const " + name + " = "
	start := strings.Index(text, marker)
	if start < 0 {
		return text, false
	}
	lineEnd := strings.IndexByte(text[start:], '\n')
	if lineEnd < 0 {
		lineEnd = len(text) - start
	}
	lineEnd += start
	replacement := marker + value + ";"
	return text[:start] + replacement + text[lineEnd:], true
}

func WriteArtifact(templatePath string, project Project, outPath string) (string, error) {
	b, err := BuildArtifact(templatePath, project)
	if err != nil {
		return "", err
	}
	if err := os.WriteFile(outPath, b, 0o644); err != nil {
		return "", err
	}
	hash := SHA256Bytes(b)
	if err := os.WriteFile(outPath+".sha256", []byte(hash+"  "+baseName(outPath)+"\n"), 0o644); err != nil {
		return "", err
	}
	return hash, nil
}

func baseName(path string) string {
	if i := strings.LastIndexAny(path, "/\\"); i >= 0 {
		return path[i+1:]
	}
	return path
}
