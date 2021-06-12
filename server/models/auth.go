package models

import (
	"context"
	"fmt"

	"github.com/txfs19260817/scopelens/server/utils/encrypt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Login struct {
	UserName string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// Register inserts new user to DB
func (d *DBDriver) Register(user User) (bool, error) {
	var err error
	user.ID = primitive.NewObjectID()
	user.Password, err = encrypt.PasswordEncrypt(user.Password)
	if err != nil {
		return false, err
	}
	// Insert
	_, err = d.DB.Collection("users").InsertOne(context.Background(), user)
	if err != nil {
		return false, err
	}
	return true, nil
}

// LoginValidate validates username and password
func (d *DBDriver) LoginValidate(loginReq Login) (bool, error) {
	// Find password
	var res bson.M
	opt := options.FindOne().SetProjection(bson.M{"_id": 0, "password": 1})
	err := d.DB.Collection("users").
		FindOne(context.Background(), bson.M{"username": loginReq.UserName}, opt).
		Decode(&res)
	if err != nil {
		return false, fmt.Errorf("user is not found")
	}

	// Verify password
	hashedPassword := res["password"].(string)
	if err := encrypt.PasswordVerification(hashedPassword, loginReq.Password); err != nil {
		return false, fmt.Errorf("password is not correct: %w", err)
	}
	return true, nil
}

// CheckUsernameAvailability checks username availability
func (d *DBDriver) CheckUsernameAvailability(username string) (bool, error) {
	// Check if the username has already existed.
	countUsers, err := d.DB.Collection("users").CountDocuments(context.Background(), bson.M{"username": username})
	if err != nil {
		return false, err
	}
	return countUsers == 0, nil
}
