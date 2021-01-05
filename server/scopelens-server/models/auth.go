package models

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"scopelens-server/utils/encrypt"
	"scopelens-server/utils/logger"
)

type Login struct {
	UserName string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// Insert new user
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

// Check username and password
func (d *DBDriver) LoginValidate(loginReq Login) (bool, error) {
	// Find password
	var res bson.M
	opt := options.FindOne().SetProjection(bson.M{"_id": 0, "password": 1})
	err := d.DB.Collection("users").
		FindOne(context.Background(), bson.M{"username": loginReq.UserName}, opt).
		Decode(&res)
	if err != nil {
		logger.SugaredLogger.Error(err)
		return false, fmt.Errorf("user is not found")
	}

	// Verify password
	hashedPassword := res["password"].(string)
	if err := encrypt.PasswordVerification(hashedPassword, loginReq.Password); err != nil {
		logger.SugaredLogger.Error(err)
		return false, fmt.Errorf("password is not correct")
	}
	return true, nil
}

// Check username availability
func (d *DBDriver) CheckUsernameAvailability(username string) (bool, error) {
	// Check if the username has already existed.
	countUsers, err := d.DB.Collection("users").
		CountDocuments(context.Background(), bson.M{"username": username})
	if err != nil {
		return false, err
	}
	return countUsers == 0, nil
}
