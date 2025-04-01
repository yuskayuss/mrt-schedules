package station

import (
	"github.com/gin-gonic/gin"
)

func Initiate(router *gin.RouterGroup){
	station :=router.Group("/station")
	station.GET("", func(c *gin.Context){

	})
}