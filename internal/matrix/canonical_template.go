// Copyright © 2026 es-3581100. ALL RIGHTS RESERVED. See LICENSE and LEGAL.md.
package matrix

import _ "embed"

// canonicalArtifactTemplate is an exact embedded copy of
// web/app-dir-matrix-core-v0.2.0.html. Tests bind the two byte-for-byte.
//
//go:embed app-dir-matrix-core-v0.2.0.html
var canonicalArtifactTemplate []byte

// BuildCanonicalArtifact renders the project with the core-owned canonical UI
// template. It performs no network access and does not mutate project truth.
func BuildCanonicalArtifact(project Project) ([]byte, error) {
	if err := project.Validate(); err != nil {
		return nil, err
	}
	return buildArtifactValidated(canonicalArtifactTemplate, project)
}

// WriteCanonicalArtifact writes the canonical HTML artifact and the same SHA-256
// sidecar format used by WriteArtifact.
func WriteCanonicalArtifact(project Project, outPath string) (string, error) {
	b, err := BuildCanonicalArtifact(project)
	if err != nil {
		return "", err
	}
	return writeArtifactBytes(b, outPath)
}
