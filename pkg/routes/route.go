package routes

import (
	"scrapping/pkg/config"
	"scrapping/pkg/logger"
	"strings"
	"scrapping/pkg/controller"
	"scrapping/pkg/handler"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
)

func setupCORS(r *gin.Engine, cfg *config.Config) {
	corsOrigins := strings.Split(cfg.ApiServer.AllowedOrigins, ";")
	r.Use(func(c *gin.Context) {
		cors.New(
			cors.Config{
				AllowOrigins: corsOrigins,
				AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "HEAD"},
				AllowHeaders: []string{
					"Origin", "Host", "Content-Type", "Content-Length", "Accept-Encoding", "Accept-Language", "Accept",
					"X-CSRF-Token", "Authorization", "X-Requested-With", "X-Access-Token",
				},
				AllowCredentials: true,
			},
		)(c)
	})
}

func NewRoutes(cfg *config.Config, logger logger.Logger) *gin.Engine {
	// programmatically set swagger info

	r := gin.New()

	pprof.Register(r)

	ctrl := controller.New( logger, cfg)
	h := handler.New( ctrl, logger, cfg)

	r.Use(
		gin.LoggerWithWriter(gin.DefaultWriter, "/healthz"),
		gin.Recovery(),
	)
	// config CORS
	setupCORS(r, cfg)

	// load API here
	v1 := r.Group("/api/v1")
	scrapRoute := v1.Group("/scrap-data")
	{
	
		scrapRoute.POST("/", h.ScrapData.Scrap)

	}

	



	return r
}