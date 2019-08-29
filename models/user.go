package models

import (
	"DiscussionBoard/utils"
	"context"
	"log"
	"orderFunc/databases"

	"golang.org/x/crypto/bcrypt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type User struct {
	ID       primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	Email    string             `json:"email" bson:"email"`
	Password string             `json:"password" bson:"password"`
	Role     string             `json:"role" bson:"role"`
}

var projectionForRemovingPassword = bson.D{
	{"password", 0},
}

func AddUser(inputUser *User) (interface{}, error) {

	result, err := databases.DB.Collection("user").InsertOne(context.TODO(), inputUser)

	return result.InsertedID, err
}

func UpdateUsers(filterDetail bson.M, updateDetail bson.M) (interface{}, error) {
	result, err := databases.DB.Collection("user").UpdateMany(context.TODO(), filterDetail, updateDetail)

	if err != nil {
		return nil, err
	}

	return result.UpsertedID, nil
}

// FindUserByID - this function will not return password
func FindUserByID(id string) (interface{}, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var user User
	err = databases.DB.Collection("user").FindOne(context.TODO(), bson.M{"_id": oid},
		options.FindOne().SetProjection(projectionForRemovingPassword)).Decode(&user)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func FindUserByEmail(email string) (interface{}, error) {
	var user User

	err := databases.DB.Collection("user").FindOne(context.TODO(), bson.M{"email": email},
		options.FindOne().SetProjection(projectionForRemovingPassword)).Decode(&user)

	if err != nil {
		return nil, err
	}

	return user, nil
}

// CheckingTheAuth - if the user is returned, then the auth is valid
func CheckingTheAuth(email string, password string) (*User, error) {
	var user User
	err := databases.DB.Collection("user").FindOne(context.TODO(), bson.M{"email": email}).Decode(&user)

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))

	if err != nil {
		return nil, err
	}
	return &user, nil
}

// FindUsers - this function will get multiple users without the password field
func FindUsers(filterDetail bson.M) []*User {
	var users []*User
	result, err := databases.DB.Collection("user").Find(context.TODO(), filterDetail,
		options.Find().SetProjection(projectionForRemovingPassword))

	if err != nil {
		log.Fatal(err)
	}
	defer result.Close(context.TODO())

	for result.Next(context.TODO()) {
		var elem User
		err := result.Decode(&elem)
		utils.ErrorChecking(err)
		users = append(users, &elem)
	}

	return users

}
