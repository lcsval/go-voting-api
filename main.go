package main

import (
	"github.com/gin-gonic/gin"
	"github.com/lcsval/go-voting-api/internal/config"
	"github.com/lcsval/go-voting-api/internal/database"
	"github.com/lcsval/go-voting-api/user"
)

func main() {
	r := gin.Default()
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "healthy"})
	})

	db, err := database.NewDB(config.NewConfig())

	if err != nil {
		panic(err)
	}

	user.RegisterRoutes(db, r)

	r.Run()
}
