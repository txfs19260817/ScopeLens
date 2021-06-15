package models

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
	jsoniter "github.com/json-iterator/go"
	"github.com/txfs19260817/scopelens/server/config"
	"github.com/txfs19260817/scopelens/server/utils/file"
	"github.com/txfs19260817/scopelens/server/utils/showdown"
	"github.com/txfs19260817/scopelens/server/utils/storage"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

// Team holds
type Team struct {
	ID primitive.ObjectID `bson:"_id" json:"id"`
	// User-defined fields
	Title       string   `bson:"title" json:"title" binding:"required"`
	Author      string   `bson:"author" json:"author" binding:"required"`
	Format      string   `bson:"format" json:"format" binding:"required"`
	Pokemon     []string `bson:"pokemon" json:"pokemon" binding:"required"`
	Showdown    string   `bson:"showdown" json:"showdown"`
	Image       string   `bson:"image" json:"image"` // upload from client: base64; store in DB: URL
	Description string   `bson:"description" json:"description"`
	// Auto-generated fields
	Uploader    string    `bson:"uploader" json:"uploader"` // User.username
	CreatedAt   time.Time `bson:"created_at" json:"created_at"`
	Likes       int       `bson:"likes" json:"likes"` // 0
	State       int       `bson:"state" json:"state"` // 0
	HasShowdown bool      `bson:"has_showdown" json:"has_showdown"`
	HasRental   bool      `bson:"has_rental" json:"has_rental"`
}

// Usage struct
type Usage struct {
	Pokemon string `bson:"_id"   json:"name"`
	Count   int    `bson:"count" json:"value"`
}

// Search criteria struct
type Search struct {
	Format      string   `bson:"format" json:"format"`
	Pokemon     []string `bson:"pokemon" json:"pokemon"`
	HasShowdown bool     `bson:"has_showdown" json:"has_showdown"`
	HasRental   bool     `bson:"has_rental" json:"has_rental"`
	OrderBy     string   `json:"order_by"`
}

// InsertTeam inserts a team
func (d *DBDriver) InsertTeam(team Team) (ok bool, err error) {
	var fileFullPath string // image save path
	ctx := context.Background()

	// use uploaded image first
	if len(team.Image) != 0 {
		// decode uploaded base64 string to file
		fileFullPath, err = file.DecodeBase64AndSave(team.Image)
		if err != nil {
			return false, err
		}
	} else {
		if len(team.Showdown) == 0 {
			return false, fmt.Errorf("Image and Showdown cannot be empty at the same time. ")
		}
		// generate image from showdown text.
		fileFullPath, err = showdown.RentalTeamMaker(team.Showdown, team.Title, team.Author)
		if err != nil {
			return false, err
		}
	}
	// resize
	if err := file.Rescale(fileFullPath); err != nil {
		return false, err
	}

	// upload to S3.
	imageFile, err := os.Open(fileFullPath)
	if err != nil {
		return false, fmt.Errorf("failed to open file %q, %v", fileFullPath, err)
	}
	tempPath := strings.Split(fileFullPath, "/")
	uploadPath := config.Aws.TeamPath + "/" + tempPath[len(tempPath)-1]
	url, err := storage.S3Client.Save(uploadPath, imageFile)
	if err != nil {
		return false, err
	}

	// fill fields
	team.ID, team.CreatedAt, team.Image = primitive.NewObjectID(), time.Now(), url

	// Insert
	if _, err := d.DB.Collection("teams").InsertOne(ctx, team); err != nil {
		return false, err
	}

	// reset hash maps for page data
	redisKeys := []string{Total, TimeOrderAll, LikesOrderAll}
	if err := Rdb.Del(ctx, redisKeys...).Err(); err != nil {
		return false, err
	}
	return true, nil
}

// GetTeamsCount counts teams
func (d *DBDriver) GetTeamsCount(filter bson.D) (int, error) {
	countTeams, err := d.DB.Collection("teams").
		CountDocuments(context.Background(), filter)
	if err != nil {
		return -1, err
	}
	return int(countTeams), nil
}

