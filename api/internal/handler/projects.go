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
	userID := c.GetString("user_id")
	projectName := c.Param("project_id")

	// First, find the project by name to get its UUID
	project, err := h.repo.GetProjectByName(c.Request.Context(), userID, projectName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get project"})
		return
	}
	if project == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "project not found"})
		return
	}

	// Query scans using the actual project UUID
	scans, err := h.repo.ListScansByProject(c.Request.Context(), project.ProjectID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get scans"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"project_id": project.ProjectID,
		"name":       project.Name,
		"scans":      scans,
	})
}

func (h *API) ListScans(c *gin.Context) {
	userID := c.GetString("user_id")
	projectName := c.Param("project_id")

	// Find project by name to get UUID
	project, err := h.repo.GetProjectByName(c.Request.Context(), userID, projectName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get project"})
		return
	}
	if project == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "project not found"})
		return
	}

	scans, err := h.repo.ListScansByProject(c.Request.Context(), project.ProjectID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list scans"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"project_id": project.ProjectID,
		"name":       project.Name,
		"scans":      scans,
	})
}
