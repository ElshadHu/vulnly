package osv

type PackageVuln struct {
	Name       string
	Version    string
	Ecosystem  string
	Vuln       *Vulnerability
	Severity   Severity
	FixVersion string
}

type Summary struct {
	TotalDeps int
	Critical  int
	High      int
	Medium    int
	Low       int
	Unknown   int
}

type ScanResult struct {
	Packages []PackageVuln
	Summary  Summary
}