// GetTeams gets a page of teams and the total number of results
func (d *DBDriver) GetTeams(pageNum, pageSize int, criteria Search, isSearching bool) (teams []Team, count int, err error) {
	ctx := context.Background()
	hKey, hField := criteria.OrderBy+":all", fmt.Sprintf("%d", pageNum)

	// support redis cache for a page of data in non-searching scenario ONLY
	if !isSearching {
		// get teams cache
		val, err := Rdb.HGet(ctx, hKey, hField).Result()
		switch {
		case err == redis.Nil || val == "":
			zap.L().Warn("no such key", zap.String("redis.key", hKey))
		case err != nil:
			zap.L().Error("an error occurred when accessing the key and field", zap.String("redis.key", hKey), zap.String("redis.field", hField), zap.Error(err))
		default:
			zap.L().Info("hit the key", zap.String("redis.key", hKey))
			if err := json.Unmarshal([]byte(val), &teams); err != nil {
				return nil, 0, err
			}
		}

		// get count cache
		count, err = Rdb.Get(ctx, Total).Int()
		switch {
		case err == redis.Nil || val == "":
			zap.L().Warn("no such key", zap.String("redis.key", Total))
		case err != nil:
			zap.L().Error("an error occurred when accessing the key", zap.String("redis.key", Total), zap.Error(err))
		default:
			zap.L().Info("hit the key", zap.String("redis.key", Total))
		}

		// return if hit the cache
		if len(teams) > 0 && count > 0 {
			zap.L().Info("hit a page", zap.Int("pageNum", pageNum), zap.String("order", criteria.OrderBy), zap.Int("total", count))
			return teams, count, nil
		}
	}

	// get skip number
	var skip int64
	if pageNum > 0 {
		skip = int64((pageNum - 1) * pageSize)
	}

	// options
	opts := options.Find()
	opts.SetLimit(int64(pageSize))
	opts.SetSkip(skip)
	if criteria.OrderBy == "likes" {
		opts.SetSort(bson.D{{Key: "likes", Value: -1}}) // order by likes dec
	} else {
		opts.SetSort(bson.D{{Key: "created_at", Value: -1}}) // order by time dec
	}

	// filter
	filter := bson.D{
		{Key: "state", Value: 1},
	}
	if isSearching {
		if len(criteria.Format) > 0 {
			filter = append(filter, bson.E{Key: "format", Value: criteria.Format})
		}
		if len(criteria.Pokemon) > 0 {
			filter = append(filter, bson.E{Key: "pokemon", Value: bson.D{{Key: "$all", Value: criteria.Pokemon}}})
		}
		if criteria.HasShowdown {
			filter = append(filter, bson.E{Key: "has_showdown", Value: true})
		}
		if criteria.HasRental {
			filter = append(filter, bson.E{Key: "has_rental", Value: true})
		}
	}

	// get the number of teams
	count, err = d.GetTeamsCount(filter)
	if err != nil {
		return nil, 0, err
	}

	// query
	paginatedCursor, err := d.DB.Collection("teams").Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, err
	}

	// unmarshal retrieved data to struct and append to list
	if err = paginatedCursor.All(ctx, &teams); err != nil {
		return nil, 0, err
	}

	// set redis cache in non-searching scenario ONLY
	if !isSearching {
		redisValue, err := json.MarshalToString(&teams)
		if err != nil {
			return nil, 0, err
		}
		_, err = Rdb.Pipelined(ctx, func(pipe redis.Pipeliner) error {
			if err := Rdb.HSet(ctx, hKey, []string{hField, redisValue}).Err(); err != nil {
				return err
			}
			zap.L().Info("set a key", zap.String("redis.key", hKey))
			if err := Rdb.Set(ctx, Total, count, 0).Err(); err != nil {
				return err
			}
			zap.L().Info("set a key", zap.String("redis.key", Total))
			return nil
		})
		if err != nil {
			return nil, 0, err
		}
	}

	return teams, count, nil
}

// GetTeamByID gets a team by `_id`
func (d *DBDriver) GetTeamByID(id string) (*Team, error) {
	var team *Team
	ctx := context.Background()

	// redis cache
	redisKey := "team:" + id
	val, err := Rdb.Get(ctx, redisKey).Result()
	switch {
	case err == redis.Nil || val == "":
		zap.L().Warn("no such key", zap.String("redis.key", redisKey))
	case err != nil:
		zap.L().Error("accessing the key error", zap.String("redis.key", redisKey), zap.Error(err))
	default:
		zap.L().Info("hit the key", zap.String("redis.key", redisKey))
		if err := json.Unmarshal([]byte(val), &team); err != nil {
			return nil, err
		}
		return team, nil
	}

	// database
	hex, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	err = d.DB.Collection("teams").
		FindOne(ctx, bson.M{"_id": hex}).
		Decode(&team)
	if err != nil {
		return nil, err
	}

	// set redis cache
	redisValue, err := json.MarshalToString(&team)
	if err != nil {
		return nil, err
	}
	if err := Rdb.SetEX(ctx, redisKey, redisValue, time.Duration(config.Redis.Expiry)*time.Second).Err(); err != nil {
		return nil, err
	}
	zap.L().Info("set a key", zap.String("redis.key", redisKey))
	return team, nil
}

