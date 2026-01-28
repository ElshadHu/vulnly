package lockfile

import (
	"bufio"
	"io"
	"path/filepath"
	"regexp"
	"strings"
)

type PipExtractor struct{}

func (e PipExtractor) Ecosystem() string {
	return "PyPI"
}

func (e PipExtractor) ShouldExtract(path string) bool {
	base := filepath.Base(path)
	return base == "requirements.txt" || base == "requirements-dev.txt"
}

// package==1.0.0, package>=1.0.0, package~=1.0.0
var pipVersionRegex = regexp.MustCompile(`^([a-zA-Z0-9_-]+)\s*([=~<>!]+)\s*([0-9][a-zA-Z0-9._-]*)`)

func (e PipExtractor) Extract(r io.Reader, _ string) ([]PackageDetails, error) {
	var packages []PackageDetails

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		if strings.HasPrefix(line, "-") {
			continue
		}

		matches := pipVersionRegex.FindStringSubmatch(line)
		if len(matches) >= 4 {
			packages = append(packages, PackageDetails{
				Name:      strings.ToLower(matches[1]),
				Version:   matches[3],
				Ecosystem: e.Ecosystem(),
			})
		}
	}

	return packages, scanner.Err()
}
