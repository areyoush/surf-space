package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type TokenChecker interface {
	IsTokenDenylisted(token string) bool
}

func Auth(jwtSecret string, checker TokenChecker) gin.HandlerFunc {
    return func(c *gin.Context) {
        header := c.GetHeader("Authorization")
        if header == "" {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "authorization header required"})
            return
        }

        tokenString := strings.TrimPrefix(header, "Bearer ")

        token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
            return []byte(jwtSecret), nil
        })

        if err != nil || !token.Valid {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
            return
        }

        if checker.IsTokenDenylisted(tokenString) {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "token revoked"})
            return
        }

        claims := token.Claims.(jwt.MapClaims)
        userID := claims["sub"].(string)
        c.Set("userID", userID)

        c.Next()
    }
}