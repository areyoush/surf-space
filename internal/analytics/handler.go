package analytics

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/areyoush/surfspace/internal/links"
)

type Handler struct {
	svc      		*Service
	linksSvc 		*links.Service		
}

func NewHandler(svc *Service, linksSvc *links.Service) *Handler {
	return &Handler{svc: svc, linksSvc: linksSvc}
}

func (h *Handler) GetAnalytics(c *gin.Context) {
	shortCode := c.Param("code")
	userID := c.MustGet("userID").(string)

	link, err := h.linksSvc.GetLink(shortCode)
	if err != nil || link == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "link not found"})
		return
	}

	if link.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "you don't own this link"})
		return
	}

	clicks, err := h.svc.GetAnalytics(link.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not fetch analytics"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"short_code":  shortCode,
		"total_clicks": len(clicks),
		"clicks":      clicks,
	})
}

func RegisterRoutes(rg *gin.RouterGroup, h *Handler) {
	rg.GET("/links/:code/analytics", h.GetAnalytics)
}