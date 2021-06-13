package models

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/txfs19260817/scopelens/server/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// Redis keys
const (
	DBTimeout = 10 * time.Second

	Total = "total" // Total refers to a key which binds to the total number of teams a.k.a. count

	TimeOrderAll  = "time:all"  // TimeOrderAll refers to a hash key that stores pages of data ordered by time
	LikesOrderAll = "likes:all" // LikesOrderAll refers to a hash key that stores pages of data ordered by likes
)

var (
	Db  *DBDriver     // Db is a database instance
	Rdb *redis.Client // Rdb is a redis client instance
)

// DBDriver is a wrapper for the mongo-go-driver.
type DBDriver struct {
	DB      *mongo.Database
	Client  *mongo.Client
	Context context.Context
}

// Close closes the mongo-go-driver connection.
func (d *DBDriver) Close() {
	_ = d.Client.Disconnect(d.Context)
}

// InitDB initializes a database instance
func InitDB() (*DBDriver, error) {
	// Define a timeout duration
	ctx, cancel := context.WithTimeout(context.Background(), DBTimeout)
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
	ctxPing, cancelPing := context.WithTimeout(context.Background(), DBTimeout)
	defer cancelPing()
	if err := client.Ping(ctxPing, readpref.Primary()); err != nil {
		return nil, err
	}

	db := client.Database(config.Database.DBName)
	return &DBDriver{DB: db, Client: client, Context: ctx}, nil
}

// InitRedis initializes a Redis instance
func InitRedis() (*redis.Client, error) {
	c := redis.NewClient(&redis.Options{
		Addr:     config.Redis.Host + ":" + config.Redis.Port,
		Password: config.Redis.Password,
		DB:       0, // use default DB
	})
	return c, ping(c)
}

func ping(client *redis.Client) error {
	ctx, cancel := context.WithTimeout(context.Background(), DBTimeout)
	defer cancel()
	pong, err := client.Ping(ctx).Result()
	if err != nil {
		return err
	}
	if pong != "PONG" {
		return fmt.Errorf("the response of PING is not PONG !!! ")
	}
	return nil
}
