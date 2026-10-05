// Copyright © 2026 es-3581100. ALL RIGHTS RESERVED. See LICENSE and LEGAL.md.
package matrix

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"
	"unicode"
)

const (
	AgentBootstrapSchema      = "app-dir-matrix.agent-bootstrap/v1"
	AgentResolveSchema        = "app-dir-matrix.agent-resolve/v1"
	AgentInspectSchema        = "app-dir-matrix.agent-inspect/v1"
	AgentExpandSchema         = "app-dir-matrix.agent-expand/v1"
	ContextCapsuleSchema      = "app-dir-matrix.context-capsule/v1"
	RetrievalReceiptSchema    = "app-dir-matrix.retrieval-receipt/v1"
	ReceiptVerificationSchema = "app-dir-matrix.receipt-verification/v1"
)

type AgentNode struct {
	ID        string `json:"id"`
	Kind      string `json:"kind"`
	Repo      string `json:"repo,omitempty"`
	Path      string `json:"path,omitempty"`
	SHA       string `json:"sha,omitempty"`
	Class     string `json:"class,omitempty"`
	Claim     string `json:"claim,omitempty"`
	Source    string `json:"source,omitempty"`
	Value     any    `json:"value,omitempty"`
	Authority string `json:"authority"`
}

type AgentMatch struct {
	Node    AgentNode `json:"node"`
	Score   int       `json:"score"`
	Matched []string  `json:"matched"`
}

type AgentIndex struct {
	Project       Project
	ProjectSHA256 string
	nodes         map[string]AgentNode
}

type AgentBootstrap struct {
	Schema        string   `json:"schema"`
	ProjectID     string   `json:"project_id"`
	ProjectSHA256 string   `json:"project_sha256"`
	CoreVersion   string   `json:"core_version"`
	Authority     string   `json:"authority"`
	Operations    []string `json:"operations"`
	Invariants    []string `json:"invariants"`
}

type AgentResolve struct {
	Schema        string       `json:"schema"`
	ProjectID     string       `json:"project_id"`
	ProjectSHA256 string       `json:"project_sha256"`
	Query         string       `json:"query"`
	Authority     string       `json:"authority"`
	Results       []AgentMatch `json:"results"`
}

type AgentInspect struct {
	Schema        string    `json:"schema"`
	ProjectID     string    `json:"project_id"`
	ProjectSHA256 string    `json:"project_sha256"`
	Authority     string    `json:"authority"`
	Node          AgentNode `json:"node"`
}

type AgentExpand struct {
	Schema        string      `json:"schema"`
	ProjectID     string      `json:"project_id"`
	ProjectSHA256 string      `json:"project_sha256"`
	Pointer       string      `json:"pointer"`
	Depth         int         `json:"depth"`
	Authority     string      `json:"authority"`
	Nodes         []AgentNode `json:"nodes"`
}

type ReceiptNode struct {
	ID         string `json:"id"`
	NodeSHA256 string `json:"node_sha256"`
}

type RetrievalReceipt struct {
	Schema        string        `json:"schema"`
	ReceiptID     string        `json:"receipt_id"`
	ProjectID     string        `json:"project_id"`
	ProjectSHA256 string        `json:"project_sha256"`
	Query         string        `json:"query"`
	Authority     string        `json:"authority"`
	Selected      []ReceiptNode `json:"selected"`
}

type ContextCapsule struct {
	Schema        string           `json:"schema"`
	CapsuleID     string           `json:"capsule_id"`
	ProjectID     string           `json:"project_id"`
	ProjectSHA256 string           `json:"project_sha256"`
	Query         string           `json:"query"`
	Authority     string           `json:"authority"`
	Included      []AgentNode      `json:"included"`
	ExcludedCount int              `json:"excluded_count"`
	Unresolved    []string         `json:"unresolved"`
	Receipt       RetrievalReceipt `json:"receipt"`
}

type ReceiptVerification struct {
	Schema        string `json:"schema"`
	ProjectID     string `json:"project_id"`
	ProjectSHA256 string `json:"project_sha256"`
	ReceiptID     string `json:"receipt_id"`
	Authority     string `json:"authority"`
	OK            bool   `json:"ok"`
}

