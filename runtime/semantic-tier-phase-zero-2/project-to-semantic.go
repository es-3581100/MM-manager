// Copyright © 2026 es-3581100. ALL RIGHTS RESERVED. See LICENSE and LEGAL.md in MM-manager.
package main

import (
    "encoding/json"
    "flag"
    "fmt"
    "os"
    "sort"
)

const projectSchema = "app-dir-matrix.project/v1"
const descriptorSchema = "memory-matrix.semantic-descriptor/v1"
const descriptorSetSchema = "memory-matrix.semantic-descriptor-set/v1"

type project struct {
    SchemaVersion string `json:"schema_version"`
    CoreVersion string `json:"core_version"`
    ProjectID string `json:"project_id"`
    Repositories []repository `json:"repositories"`
}
type repository struct {
    Repo string `json:"repo"`
    Branch string `json:"branch"`
    TreeSHA string `json:"tree_sha"`
    SourceClass string `json:"source_class"`
    Tree []treeEntry `json:"tree"`
    Truncated bool `json:"truncated"`
}
type treeEntry struct {
    Path string `json:"path"`
    Type string `json:"type"`
    Size *int64 `json:"size"`
    SHA string `json:"sha"`
}
type provenance struct {
    Producer string `json:"producer"`
    SemanticType string `json:"semantic_type"`
    Source string `json:"source"`
    ObservedAtNanos int64 `json:"observed_at_nanos"`
    Detail string `json:"detail"`
    Authority string `json:"authority"`
    SourceLocator string `json:"source_locator"`
}
type descriptor struct {
    SchemaVersion string `json:"schema_version"`
    ObjectID string `json:"object_id"`
    Kind string `json:"kind"`
    Dependencies []string `json:"dependencies"`
    Mutability string `json:"mutability"`
    Recomputability string `json:"recomputability"`
    Confidence float64 `json:"confidence"`
    Provenance []provenance `json:"provenance"`
}
type descriptorSet struct {
    SchemaVersion string `json:"schema_version"`
    SourceProjectSchema string `json:"source_project_schema"`
    SourceProjectID string `json:"source_project_id"`
    Adapter string `json:"adapter"`
    InferenceBoundary []string `json:"inference_boundary"`
    Descriptors []descriptor `json:"descriptors"`
}

func main() {
    in := flag.String("in", "", "app-dir-matrix.project/v1 JSON")
    out := flag.String("out", "", "semantic descriptor set JSON")
    flag.Parse()
    if *in == "" || *out == "" { fatal("--in and --out are required") }
    raw, err := os.ReadFile(*in)
    if err != nil { fatal(err.Error()) }
    var p project
    if err := json.Unmarshal(raw, &p); err != nil { fatal("parse project: " + err.Error()) }
    if p.SchemaVersion != projectSchema {
        fatal(fmt.Sprintf("unsupported project schema %q; expected %q", p.SchemaVersion, projectSchema))
    }

    set := descriptorSet{
        SchemaVersion: descriptorSetSchema,
        SourceProjectSchema: p.SchemaVersion,
        SourceProjectID: p.ProjectID,
        Adapter: "mm-manager.project-to-semantic/v1",
        InferenceBoundary: []string{
            "tree membership and blob/tree identity are structural evidence only",
            "path adjacency does not create semantic dependencies",
            "mutability and recomputability remain UNKNOWN without explicit evidence",
            "all imported evidence remains authority=none",
        },
    }
    for _, repo := range p.Repositories {
        if repo.Truncated { fatal("refuse truncated repository tree: " + repo.Repo) }
        for _, entry := range repo.Tree {
            kind := "UNKNOWN"
            switch entry.Type { case "blob": kind = "FILE"; case "tree": kind = "DIRECTORY" }
            set.Descriptors = append(set.Descriptors, descriptor{
                SchemaVersion: descriptorSchema,
                ObjectID: repo.Repo + ":" + entry.Path,
                Kind: kind,
                Dependencies: []string{},
                Mutability: "UNKNOWN",
                Recomputability: "UNKNOWN",
                Confidence: 1.0,
                Provenance: []provenance{{
                    Producer: "mm-manager.project-to-semantic/v1",
                    SemanticType: "PROGRAM",
                    Source: "STRUCTURAL_GRAPH",
                    ObservedAtNanos: 0,
                    Detail: fmt.Sprintf("pinned tree entry type=%s entry_sha=%s tree_sha=%s", entry.Type, entry.SHA, repo.TreeSHA),
                    Authority: "none",
                    SourceLocator: fmt.Sprintf("github:%s@%s:%s#%s", repo.Repo, repo.Branch, entry.Path, entry.SHA),
                }},
            })
        }
    }
    sort.Slice(set.Descriptors, func(i, j int) bool { return set.Descriptors[i].ObjectID < set.Descriptors[j].ObjectID })
    enc, err := json.MarshalIndent(set, "", "  ")
    if err != nil { fatal(err.Error()) }
    enc = append(enc, '\n')
    if err := os.WriteFile(*out, enc, 0644); err != nil { fatal(err.Error()) }
}
func fatal(msg string) { fmt.Fprintln(os.Stderr, "error:", msg); os.Exit(1) }
