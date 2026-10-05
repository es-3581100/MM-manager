// Copyright © 2026 es-3581100. ALL RIGHTS RESERVED. See LICENSE and LEGAL.md.
package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"appdirmatrix.local/core/internal/matrix"
)

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(2)
	}
	var err error
	switch os.Args[1] {
	case "build":
		err = cmdBuild(os.Args[2:])
	case "verify":
		err = cmdVerify(os.Args[2:])
	case "compare":
		err = cmdCompare(os.Args[2:])
	case "normalize":
		err = cmdNormalize(os.Args[2:])
	case "pin":
		err = cmdPin(os.Args[2:])
	case "agent":
		err = cmdAgent(os.Args[2:])
	default:
		usage()
		os.Exit(2)
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}

func usage() {
	fmt.Fprintln(os.Stderr, "appdir-matrix <normalize|pin|build|verify|compare|agent> [flags]")
	fmt.Fprintln(os.Stderr, "appdir-matrix agent <bootstrap|resolve|inspect|expand|scope|verify-receipt> [flags]")
}

func cmdBuild(args []string) error {
	fs := flag.NewFlagSet("build", flag.ContinueOnError)
	projectPath := fs.String("project", "", "project JSON")
	template := fs.String("template", "", "core HTML template")
	out := fs.String("out", "", "output HTML")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if *projectPath == "" || *template == "" || *out == "" {
		return fmt.Errorf("--project, --template and --out are required")
	}
	p, err := matrix.LoadProject(*projectPath)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(*out), 0o755); err != nil {
		return err
	}
	h, err := matrix.WriteArtifact(*template, p, *out)
	if err != nil {
		return err
	}
	fmt.Printf("built %s\nsha256 %s\n", *out, h)
	return nil
}

func cmdVerify(args []string) error {
	fs := flag.NewFlagSet("verify", flag.ContinueOnError)
	artifact := fs.String("artifact", "", "artifact HTML")
	sha := fs.String("sha", "", "optional sha256 sidecar")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if *artifact == "" {
		return fmt.Errorf("--artifact is required")
	}
	r, err := matrix.VerifyArtifact(*artifact)
	if err != nil {
		return err
	}
	if *sha != "" {
		if err := matrix.VerifySHAFile(*artifact, *sha); err != nil {
			return err
		}
		r.Checks = append(r.Checks, "sha256 sidecar matches artifact")
	}
	return printJSON(r)
}

func cmdCompare(args []string) error {
	fs := flag.NewFlagSet("compare", flag.ContinueOnError)
	left := fs.String("left", "", "left project JSON")
	right := fs.String("right", "", "right project JSON")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if *left == "" || *right == "" {
		return fmt.Errorf("--left and --right are required")
	}
	a, err := matrix.LoadProject(*left)
	if err != nil {
		return err
	}
	b, err := matrix.LoadProject(*right)
	if err != nil {
		return err
	}
	return printJSON(matrix.CompareProjects(a, b))
}

func cmdNormalize(args []string) error {
	fs := flag.NewFlagSet("normalize", flag.ContinueOnError)
	in := fs.String("in", "", "input ref glob")
	out := fs.String("out", "", "optional output JSON")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if *in == "" {
		return fmt.Errorf("--in is required")
	}
	b, err := os.ReadFile(*in)
	if err != nil {
		return err
	}
	r := matrix.NormalizeRefGlob(string(b))
	j, _ := json.MarshalIndent(r, "", "  ")
	if *out != "" {
		return os.WriteFile(*out, append(j, '\n'), 0o644)
	}
	fmt.Println(string(j))
	return nil
}

func cmdPin(args []string) error {
	fs := flag.NewFlagSet("pin", flag.ContinueOnError)
	in := fs.String("in", "", "input ref glob")
	out := fs.String("out", "", "output project JSON")
	projectID := fs.String("project-id", "imported-project", "project id")
	buildID := fs.String("build-id", "pin-001", "deterministic build id")
	api := fs.String("github-api", "https://api.github.com", "GitHub API base URL")
	fixtureDir := fs.String("github-fixture-dir", "", "offline GitHub API fixture directory containing capture.json")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if *in == "" || *out == "" {
		return fmt.Errorf("--in and --out are required")
	}
	b, err := os.ReadFile(*in)
	if err != nil {
		return err
	}
	refs := matrix.NormalizeRefGlob(string(b))
	r := matrix.GitHubResolver{BaseURL: *api, Token: os.Getenv("GITHUB_TOKEN")}
	var replayFixture *matrix.GitHubFixture
	if *fixtureDir != "" {
		client, fixture, err := matrix.NewGitHubFixtureClient(*fixtureDir)
		if err != nil {
			return err
		}
		replayFixture = &fixture
		r.Client = client
		r.Token = ""
		fmt.Fprintf(os.Stderr, "offline GitHub fixture: %s captured %s via %s\n", fixture.Repository, fixture.CapturedAt, fixture.CaptureMethod)
	}
	p, err := r.Pin(context.Background(), refs, *projectID, *buildID)
	if err != nil {
		return err
	}
	if replayFixture != nil {
		pack, ok := p.RefPack.(map[string]any)
		if !ok || pack == nil {
			pack = map[string]any{}
		}
		pack["github_fixture"] = map[string]any{
			"schema":         replayFixture.Schema,
			"captured_at":    replayFixture.CapturedAt,
			"source":         replayFixture.Source,
			"repository":     replayFixture.Repository,
			"capture_method": replayFixture.CaptureMethod,
		}
		p.RefPack = pack
	}
	j, err := json.MarshalIndent(p, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(*out, append(j, '\n'), 0o644)
}

func cmdAgent(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("agent subcommand is required")
	}
	switch args[0] {
	case "bootstrap":
		return cmdAgentBootstrap(args[1:])
	case "resolve":
		return cmdAgentResolve(args[1:])
	case "inspect":
		return cmdAgentInspect(args[1:])
	case "expand":
		return cmdAgentExpand(args[1:])
	case "scope":
		return cmdAgentScope(args[1:])
	case "verify-receipt":
		return cmdAgentVerifyReceipt(args[1:])
	default:
		return fmt.Errorf("unknown agent subcommand %q", args[0])
	}
}

