package main

import (
	"github.com/gin-gonic/gin"
	"github.com/globalsign/mgo/bson"
	"log"
	"time"
)

func uploadImage(c *gin.Context) {
	uploadDir := globalConf.ResDir + "/pics/"
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
	c.JSON(200, gin.H{"status": true, "path": globalConf.ResRef + "/pics/"})
	return
}
func addCircle(c *gin.Context) {
	var collection outCircleModel
	c.ShouldBind(&collection)
	if len(collection.Pics) == 0 && len(collection.Content) == 0 {
		c.JSON(200, gin.H{"status": false, "msg": "信息不完整"})
		return
	}
	if len(collection.Pics) > 4 {
		c.JSON(200, gin.H{"status": false, "msg": "图片数过多"})
		return
	}
	collection.ID = bson.NewObjectId()
	collection.Date = time.Now().Unix()
	collection.Owner = bson.ObjectIdHex(c.MustGet("auth").(string))
	err := insertC("community", collection)
	if err != nil {
		log.Println(err.Error())
		c.JSON(200, gin.H{"status": false, "msg": "数据库错误"})
		return
	}
	c.JSON(200, gin.H{"status": true, "cid": collection.ID, "msg": "发布成功"})
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
	commentslength := bson.M{"$addFields": bson.M{"comments_length": bson.M{"$size": bson.M{"$ifNull": []interface{}{"$comments", []string{}}}}, "likes_length": bson.M{"$size": bson.M{"$ifNull": []interface{}{"$likes", []string{}}}}, "isLiked": bson.M{"$in": []interface{}{c.MustGet("auth").(string), bson.M{"$ifNull": []interface{}{"$likes", []string{}}}}}}}
	err := latestC("community", []bson.M{commentslength}, *params.Start*params.Size, params.Size*(*params.Start+1), &re)
	if err != nil {
		c.JSON(200, false)
		return
	}
	c.JSON(200, re)
	return
}
func commAll(c *gin.Context) {
	id := c.Param("id")
	if len(id) != 24 {
		c.JSON(200, "wrong id")
		return
	}
	var re []ownCircleModel
	likeslength := bson.M{"$addFields": bson.M{"comments_length": bson.M{"$size": bson.M{"$ifNull": []interface{}{"$comments", []string{}}}}, "likes_length": bson.M{"$size": bson.M{"$ifNull": []interface{}{"$likes", []string{}}}}, "isLiked": bson.M{"$in": []interface{}{c.MustGet("auth").(string), bson.M{"$ifNull": []interface{}{"$likes", []string{}}}}}}}
	line := []bson.M{
		{"$match": bson.M{"owner": bson.ObjectIdHex(id)}},
		bson.M{"$sort": bson.M{"_id": -1}},
		likeslength,
	}
	err := pipiC("community", line, &re, true)
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
	// log.Println(re)
	// for _, val := range re.Pics {
	// 	os.Remove(globalConf.ResDir + "/community/" + filepath.Base(val))
	// 	os.Remove(globalConf.ResDir + "/community/" + filepath.Base(val) + "_thb.jpeg")
	// }
	c.JSON(200, gin.H{"status": true, "msg": "删除状态成功"})
	return
}
