package internal

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"net/http"
	"strings"
	"synthomat/minimi/internal/db"
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
	AUTH_SECRET = "abcde"
)

type LinkDto struct {
	Slug string `form:"slug" binding:"required,alphanum,min=3,max=20"`
	Url  string `form:"url" binding:"required,url"`
}

func RunServer(gdb *gorm.DB) {
	server := gin.Default()
	sessSecret, _ := uuid.NewRandom()

	store := cookie.NewStore([]byte(sessSecret.String()[:32]))
	server.Use(sessions.Sessions("sess", store))

	enLang := en.New()
	uni := ut.New(enLang)
	trans, _ := uni.GetTranslator("en")

	b, _ := binding.Validator.Engine().(*validator.Validate)
	en_translations.RegisterDefaultTranslations(b, trans)

	ginHtmlRenderer := server.HTMLRender
	server.HTMLRender = &HTMLTemplRenderer{FallbackHtmlRenderer: ginHtmlRenderer}
	server.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusPermanentRedirect, "/a/")
	})
	server.GET("/:slug", func(c *gin.Context) {
		slugName := c.Param("slug")

		slug, err := db.LinkBySlug(gdb, slugName)

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
				r := New(c, internal.Login(next))
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

			r := New(c, internal.Login(next))
			c.Render(http.StatusOK, r)
		})

		authGroup.POST("/logout", func(c *gin.Context) {
			session := sessions.Default(c)
			session.Clear()
			session.Options(sessions.Options{Path: "/", MaxAge: -1})
			session.Save()

			c.Redirect(http.StatusFound, "/a/")
		})
	}

	adminGroup := server.Group("/a").Use(AuthMiddleware())
	{
		adminGroup.GET("/", func(c *gin.Context) {
			c.Set("authed", true)
			var links []db.Link
			gdb.Order("slug asc, created_at asc").Find(&links)

			r := New(c, internal.Links(links))
			c.Render(http.StatusOK, r)
		})

		adminGroup.Any("/links/new", func(c *gin.Context) {

			if c.Request.Method == http.MethodPost {
				slug := c.PostForm("slug")
				url := c.PostForm("url")

				link := db.Link{Slug: slug, OriginalUrl: url}

				if err := gdb.Create(&link).Error; err != nil {

				}

				c.Redirect(http.StatusFound, "/")
				return
			}

			r := New(c, internal.NewLinkLayout())
			c.Render(http.StatusOK, r)
		})

		adminGroup.Any("/links/:linkId/edit", func(c *gin.Context) {

			linkId := c.Param("linkId")
			link, _ := db.LinkById(gdb, linkId)

			if c.Request.Method != http.MethodPost {
				r := New(c, internal.EditLinkLayout(*link, nil))
				c.Render(http.StatusOK, r)
				return
			}

			var linkDto LinkDto
			err := c.ShouldBind(&linkDto)

			if err != nil {
				validationErrors := err.(validator.ValidationErrors)
				errors := make(map[string]string)
				for _, v := range validationErrors {
					errors[strings.ToLower(v.Field())] = v.Translate(trans)
				}

				r := New(c, internal.EditLinkLayout(*link, errors))
				c.Render(http.StatusOK, r)
				return
			}

			link.Slug = linkDto.Slug
			link.OriginalUrl = linkDto.Url

			if err := gdb.Save(&link).Error; err != nil {
				r := New(c, internal.EditLinkLayout(*link, nil))
				c.Render(http.StatusOK, r)
				return
			}

			c.Redirect(http.StatusFound, "/a/")
		})
	}

	server.Run(":8000")
}
