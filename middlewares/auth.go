package middlewares

import (
	"fmt"
	"log"
	"net/http"
	"orderFunc/models"
	"os"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func LoginAuth() gin.HandlerFunc {
	return func(c *gin.Context) {

		tokenStr := c.GetHeader("Authorization")

		// log.Printf("tokenStr: %+v", tokenStr)

		if tokenStr == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "No token Provided",
				"msg":   "Token is needed",
			})
			return
		}

		// Check if it use Bearer

		if s := strings.Split(tokenStr, " "); len(s) == 2 {
			tokenStr = s[1]
		}

		// Problem: token is generatted but not valid

		// Add Claim Latter
		token, _ := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid Token", "msg": "Cannot parse the given token"})
				return nil, fmt.Errorf("Invalid Token")
			}

			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		// fmt.Printf("The token:\n %+v\n", token)

		// fmt.Printf("Claims:\n%+v\n", token.Claims.(jwt.MapClaims))

		// Find the User and store in c

		claims := token.Claims.(jwt.MapClaims)
		inputClaim := claims["_id"].(string)
		user, err := models.FindUserByID(inputClaim)

		if err != nil {
			log.Print(err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "User Not Authorised",
				"msg":   "Cannot find this user",
			})
			return
		}

		c.Set("user", user)

		if !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"getToken": tokenStr,
				"error":    "Token is not valid",
				"msg":      "The token is not valid",
			})
			return
		}

		c.Next()
		// t := time.Now()
	}
}
