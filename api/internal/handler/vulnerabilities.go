package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

const (
	DefaultVulnLimit = 100
	MaxVulnLimit     = 500
)

var validSeverities = map[string]bool{
	"CRITICAL": true,
	"HIGH":     true,
	"MEDIUM":   true,
	"LOW":      true,
}

// ListVulnerabilities returns vulnerabilities with optional filters
// GET /api/vulnerabilities?scanId=xxx&severity=HIGH&package=lodash&limit=50
func (h *API) ListVulnerabilities(c *gin.Context) {
	scanID := c.Query("scanId")
	if scanID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "scanId is required"})
		return
	}
	userID := c.GetString("user_id")

	// Get scan to find its project
	scan, err := h.repo.GetScanByID(c.Request.Context(), scanID)
	if err != nil || scan == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "scan not found"})
		return
	}

	// Verify project belongs to user
	project, err := h.repo.GetProjectByID(c.Request.Context(), userID, scan.ProjectID)
	if err != nil || project == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "scan not found"})
		return
	}
	severity := c.Query("severity")
	if severity != "" && !validSeverities[severity] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid severity, must be CRITICAL, HIGH, MEDIUM, or LOW"})
		return
	}

	packageName := c.Query("package")
	limit := DefaultVulnLimit
	if limitStr := c.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= MaxVulnLimit {
			limit = l
		}
	}

	vulns, err := h.repo.ListVulnerabilitiesByScan(
		c.Request.Context(),
		scanID,
		severity,
		packageName,
		limit,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to query vulnerabilities"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"vulnerabilities": vulns,
		"count":           len(vulns),
	})
}
