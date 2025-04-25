package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/RasoulZamani/motivational-api/storage"
)

func main() {
	store, err := storage.NewSQLiteStorage("./storage/quotes.db")
	if err != nil {
		log.Fatal("Failed to initialize storage:", err)
	}

	if err := store.SyncFromJSON("./storage/quotes.json"); err != nil {
		log.Fatal("Failed to sync quotes:", err)
	}

	// Start server
	r := gin.Default()
	r.GET("/quote", func(c *gin.Context) {
		quote, err := store.GetRandomQuote()
		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to fetch quote"})
			return
		}
		c.JSON(200, gin.H{"quote": quote})
	})
	r.Run(":8042")
}