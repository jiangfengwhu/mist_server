package main

import (
	"log"

	"github.com/gin-gonic/gin"
)

func init() {
	getConfig()
	dial()
}
func main() {
	log.Println(globalConf.RecapSecure)
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	hub := newHub()
	go hub.run()
	api := r.Group("/api")
	api.POST("regist", regist)
	api.POST("login", login)
	api.GET("logstatus", logstatus)
	api.GET("logout", logout)
	api.GET("activeAccount", activeAccount)
	api.PUT("updateinfo", Auth(), updateInfo)
	api.GET("user/:id", getUser)
	api.POST("createCollection", Auth(), createCollection)
	api.GET("uploadfile/:id", Auth(), checkBP)
	api.POST("uploadfile", Auth(), uploadFile)
	api.POST("addVideo", Auth(), addVideo)
	api.POST("changeVC", Auth(), changeCover)
	api.GET("getVideo", latestVideo)
	api.GET("getVideo/:id", CheckGuest(), getVideo)
	api.POST("changeAvatar", Auth(), changepic)
	api.GET("myvideo", Auth(), myVideo)
	api.PUT("delvideoc", Auth(), delvideoc)
	api.GET("checkvown/:id", Auth(), checkOWN)
	api.PUT("updatevc", Auth(), updateVC)
	api.PUT("updatevideo", Auth(), updateVideo)
	api.DELETE("delvideo", Auth(), delvideo)
	api.GET("videoall/:id", videoAll)
	api.POST("uploadImage", Auth(), uploadImage)
	api.POST("addCircle", Auth(), addCircle)
	api.GET("getCircles", CheckGuest(), latestCircle)
	api.POST("addGallery", Auth(), addGallery)
	api.GET("getGallery", CheckGuest(), latesetGallery)
	api.GET("oneGallery/:id", CheckGuest(), getGallery)
	api.GET("commall/:id", CheckGuest(), commAll)
	api.GET("galleryall/:id", galleryAll)
	api.DELETE("delcomms/:id", Auth(), delcomms)
	api.DELETE("delgallery/:id", Auth(), delGallery)
	api.POST("addComment", Auth(), addComment)
	api.GET("getcomments", CheckGuest(), getComments)
	api.POST("like", Auth(), setLike)
	api.GET("chat/:id", CheckGuest(), func(c *gin.Context) {
		serveWs(hub, c)
	})
	api.GET("listAll/:id", listAll)
	api.POST("newlist", Auth(), newList)
	api.PUT("addtolist", Auth(), addtoList)
	api.PUT("refromlist", Auth(), removeFromList)
	api.GET("searchVideo", searchVideo)
	r.Run(":3030")
}
