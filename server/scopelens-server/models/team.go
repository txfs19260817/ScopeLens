package models

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

// The Team holds
type Team struct {
	ID          primitive.ObjectID `bson:"_id" json:"id"`
	// User-defined
	Title       string             `bson:"title" json:"title" binding:"required"`
	Author      string             `bson:"author" json:"author" binding:"required"`
	Format      string             `bson:"format" json:"format" binding:"required"`
	Pokemon     []string           `bson:"pokemon" json:"pokemon" binding:"required"`
	Showdown    string             `bson:"showdown" json:"showdown"`
	Image       string             `bson:"image" json:"image"` // upload from client: base64; store in DB: URL
	Description string             `bson:"description" json:"description"`
	// Auto-generated
	Uploader    string             `bson:"uploader" json:"uploader"` // User.username
	CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
	Likes       int                `bson:"likes" json:"likes"` // 0
	State       int                `bson:"state" json:"state"` // 0
}

// Insert a team
func (d *DBDriver) InsertTeam(team Team) (bool, error){
	// Replenish fields
	team.ID = primitive.NewObjectID()
	team.CreatedAt = time.Now()

	// Insert
	if _, err := d.DB.Collection("teams").InsertOne(context.Background(), team); err != nil {
		return false, err
	}
	return true, nil
}