package models

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"
	"scopelens-server/config"
	"scopelens-server/utils/file"
	"scopelens-server/utils/showdown"
	"scopelens-server/utils/storage"
	"strings"
	"time"
)

// The Team holds
type Team struct {
	ID primitive.ObjectID `bson:"_id" json:"id"`
	// User-defined
	Title       string   `bson:"title" json:"title" binding:"required"`
	Author      string   `bson:"author" json:"author" binding:"required"`
	Format      string   `bson:"format" json:"format" binding:"required"`
	Pokemon     []string `bson:"pokemon" json:"pokemon" binding:"required"`
	Showdown    string   `bson:"showdown" json:"showdown"`
	Image       string   `bson:"image" json:"image"` // upload from client: base64; store in DB: URL
	Description string   `bson:"description" json:"description"`
	// Auto-generated
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

// Insert a team
func (d *DBDriver) InsertTeam(team Team) (bool, error) {
	var err error
	var filename string // image save path

	// Replenish fields
	team.ID = primitive.NewObjectID()
	team.CreatedAt = time.Now()
	// Use uploaded image first
	if len(team.Image) != 0 {
		// decode uploaded base64 string to file
		filename, err = file.DecodeBase64AndSave(team.Image)
		if err != nil {
			return false, err
		}
	} else {
		if len(team.Showdown) == 0 {
			return false, fmt.Errorf("Image and Showdown cannot be empty at the same time. ")
		}
		// generate image from showdown text.
		filename, err = showdown.RentalTeamMaker(team.Showdown, team.Title, team.Author)
		if err != nil {
			return false, err
		}
	}
	// then upload to S3.
	s3, err := storage.NewAmazonS3(config.Aws.AccessKey, config.Aws.SecretKey, config.Aws.Region, config.Aws.Bucket)
	if err != nil {
		return false, err
	}
	f, err := os.Open(filename)
	if err != nil {
		return false, fmt.Errorf("failed to open file %q, %v", filename, err)
	}
	tempPath := strings.Split(filename, "/")
	uploadPath := config.Aws.TeamPath + "/" + tempPath[len(tempPath)-1]
	url, err := s3.Save(uploadPath, f)
	if err != nil {
		return false, err
	}
	// store image URL
	team.Image = url

	// Insert
	if _, err := d.DB.Collection("teams").InsertOne(context.Background(), team); err != nil {
		return false, err
	}
	return true, nil
}

// Count teams
func (d *DBDriver) GetTeamsCount(filter bson.D) (int, error) {
	countTeams, err := d.DB.Collection("teams").
		CountDocuments(context.Background(), filter)
	if err != nil {
		return -1, err
	}
	return int(countTeams), nil
}

// Get teams
func (d *DBDriver) GetTeams(pageNum, pageSize int, s Search) ([]Team, int, error) {
	// get skip number
	var skip int64
	if pageNum > 0 {
		skip = int64((pageNum - 1) * pageSize)
	}

	// options
	opts := options.Find()
	opts.SetLimit(int64(pageSize))
	opts.SetSkip(skip)
	if s.OrderBy == "likes" {
		opts.SetSort(bson.D{{"likes", -1}}) // order by likes dec
	} else {
		opts.SetSort(bson.D{{"created_at", -1}}) // order by time dec
	}

	// filter
	filter := bson.D{
		{"state", 1},
	}
	if len(s.Format) != 0 {
		filter = append(filter, bson.E{Key: "format", Value: s.Format})
	}
	if len(s.Pokemon) != 0 {
		filter = append(filter, bson.E{Key: "pokemon", Value: bson.D{{"$all", s.Pokemon}}})
	}
	if s.HasShowdown {
		filter = append(filter, bson.E{Key: "has_showdown", Value: true})
	}
	if s.HasRental {
		filter = append(filter, bson.E{Key: "has_rental", Value: true})
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

// Get a team by _id
func (d *DBDriver) GetTeamByID(id string) (*Team, error) {
	hex, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	var team *Team
	err = d.DB.Collection("teams").
		FindOne(context.Background(), bson.M{"_id": hex}).
		Decode(&team)
	if err != nil {
		return nil, err
	}
	return team, nil
}

// Get pokemon usage by format
func (d *DBDriver) GetPokemonUsageByFormat(format string) ([]Usage, error) {
	unwindStage := bson.D{{"$unwind", bson.D{{"path", "$pokemon"}, {"preserveNullAndEmptyArrays", false}}}}
	matchStage0 := bson.D{{"$match", bson.D{{"format", format}}}}
	groupStage := bson.D{{"$group", bson.D{{"_id", "$pokemon"}, {"count", bson.D{{"$sum", 1}}}}}}
	matchStage1 := bson.D{{"$match", bson.D{{"count", bson.D{{"$gt", 1}}}}}}
	sortStage := bson.D{{"$sort", bson.D{{"count", -1}}}}

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
		{"state", 1},
		{"_id", bson.D{{"$in", like}}},
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

// Get uploaded teams
func (d *DBDriver) GetUploadedTeamsByUsername(pageNum, pageSize int, username string, ) ([]Team, int, error) {
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
		{"state", 1},
		{"uploader", username},
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
