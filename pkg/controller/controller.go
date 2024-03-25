package controller

import (
	"scrapping/pkg/config"
	scrapdata "scrapping/pkg/controller/scrap"
	"scrapping/pkg/logger"
)

type Controller struct {
	ScrapData scrapdata.IController
}

func New(logger logger.Logger, cfg *config.Config) *Controller {
	return &Controller{

		ScrapData: scrapdata.New(logger, cfg),
	}
}
