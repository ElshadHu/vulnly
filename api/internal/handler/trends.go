package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

const (
	DefaultTrendDays = 30
	MaxTrendDays     = 90
)

// TrendDataPoint represents a single data point for trend charts
type TrendDataPoint struct {
	Date     string `json:"date"`
	Critical int    `json:"critical"`
	High     int    `json:"high"`
	Medium   int    `json:"medium"`
	Low      int    `json:"low"`
	Total    int    `json:"total"`
	ScanID   string `json:"scanId"`
}

// GetTrends returns vulnerability trend data for charts
// GET /api/trends?projectId=xxx&days=30
func (h *API) GetTrends(c *gin.Context) {
	projectID := c.Query("projectId")
	if projectID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "projectId is required"})
		return
	}

	// Get user ID from context
	userID := c.GetString("user_id")
	// Verify project belongs to user
	project, err := h.repo.GetProjectByName(c.Request.Context(), userID, projectID)
	if err != nil || project == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "project not found"})
		return
	}
	days := DefaultTrendDays
	if daysStr := c.Query("days"); daysStr != "" {
		if d, err := strconv.Atoi(daysStr); err == nil && d > 0 && d <= MaxTrendDays {
			days = d
		}
	}

	scans, err := h.repo.GetRecentScans(c.Request.Context(), projectID, days)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get trend data"})
		return
	}

	dataPoints := make([]TrendDataPoint, len(scans))
	for i, scan := range scans {
		dataPoints[i] = TrendDataPoint{
			Date:     scan.CreatedAt.Format("2006-01-02"),
			Critical: scan.Summary.Critical,
			High:     scan.Summary.High,
			Medium:   scan.Summary.Medium,
			Low:      scan.Summary.Low,
			Total:    scan.Summary.Critical + scan.Summary.High + scan.Summary.Medium + scan.Summary.Low,
			ScanID:   scan.ScanID,
		}
	}

	c.JSON(http.StatusOK, gin.H{"dataPoints": dataPoints})
}
