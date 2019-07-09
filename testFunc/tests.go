package testFunc

import (
	"log"
	"orderFunc/models"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func AllTestingRun() {
	resultID, err := InitMemberTest()

	if err != nil {
		log.Fatal(err)
	}

	inputID := resultID.(primitive.ObjectID).Hex()

	log.Printf("This added id is : %+v", inputID)

	upsertID, err := models.UpdateMemberByID(inputID, bson.M{
		"$set": bson.M{"remark": "newRemark"},
	})

	if err != nil {
		log.Fatal(err)
	}

	log.Printf("The upsertID is: %+v", upsertID)

	member, err := models.FindMemberByID(inputID)

	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Retrieve Test: %+v", member)

	inputProductID := primitive.NewObjectID()

	inputProduct := map[string]interface{}{}
	inputProduct["_id"] = inputProductID
	inputProduct["name"] = "Product1"
	inputProduct["quantity"] = "q1"
	inputProduct["buyPriceAUD"] = 20
	inputProduct["sellPriceTWD"] = 550
	inputProduct["seller"] = "chemistwarehouse"
	inputProduct["paid"] = false
	inputProduct["bought"] = false
	inputProduct["received"] = false
	inputProduct["createdAt"] = time.Now()

	productUpsertID, err := models.AddProductToMemberByID(inputID, inputProduct)

	if err != nil {
		log.Fatal(err)
	}

	log.Printf("productUpsertID is: %+v", productUpsertID)

	updatedProduct := map[string]interface{}{}
	// updatedProduct["_id"] = inputProductID
	// updatedProduct["quantity"] = "q1"
	// updatedProduct["buyPriceAUD"] = 20
	updatedProduct["sellPriceTWD"] = 55055
	updatedProduct["seller"] = "chemistwarehouse_2"
	// updatedProduct["bough"] = false
	// updatedProduct["received"] = false
	// updatedProduct["createdAt"] = time.Now()
	updatedProduct["name"] = "Product1Update"
	updatedProduct["paid"] = true

	productUpdateID, err := models.UpdateProductToMemberByID(inputID, inputProductID.Hex(), updatedProduct)

	if err != nil {
		log.Fatal(err)
	}

	log.Printf("productUpdateID is: %+v", productUpdateID)

	productDeletedID, err := models.DeleteProductFromMemberByID(inputID, inputProductID.Hex())

	if err != nil {
		log.Fatal(err)
	}

	log.Printf("productDeletedID is: %+v", productDeletedID)

	err = models.DeleteMemberByID(inputID)

	if err != nil {
		log.Fatal(err)
	} else {
		log.Println("Successfully remove a member")
	}

}

func InitMemberTest() (interface{}, error) {
	initProduct := []map[string]interface{}{}

	resultID, err := models.AddMember(models.Member{
		Name:     "Szu",
		Remark:   "",
		Products: initProduct,
	})

	if err != nil {
		return nil, err
	}

	return resultID, nil
}
