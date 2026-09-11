package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"ai-audio-tools/internal/store"
)

// RegisterRequest represents a register request
type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

// LoginRequest represents a login request
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// UserInfo represents current user info
type UserInfo struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	IsPro     bool   `json:"is_pro"`
	UsageMin  int    `json:"usage_minutes"`
	FreeQuota int    `json:"free_quota"`
}

// SetupRoutes sets up auth routes
func SetupRoutes(r *gin.Engine, s *store.Store) {
	auth := r.Group("/api/auth")
	{
		auth.POST("/register", func(c *gin.Context) {
			var req RegisterRequest
			if err := c.ShouldBindJSON(&req); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			// TODO: implement actual registration with bcrypt
			c.JSON(http.StatusCreated, gin.H{"message": "registered"})
		})
		auth.POST("/login", func(c *gin.Context) {
			var req LoginRequest
			if err := c.ShouldBindJSON(&req); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			// TODO: implement actual login
			c.JSON(http.StatusOK, gin.H{"message": "login"})
		})
		auth.GET("/me", func(c *gin.Context) {
			// TODO: implement user info
			c.JSON(http.StatusOK, gin.H{"message": "auth"})
		})
	}

	user := r.Group("/api/user")
	{
		user.GET("/usage", func(c *gin.Context) {
			allowed, remaining, err := s.CheckQuota("", "")
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusOK, gin.H{
				"allowed":       allowed,
				"remaining_min": remaining,
				"free_quota":    store.FreeQuota,
			})
		})
	}
}
