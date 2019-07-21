package apis

import (
	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	router := gin.Default()
	UserApiInit(router)
	OrderApiInit(router)
	MemberApiInit(router)

	return router
}