func cmdAgentBootstrap(args []string) error {
	fs := flag.NewFlagSet("agent bootstrap", flag.ContinueOnError)
	project := fs.String("project", "", "project JSON")
	if err := fs.Parse(args); err != nil {
		return err
	}
	idx, err := loadAgentIndex(*project)
	if err != nil {
		return err
	}
	return printJSON(idx.Bootstrap())
}

func cmdAgentResolve(args []string) error {
	fs := flag.NewFlagSet("agent resolve", flag.ContinueOnError)
	project := fs.String("project", "", "project JSON")
	query := fs.String("query", "", "query text")
	limit := fs.Int("limit", 8, "maximum results (1-32)")
	if err := fs.Parse(args); err != nil {
		return err
	}
	idx, err := loadAgentIndex(*project)
	if err != nil {
		return err
	}
	result, err := idx.Resolve(*query, *limit)
	if err != nil {
		return err
	}
	return printJSON(result)
}

func cmdAgentInspect(args []string) error {
	fs := flag.NewFlagSet("agent inspect", flag.ContinueOnError)
	project := fs.String("project", "", "project JSON")
	pointer := fs.String("pointer", "", "exact node pointer")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if *pointer == "" {
		return fmt.Errorf("--pointer is required")
	}
	idx, err := loadAgentIndex(*project)
	if err != nil {
		return err
	}
	result, err := idx.Inspect(*pointer)
	if err != nil {
		return err
	}
	return printJSON(result)
}

func cmdAgentExpand(args []string) error {
	fs := flag.NewFlagSet("agent expand", flag.ContinueOnError)
	project := fs.String("project", "", "project JSON")
	pointer := fs.String("pointer", "", "exact node pointer")
	depth := fs.Int("depth", 1, "structural expansion depth (1-2)")
	limit := fs.Int("limit", 8, "maximum nodes (1-32)")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if *pointer == "" {
		return fmt.Errorf("--pointer is required")
	}
	idx, err := loadAgentIndex(*project)
	if err != nil {
		return err
	}
	result, err := idx.Expand(*pointer, *depth, *limit)
	if err != nil {
		return err
	}
	return printJSON(result)
}

func cmdAgentScope(args []string) error {
	fs := flag.NewFlagSet("agent scope", flag.ContinueOnError)
	project := fs.String("project", "", "project JSON")
	query := fs.String("query", "", "task/query text")
	limit := fs.Int("limit", 12, "maximum included nodes (1-32)")
	if err := fs.Parse(args); err != nil {
		return err
	}
	idx, err := loadAgentIndex(*project)
	if err != nil {
		return err
	}
	result, err := idx.Scope(*query, *limit)
	if err != nil {
		return err
	}
	return printJSON(result)
}

func cmdAgentVerifyReceipt(args []string) error {
	fs := flag.NewFlagSet("agent verify-receipt", flag.ContinueOnError)
	project := fs.String("project", "", "project JSON")
	receiptPath := fs.String("receipt", "", "retrieval receipt JSON")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if *receiptPath == "" {
		return fmt.Errorf("--receipt is required")
	}
	idx, err := loadAgentIndex(*project)
	if err != nil {
		return err
	}
	receipt, err := matrix.LoadRetrievalReceipt(*receiptPath)
	if err != nil {
		return err
	}
	result, err := idx.VerifyReceipt(receipt)
	if err != nil {
		return err
	}
	return printJSON(result)
}

func loadAgentIndex(project string) (*matrix.AgentIndex, error) {
	if project == "" {
		return nil, fmt.Errorf("--project is required")
	}
	return matrix.NewAgentIndex(project)
}

func printJSON(v any) error {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(b))
	return nil
}
