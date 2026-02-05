package handler

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/ElshadHu/vulnly/api/internal/repository"
	"github.com/gin-gonic/gin"
)

// DependencyInput represents a package dependency from the scan
type DependencyInput struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

// VulnerabilityInput represents a vulnerability found during scan
type VulnerabilityInput struct {
	ID           string `json:"id"`
	Package      string `json:"package"`
	Version      string `json:"version"`
	Severity     string `json:"severity"`
	FixedVersion string `json:"fixed_version"`
	Description  string `json:"description"`
}

// IngestResponse is returned after successful scan ingestion
type IngestResponse struct {
	ScanID    string    `json:"scan_id"`
	CreatedAt time.Time `json:"created_at"`
}

// IngestRequest is the payload sent by CLI after scanning
type IngestRequest struct {
	Project         string                 `json:"project" binding:"required"`
	Commit          string                 `json:"commit"`
	Branch          string                 `json:"branch"`
	Ecosystem       string                 `json:"ecosystem" binding:"required"`
	Dependencies    []DependencyInput      `json:"dependencies"`
	Vulnerabilities []VulnerabilityInput   `json:"vulnerabilities"`
	Summary         repository.VulnSummary `json:"summary"`
}

func (h *API) Ingest(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		userID = "anonymous"
	}
	var req IngestRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("ingest bind error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	project, err := h.repo.GetOrCreateProject(c.Request.Context(), userID, req.Project)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create a project"})
		return
	}
	// scanID uses nanosecond timestamp to ensure uniqueness within a project
	scanID := fmt.Sprintf("scan_%d", time.Now().UnixNano())
	now := time.Now().UTC()

	scan := &repository.Scan{
		ProjectID: project.ProjectID,
		ScanID:    scanID,
		Commit:    req.Commit,
		Branch:    req.Branch,
		Ecosystem: req.Ecosystem,
		TotalDeps: len(req.Dependencies),
		Summary:   req.Summary,
		CreatedAt: now,
	}
	if err := h.repo.CreateScan(c.Request.Context(), scan); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create scan"})
		return
	}
	var vulns []repository.Vulnerability
	for _, v := range req.Vulnerabilities {
		vulns = append(vulns, repository.Vulnerability{
			ScanID:         scanID,
			VulnID:         v.ID,
			PackageName:    v.Package,
			PackageVersion: v.Version,
			Severity:       v.Severity,
			FixedVersion:   v.FixedVersion,
			Description:    v.Description,
		})
	}

	// CreateVulnerabilities handles chunking into 25-item batches (DynamoDB limit)
	if err := h.repo.CreateVulnerabilities(c.Request.Context(), vulns); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create vulnerabilities"})
		return
	}

	c.JSON(http.StatusCreated, IngestResponse{
		ScanID:    scanID,
		CreatedAt: now,
	})
}
