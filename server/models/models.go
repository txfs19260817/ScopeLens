package models

import (
	"context"
	"time"

	"github.com/txfs19260817/scopelens/server/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// Database instance
var Db *DBDriver

// DBDriver is a wrapper for the mongo-go-driver.
type DBDriver struct {
	DB      *mongo.Database
	Client  *mongo.Client
	Context context.Context
}

// Close closes the mongo-go-driver connection.
func (d *DBDriver) Close() {
	d.Client.Disconnect(d.Context)
}

func InitDB() (*DBDriver, error) {
	// Define a timeout duration
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Configure our client to use the correct URI, but we're not yet connecting to it.
	client, err := mongo.NewClient(options.Client().ApplyURI(config.Database.AtlasURI))
	if err != nil {
		return nil, err
	}

	// Try to connect using the defined timeout duration
	if err := client.Connect(ctx); err != nil {
		return nil, err
	}

	// Ping the cluster to ensure we're already connected
	ctxPing, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := client.Ping(ctxPing, readpref.Primary()); err != nil {
		return nil, err
	}

	db := client.Database(config.Database.DBName)
	return &DBDriver{DB: db, Client: client, Context: ctx}, nil
}
