// Copyright © 2026 es-3581100. ALL RIGHTS RESERVED. See LICENSE and LEGAL.md.
package matrix

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

const GitHubFixtureSchemaVersion = "app-dir-matrix.github-fixture/v1"

type GitHubFixture struct {
	Schema        string                 `json:"schema"`
	CapturedAt    string                 `json:"captured_at"`
	Source        string                 `json:"source"`
	Repository    string                 `json:"repository"`
	CaptureMethod string                 `json:"capture_method"`
	Requests      []GitHubFixtureRequest `json:"requests"`
}

type GitHubFixtureRequest struct {
	Method       string `json:"method"`
	RequestURI   string `json:"request_uri"`
	Status       int    `json:"status"`
	ResponseFile string `json:"response_file"`
	SHA256       string `json:"sha256"`
}

type fixtureResponse struct {
	status int
	body   []byte
}

type fixtureTransport struct {
	responses map[string]fixtureResponse
}

func NewGitHubFixtureClient(dir string) (*http.Client, GitHubFixture, error) {
	fixturePath := filepath.Join(dir, "capture.json")
	raw, err := os.ReadFile(fixturePath)
	if err != nil {
		return nil, GitHubFixture{}, fmt.Errorf("read GitHub fixture manifest: %w", err)
	}
	var fixture GitHubFixture
	if err := json.Unmarshal(raw, &fixture); err != nil {
		return nil, GitHubFixture{}, fmt.Errorf("parse GitHub fixture manifest: %w", err)
	}
	if fixture.Schema != GitHubFixtureSchemaVersion {
		return nil, GitHubFixture{}, fmt.Errorf("GitHub fixture schema %q != %q", fixture.Schema, GitHubFixtureSchemaVersion)
	}
	if len(fixture.Requests) == 0 {
		return nil, GitHubFixture{}, fmt.Errorf("GitHub fixture has no recorded requests")
	}

	responses := make(map[string]fixtureResponse, len(fixture.Requests))
	for _, item := range fixture.Requests {
		method := strings.ToUpper(strings.TrimSpace(item.Method))
		if method == "" || item.RequestURI == "" || item.ResponseFile == "" || item.Status == 0 || item.SHA256 == "" {
			return nil, GitHubFixture{}, fmt.Errorf("invalid GitHub fixture request: %+v", item)
		}
		full := filepath.Join(dir, filepath.FromSlash(item.ResponseFile))
		cleanRoot, err := filepath.Abs(dir)
		if err != nil {
			return nil, GitHubFixture{}, err
		}
		cleanFull, err := filepath.Abs(full)
		if err != nil {
			return nil, GitHubFixture{}, err
		}
		prefix := cleanRoot + string(os.PathSeparator)
		if cleanFull != cleanRoot && !strings.HasPrefix(cleanFull, prefix) {
			return nil, GitHubFixture{}, fmt.Errorf("fixture response escapes root: %s", item.ResponseFile)
		}
		body, err := os.ReadFile(cleanFull)
		if err != nil {
			return nil, GitHubFixture{}, fmt.Errorf("read fixture response %s: %w", item.ResponseFile, err)
		}
		sum := sha256.Sum256(body)
		got := hex.EncodeToString(sum[:])
		if !strings.EqualFold(got, item.SHA256) {
			return nil, GitHubFixture{}, fmt.Errorf("fixture response hash mismatch for %s: got %s want %s", item.ResponseFile, got, item.SHA256)
		}
		key := method + " " + item.RequestURI
		if _, exists := responses[key]; exists {
			return nil, GitHubFixture{}, fmt.Errorf("duplicate GitHub fixture request %s", key)
		}
		responses[key] = fixtureResponse{status: item.Status, body: body}
	}

	return &http.Client{Transport: fixtureTransport{responses: responses}}, fixture, nil
}

func (t fixtureTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	key := strings.ToUpper(req.Method) + " " + req.URL.RequestURI()
	hit, ok := t.responses[key]
	if !ok {
		return nil, fmt.Errorf("offline GitHub fixture has no recorded response for %s", key)
	}
	return &http.Response{
		StatusCode: hit.status,
		Status:     fmt.Sprintf("%d %s", hit.status, http.StatusText(hit.status)),
		Header:     http.Header{"Content-Type": []string{"application/json"}},
		Body:       io.NopCloser(bytes.NewReader(hit.body)),
		Request:    req,
	}, nil
}
