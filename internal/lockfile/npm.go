package lockfile

import (
	"encoding/json"
	"fmt"
	"io"
	"path/filepath"
	"strings"
)

type NpmExtractor struct{}

func (e NpmExtractor) Ecosystem() string {
	return "npm"
}

func (e NpmExtractor) ShouldExtract(path string) bool {
	return filepath.Base(path) == "package-lock.json"
}

type npmLockFile struct {
	Packages map[string]npmPackage `json:"packages"`
}
type npmPackage struct {
	Version string `json:"version"`
}

func (e NpmExtractor) Extract(r io.Reader, path string) ([]PackageDetails, error) {
	var lockfile npmLockFile

	decoder := json.NewDecoder(r)
	if err := decoder.Decode(&lockfile); err != nil {
		return nil, fmt.Errorf("failed to decode %s: %w", path, err)
	}

	var packages []PackageDetails
	for pkgPath, pkg := range lockfile.Packages {
		if pkgPath == "" || pkg.Version == "" {
			continue
		}
		name := strings.TrimPrefix(pkgPath, "node_modules/")
		if strings.Contains(name, "node_modules/") {
			parts := strings.Split(name, "node_modules/")
			name = parts[len(parts)-1]
		}

		packages = append(packages, PackageDetails{
			Name:      name,
			Version:   pkg.Version,
			Ecosystem: e.Ecosystem(),
		})

	}
	return packages, nil
}
