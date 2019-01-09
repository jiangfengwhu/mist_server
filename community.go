package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/globalsign/mgo/bson"
	"log"
	"os"
	"path/filepath"
	"time"
)

func uploadImage(c *gin.Context) {
	uploadDir := globalConf.ResDir + "/community/"
	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(200, gin.H{"status": false, "msg": err.Error()})
		return
	}
	files := form.File["pics"]
	for _, file := range files {
		if err := c.SaveUploadedFile(file, uploadDir+file.Filename); err != nil {
			c.JSON(200, gin.H{"status": false, "msg": err.Error()})
			return
		}
	}
	c.JSON(200, gin.H{"status": true, "path": globalConf.ResRef + "/community/"})
	return
}
func addCircle(c *gin.Context) {
	var collection circleModel
	if err := c.ShouldBind(&collection); err != nil {
		log.Println(err)
		c.JSON(200, gin.H{"status": false, "msg": "信息不完整"})
		return
	}
	id := bson.NewObjectId()
	err := insertC("community", bson.M{"_id": id, "cont": collection.Content, "owner": bson.ObjectIdHex(c.MustGet("auth").(string)), "pics": collection.Pics, "date": fmt.Sprintf("%d", time.Now().Unix())})

	if err != nil {
		log.Println(err.Error())
		c.JSON(200, gin.H{"status": false, "msg": "数据库错误"})
		return
	}
	c.JSON(200, gin.H{"status": true, "cid": id, "msg": "发布成功"})
	return
}

func latestCircle(c *gin.Context) {
	var params getVModel
	if err := c.ShouldBind(&params); err != nil {
		c.JSON(200, gin.H{"status": false, "msg": err.Error()})
		return
	}
	var re []outCircleModel
	log.Println(*params.Start*params.Size, *params.Start*(params.Size+1))
	err := latestC("community", []bson.M{}, *params.Start*params.Size, params.Size*(*params.Start+1), &re)
	if err != nil {
		c.JSON(200, false)
		return
	}
	c.JSON(200, re)
	return
}
func myCommAll(c *gin.Context) {
	var re []myCircleModel
	err := findC("community", bson.M{"owner": bson.ObjectIdHex(c.MustGet("auth").(string))}, true, &re)
	if err != nil {
		c.JSON(200, false)
		return
	}
	c.JSON(200, re)
	return
}
func delcomms(c *gin.Context) {
	id := c.Param("id")
	if len(id) != 24 {
		log.Println(id)
		c.JSON(200, gin.H{"status": false, "msg": "not correct id"})
		return
	}
	var re circleModel
	err := findAndModify("community", bson.M{"_id": bson.ObjectIdHex(id), "owner": bson.ObjectIdHex(c.MustGet("auth").(string))}, nil, false, true, false, &re)
	if err != nil {
		c.JSON(200, gin.H{"status": false, "msg": err.Error()})
		return
	}
	log.Println(re)
	for _, val := range re.Pics {
		os.Remove(globalConf.ResDir + "/community/" + filepath.Base(val))
		os.Remove(globalConf.ResDir + "/community/" + filepath.Base(val) + "_thb.jpeg")
	}
	c.JSON(200, gin.H{"status": true, "msg": "删除状态成功"})
	return
}
