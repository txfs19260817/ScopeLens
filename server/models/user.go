package models

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// The User holds
type User struct {
	ID       primitive.ObjectID `bson:"_id"      json:"id"`
	UserName string             `bson:"username" json:"username" binding:"required"`
	Email    string             `bson:"email"    json:"email"    binding:"required,email"`
	Password string             `bson:"password" json:"password" binding:"required"`
	Like     []string           `bson:"like"     json:"like"`
}

// GetUserByUsername returns the User struct pointer by the given name.
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

// InsertLikeByUsername adds one like onto the team according to the username
func (d *DBDriver) InsertLikeByUsername(username, id string) error {
	ctx := context.Background()

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
		return fmt.Errorf("You have already liked this team. ")
	}

	// clear redis cache
	redisKeys := []string{LikesOrderAll}
	if err := Rdb.Del(ctx, redisKeys...).Err(); err != nil {
		return err
	}

	// update user's likes list
	_, err = d.DB.Collection("users").
		UpdateOne(ctx,
			bson.D{{"username", username}},
			bson.D{{"$push", bson.D{{"like", id}}}})
	if err != nil {
		_, err = d.DB.Collection("users").
			UpdateOne(ctx,
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
		UpdateOne(ctx,
			bson.D{{"_id", _id}},
			bson.D{{"$inc", bson.D{{"likes", 1}}}})
	if err != nil {
		return err
	}
	return nil
}
