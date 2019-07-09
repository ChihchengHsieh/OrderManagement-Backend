package models

import (
	"context"
	"orderFunc/databases"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Member struct {
	ID       primitive.ObjectID       `json:"_id" bson:"_id,omitempty"`
	Name     string                   `json:"name" bson:"name"`
	Remark   string                   `json:"remark" bson:"remark"`
	Products []map[string]interface{} `json:"products" bson:"products"`
}

type Product struct {
	ID           primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	Name         string             `json:"name" bson:"name"`
	Quantity     int                `json:"quantity" bson:"quantity"`
	BuyPriceAUD  float32            `json:"buyPriceAUD" bson:"buyPriceAUD"`
	SellPriceTWD float32            `json:"sellPriceTWD" bson:"sellPriceTWD"`
	Seller       string             `jsonp:"seller" bson:"seller"`
	Paid         bool               `json:"paid" bson:"paid"`
	Bought       bool               `json:"bought" bson:"bought"`
	Received     bool               `json:"received" bson:"received"`
	CreatedAt    time.Time          `json:"createdAt" bson:"createdAt"`
}

func AddMember(inputMember Member) (interface{}, error) {
	result, err := databases.DB.Collection("member").InsertOne(context.TODO(), inputMember)

	if err != nil {
		return nil, err
	}
	return result.InsertedID, nil
}

func UpdateMemberByID(id string, updateDetail bson.M) (interface{}, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	result, err := databases.DB.Collection("member").UpdateOne(context.TODO(), bson.M{"_id": oid}, updateDetail)

	if err != nil {
		return nil, err
	}

	return result.UpsertedID, nil
}

func DeleteMemberByID(id string) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = databases.DB.Collection("member").DeleteOne(context.TODO(), bson.M{"_id": oid})

	if err != nil {
		return err
	}

	return nil
}

func FindMemberByID(id string) (*Member, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	var member Member
	err = databases.DB.Collection("member").FindOne(context.TODO(), bson.M{"_id": oid}).Decode(&member)

	if err != nil {
		return nil, err
	}

	return &member, nil
}

func FindMembers(filterDetail bson.M) ([]*Member, error) {
	var members []*Member
	result, err := databases.DB.Collection("member").Find(context.TODO(), filterDetail)
	defer result.Close(context.TODO())

	if err != nil {
		return nil, err
	}

	for result.Next(context.TODO()) {
		var elem Member
		err := result.Decode(&elem)
		if err != nil {
			return nil, err
		}
		members = append(members, &elem)
	}

	return members, nil
}

func AddProductToMemberByID(id string, inputProduct map[string]interface{}) (interface{}, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	result, err := databases.DB.Collection("member").UpdateOne(context.TODO(),
		bson.M{"_id": oid},
		bson.M{"$push": bson.M{"products": inputProduct}})
	if err != nil {
		return nil, err
	}
	return result.UpsertedID, nil
}

func UpdateProductToMemberByID(memberID string, productID string, newProduct map[string]interface{}) (interface{}, error) {
	memebrOid, err := primitive.ObjectIDFromHex(memberID)
	if err != nil {
		return nil, err
	}

	productOid, err := primitive.ObjectIDFromHex(productID)
	if err != nil {
		return nil, err
	}

	updateFields := bson.M{}

	for k, v := range newProduct {
		updateFields["products.$."+k] = v
	}

	result, err := databases.DB.Collection("member").UpdateOne(context.TODO(),
		bson.M{"_id": memebrOid, "products._id": productOid},
		bson.M{"$set": updateFields})
	// bson.M{"$set": bson.M{"products.$": newProduct}})
	if err != nil {
		return nil, err
	}

	return result.UpsertedID, nil

}

func DeleteProductFromMemberByID(memberID string, productID string) (interface{}, error) {
	memebrOid, err := primitive.ObjectIDFromHex(memberID)
	if err != nil {
		return nil, err
	}

	productOid, err := primitive.ObjectIDFromHex(productID)
	if err != nil {
		return nil, err
	}

	result, err := databases.DB.Collection("member").UpdateOne(context.TODO(),
		bson.M{"_id": memebrOid},
		bson.M{"$pull": bson.M{"products": bson.M{"_id": productOid}}})

	if err != nil {
		return nil, err
	}

	return result.UpsertedID, nil
}
