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
	zap.L().Info("loaded config file", zap.String("configFilePath", config.CfgPath))

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
		zap.L().Panic("Failed to connect database",zap.Error(err), zap.String("databaseType", config.Database.Type))
	}
	zap.L().Info("Database connected", zap.String("databaseType", config.Database.Type))
	defer models.Db.Close()

	// Redis Connection
	models.Rdb, err = models.InitRedis()
	if err != nil {
		zap.L().Panic("Failed to connect Redis",zap.Error(err), zap.String("databaseType", "Redis"))
	}
	zap.L().Info("Redis Connected", zap.String("redisHost", config.Redis.Host), zap.String("redisPort", config.Redis.Port))
	defer models.Rdb.Close()

	// S3 session establishing
	storage.S3Client, err = storage.NewAmazonS3(config.Aws.AccessKey, config.Aws.SecretKey, config.Aws.Region, config.Aws.Bucket)
	if err != nil {
		zap.L().Panic("Failed to connect AWS S3", zap.Error(err),zap.String("s3Bucket", config.Aws.Bucket), zap.String("s3Region", config.Aws.Region))
	}
	zap.L().Info("AWS S3 session established. ", zap.String("s3Bucket", config.Aws.Bucket), zap.String("s3Region", config.Aws.Region))

	// Start server
	if config.Server.EnableHttps {
		if err := s.ListenAndServeTLS(config.Server.HttpsCrt, config.Server.HttpsKey); err != nil {
			zap.L().Panic("Failed to establish HTTPs server", zap.Error(err), zap.Bool("TLS", config.Server.EnableHttps))
		}
	} else {
		if err := s.ListenAndServe(); err != nil {
			zap.L().Panic("Failed to establish HTTP server", zap.Error(err), zap.Bool("TLS", config.Server.EnableHttps))
		}
	}
}
