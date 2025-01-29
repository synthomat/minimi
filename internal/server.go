package internal

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
)

var reservedSlugs = []string{
	"admin",
	"settings",
}

func RunServer(db *gorm.DB) {
	server := gin.Default()

	server.GET("/:slug", func(c *gin.Context) {
		slugName := c.Param("slug")

		slug, _ := LinkBySlug(db, slugName)

		/*
			c.JSON(200, gin.H{
				"data": slug,
			})
		*/

		c.Redirect(http.StatusFound, slug.OriginalUrl)
	})

	server.Run(":8000")
}
