package scrapdata

import (
	"scrapping/pkg/config"
	"scrapping/pkg/logger"

	"github.com/gin-gonic/gin"
)

type controller struct {
	logger logger.Logger
	config *config.Config
}

func New(logger logger.Logger, cfg *config.Config) IController {
	return &controller{
		logger: logger,
		config: cfg,
	}
}

type IController interface {
	Scrap(c *gin.Context, url string) (err error)
}
