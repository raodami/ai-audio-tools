package api

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"ai-audio-tools/internal/auth"
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

// setupMiddleware sets up global middleware
func setupMiddleware(r *gin.Engine) {
	// CORS
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Authorization, Content-Type")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})
}

// getLoggedInUser extracts user ID from JWT token
func getLoggedInUser(c *gin.Context) (string, bool) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing authorization header"})
		c.Abort()
		return "", false
	}

	tokenStr, err := auth.ExtractTokenFromHeader(authHeader)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token format"})
		c.Abort()
		return "", false
	}

	userID, err := auth.ValidateToken(tokenStr)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
		c.Abort()
		return "", false
	}

	c.Set("user_id", userID)
	return userID, true
}

// SetupRoutes sets up all API routes
func SetupRoutes(r *gin.Engine, s *store.Store) {
	setupMiddleware(r)

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok", "timestamp": time.Now().Unix()})
	})

	// Auth routes
	authGroup := r.Group("/api/auth")
	{
		authGroup.POST("/register", func(c *gin.Context) {
			var req RegisterRequest
			if err := c.ShouldBindJSON(&req); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}

			// Check if user exists
			_, err := s.GetUserByEmail(req.Email)
			if err == nil {
				c.JSON(http.StatusConflict, gin.H{"error": "email already registered"})
				return
			}

			// Hash password
			hashedPassword, err := auth.HashPassword(req.Password)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to hash password"})
				return
			}

			// Create user
			userID := uuid.New().String()
			if err := s.CreateUser(userID, req.Email, hashedPassword); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create user"})
				return
			}

			// Generate JWT token
			token, err := auth.GenerateToken(userID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
				return
			}

			c.JSON(http.StatusCreated, gin.H{
				"token": token,
				"user": UserInfo{
					ID:        userID,
					Email:     req.Email,
					IsPro:     false,
					UsageMin:  0,
					FreeQuota: store.FreeQuota,
				},
			})
		})

		authGroup.POST("/login", func(c *gin.Context) {
			var req LoginRequest
			if err := c.ShouldBindJSON(&req); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}

			user, err := s.GetUserByEmail(req.Email)
			if err != nil {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
				return
			}

			if err := auth.CheckPasswordHash(req.Password, user.Password); err != nil {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
				return
			}

			token, err := auth.GenerateToken(user.ID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
				return
			}

			c.JSON(http.StatusOK, gin.H{
				"token": token,
				"user": UserInfo{
					ID:        user.ID,
					Email:     user.Email,
					IsPro:     user.IsPro,
					UsageMin:  user.UsageMin,
					FreeQuota: store.FreeQuota,
				},
			})
		})

		authGroup.GET("/me", func(c *gin.Context) {
			userID, ok := getLoggedInUser(c)
			if !ok {
				return
			}

			user, err := s.GetUserByID(userID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get user"})
				return
			}

			c.JSON(http.StatusOK, gin.H{
				"id":           user.ID,
				"email":        user.Email,
				"is_pro":       user.IsPro,
				"usage_min":    user.UsageMin,
				"free_quota":   store.FreeQuota,
				"remaining_min": getMax(0, store.FreeQuota-user.UsageMin),
			})
		})
	}

	// User routes
	userGroup := r.Group("/api/user")
	{
		userGroup.GET("/usage", func(c *gin.Context) {
			userID, ok := getLoggedInUser(c)
			if !ok {
				return
			}

			user, err := s.GetUserByID(userID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get user"})
				return
			}

			c.JSON(http.StatusOK, gin.H{
				"id":            user.ID,
				"email":         user.Email,
				"is_pro":        user.IsPro,
				"usage_minutes": user.UsageMin,
				"free_quota":    store.FreeQuota,
				"remaining":     getMax(0, store.FreeQuota-user.UsageMin),
			})
		})

		userGroup.PUT("/subscribe", func(c *gin.Context) {
			// TODO: integrate Stripe checkout
			c.JSON(http.StatusOK, gin.H{"message": "Stripe integration coming soon"})
		})

		// GET /api/user/analytics — get usage statistics
		userGroup.GET("/analytics", func(c *gin.Context) {
			userID, ok := getLoggedInUser(c)
			if !ok {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
				return
			}

			analytics, err := s.GetAnalytics(userID, 30)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get analytics"})
				return
			}

			c.JSON(http.StatusOK, analytics)
		})
	}

	// Audio routes
	audioGroup := r.Group("/api/audio")
	{
		audioGroup.POST("/upload", func(c *gin.Context) {
			userID, _ := getLoggedInUser(c)
			cookieHash := c.GetHeader("X-Session-ID")

			// Check quota
			allowed, remaining, err := s.CheckQuota(userID, cookieHash)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			if !allowed {
				c.JSON(http.StatusForbidden, gin.H{"error": "quota exceeded", "remaining_min": remaining})
				return
			}

			// Handle multipart upload
			file, err := c.FormFile("audio")
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "file is required"})
				return
			}

			// Validate file type
			ext := getFileExtension(file.Filename)
			allowedExts := map[string]bool{
				".mp3": true, ".wav": true, ".m4a": true, ".flac": true,
				".ogg": true, ".aac": true, ".mp4": true,
			}
			if !allowedExts[ext] {
				c.JSON(http.StatusBadRequest, gin.H{"error": "unsupported file type"})
				return
			}

			// Create job record first
			jobID := uuid.New().String()
			if err := s.CreateJob(jobID, userID, cookieHash, file.Filename, file.Size); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create job"})
				return
			}

			// Save file to uploads directory
			os.MkdirAll("uploads", 0755)
			dst := fmt.Sprintf("uploads/%s%s", jobID, ext)
			if err := c.SaveUploadedFile(file, dst); err != nil {
				s.UpdateJobStatus(jobID, "failed", err.Error())
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save file"})
				return
			}

			// TODO: trigger async processing with Deepgram ASR
			go processAudioAsync(s, jobID, dst, userID, cookieHash)

			c.JSON(http.StatusAccepted, gin.H{
				"job_id":    jobID,
				"status":    "processing",
				"file_name": file.Filename,
			})
		})

		audioGroup.GET("/:id", func(c *gin.Context) {
			jobID := c.Param("id")
			job, err := s.GetJob(jobID)
			if err != nil {
				c.JSON(http.StatusNotFound, gin.H{"error": "job not found"})
				return
			}
			c.JSON(http.StatusOK, job)
		})
	}
}

// processAudioAsync handles async audio processing
func processAudioAsync(s *store.Store, jobID, filePath, userID, cookieHash string) {
	// TODO: implement Deepgram transcription + DeepSeek summarization
	// For now, just mark as completed with placeholder
	result := `{"transcript": "Processing... Connect Deepgram API for real transcription"}`
	s.UpdateJobStatus(jobID, "completed", result)
}

// Helper functions
func getFileExtension(filename string) string {
	for i := len(filename) - 1; i >= 0; i-- {
		if filename[i] == '.' {
			return filename[i:]
		}
	}
	return ""
}

func getMax(a, b int) int {
	if a > b {
		return a
	}
	return b
}
