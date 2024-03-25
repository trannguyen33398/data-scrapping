package handler

import (
	"scrapping/pkg/config"
	"scrapping/pkg/controller"
	scrapdata "scrapping/pkg/handler/scrap"
	"scrapping/pkg/logger"
)

type Handler struct {
	ScrapData scrapdata.IHandler
}

func New(ctrl *controller.Controller, logger logger.Logger, cfg *config.Config) *Handler {
	return &Handler{

		ScrapData: scrapdata.New(ctrl, logger, cfg),
	}
}
