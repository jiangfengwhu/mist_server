package main

import (
	"github.com/gin-gonic/gin"
	"log"
)

func init() {
	getConfig()
	dial()
}
func main() {
	log.Println(globalConf.RecapSecure)
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	api := r.Group("/api")
	api.POST("regist", regist)
	api.POST("login", login)
	api.GET("logstatus", logstatus)
	api.GET("logout", logout)
	api.GET("activeAccount", activeAccount)
	api.PUT("updateinfo", Auth(), updateInfo)
	api.POST("createCollection", Auth(), createCollection)
	api.OPTIONS("uploadfile/:id", Auth(), checkBP)
	api.POST("uploadfile", Auth(), uploadFile)
	api.POST("addVideo", Auth(), addVideo)
	api.POST("changeVC", Auth(), changeCover)
	api.GET("getVideo", latestVideo)
	api.GET("getVideo/:id", getVideo)
	api.POST("changeAvatar", Auth(), changepic)
	api.GET("myvideo", Auth(), myVideo)
	api.PUT("delvideoc", Auth(), delvideoc)
	api.PUT("updatevc", Auth(), updateVC)
	api.PUT("updatevideo", Auth(), updateVideo)
	api.DELETE("delvideo", Auth(), delvideo)
	api.GET("myvideoall", Auth(), myVideoAll)
	api.POST("uploadImage", Auth(), uploadImage)
	api.POST("addCircle", Auth(), addCircle)
	api.GET("getCircles", latestCircle)
	api.GET("mycommall", Auth(), myCommAll)
	api.DELETE("delcomms/:id", Auth(), delcomms)
	r.Run(":3000")
}
