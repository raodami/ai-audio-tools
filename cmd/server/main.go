package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"ai-audio-tools/internal/api"
	"ai-audio-tools/internal/store"
	"ai-audio-tools/internal/webhook"
	"ai-audio-tools/internal/ws"
)

func main() {
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "data/app.db"
	}

	db, err := store.NewDB(dbPath)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Initialize WebSocket server
	wsServer := ws.NewWSServer()
	go wsServer.Run()

	r := gin.Default()

	api.SetupRoutes(r, db)
	webhook.SetupWebhookRoutes(r, db)

	// WebSocket endpoint
	r.GET("/ws", func(c *gin.Context) {
		wsServer.ServeHTTP(c.Writer, c.Request)
	})

	log.Printf("🎵 AI Audio Tools Server starting on :%s", port)
	log.Fatal(r.Run(":" + port))
}
