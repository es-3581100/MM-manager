// Copyright © 2026 es-3581100. ALL RIGHTS RESERVED. See LICENSE and LEGAL.md.
package matrix

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func fixtureProject(t *testing.T) Project {
	t.Helper()
	path := filepath.Join("..", "..", "fixtures", "project-v1.json")
	p, err := LoadProject(path)
	if err != nil {
		t.Fatal(err)
	}
	return p
}

func TestNormalizeRefGlob(t *testing.T) {
	input := `<ref-block>
[alpha](https://github.com/fixture/alpha)
https\://github.com/fixture/alpha/tree/main/cmd/app
https://github.com/fixture/alpha
https://example.invalid/reference
Note: descriptive only.
</ref-block>`
	r := NormalizeRefGlob(input)
	if len(r.Repositories) != 1 || r.Repositories[0].Repo != "fixture/alpha" {
		t.Fatalf("repos=%+v", r.Repositories)
	}
	if len(r.Subpaths) != 1 || r.Subpaths[0].Path != "cmd/app" || r.Subpaths[0].Ref != "main" {
		t.Fatalf("subpaths=%+v", r.Subpaths)
	}
	if len(r.ExternalDocs) != 1 || r.ExternalDocs[0].Host != "example.invalid" {
		t.Fatalf("external=%+v", r.ExternalDocs)
	}
	if len(r.Annotations) != 1 || r.Annotations[0] != "Note: descriptive only." {
		t.Fatalf("annotations=%+v", r.Annotations)
	}
}

func TestValidateRejectsAuthorityInflation(t *testing.T) {
	p := fixtureProject(t)
	p.Evidence[0].Authority = "execute"
	if err := p.Validate(); err == nil || !strings.Contains(err.Error(), "authority inflation") {
		t.Fatalf("unexpected err=%v", err)
	}
}

func TestValidateRejectsTruncatedTree(t *testing.T) {
	p := fixtureProject(t)
	p.Repositories[0].Truncated = true
	if err := p.Validate(); err == nil || !strings.Contains(err.Error(), "truncated") {
		t.Fatalf("unexpected err=%v", err)
	}
}

func TestDerivedLevelContract(t *testing.T) {
	cases := map[string]string{"README.md": "root", "cmd/app": "core", "cmd/app/main.go": "trunk", "a/b/c/d.go": "branch", "a/b/c/d/e.go": "leaves"}
	for path, want := range cases {
		if got := LevelForPath(path); got != want {
			t.Fatalf("%s got=%s want=%s", path, got, want)
		}
	}
}

func TestBuildVerifyAndDeterminism(t *testing.T) {
	p := fixtureProject(t)
	template := filepath.Join("..", "..", "web", "app-dir-matrix-core-v0.2.0.html")
	a, err := BuildArtifact(template, p)
	if err != nil {
		t.Fatal(err)
	}
	b, err := BuildArtifact(template, p)
	if err != nil {
		t.Fatal(err)
	}
	if string(a) != string(b) {
		t.Fatal("artifact build is not deterministic")
	}
	if strings.Contains(string(a), "donCannoli-burns/kol-agent-sandbox") {
		t.Fatal("template repository data leaked into generated artifact")
	}
	if strings.Contains(string(a), `</script>\n<script id="embeddedProjectDocument"`) {
		t.Fatal("artifact contains literal backslash-n between embedded scripts")
	}
	d := t.TempDir()
	out := filepath.Join(d, "matrix.html")
	if err := os.WriteFile(out, a, 0o644); err != nil {
		t.Fatal(err)
	}
	r, err := VerifyArtifact(out)
	if err != nil {
		t.Fatal(err)
	}
	if !r.OK || r.EntryCount != 6 || r.RepositoryCount != 1 {
		t.Fatalf("report=%+v", r)
	}
}

func TestCompareProjectsDetectsDrift(t *testing.T) {
	a := fixtureProject(t)
	b := fixtureProject(t)
	b.Repositories[0].TreeSHA = "tree-alpha-002"
	b.Repositories[0].Tree[0].SHA = "blob-readme-002"
	b.Repositories[0].Tree = append(b.Repositories[0].Tree, TreeEntry{Path: "docs", Type: "tree", SHA: "tree-docs-001"})
	r := CompareProjects(a, b)
	if r.Equal || len(r.RepoChanges) != 1 {
		t.Fatalf("report=%+v", r)
	}
	if len(r.RepoChanges[0].Changed) != 1 || r.RepoChanges[0].Changed[0] != "README.md" {
		t.Fatalf("changed=%v", r.RepoChanges[0].Changed)
	}
	if len(r.RepoChanges[0].Added) != 1 || r.RepoChanges[0].Added[0] != "docs" {
		t.Fatalf("added=%v", r.RepoChanges[0].Added)
	}
}

