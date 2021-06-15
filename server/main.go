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
	"go.uber.org/zap"
)

func main() {
	// Setup global logger
	l := logger.CreateGlobalLogger()
	defer l.Sync()
	zap.L().Info("loaded config file", zap.String("path", config.CfgPath))

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
		zap.L().Panic("Failed to connect database", zap.Error(err), zap.String("type", config.Database.Type), zap.String("name", config.Database.DBName))
	}
	zap.L().Info("Database connected", zap.String("type", config.Database.Type), zap.String("name", config.Database.DBName))
	defer models.Db.Close()

	// Redis Connection
	models.Rdb, err = models.InitRedis()
	if err != nil {
		zap.L().Panic("Failed to connect Redis", zap.Error(err), zap.String("host", config.Redis.Host), zap.String("port", config.Redis.Port), zap.String("type", "Redis"))
	}
	zap.L().Info("Redis Connected", zap.String("host", config.Redis.Host), zap.String("port", config.Redis.Port), zap.String("type", "Redis"))
	defer models.Rdb.Close()

	// S3 session establishing
	storage.S3Client, err = storage.NewAmazonS3(config.Aws.AccessKey, config.Aws.SecretKey, config.Aws.Region, config.Aws.Bucket)
	if err != nil {
		zap.L().Panic("Failed to connect AWS S3", zap.Error(err), zap.String("bucket", config.Aws.Bucket), zap.String("region", config.Aws.Region))
	}
	zap.L().Info("AWS S3 session established. ", zap.String("bucket", config.Aws.Bucket), zap.String("region", config.Aws.Region))

	// Start server
	if config.Server.EnableHttps {
		if err := s.ListenAndServeTLS(config.Server.HttpsCrt, config.Server.HttpsKey); err != nil {
			zap.L().Panic("Failed to establish HTTPs server", zap.Error(err), zap.Bool("tls", config.Server.EnableHttps))
		}
	} else {
		if err := s.ListenAndServe(); err != nil {
			zap.L().Panic("Failed to establish HTTP server", zap.Error(err), zap.Bool("tls", config.Server.EnableHttps))
		}
	}
}
