package lockfile

import "io"

// PackageDetails represents a single dependency

type PackageDetails struct {
	Name      string `json:"name"`
	Version   string `json:"version"`
	Ecosystem string `json:"ecosystem"`
}

// Extractor interface based on osv-scanner pattern
type Extractor interface {
	ShouldExtract(path string) bool
	Extract(r io.Reader, path string) ([]PackageDetails, error)
	Ecosystem() string
}

// LocalFile for local filesystem
type LocalFile struct {
	io.Reader
	FilePath string
}

func (f LocalFile) Path() string {
	return f.FilePath
}