func NewAgentIndex(projectPath string) (*AgentIndex, error) {
	p, err := LoadProject(projectPath)
	if err != nil {
		return nil, err
	}
	h, err := SHA256File(projectPath)
	if err != nil {
		return nil, err
	}
	idx := &AgentIndex{
		Project:       p,
		ProjectSHA256: h,
		nodes:         map[string]AgentNode{},
	}
	idx.buildNodes()
	return idx, nil
}

func (a *AgentIndex) buildNodes() {
	for _, r := range a.Project.Repositories {
		repoID := repoPointer(r.Repo)
		a.nodes[repoID] = AgentNode{
			ID: repoID, Kind: "repository", Repo: r.Repo, SHA: r.TreeSHA,
			Claim: r.Description, Source: r.Branch, Value: map[string]any{
				"deployment": r.Deployment, "source_class": r.SourceClass, "boundary": r.Boundary,
			}, Authority: "none",
		}
		for _, e := range r.Tree {
			id := pathPointer(r.Repo, e.Path)
			a.nodes[id] = AgentNode{
				ID: id, Kind: "path", Repo: r.Repo, Path: e.Path, SHA: e.SHA,
				Value: map[string]any{"type": e.Type, "size": e.Size, "level": LevelForPath(e.Path)},
				Authority: "none",
			}
		}
	}
	for _, e := range a.Project.Evidence {
		id := evidencePointer(e.ID)
		a.nodes[id] = AgentNode{
			ID: id, Kind: "evidence", Repo: e.Repo, Path: e.Path, Class: e.Class,
			Claim: e.Claim, Source: e.Source, Value: e.Value, Authority: "none",
		}
	}
}

func repoPointer(repo string) string { return "mem://repo/" + repo }
func pathPointer(repo, path string) string { return "mem://path/" + repo + "/" + strings.TrimPrefix(path, "/") }
func evidencePointer(id string) string { return "mem://evidence/" + id }

func (a *AgentIndex) Bootstrap() AgentBootstrap {
	return AgentBootstrap{
		Schema: AgentBootstrapSchema, ProjectID: a.Project.ProjectID, ProjectSHA256: a.ProjectSHA256,
		CoreVersion: a.Project.CoreVersion, Authority: "none",
		Operations: []string{"bootstrap", "resolve", "inspect", "expand", "scope", "verify-receipt"},
		Invariants: []string{
			"context-is-not-permission",
			"evidence-authority-is-none",
			"unknown-pointers-fail-closed",
			"no-automatic-scope-widening",
		},
	}
}

func (a *AgentIndex) Resolve(query string, limit int) (AgentResolve, error) {
	query = strings.TrimSpace(query)
	if query == "" {
		return AgentResolve{}, fmt.Errorf("query is required")
	}
	if limit < 1 || limit > 32 {
		return AgentResolve{}, fmt.Errorf("limit must be between 1 and 32")
	}
	tokens := tokenize(query)
	if len(tokens) == 0 {
		return AgentResolve{}, fmt.Errorf("query has no searchable tokens")
	}
	matches := make([]AgentMatch, 0)
	for _, node := range a.sortedNodes() {
		nodeTokens := tokenize(nodeSearchText(node))
		matched := intersectTokens(tokens, nodeTokens)
		if len(matched) == 0 {
			continue
		}
		matches = append(matches, AgentMatch{Node: node, Score: len(matched), Matched: matched})
	}
	sort.Slice(matches, func(i, j int) bool {
		if matches[i].Score != matches[j].Score {
			return matches[i].Score > matches[j].Score
		}
		return matches[i].Node.ID < matches[j].Node.ID
	})
	if len(matches) > limit {
		matches = matches[:limit]
	}
	return AgentResolve{
		Schema: AgentResolveSchema, ProjectID: a.Project.ProjectID, ProjectSHA256: a.ProjectSHA256,
		Query: query, Authority: "none", Results: matches,
	}, nil
}

func (a *AgentIndex) Inspect(pointer string) (AgentInspect, error) {
	node, ok := a.nodes[pointer]
	if !ok {
		return AgentInspect{}, fmt.Errorf("unknown pointer %q", pointer)
	}
	return AgentInspect{
		Schema: AgentInspectSchema, ProjectID: a.Project.ProjectID, ProjectSHA256: a.ProjectSHA256,
		Authority: "none", Node: node,
	}, nil
}

