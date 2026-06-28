package links

import (

	"time"
	"log"
	
	"net/http"
	"github.com/gin-gonic/gin"
)

type ClickRecorder interface {
    RecordClick(linkID, referrer, userAgent string) error
}

type Handler struct {
	svc 			*Service
	clickRecorder	ClickRecorder	
}

func NewHandler(svc *Service, clickRecorder ClickRecorder) *Handler {
	return &Handler{svc: svc, clickRecorder: clickRecorder}
}

type ShortenInput struct {
	OriginalURL string `json:"original_url"`
	ExpiresAt	*time.Time `json:"expires_at"`
}

func (h *Handler) ShortenURL(c *gin.Context) {
	userID, _ := c.Get("userID")

	var input ShortenInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	link, err := h.svc.ShortenURL(userID.(string), input.OriginalURL, input.ExpiresAt)
	if err != nil {
    	log.Println("ShortenURL error:", err)
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

	if link.ExpiresAt != nil && time.Now().After(*link.ExpiresAt) {
		c.JSON(http.StatusGone, gin.H{"error": "this link has expired"})
		return
	}

	go func() {
    	if err := h.svc.IncrementClickCount(shortCode); err != nil {
        	log.Println("failed to increment click count:", err)
     	}
      	if h.clickRecorder != nil {
        	referrer := c.GetHeader("Referer")
         	userAgent := c.GetHeader("User-Agent")
          	if err := h.clickRecorder.RecordClick(link.ID, referrer, userAgent); err != nil {
            	log.Println("failed to record click:", err)
           	}
       	}
	}()


	c.Redirect(http.StatusMovedPermanently, link.OriginalURL)
}

func RegisterRoutes(rg *gin.RouterGroup, h *Handler, publicRg *gin.RouterGroup) {
	rg.POST("/links", h.ShortenURL)
	publicRg.GET("/:code", h.Redirect)
}