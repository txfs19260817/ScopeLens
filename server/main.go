package main

import (
	"fmt"
	"net/http"
	"time"

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
	logger.SugaredLogger.Infof("Loaded config from: %s", config.CfgPath)

	var err error
	// Database Connection
	models.Db, err = models.InitDB()
	if err != nil {
		panic(err)
	}
	logger.SugaredLogger.Infof("Database %s Connected. ", config.Database.Type)
	defer models.Db.Close()

	// Redis Connection
	models.Rdb, err = models.InitRedis()
	if err != nil {
		panic(err)
	}
	logger.SugaredLogger.Info("Redis Connected. ")
	defer models.Rdb.Close()

	// S3 session establishing
	storage.S3Client, err = storage.NewAmazonS3(config.Aws.AccessKey, config.Aws.SecretKey, config.Aws.Region, config.Aws.Bucket)
	if err != nil {
		panic(err)
	}
	logger.SugaredLogger.Info("AWS S3 session established. ")

	// Start server
	if config.Server.EnableHttps {
		if err := s.ListenAndServeTLS(config.Server.HttpsCrt, config.Server.HttpsKey); err != nil {
			panic(err.Error())
		}
	} else {
		if err := s.ListenAndServe(); err != nil {
			panic(err)
		}
	}
}