func (a *AgentIndex) Expand(pointer string, depth, limit int) (AgentExpand, error) {
	if depth < 1 || depth > 2 {
		return AgentExpand{}, fmt.Errorf("depth must be 1 or 2")
	}
	if limit < 1 || limit > 32 {
		return AgentExpand{}, fmt.Errorf("limit must be between 1 and 32")
	}
	if _, ok := a.nodes[pointer]; !ok {
		return AgentExpand{}, fmt.Errorf("unknown pointer %q", pointer)
	}

	seen := map[string]bool{pointer: true}
	frontier := []string{pointer}
	out := make([]AgentNode, 0, limit)
	for level := 0; level < depth && len(frontier) > 0 && len(out) < limit; level++ {
		next := make([]string, 0)
		for _, current := range frontier {
			for _, id := range a.neighborIDs(current) {
				if seen[id] {
					continue
				}
				seen[id] = true
				next = append(next, id)
			}
		}
		sort.Strings(next)
		for _, id := range next {
			if len(out) >= limit {
				break
			}
			out = append(out, a.nodes[id])
		}
		frontier = next
	}

	return AgentExpand{
		Schema: AgentExpandSchema, ProjectID: a.Project.ProjectID, ProjectSHA256: a.ProjectSHA256,
		Pointer: pointer, Depth: depth, Authority: "none", Nodes: out,
	}, nil
}

func (a *AgentIndex) Scope(query string, limit int) (ContextCapsule, error) {
	resolved, err := a.Resolve(query, limit)
	if err != nil {
		return ContextCapsule{}, err
	}
	included := make([]AgentNode, 0, len(resolved.Results))
	selected := make([]ReceiptNode, 0, len(resolved.Results))
	ids := make([]string, 0, len(resolved.Results))
	for _, match := range resolved.Results {
		included = append(included, match.Node)
		h, err := nodeHash(match.Node)
		if err != nil {
			return ContextCapsule{}, err
		}
		selected = append(selected, ReceiptNode{ID: match.Node.ID, NodeSHA256: h})
		ids = append(ids, match.Node.ID)
	}
	receipt := RetrievalReceipt{
		Schema: RetrievalReceiptSchema, ProjectID: a.Project.ProjectID, ProjectSHA256: a.ProjectSHA256,
		Query: resolved.Query, Authority: "none", Selected: selected,
	}
	receipt.ReceiptID = receiptID(receipt)
	capsuleID := SHA256Bytes([]byte(a.ProjectSHA256 + "\n" + resolved.Query + "\n" + strings.Join(ids, "\n")))
	return ContextCapsule{
		Schema: ContextCapsuleSchema, CapsuleID: capsuleID, ProjectID: a.Project.ProjectID,
		ProjectSHA256: a.ProjectSHA256, Query: resolved.Query, Authority: "none",
		Included: included, ExcludedCount: len(a.nodes) - len(included), Unresolved: []string{},
		Receipt: receipt,
	}, nil
}

func LoadRetrievalReceipt(path string) (RetrievalReceipt, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return RetrievalReceipt{}, err
	}
	var r RetrievalReceipt
	if err := json.Unmarshal(b, &r); err != nil {
		return RetrievalReceipt{}, fmt.Errorf("decode receipt: %w", err)
	}
	return r, nil
}

func (a *AgentIndex) VerifyReceipt(r RetrievalReceipt) (ReceiptVerification, error) {
	if r.Schema != RetrievalReceiptSchema {
		return ReceiptVerification{}, fmt.Errorf("receipt schema %q != %q", r.Schema, RetrievalReceiptSchema)
	}
	if r.Authority != "none" {
		return ReceiptVerification{}, fmt.Errorf("receipt attempts authority inflation: %q", r.Authority)
	}
	if r.ProjectID != a.Project.ProjectID || r.ProjectSHA256 != a.ProjectSHA256 {
		return ReceiptVerification{}, fmt.Errorf("receipt project identity/hash mismatch")
	}
	if r.ReceiptID == "" || r.ReceiptID != receiptID(r) {
		return ReceiptVerification{}, fmt.Errorf("receipt id mismatch")
	}
	seen := map[string]bool{}
	for _, selected := range r.Selected {
		if seen[selected.ID] {
			return ReceiptVerification{}, fmt.Errorf("duplicate receipt node %q", selected.ID)
		}
		seen[selected.ID] = true
		node, ok := a.nodes[selected.ID]
		if !ok {
			return ReceiptVerification{}, fmt.Errorf("receipt references unknown node %q", selected.ID)
		}
		h, err := nodeHash(node)
		if err != nil {
			return ReceiptVerification{}, err
		}
		if h != selected.NodeSHA256 {
			return ReceiptVerification{}, fmt.Errorf("receipt node hash mismatch for %q", selected.ID)
		}
	}
	return ReceiptVerification{
		Schema: ReceiptVerificationSchema, ProjectID: a.Project.ProjectID, ProjectSHA256: a.ProjectSHA256,
		ReceiptID: r.ReceiptID, Authority: "none", OK: true,
	}, nil
}

