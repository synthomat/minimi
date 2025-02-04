package internal

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	db2 "synthomat/minimi/internal/db"
	internal "synthomat/minimi/internal/templates"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		authed := session.Get("authed")

		if authed != true {
			c.Redirect(http.StatusFound, "/auth/login?next="+c.Request.RequestURI)
			return
		}

		c.Set("authed", true)
		c.Next()
	}
}

// INSECURE // INSECURE // INSECURE // INSECURE // INSECURE // INSECURE // INSECURE // INSECURE // INSECURE // INSECURE
// DO NOT, UNDER ANY CIRCUMSTANCES, USE THIS FOR PRODUCTION!
// Will change it later
// INSECURE // INSECURE // INSECURE // INSECURE // INSECURE // INSECURE // INSECURE // INSECURE // INSECURE // INSECURE
var (
	SESS_SECRET = "sessionsecret"
	AUTH_SECRET = "abcde"
)

type LinkDto struct {
	Slug string `form:"slug"`
	Url  string `form:"url"`
}

func RunServer(db *gorm.DB) {
	server := gin.Default()
	store := cookie.NewStore([]byte(SESS_SECRET))
	server.Use(sessions.Sessions("mysession", store))

	ginHtmlRenderer := server.HTMLRender
	server.HTMLRender = &HTMLTemplRenderer{FallbackHtmlRenderer: ginHtmlRenderer}

	server.GET("/:slug", func(c *gin.Context) {
		slugName := c.Param("slug")

		slug, err := db2.LinkBySlug(db, slugName)

		if err != nil {
			c.JSON(http.StatusNotFound, "not found")
			return
		}

		c.Redirect(http.StatusFound, slug.OriginalUrl)
	})

	authGroup := server.Group("/auth")
	{
		authGroup.Any("/login", func(c *gin.Context) {
			next := c.Query("next")

			if c.Request.Method != http.MethodPost {
				r := New(c, http.StatusOK, internal.Login(next))
				c.Render(http.StatusOK, r)
				return
			}

			password := c.PostForm("password")

			if password == AUTH_SECRET {
				session := sessions.Default(c)
				session.Set("authed", true)
				err := session.Save()
				if err != nil {
					panic(err)
					return
				}

				if next == "" {
					next = "/a/"
				}

				c.Redirect(http.StatusFound, next)
				return
			}

			r := New(c, http.StatusOK, internal.Login(next))
			c.Render(http.StatusOK, r)
		})
	}

	adminGroup := server.Group("/a").Use(AuthMiddleware())
	{
		adminGroup.GET("/", func(c *gin.Context) {
			c.Set("authed", true)
			var links []db2.Link
			db.Order("slug asc, created_at asc").Find(&links)

			r := New(c, http.StatusOK, internal.Links(links))
			c.Render(http.StatusOK, r)
		})

		adminGroup.Any("/links/new", func(c *gin.Context) {

			if c.Request.Method == http.MethodPost {
				slug := c.PostForm("slug")
				url := c.PostForm("url")

				link := db2.Link{Slug: slug, OriginalUrl: url}

				if err := db.Create(&link).Error; err != nil {

				}

				c.Redirect(http.StatusFound, "/")
				return
			}

			r := New(c, http.StatusOK, internal.NewLinkLayout())
			c.Render(http.StatusOK, r)
		})

		adminGroup.Any("/links/:linkId/edit", func(c *gin.Context) {

			linkId := c.Param("linkId")
			link, _ := db2.LinkById(db, linkId)

			if c.Request.Method != http.MethodPost {
				r := New(c, http.StatusOK, internal.EditLinkLayout(*link))
				c.Render(http.StatusOK, r)
				return
			}

			var linkDto LinkDto
			c.Bind(&linkDto)

			link.Slug = linkDto.Slug
			link.OriginalUrl = linkDto.Url

			if err := db.Save(&link).Error; err != nil {
				r := New(c, http.StatusOK, internal.EditLinkLayout(*link))
				c.Render(http.StatusOK, r)
				return
			}

			c.Redirect(http.StatusFound, "/a/")
		})
	}

	server.Run(":8000")
}
