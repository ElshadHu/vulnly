package lockfile

import (
	"os"
	"path/filepath"
)

var extractors = []Extractor{
	NpmExtractor{},
	GoModExtractor{},
	PipExtractor{},
}

// FindLockfiles walks the directory and finds all supported lockfiles

func FindLockfiles(root string) ([]string, error) {
	var files []string

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			// Skip node_modules and vendor directories
			if info.Name() == "node_modules" || info.Name() == "vendor" || info.Name() == ".git" {
				return filepath.SkipDir
			}
			return nil
		}

		for _, e := range extractors {
			if e.ShouldExtract(path) {
				files = append(files, path)
				break
			}
		}

		return nil
	})

	return files, err
}

func ExtractAll(root string) ([]PackageDetails, error) {
	files, err := FindLockfiles(root)
	if err != nil {
		return nil, err
	}

	var allPackages []PackageDetails
	for _, path := range files {
		packages, err := extractFromFile(path)
		if err != nil {
			return nil, err
		}
		allPackages = append(allPackages, packages...)
	}

	return allPackages, nil
}

func extractFromFile(path string) ([]PackageDetails, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer func() { _ = file.Close() }()

	for _, e := range extractors {
		if e.ShouldExtract(path) {
			return e.Extract(file, path)
		}
	}

	return nil, nil
}
