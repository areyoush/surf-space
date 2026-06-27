package auth 

import (
	"net/http"
	"errors"

	"github.com/gin-gonic/gin"
)


type Handler struct {
	svc	*Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

type RegisterInput struct {
	Email		string	`json:"email" binding:"required,email"`
	Password	string	`json:"password" binding:"required,min=8"`
}

type LoginInput struct {
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required"`
}

func (h *Handler) Register(c *gin.Context) {
	var input RegisterInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if len(input.Password) > 72 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "password must be between 8 and 72 characters"})
			return
	}

	token, err := h.svc.Register(input.Email, input.Password)
	if err != nil {
		if errors.Is(err, ErrEmailTaken) {
			c.JSON(http.StatusConflict, gin.H{"error": "email already registered"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "registration failed"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"token": token})
}

func (h *Handler) Login(c *gin.Context) {
	var input LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	token, err := h.svc.Login(input.Email, input.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "login failed. please check your details"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": token})
}

func RegisterRoutes(rg *gin.RouterGroup, h *Handler) {
	rg.POST("/register", h.Register)
	rg.POST("/login", h.Login)
}