// GetPokemonUsageByFormat gets pokemon usage by format
func (d *DBDriver) GetPokemonUsageByFormat(format string) ([]Usage, error) {
	unwindStage := bson.D{{Key: "$unwind", Value: bson.D{{Key: "path", Value: "$pokemon"}, {Key: "preserveNullAndEmptyArrays", Value: false}}}}
	matchStage0 := bson.D{{Key: "$match", Value: bson.D{{Key: "format", Value: format}}}}
	groupStage := bson.D{{Key: "$group", Value: bson.D{{Key: "_id", Value: "$pokemon"}, {Key: "count", Value: bson.D{{Key: "$sum", Value: 1}}}}}}
	matchStage1 := bson.D{{Key: "$match", Value: bson.D{{Key: "count", Value: bson.D{{Key: "$gt", Value: 1}}}}}}
	sortStage := bson.D{{Key: "$sort", Value: bson.D{{Key: "count", Value: -1}}}}

	cursor, err := d.DB.Collection("teams").Aggregate(context.Background(), mongo.Pipeline{unwindStage, matchStage0, groupStage, matchStage1, sortStage})
	if err != nil {
		return nil, err
	}
	var usages []Usage
	if err = cursor.All(context.Background(), &usages); err != nil {
		return nil, err
	}
	return usages, nil
}

// GetLikedTeamsByUsername gets a user's liked teams
func (d *DBDriver) GetLikedTeamsByUsername(pageNum, pageSize int, username string) ([]Team, int, error) {
	// get skip number
	var skip int64
	if pageNum > 0 {
		skip = int64((pageNum - 1) * pageSize)
	}

	// options
	opts := options.Find()
	opts.SetLimit(int64(pageSize))
	opts.SetSkip(skip)

	// get liked teams by user
	user, err := d.GetUserByUsername(username)
	if err != nil {
		return nil, -1, err
	}
	var like []primitive.ObjectID
	for _, l := range user.Like {
		hex, err := primitive.ObjectIDFromHex(l)
		if err != nil {
			continue
		}
		like = append(like, hex)
	}

	// filter
	filter := bson.D{
		{Key: "state", Value: 1},
		{Key: "_id", Value: bson.D{{Key: "$in", Value: like}}},
	}

	// get the number of teams
	count, err := d.GetTeamsCount(filter)
	if err != nil {
		return nil, -1, err
	}

	// query
	paginatedCursor, err := d.DB.Collection("teams").Find(context.Background(), filter, opts)
	if err != nil {
		return nil, -1, err
	}

	// unmarshal retrieved data to struct and append to list
	var res []Team
	if err = paginatedCursor.All(context.Background(), &res); err != nil {
		return nil, -1, err
	}
	return res, count, nil
}

// GetUploadedTeamsByUsername gets uploaded teams by user name
func (d *DBDriver) GetUploadedTeamsByUsername(pageNum, pageSize int, username string) ([]Team, int, error) {
	// get skip number
	var skip int64
	if pageNum > 0 {
		skip = int64((pageNum - 1) * pageSize)
	}

	// options
	opts := options.Find()
	opts.SetLimit(int64(pageSize))
	opts.SetSkip(skip)

	// filter
	filter := bson.D{
		{Key: "state", Value: 1},
		{Key: "uploader", Value: username},
	}

	// get the number of teams
	count, err := d.GetTeamsCount(filter)
	if err != nil {
		return nil, -1, err
	}

	// query
	paginatedCursor, err := d.DB.Collection("teams").Find(context.Background(), filter, opts)
	if err != nil {
		return nil, -1, err
	}

	// unmarshal retrieved data to struct and append to list
	var res []Team
	if err = paginatedCursor.All(context.Background(), &res); err != nil {
		return nil, -1, err
	}
	return res, count, nil
}
