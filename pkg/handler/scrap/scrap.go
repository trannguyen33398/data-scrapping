package scrapdata

import (
	"net/http"
	"scrapping/pkg/config"
	"scrapping/pkg/handler/scrap/request"
	"scrapping/pkg/controller"
	"scrapping/pkg/logger"
	"scrapping/pkg/utils/validate"
	"github.com/gin-gonic/gin"
)

type handler struct {
	logger     logger.Logger
	config     *config.Config
	controller *controller.Controller
}

// Scrap implements IHandler.

// New returns a handler
func New(controller *controller.Controller, logger logger.Logger, cfg *config.Config) IHandler {
	return &handler{
		logger:     logger,
		config:     cfg,
		controller: controller,
	}
}

func (h *handler) Scrap(c *gin.Context) {
	input := request.UrlRequest{}

	if err := c.ShouldBindJSON(&input); err != nil {

		c.JSON(http.StatusBadRequest, err.Error())
		return 
	}
	l := h.logger.Fields(logger.Fields{
		"handler": "scrap-data",
		"method":  "POST",
		"url":     input.Url,
	})

	validateError := validate.IsWikipediaUrl(input.Url)

	if validateError != nil {
		c.JSON(http.StatusBadRequest, validateError.Error())
		return 
	}

	err := h.controller.ScrapData.Scrap(c, input.Url)

	if err != nil {
		l.Error( err,"failed to craw data")
		c.JSON(http.StatusInternalServerError, err)
		return 
	}

	c.Status(http.StatusCreated)
	return 
}
