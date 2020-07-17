package models

import (
	"context"
	"errors"
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

func (d *DBDriver) InsertLikeByUsername(username, id string) error  {
	// filter
	filter := bson.D{
		{"username", username},
		{"like", id},
	}
	count, err := d.DB.Collection("users").CountDocuments(context.Background(), filter)
	if err != nil {
		return err
	}
	if count != 0 {
		return errors.New("You have already liked this team. ")
	}

	// update
	_, err = d.DB.Collection("users").
		UpdateOne(context.Background(),
			bson.D{{"username", username}},
			bson.D{{"$push", bson.D{{"like", id}}}})
	if err != nil {
		_, err = d.DB.Collection("users").
			UpdateOne(context.Background(),
				bson.D{{"username", username}},
				bson.D{{"$set", bson.D{{"like", bson.A{id}}}}})
		if err != nil {
			return err
		}
	}

	// add 1 like to the team
	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	_, err = d.DB.Collection("teams").
		UpdateOne(context.Background(),
			bson.D{{"_id",_id}},
			bson.D{{"$inc", bson.D{{"likes", 1}}}})
	if err != nil {
		return err
	}
	return nil
}