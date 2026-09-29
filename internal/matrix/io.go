// Copyright © 2026 es-3581100. ALL RIGHTS RESERVED. See LICENSE and LEGAL.md.
package matrix

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
)

func LoadProject(path string) (Project, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return Project{}, err
	}
	var p Project
	if err := json.Unmarshal(b, &p); err != nil {
		return Project{}, fmt.Errorf("decode project: %w", err)
	}
	if err := p.Validate(); err != nil {
		return Project{}, err
	}
	return p, nil
}

func SHA256Bytes(b []byte) string {
	s := sha256.Sum256(b)
	return hex.EncodeToString(s[:])
}

func SHA256File(path string) (string, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return SHA256Bytes(b), nil
}
