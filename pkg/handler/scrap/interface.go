package scrapdata

import (

	"github.com/gin-gonic/gin"

)

type IHandler interface {
	Scrap(c *gin.Context)
}