package ginadapter

import (
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type RouterConfig struct {
	AllowedOrigin string
}

func NewRouter(handler *ProteinHandler, config RouterConfig) *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery())
	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{config.AllowedOrigin},
		AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodOptions},
		AllowHeaders: []string{"Origin", "Content-Type", "Accept"},
		MaxAge:       12 * time.Hour,
	}))

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	api := router.Group("/api/v1")
	proteins := api.Group("/proteins")
	{
		proteins.GET("", handler.List)
		proteins.POST("", handler.Create)
		proteins.GET("/:accession", handler.GetDetails)
		proteins.PUT("/:accession", handler.Update)
		proteins.DELETE("/:accession", handler.Delete)
	}

	return router
}
