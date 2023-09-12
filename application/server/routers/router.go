package routers

import (
	v1 "application/api/v1"
	"github.com/gin-gonic/gin"
)

// InitRouter 初始化路由信息
func InitRouter() *gin.Engine {
	r := gin.Default()

	apiV1 := r.Group("/api/v1")
	{
		apiV1.GET("/hello", v1.Hello)
		apiV1.POST("/queryAccountList", v1.QueryAccountList)
		apiV1.POST("/buySongrong", v1.BuySongrong)
		apiV1.POST("/confirmSongrong", v1.ConfirmSongrong)
		apiV1.POST("/querySellingBuyList", v1.QuerySellingBuyList)
		apiV1.POST("/querySellingConfirmList", v1.QuerySellingConfirmList)
		apiV1.POST("/sellSongrong", v1.SellSongrong)
		apiV1.POST("/querySellSongrong", v1.QuerySellSongrong)
		apiV1.POST("/uploadSongrong", v1.UploadSongrong)
		apiV1.POST("/queryUploadSongrong", v1.QueryUploadSongrong)
		apiV1.POST("/packingSongrong", v1.PackingSongrong)
		apiV1.POST("/queryPackingSongrong", v1.QueryPackingSongrong)
	}
	return r
}
