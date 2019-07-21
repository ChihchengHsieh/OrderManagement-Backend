package apis

import (
	"encoding/json"
	"net/http"
	"orderFunc/middlewares"
	"orderFunc/models"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

func MemberApiInit(router *gin.Engine) {
	memberRouter := router.Group("/member")
	memberRouter.Use(middlewares.LoginAuth())
	{
		// Get all the members
		memberRouter.GET("/", func(c *gin.Context) {

			members, err := models.FindMembers(bson.M{})

			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": err,
					"msg":   "Cannot retrieve all members",
				})
				return
			}

			c.JSON(http.StatusOK, gin.H{
				"members": members,
			})

		})

		// Get Single member
		memberRouter.GET("/:uid", func(c *gin.Context) {

			member, err := models.FindMemberByID(c.Param("uid"))

			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": err,
					"msg":   "Cannot retrieve this user",
				})
				return
			}

			c.JSON(http.StatusOK, gin.H{
				"member": member,
			})

		})

		// Create a memeber
		memberRouter.POST("/", func(c *gin.Context) {

			newMember := models.Member{
				Name:     c.PostForm("name"),
				Remark:   "",
				Products: []map[string]interface{}{},
			}
			insertID, err := models.AddMember(newMember)

			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": err,
					"msg":   "Cannot create this user",
				})
				return
			}

			c.JSON(http.StatusOK, gin.H{
				"insertID":     insertID,
				"insertMember": newMember,
			})

		})

		// Update a memeber
		memberRouter.PUT("/:uid", func(c *gin.Context) {

			updateJSON := c.PostForm("member")
			var updateMember map[string]interface{}
			err := json.Unmarshal([]byte(updateJSON), &updateMember)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": err,
					"msg":   "Cannot Unmarshal the input",
				})
				return
			}

			// memberAcceoptFields := []string{"name", "remark"}

			// updateingField := bson.M{}

			// for _, f := range memberAcceoptFields {
			// 	if k := c.PostForm(f); k != "" {
			// 		updateingField[f] = k
			// 	}
			// }

			upsertID, err := models.UpdateMemberByID(c.Param("uid"), bson.M{"$set": updateMember})

			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": err,
					"msg":   "Cannot update the member",
				})
				return
			}

			c.JSON(http.StatusOK, gin.H{
				"updatingFields": updateMember,
				"upsertID":       upsertID,
			})

		})

		// Delete a member
		memberRouter.DELETE("/:uid", func(c *gin.Context) {
			err := models.DeleteMemberByID(c.Param("uid"))

			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": err,
					"msg":   "Cannot delete the member",
				})
				return
			}

			c.JSON(http.StatusOK, gin.H{
				"removedID": c.Param("uid"),
			})

		})

		productRouter := memberRouter.Group("/:uid/product")

		{
			// Don't need GET (find) since we can get the product from the member

			// Add a new product to a user
			productRouter.POST("/", func(c *gin.Context) {
				prodcutJSON := c.PostForm("product")

				var product map[string]interface{}

				err := json.Unmarshal([]byte(prodcutJSON), &product)

				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{
						"error": err,
						"msg":   "Cannot get the product",
					})
					return
				}

				product["createdAt"] = time.Now()
				product["_id"] = primitive.NewObjectID()

				upsertID, err := models.AddProductToMemberByID(c.Param("uid"), product)

				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{
						"error": err,
						"msg":   "Cannot add the product properly",
					})
					return
				}

				c.JSON(http.StatusOK, gin.H{
					"upsertedID":   upsertID,
					"addedProduct": product,
					"uid":          c.Param("uid"),
				})

			})

			// update a product for a user
			productRouter.PUT("/:pid", func(c *gin.Context) {
				updateJSON := c.PostForm("product")
				var updateProduct map[string]interface{}
				err := json.Unmarshal([]byte(updateJSON), &updateProduct)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{
						"error": err,
					})
					return
				}

				upsertID, err := models.UpdateProductToMemberByID(c.Param("uid"), c.Param("pid"), updateProduct)

				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{
						"error": err,
					})
					return
				}

				c.JSON(http.StatusOK, gin.H{
					"upsertID":       upsertID,
					"productID":      c.Param("pid"),
					"userID":         c.Param("uid"),
					"updatedProduct": updateProduct,
				})

			})

			// delete a product for a user
			productRouter.DELETE("/:pid", func(c *gin.Context) {
				upsertID, err := models.DeleteProductFromMemberByID(c.Param("uid"), c.Param("pid"))
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{
						"upsertID": upsertID,
					})
					return
				}

				c.JSON(http.StatusOK, gin.H{
					"upsertID":  upsertID,
					"productID": c.Param("pid"),
					"userID":    c.Param("uid"),
				})

			})

		}
	}
}
