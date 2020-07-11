package models

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// The User holds
type User struct {
	ID       primitive.ObjectID `bson:"_id"      json:"id"`
	UserName string             `bson:"username" json:"username" binding:"required"`
	Email    string             `bson:"email"    json:"email"    binding:"required,email"`
	Password string             `bson:"password" json:"password" binding:"required"`
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

// GetEmailByUsername returns the Email string by the given name.
func (d *DBDriver) GetEmailByUsername(username string) (string, error) {
	var res bson.M
	opt := options.FindOne().SetProjection(bson.M{"_id": 0, "email": 1})
	err := d.DB.Collection("users").
		FindOne(context.Background(), bson.D{{Key: "username", Value: username}}, opt).
		Decode(&res)
	if err != nil {
		return "", err
	}
	return res["email"].(string), nil
}
