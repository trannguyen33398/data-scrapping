package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"scrapping/pkg/config"
	"scrapping/pkg/logger"
	"scrapping/pkg/routes"
)

func main() {
	
	cfg := config.LoadConfig()
	log := logger.NewLogrusLogger()
	log.Infof("Server starting at port %v", cfg.ApiServer.Port)

	_, cancel := context.WithCancel(context.Background())
	router := routes.NewRoutes(cfg,  log)
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.ApiServer.Port),
		Handler: router,
	}

	// serve http server
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err, "failed to listen and serve")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	<-quit

	cancel()

	shutdownServer(srv, log)

}

func shutdownServer(srv *http.Server, l logger.Logger) {
	l.Info("Server Shutting Down")
	if err := srv.Shutdown(context.Background()); err != nil {
		l.Error(err, "failed to shutdown server")
	}

	l.Info("Server Exit")
}