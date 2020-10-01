package main

import (
	"fmt"
	"net/http"
	"scopelens-server/config"
	"scopelens-server/models"
	"scopelens-server/routers"
	"scopelens-server/utils/logger"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	// Routers
	router := routers.InitRouters()

	// Server setup
	s := &http.Server{
		Addr:           fmt.Sprintf(":%d", config.Server.Port),
		Handler:        router,
		ReadTimeout:    config.Server.ReadTimeout * time.Second,
		WriteTimeout:   config.Server.WriteTimeout * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	// Database Connection
	var err error
	models.Db, err = models.InitDB()
	if err != nil {
		panic(err)
	}
	logger.SugaredLogger.Info("Database Connected. ")
	defer models.Db.Close()

	// Start server depending on running mode
	switch config.Mode {
	case "debug":
		gin.SetMode(config.Mode)
		s.ListenAndServe()
	case "release":
		gin.SetMode(config.Mode)
		s.ListenAndServeTLS(config.Server.HttpsCrt, config.Server.HttpsKey)
	default:
		panic("Running mode %v is not available: " + config.Mode)
	}
}
