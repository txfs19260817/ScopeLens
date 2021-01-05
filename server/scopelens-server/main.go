package main

import (
	"fmt"
	"net/http"
	"scopelens-server/config"
	"scopelens-server/models"
	"scopelens-server/routers"
	"scopelens-server/utils/logger"
	"scopelens-server/utils/storage"
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
		panic(err.Error())
	}
	logger.SugaredLogger.Info("Database Connected. ")
	defer models.Db.Close()

	// S3 session establishing
	storage.S3Client, err = storage.NewAmazonS3(config.Aws.AccessKey, config.Aws.SecretKey, config.Aws.Region, config.Aws.Bucket)
	if err != nil {
		panic(err.Error())
	}

	// Start server depending on running mode
	switch config.Mode {
	case "debug":
		gin.SetMode(config.Mode)
		if err := s.ListenAndServe(); err != nil {
			panic(err.Error())
		}
	case "release":
		gin.SetMode(config.Mode)
		if err := s.ListenAndServeTLS(config.Server.HttpsCrt, config.Server.HttpsKey); err != nil {
			panic(err.Error())
		}
	default:
		panic("Running mode %v is not available: " + config.Mode)
	}
}
