package apis

import (
	"encoding/json"
	"log"
	"net/http"
	"orderFunc/models"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

/*
type Order struct {
	ID          primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	Receiver    string             `json:"receiver" bson:"receiver"`
	CreatedAt   time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt   time.Time          `json:"updatedAt" bson:"updatedAt"`
	Products    []interface{}      `json:"products" bson:"products"`
	Paid        bool               `json:"paid" bson:"paid"`
	SendingDate time.Time          `json:"sendingDate" bson:"sendingDate"`
}
*/

// Try array of products

func OrderApiInit(router *gin.Engine) {
	orderRouter := router.Group("/order")
	{
		// Adding Comment
		orderRouter.POST("/", func(c *gin.Context) {
			receiver, productsJSON := c.PostForm("receiver"), c.PostForm("products")
			log.Printf("The receiver is %s and product is %s", receiver, productsJSON)

			var products []map[string]interface{}
			// var products []models.Product

			err := json.Unmarshal([]byte(productsJSON), &products)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": err,
					"msg":   "Cannot get products array properly",
				})
				return
			}

			log.Printf("The receiver: %s", receiver)

			insertOrder := models.Order{
				Receiver:  receiver,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
				Products:  products,
			}

			insertID, err := models.AddOrder(insertOrder)

			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": err,
					"msg":   "Cannot add the order",
				})
				return
			}

			c.JSON(http.StatusOK, gin.H{
				"insertID":     insertID,
				"insertedOder": insertOrder,
			})

		})
	}

	// Updating

	orderRouter.PUT("/:id", func(c *gin.Context) {

		id, receiver, productsJSON := c.Param("id"), c.PostForm("receiver"), c.PostForm("products")

		var products []map[string]interface{}

		err := json.Unmarshal([]byte(productsJSON), &products)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err,
				"msg":   "Cannot get products array properly",
			})
		}

		resultID, err := models.UpadteOrderByID(id, bson.M{"$set": bson.M{"receiver": receiver, "products": products, "updatedAt": time.Now()}})

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err,
				"msg":   "Cannot update the order",
			})
		}

		c.JSON(http.StatusOK, gin.H{
			"updatedID": resultID,
			"orderID":   id,
		})

	})

	// Deleting

	orderRouter.DELETE("/:id", func(c *gin.Context) {
		id := c.Param("id")
		err := models.DeleteOrderByID(id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err,
				"msg":   "Cannot Delete the order",
			})
		}

		c.JSON(http.StatusOK, gin.H{
			"deletedID": id,
		})

	})

	// Finding All

	orderRouter.GET("/", func(c *gin.Context) {
		orders, err := models.FindOrders(bson.M{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err,
				"msg":   "Cannot get all orders properly",
			})
		}
		// fmt.Printf("The orders: \n %+v", orders[0])

		c.JSON(http.StatusOK, gin.H{
			"orders": orders,
		})
	})

	// Finding SingleOne

	orderRouter.GET("/:id", func(c *gin.Context) {
		order, err := models.FindOrderByID(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err,
				"msg":   "Cannot find the order",
			})
		}

		c.JSON(http.StatusOK, gin.H{
			"order": order,
		})
	})

}
