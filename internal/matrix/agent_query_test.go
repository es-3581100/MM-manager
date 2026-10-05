// Copyright © 2026 es-3581100. ALL RIGHTS RESERVED. See LICENSE and LEGAL.md.
package matrix

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func fixtureAgentIndex(t *testing.T) *AgentIndex {
	t.Helper()
	path := filepath.Join("..", "..", "fixtures", "project-v1.json")
	idx, err := NewAgentIndex(path)
	if err != nil {
		t.Fatal(err)
	}
	return idx
}

func TestAgentBootstrapIsReadOnlyAndBounded(t *testing.T) {
	idx := fixtureAgentIndex(t)
	got := idx.Bootstrap()
	if got.Authority != "none" {
		t.Fatalf("authority=%q", got.Authority)
	}
	if len(got.Operations) != 6 {
		t.Fatalf("operations=%v", got.Operations)
	}
	if got.ProjectSHA256 == "" {
		t.Fatal("missing project hash")
	}
}

func TestAgentResolveIsDeterministicAndNonAuthoritative(t *testing.T) {
	idx := fixtureAgentIndex(t)
	a, err := idx.Resolve("main model", 8)
	if err != nil {
		t.Fatal(err)
	}
	b, err := idx.Resolve("main model", 8)
	if err != nil {
		t.Fatal(err)
	}
	aj, _ := json.Marshal(a)
	bj, _ := json.Marshal(b)
	if string(aj) != string(bj) {
		t.Fatal("resolve is not deterministic")
	}
	if a.Authority != "none" || len(a.Results) == 0 {
		t.Fatalf("resolve=%+v", a)
	}
	for _, match := range a.Results {
		if match.Node.Authority != "none" {
			t.Fatalf("node authority=%q", match.Node.Authority)
		}
	}
}

func TestAgentInspectUnknownFailsClosed(t *testing.T) {
	idx := fixtureAgentIndex(t)
	if _, err := idx.Inspect("mem://path/fixture/alpha/nope"); err == nil {
		t.Fatal("expected unknown pointer failure")
	}
}

func TestAgentExpandUsesStructuralNeighborsOnly(t *testing.T) {
	idx := fixtureAgentIndex(t)
	got, err := idx.Expand("mem://path/fixture/alpha/cmd/app", 1, 8)
	if err != nil {
		t.Fatal(err)
	}
	want := map[string]bool{
		"mem://path/fixture/alpha/cmd":             true,
		"mem://path/fixture/alpha/cmd/app/main.go": true,
	}
	for id := range want {
		found := false
		for _, node := range got.Nodes {
			if node.ID == id {
				found = true
			}
		}
		if !found {
			t.Fatalf("missing neighbor %s in %+v", id, got.Nodes)
		}
	}
}

func TestAgentScopeReceiptRoundTrip(t *testing.T) {
	idx := fixtureAgentIndex(t)
	capsule, err := idx.Scope("main model", 8)
	if err != nil {
		t.Fatal(err)
	}
	if capsule.Authority != "none" || capsule.Receipt.Authority != "none" {
		t.Fatalf("authority inflation: %+v", capsule)
	}
	verify, err := idx.VerifyReceipt(capsule.Receipt)
	if err != nil {
		t.Fatal(err)
	}
	if !verify.OK || verify.ReceiptID != capsule.Receipt.ReceiptID {
		t.Fatalf("verify=%+v", verify)
	}
}

func TestAgentReceiptRejectsProjectHashDrift(t *testing.T) {
	idx := fixtureAgentIndex(t)
	capsule, err := idx.Scope("model", 4)
	if err != nil {
		t.Fatal(err)
	}
	capsule.Receipt.ProjectSHA256 = "tampered"
	capsule.Receipt.ReceiptID = receiptID(capsule.Receipt)
	if _, err := idx.VerifyReceipt(capsule.Receipt); err == nil {
		t.Fatal("expected project hash mismatch")
	}
}

func TestLoadRetrievalReceipt(t *testing.T) {
	idx := fixtureAgentIndex(t)
	capsule, err := idx.Scope("model", 4)
	if err != nil {
		t.Fatal(err)
	}
	dir := t.TempDir()
	path := filepath.Join(dir, "receipt.json")
	b, _ := json.Marshal(capsule.Receipt)
	if err := os.WriteFile(path, b, 0o644); err != nil {
		t.Fatal(err)
	}
	got, err := LoadRetrievalReceipt(path)
	if err != nil {
		t.Fatal(err)
	}
	if got.ReceiptID != capsule.Receipt.ReceiptID {
		t.Fatalf("receipt id=%q want=%q", got.ReceiptID, capsule.Receipt.ReceiptID)
	}
}
