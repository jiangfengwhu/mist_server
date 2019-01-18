package main

import (
	"github.com/gin-gonic/gin"
	"github.com/globalsign/mgo/bson"
	"time"
)

func getComments(c *gin.Context) {
	var par getCommentsModel
	if err := c.ShouldBindQuery(&par); err != nil {
		c.JSON(200, false)
		return
	}
	dbtype := ""
	switch par.Type {
	case "1":
		dbtype = "video"
	case "2":
		dbtype = "community"
	case "reply":
		dbtype = "comment"
	default:
		c.JSON(200, false)
		return
	}
	pipeline := []bson.M{
		{"$match": bson.M{"_id": bson.ObjectIdHex(par.ID)}},
		lookcomments,
		lookreplycount,
		lookcomowner,
	}
	var re outCommentModel
	err := pipiC(dbtype, pipeline, &re, false)
	if err != nil {
		c.JSON(200, false)
		return
	}
	c.JSON(200, re)
	return
}
func addComment(c *gin.Context) {
	var par postCommentModel
	if err := c.ShouldBind(&par); err != nil {
		c.JSON(200, gin.H{"msg": err.Error(), "status": false})
		return
	}
	var doc commentModel
	doc.ID = bson.NewObjectId()
	doc.Content = par.Content
	doc.At = par.At
	doc.Date = time.Now().Unix()
	doc.Owner = bson.ObjectIdHex(c.MustGet("auth").(string))
	err := insertC("comment", doc)
	if err != nil {
		c.JSON(200, gin.H{"msg": err.Error(), "status": false})
		return
	}
	dbtype := ""
	switch par.Type {
	case "1":
		dbtype = "video"
	case "2":
		dbtype = "community"
	case "reply":
		dbtype = "comment"
	default:
		c.JSON(200, gin.H{"status": false, "msg": "类型错误"})
		return
	}
	err = updateC(dbtype, bson.M{"_id": bson.ObjectIdHex(par.ID)}, bson.M{"$push": bson.M{"comments": doc.ID}})
	if err != nil {
		c.JSON(200, gin.H{"msg": err.Error(), "status": false})
		return
	}
	c.JSON(200, gin.H{"status": true, "msg": "发表评论成功", "doc": doc})
	return
}
