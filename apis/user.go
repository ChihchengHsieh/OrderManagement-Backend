package apis

import (
	"net/http"
	"orderFunc/models"
	"orderFunc/utils"
	"os"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

func UserApiInit(router *gin.Engine) {
	userRouter := router.Group("/user")
	{
		userRouter.POST("/signup", func(c *gin.Context) {

			// Check if the email already exist

			code := c.PostForm("code")

			var role string

			if !(code == os.Getenv("REGISTER_CODE") || code == os.Getenv("ADMIN_REGISTER_CODE")) {
				c.JSON(http.StatusUnauthorized, gin.H{
					"error": "The register code is not correct",
					"msg":   "The register code is not correct",
				})
				return
			} else {
				if code == os.Getenv("REGISTER_CODE") {
					role = "normal"
				} else if code == os.Getenv("ADMIN_REGISTER_CODE") {
					role = "admin"
				}
			}

			if c.PostForm("email") == "" || c.PostForm("password") == "" {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": "Must Provide Email and Password",
					"msg":   "Must Provide Email and Password",
				})
				return
			}

			if exist := len(models.FindUsers(bson.M{"email": c.PostForm("email")})); exist > 0 {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": "This Eamil has been used",
					"msg":   "This Eamil has been used",
				})
				return
			}

			// hash the password
			hashedPassword, err := bcrypt.GenerateFromPassword([]byte(c.PostForm("password")), bcrypt.DefaultCost)

			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": err,
					"msg":   "Cannot hash the password properly",
				})
				return
			}

			registerUser := models.User{
				Email:    c.PostForm("email"),
				Password: string(hashedPassword),
				Role:     role,
			}

			insertID, err := models.AddUser(&registerUser)

			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err, "msg": "Fail to register the user"})
				return
			}

			registerUser.ID = insertID.(primitive.ObjectID)
			authToken, err := utils.GenerateAuthToken(registerUser.ID.Hex())

			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": err,
					"msg":   "Cannot get the Auth Token",
				})
				return
			}

			registerUser.Password = "" // remove the password field

			c.JSON(http.StatusOK, gin.H{
				"user":  registerUser,
				"token": authToken,
			})

		})

		// Login Route

		userRouter.POST("/login", func(c *gin.Context) {
			/*
				Require Field
					Email
					Password (Hashed)
			*/

			inputEmail := c.PostForm("email")
			inputPassword := c.PostForm("password") // Expect Hased Password

			if inputEmail == "" || inputPassword == "" {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": "Must Provide Email and Password",
					"msg":   "Must Provide Email and Password",
				})
				return
			}

			// How can we check if the user exist or not
			user, err := models.CheckingTheAuth(inputEmail, inputPassword)

			if err != nil {
				c.JSON(http.StatusUnauthorized, gin.H{
					"error": err,
					"msg":   "Email or Password is not correct",
					"user":  user,
				})
				return
			}

			authToken, err := utils.GenerateAuthToken(user.ID.Hex())

			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": err,
					"msg":   "Cannot get the Auth Token",
				})
			}

			c.JSON(http.StatusOK, gin.H{
				"token": authToken,
				"user":  user,
			})

		})

		// Get User through id
		userRouter.GET("/id/:id", func(c *gin.Context) {
			id := c.Param("id")
			user, err := models.FindUserByID(id)

			if err != nil {
				c.JSON(http.StatusNotFound, gin.H{
					"error": err,
					"msg":   "User Not Found",
				})
				return
			}

			c.JSON(http.StatusOK, gin.H{
				"user": user,
			})
		})

		// Get User through Email

		userRouter.GET("/email/:email", func(c *gin.Context) {
			email := c.Param("email")
			user, err := models.FindUserByEmail(email)

			if err != nil {
				c.JSON(http.StatusNotFound, gin.H{
					"error": err,
					"msg":   "User not found",
				})
				return
			}

			c.JSON(http.StatusOK, gin.H{
				"user": user,
			})
		})
	}

}
