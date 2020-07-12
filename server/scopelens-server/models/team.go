package models

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
	Uploader  string    `bson:"uploader" json:"uploader"` // User.username
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	Likes     int       `bson:"likes" json:"likes"` // 0
	State     int       `bson:"state" json:"state"` // 0
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
		filename, err = showdown.RentalTeamMaker(team.Showdown)
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
