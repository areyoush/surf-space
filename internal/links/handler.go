package links

import (
	"net/http"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

type ShortenInput struct {
	OriginalURL string `json:"original_url"`
}

func (h *Handler) ShortenURL(c *gin.Context) {
	userID, _ := c.Get("userID")

	var input ShortenInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	link, err := h.svc.ShortenURL(userID.(string), input.OriginalURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not shorten URL"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"short_code":   link.ShortCode,
		"original_url": link.OriginalURL,
		"short_url":    "http://localhost:8080/" + link.ShortCode,
	})
}

func (h *Handler) Redirect(c *gin.Context) {
	shortCode := c.Param("code")

	link, err := h.svc.GetLink(shortCode)
	if err != nil || link == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "link not found"})
		return
	}

	go h.svc.repo.IncrementClickCount(shortCode)

	c.Redirect(http.StatusMovedPermanently, link.OriginalURL)
}

func RegisterRoutes(rg *gin.RouterGroup, h *Handler, publicRg *gin.RouterGroup) {
	rg.POST("/links", h.ShortenURL)
	publicRg.GET("/:code", h.Redirect)
}