func (a *AgentIndex) sortedNodes() []AgentNode {
	ids := make([]string, 0, len(a.nodes))
	for id := range a.nodes {
		ids = append(ids, id)
	}
	sort.Strings(ids)
	out := make([]AgentNode, 0, len(ids))
	for _, id := range ids {
		out = append(out, a.nodes[id])
	}
	return out
}

func (a *AgentIndex) neighborIDs(pointer string) []string {
	node := a.nodes[pointer]
	set := map[string]bool{}
	switch node.Kind {
	case "repository":
		for id, candidate := range a.nodes {
			if candidate.Kind == "path" && candidate.Repo == node.Repo && !strings.Contains(candidate.Path, "/") {
				set[id] = true
			}
			if candidate.Kind == "evidence" && candidate.Repo == node.Repo && candidate.Path == "" {
				set[id] = true
			}
		}
	case "path":
		parent := parentPath(node.Path)
		if parent == "" {
			set[repoPointer(node.Repo)] = true
		} else if _, ok := a.nodes[pathPointer(node.Repo, parent)]; ok {
			set[pathPointer(node.Repo, parent)] = true
		}
		for id, candidate := range a.nodes {
			if candidate.Kind == "path" && candidate.Repo == node.Repo && parentPath(candidate.Path) == node.Path {
				set[id] = true
			}
			if candidate.Kind == "evidence" && candidate.Repo == node.Repo && candidate.Path == node.Path {
				set[id] = true
			}
		}
	case "evidence":
		if node.Repo != "" && node.Path != "" {
			if id := pathPointer(node.Repo, node.Path); a.hasNode(id) {
				set[id] = true
			}
		} else if node.Repo != "" {
			if id := repoPointer(node.Repo); a.hasNode(id) {
				set[id] = true
			}
		}
	}
	ids := make([]string, 0, len(set))
	for id := range set {
		if a.hasNode(id) {
			ids = append(ids, id)
		}
	}
	sort.Strings(ids)
	return ids
}

func (a *AgentIndex) hasNode(id string) bool {
	_, ok := a.nodes[id]
	return ok
}

func parentPath(path string) string {
	path = strings.Trim(path, "/")
	i := strings.LastIndex(path, "/")
	if i < 0 {
		return ""
	}
	return path[:i]
}

func nodeSearchText(node AgentNode) string {
	b, _ := json.Marshal(node.Value)
	return strings.Join([]string{
		node.ID, node.Kind, node.Repo, node.Path, node.SHA, node.Class,
		node.Claim, node.Source, string(b),
	}, " ")
}

func tokenize(s string) []string {
	normalized := strings.Map(func(r rune) rune {
		if unicode.IsLetter(r) || unicode.IsNumber(r) {
			return unicode.ToLower(r)
		}
		return ' '
	}, s)
	seen := map[string]bool{}
	out := make([]string, 0)
	for _, token := range strings.Fields(normalized) {
		if seen[token] {
			continue
		}
		seen[token] = true
		out = append(out, token)
	}
	sort.Strings(out)
	return out
}

func intersectTokens(query, candidate []string) []string {
	set := map[string]bool{}
	for _, token := range candidate {
		set[token] = true
	}
	out := make([]string, 0)
	for _, token := range query {
		if set[token] {
			out = append(out, token)
		}
	}
	sort.Strings(out)
	return out
}

func nodeHash(node AgentNode) (string, error) {
	b, err := json.Marshal(node)
	if err != nil {
		return "", err
	}
	return SHA256Bytes(b), nil
}

func receiptID(r RetrievalReceipt) string {
	copy := r
	copy.ReceiptID = ""
	b, _ := json.Marshal(copy)
	return SHA256Bytes(b)
}
