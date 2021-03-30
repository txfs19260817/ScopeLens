package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/txfs19260817/scopelens/server/config"
	"github.com/txfs19260817/scopelens/server/models"
	"github.com/txfs19260817/scopelens/server/routers"
	"github.com/txfs19260817/scopelens/server/utils/logger"
	"github.com/txfs19260817/scopelens/server/utils/storage"
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

	var err error
	// Database Connection
	models.Db, err = models.InitDB()
	if err != nil {
		panic(err.Error())
	}
	logger.SugaredLogger.Info("Database Connected. ")
	defer models.Db.Close()

	// Redis Connection
	models.Rdb = models.InitRedis()
	defer models.Rdb.Close()

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
