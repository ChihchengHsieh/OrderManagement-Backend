package models

import (
	"context"
	"orderFunc/databases"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Order - Struct for make an order
type Order struct {
	ID          primitive.ObjectID       `json:"_id" bson:"_id,omitempty"`
	Receiver    string                   `json:"receiver" bson:"receiver"`
	CreatedAt   time.Time                `json:"createdAt" bson:"createdAt"`
	UpdatedAt   time.Time                `json:"updatedAt" bson:"updatedAt"`
	Products    []map[string]interface{} `json:"products" bson:"products"`
	Paid        bool                     `json:"paid" bson:"paid"`
	SendingDate time.Time                `json:"sendingDate" bson:"sendingDate"`
}

// Defineing here

// THe map[string]string should have these fields

// type Product struct {
// 	Name         string
// 	Quantity     int
// 	BuyPriceAUD  float32
// 	SellPriceTWD float32
// 	Seller       string
// 	EarningTWD   float32
// }

// The Products will be an array storing map

// AddOrder - Adding one product by give a Order instance
func AddOrder(inputOrder Order) (interface{}, error) {
	result, err := databases.DB.Collection("order").InsertOne(context.TODO(), inputOrder)

	if err != nil {
		return nil, err
	}

	return result.InsertedID, nil
}

// UpadteOrderByID - upadte the Order through id
func UpadteOrderByID(id string, updateDetail bson.M) (interface{}, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	result, err := databases.DB.Collection("order").UpdateOne(context.TODO(), bson.M{"_id": oid}, updateDetail)

	if err != nil {
		return nil, err
	}

	return result.UpsertedID, nil
}

// DeleteOrderByID - Delete an Order Through the id
func DeleteOrderByID(id string) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = databases.DB.Collection("order").DeleteOne(context.TODO(), bson.M{"_id": oid})

	if err != nil {
		return err
	}

	return nil
}

// FindOrders - Find the comments correspond to the filterDetail
func FindOrders(filterDetail bson.M) ([]*Order, error) {
	var orders []*Order
	result, err := databases.DB.Collection("order").Find(context.TODO(), filterDetail)
	defer result.Close(context.TODO())

	if err != nil {
		return nil, err
	}

	for result.Next(context.TODO()) {
		var elem Order
		err := result.Decode(&elem)
		if err != nil {
			return nil, err
		}
		orders = append(orders, &elem)
	}

	return orders, nil
}

// FindOrderByID - Find an order through id
func FindOrderByID(id string) (*Order, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	var order Order
	err = databases.DB.Collection("order").FindOne(context.TODO(), bson.M{"_id": oid}).Decode(&order)

	if err != nil {
		return nil, err
	}

	return &order, nil
}