func TestGitHubResolverPinsDefaultBranchAndSubpath(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/repos/fixture/alpha", func(w http.ResponseWriter, r *http.Request) { io.WriteString(w, `{"default_branch":"dev"}`) })
	mux.HandleFunc("/repos/fixture/alpha/git/trees/dev", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("recursive") != "1" {
			t.Fatalf("missing recursive=1")
		}
		io.WriteString(w, `{"sha":"tree-dev-1","truncated":false,"tree":[{"path":"cmd","type":"tree","sha":"t1"},{"path":"cmd/app","type":"tree","sha":"t2"},{"path":"cmd/app/main.go","type":"blob","size":9,"sha":"b1"}]}`)
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()
	refs := NormalizeRefGlob("https://github.com/fixture/alpha\nhttps://github.com/fixture/alpha/tree/dev/cmd/app\nNote: context only")
	p, err := (GitHubResolver{BaseURL: srv.URL, Client: srv.Client()}).Pin(context.Background(), refs, "pin-test", "pin-001")
	if err != nil {
		t.Fatal(err)
	}
	if len(p.Repositories) != 1 || p.Repositories[0].Branch != "dev" || p.Repositories[0].TreeSHA != "tree-dev-1" {
		t.Fatalf("project=%+v", p)
	}
	found := false
	for _, e := range p.Evidence {
		if e.Class == "verified_subpath" && e.Path == "cmd/app" {
			found = true
		}
	}
	if !found {
		t.Fatal("missing verified_subpath evidence")
	}
}

func TestGitHubResolverRejectsTruncation(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/repos/fixture/alpha", func(w http.ResponseWriter, r *http.Request) { io.WriteString(w, `{"default_branch":"main"}`) })
	mux.HandleFunc("/repos/fixture/alpha/git/trees/main", func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, `{"sha":"x","truncated":true,"tree":[]}`)
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()
	refs := NormalizeRefGlob("https://github.com/fixture/alpha")
	_, err := (GitHubResolver{BaseURL: srv.URL, Client: srv.Client()}).Pin(context.Background(), refs, "pin-test", "pin-001")
	if err == nil || !strings.Contains(err.Error(), "truncated") {
		t.Fatalf("err=%v", err)
	}
}

func TestGitHubResolverRejectsMissingNamedSubpath(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/repos/fixture/alpha", func(w http.ResponseWriter, r *http.Request) { io.WriteString(w, `{"default_branch":"main"}`) })
	mux.HandleFunc("/repos/fixture/alpha/git/trees/main", func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, `{"sha":"x","truncated":false,"tree":[{"path":"README.md","type":"blob","sha":"b"}]}`)
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()
	refs := NormalizeRefGlob("https://github.com/fixture/alpha/tree/main/missing/path")
	_, err := (GitHubResolver{BaseURL: srv.URL, Client: srv.Client()}).Pin(context.Background(), refs, "pin-test", "pin-001")
	if err == nil || !strings.Contains(err.Error(), "verified subpath missing") {
		t.Fatalf("err=%v", err)
	}
}

func TestCapturedGitHubFixtureReplaysOffline(t *testing.T) {
	fixtureDir := filepath.Join("..", "..", "fixtures", "github-live")
	client, capture, err := NewGitHubFixtureClient(fixtureDir)
	if err != nil {
		t.Fatal(err)
	}
	if capture.Repository != "octocat/Hello-World" || capture.CaptureMethod == "" {
		t.Fatalf("capture=%+v", capture)
	}
	b, err := os.ReadFile(filepath.Join(fixtureDir, "ref-glob.txt"))
	if err != nil {
		t.Fatal(err)
	}
	refs := NormalizeRefGlob(string(b))
	p, err := (GitHubResolver{BaseURL: "https://api.github.com", Client: client}).Pin(context.Background(), refs, "github-live-replay", "capture-20260929")
	if err != nil {
		t.Fatal(err)
	}
	if len(p.Repositories) != 1 {
		t.Fatalf("repositories=%d", len(p.Repositories))
	}
	r := p.Repositories[0]
	if r.Repo != "octocat/Hello-World" || r.Branch != "master" || r.TreeSHA != "7fd1a60b01f91b314f59955a4e4d4e80d8edf11d" {
		t.Fatalf("repo=%+v", r)
	}
	if len(r.Tree) != 1 || r.Tree[0].Path != "README" || r.Tree[0].SHA != "980a0d5f19a64b4b30a87d4206aade58726b60e3" {
		t.Fatalf("tree=%+v", r.Tree)
	}
	verified := false
	for _, e := range p.Evidence {
		if e.Class == "verified_subpath" && e.Repo == "octocat/Hello-World" && e.Path == "README" {
			verified = true
		}
	}
	if !verified {
		t.Fatal("captured fixture replay did not preserve named-subpath verification")
	}
}

func TestGitHubFixtureHasNoNetworkFallback(t *testing.T) {
	fixtureDir := filepath.Join("..", "..", "fixtures", "github-live")
	client, _, err := NewGitHubFixtureClient(fixtureDir)
	if err != nil {
		t.Fatal(err)
	}
	req, err := http.NewRequest(http.MethodGet, "https://api.github.com/repos/octocat/Hello-World/git/trees/not-recorded?recursive=1", nil)
	if err != nil {
		t.Fatal(err)
	}
	_, err = client.Do(req)
	if err == nil || !strings.Contains(err.Error(), "no recorded response") {
		t.Fatalf("expected offline fixture miss, got %v", err)
	}
}

func TestGitHubFixtureRejectsTamperedResponse(t *testing.T) {
	src := filepath.Join("..", "..", "fixtures", "github-live")
	dst := t.TempDir()
	if err := os.CopyFS(dst, os.DirFS(src)); err != nil {
		t.Fatal(err)
	}
	response := filepath.Join(dst, "octocat", "Hello-World", "tree-master.json")
	f, err := os.OpenFile(response, os.O_APPEND|os.O_WRONLY, 0)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := f.WriteString(" "); err != nil {
		f.Close()
		t.Fatal(err)
	}
	if err := f.Close(); err != nil {
		t.Fatal(err)
	}
	_, _, err = NewGitHubFixtureClient(dst)
	if err == nil || !strings.Contains(err.Error(), "hash mismatch") {
		t.Fatalf("expected hash mismatch, got %v", err)
	}
}
