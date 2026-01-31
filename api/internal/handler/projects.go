package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *API) ListProjects(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		userID = "anonymous"
	}
	projects, err := h.repo.ListProjectsByUser(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list projects"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"projects": projects})
}

func (h *API) GetProject(c *gin.Context) {
	projectID := c.Param("project_id")
	scans, err := h.repo.ListScansByProject(c.Request.Context(), projectID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get project"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"project_id": projectID,
		"scans":      scans,
	})
}

func (h *API) ListScans(c *gin.Context) {
	projectID := c.Param("project_id")

	scans, err := h.repo.ListScansByProject(c.Request.Context(), projectID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list scans"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"project_id": projectID,
		"scans":      scans,
	})
}
