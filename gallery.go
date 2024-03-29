package main

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/globalsign/mgo/bson"
)

func addGallery(c *gin.Context) {
	var collection newGalleryModel
	c.ShouldBind(&collection)
	if len(collection.Pics) == 0 {
		c.JSON(200, gin.H{"status": false, "msg": "至少添加一个图片"})
		return
	}
	collection.ID = bson.NewObjectId()
	collection.Date = time.Now().Unix()
	collection.Owner = bson.ObjectIdHex(c.MustGet("auth").(string))
	err := insertC("gallery", collection)
	if err != nil {
		log.Println(err.Error())
		c.JSON(200, gin.H{"status": false, "msg": "数据库错误"})
		return
	}
	c.JSON(200, gin.H{"status": true, "cid": collection.ID, "msg": "发布成功"})
	return
}

func latesetGallery(c *gin.Context) {
	var params getOModel
	if err := c.ShouldBind(&params); err != nil {
		c.JSON(200, gin.H{"status": false, "msg": err.Error()})
		return
	}
	var re []faceGalleryModel
	log.Println(*params.Start*params.Size, *params.Start*(params.Size+1))
	err := pipiC("gallery", []bson.M{{"$sort": bson.M{"_id": -1}}, bson.M{"$limit": params.Size * (*params.Start + 1)}, bson.M{"$skip": *params.Start * params.Size}}, &re, true)
	if err != nil {
		c.JSON(200, false)
		return
	}
	c.JSON(200, re)
	return
}
func galleryAll(c *gin.Context) {
	id := c.Param("id")
	if len(id) != 24 {
		c.JSON(200, "wrong id")
		return
	}
	var re []outGalleryModel
	err := findC("gallery", bson.M{"owner": bson.ObjectIdHex(id)}, true, &re)
	if err != nil {
		c.JSON(200, false)
		return
	}
	c.JSON(200, re)
	return
}
func getGallery(c *gin.Context) {
	id := c.Param("id")
	if len(id) != 24 {
		c.JSON(200, "wrong id")
		return
	}
	var re outGalleryModel
	likeslength := bson.M{"$addFields": bson.M{"comments_length": bson.M{"$size": bson.M{"$ifNull": []interface{}{"$comments", []string{}}}}, "likes_length": bson.M{"$size": bson.M{"$ifNull": []interface{}{"$likes", []string{}}}}, "isLiked": bson.M{"$in": []interface{}{c.MustGet("auth").(string), bson.M{"$ifNull": []interface{}{"$likes", []string{}}}}}}}
	line := []bson.M{
		{"$match": bson.M{"_id": bson.ObjectIdHex(id)}},
		lookowner,
		unwind,
		likeslength,
	}
	err := pipiC("gallery", line, &re, false)
	if err != nil {
		c.JSON(200, false)
		return
	}
	c.JSON(200, re)
}
func delGallery(c *gin.Context) {
	id := c.Param("id")
	if len(id) != 24 {
		log.Println(id)
		c.JSON(200, gin.H{"status": false, "msg": "not correct id"})
		return
	}
	var re circleModel
	err := findAndModify("gallery", bson.M{"_id": bson.ObjectIdHex(id), "owner": bson.ObjectIdHex(c.MustGet("auth").(string))}, nil, false, true, false, &re)
	if err != nil {
		c.JSON(200, gin.H{"status": false, "msg": err.Error()})
		return
	}
	// log.Println(re)
	// for _, val := range re.Pics {
	// 	os.Remove(globalConf.ResDir + "/gallery/" + filepath.Base(val))
	// 	os.Remove(globalConf.ResDir + "/gallery/" + filepath.Base(val) + "_thb.jpeg")
	// }
	c.JSON(200, gin.H{"status": true, "msg": "删除图片成功"})
	return
}
