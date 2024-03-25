package test

import (
	"bytes"
	"encoding/json"
	"os"

	"net/http"
	"net/http/httptest"
	"scrapping/pkg/config"
	"scrapping/pkg/controller"
	"scrapping/pkg/handler/scrap/request"

	"scrapping/pkg/logger"
	"scrapping/pkg/routes"
	"testing"

	"github.com/gin-gonic/gin"

	"github.com/stretchr/testify/assert"
)

func router() *gin.Engine {

	cfg := &config.Config{
		Debug: true,

		ApiServer: config.ApiServer{
			Port:           "3000",
			AllowedOrigins: "*",
		},
	}
	log := logger.NewLogrusLogger()
	router := routes.NewRoutes(cfg, log)
	return router
}

func makeRequest(method, url string, body interface{}) *httptest.ResponseRecorder {

	requestBody, _ := json.Marshal(body)
	request, _ := http.NewRequest(method, url, bytes.NewBuffer(requestBody))
	request.Header.Set("Content-Type", "application/json")
	writer := httptest.NewRecorder()
	r := router()
	r.ServeHTTP(writer, request)

	return writer
}

func TestScrapDataWithInvalidUrl(t *testing.T) {
	data := request.UrlRequest{
		Url: "http://google.com",
	}
	body, _ := json.Marshal(data)
	reader := bytes.NewReader(body)
	writer := makeRequest("POST", "http://localhost:3000/api/v1/scrap-data/", reader)

	assert.Equal(t, http.StatusBadRequest, writer.Code)
}

func TestScrapDataWithValidUrl(t *testing.T) {
	data := request.UrlRequest{
		Url: "https://en.wikipedia.org/wiki/Women%27s_100_metres_world_record_progression",
	}
	body, _ := json.Marshal(data)
	reader := bytes.NewReader(body)
	writer := makeRequest("POST", "http://localhost:3000/api/v1/scrap-data/", reader)

	assert.Equal(t, http.StatusBadRequest, writer.Code)
}

func TestScrapDataController(t *testing.T) {
	log := logger.NewLogrusLogger()

	cfg := &config.Config{
		Debug: true,

		ApiServer: config.ApiServer{
			Port:           "3000",
			AllowedOrigins: "*",
		},
	}
	c := controller.New(log, cfg)
	g, _ := gin.CreateTestContext(httptest.NewRecorder())
	url := "https://en.wikipedia.org/wiki/Women%27s_100_metres_world_record_progression"
	err := c.ScrapData.Scrap(g, url)
	assert.NoError(t, err, "Expected no error with valid input")

	filenames := []string{"image1.png", "image2.png", "image3.png"}
	for _, filename := range filenames {
		_, err := os.Stat(filename)
		assert.Nil(t, err, "Error checking file %s", filename)
	}

}
