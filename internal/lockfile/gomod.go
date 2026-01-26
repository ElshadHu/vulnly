package lockfile

import (
	"bufio"
	"io"
	"path/filepath"
	"strings"
)

type GoModExtractor struct{}

func (e GoModExtractor) Ecosystem() string {
	return "Go"
}

func (e GoModExtractor) ShouldExtract(path string) bool {
	return filepath.Base(path) == "go.sum"
}

func (e GoModExtractor) Extract(r io.Reader, path string) ([]PackageDetails, error) {
	var packages []PackageDetails
	seen := make(map[string]bool)

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		parts := strings.Fields(line)
		if len(parts) < 2 {
			continue
		}

		name := parts[0]
		version := parts[1]

		// Skip /go.mod entries
		if strings.HasSuffix(version, "/go.mod") {
			continue
		}

		version = strings.TrimSuffix(version, "+incompatible")

		// Deduplicate
		key := name + "@" + version
		if seen[key] {
			continue
		}
		seen[key] = true

		packages = append(packages, PackageDetails{
			Name:      name,
			Version:   version,
			Ecosystem: e.Ecosystem(),
		})
	}

	return packages, scanner.Err()
}
