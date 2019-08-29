package apis

import (
	"encoding/json"
	"net/http"
	"orderFunc/middlewares"
	"orderFunc/models"
	"strconv"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

func MemberApiInit(router *gin.Engine) {
	memberRouter := router.Group("/member")
	memberRouter.Use(middlewares.LoginAuth())
	{
		// Get all the members
		memberRouter.GET("/", func(c *gin.Context) {

			skip, limit, sort := c.Query("skip"), c.Query("limit"), c.Query("sort")

			findOptions := options.Find()
			if strings.TrimSpace(skip) != "" {
				inputSkip, err := strconv.ParseInt(skip, 10, 64)
				if err != nil {
					c.JSON(http.StatusBadRequest, gin.H{
						"err":  err,
						"msg":  "Cannot setup skip",
						"skip": skip,
					})
					return
				}
				findOptions.SetSkip(inputSkip)
			}

			if strings.TrimSpace(limit) != "" {
				inputLimit, err := strconv.ParseInt(limit, 10, 64)
				if err != nil {
					c.JSON(http.StatusBadRequest, gin.H{
						"err":   err,
						"msg":   "Cannot setup limit",
						"limit": limit,
					})
					return
				}
				findOptions.SetLimit(inputLimit)
			}

			sortMap := map[string]int{}
			if strings.TrimSpace(sort) != "" {
				if s := strings.Split(sort, "_"); len(s) == 2 {
					sortOrd, err := strconv.Atoi(s[1])
					if err != nil {
						c.JSON(http.StatusBadRequest, gin.H{
							"err":  err,
							"msg":  "Cannot get the sort order",
							"s[1]": s[1],
							"s[0]": s[0],
						})
						return
					}
					// fmt.Printf("s is %+v\n", s)
					// fmt.Printf("s[0] is %+v\n", s[0])
					// fmt.Printf("s[1] is %+v\n", s[1])
					// fmt.Printf("sortOrd is %+v\n", sortOrd)

					sortMap[s[0]] = sortOrd
				} else {
					sortMap[sort] = -1
				}

				findOptions.SetSort(sortMap)
			}

			members, err := models.FindMembers(bson.M{}, findOptions)

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
