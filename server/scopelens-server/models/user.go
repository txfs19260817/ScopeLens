package models

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// The User holds
type User struct {
	ID       primitive.ObjectID   `bson:"_id"      json:"id"`
	UserName string               `bson:"username" json:"username" binding:"required"`
	Email    string               `bson:"email"    json:"email"    binding:"required,email"`
	Password string               `bson:"password" json:"password" binding:"required"`
	Like     []string             `bson:"like"     json:"like"`
}

// GetUserByName returns the User struct pointer by the given name.
func (d *DBDriver) GetUserByUsername(username string) (*User, error) {
	var user *User
	err := d.DB.Collection("users").
		FindOne(context.Background(), bson.D{{Key: "username", Value: username}}).
		Decode(&user)
	if err != nil {
		return nil, err
	}
	return user, nil
}
