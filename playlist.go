package main

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/globalsign/mgo/bson"
)

type basicListModel struct {
	outListModel `bson:",inline"`
	ID           bson.ObjectId  `bson:"_id" json:"id"`
	Cover        string         `bson:"cover" json:"cover"`
	Owner        bson.ObjectId  `bson:"owner" json:"-"`
	OwnerDoc     *basicUserModel `bson:"owner_doc,omitempty" json:"owner,omitempty"`
}
type outListModel struct {
	Title string `bson:"title" json:"title,omitempty" binding:"required"`
	Date  int64  `bson:"date" json:"date,omitempty"`
}
type updateListModel struct {
	ID     bson.ObjectId   `bson:"_id" json:"id" binding:"required"`
	Videos []bson.ObjectId `bson:"videos" json:"videos" binding:"required"`
}
type removeListModel struct {
	Videos []bson.ObjectId `bson:"videos" json:"videos" binding:"required"`
}
func newList(c *gin.Context) {
	var list basicListModel
	if err := c.ShouldBind(&list); err != nil {
		log.Println(err)
		c.JSON(200, gin.H{"status": false, "msg": err.Error()})
		return
	}
	log.Println(list)
	list.Date = time.Now().Unix()
	list.ID = bson.NewObjectId()
	list.Owner = bson.ObjectIdHex(c.MustGet("auth").(string))
	if err := insertC("playlist", list); err != nil {
		c.JSON(200, gin.H{"status": false, "msg": err.Error()})
		return
	}
	c.JSON(200, gin.H{"status": true, "msg": "添加成功", "list": list})
	return
}

func listAll(c *gin.Context) {
	id := c.Param("id")
	if len(id) != 24 {
		c.JSON(200, "wrong id")
		return
	}
	var re []basicListModel
	err := findC("playlist", bson.M{"owner": bson.ObjectIdHex(id)}, true, &re)
	if err != nil {
		c.JSON(200, false)
		return
	}
	c.JSON(200, re)
	return
}

func addtoList(c *gin.Context) {
	var update updateListModel
	if err := c.ShouldBind(&update); err != nil {
		log.Println(err.Error())
		c.JSON(200, gin.H{"status": false, "msg": err.Error()})
		return
	}
	info, err := updateAllC("video", bson.M{"_id": bson.M{"$in": update.Videos}, "owner": bson.ObjectIdHex(c.MustGet("auth").(string))}, bson.M{"$set": bson.M{"playlist": update.ID}})
	log.Println(info)
	if err != nil {
		c.JSON(200, gin.H{"status": false, "msg": err.Error()})
		return
	}
	c.JSON(200, gin.H{"status": true, "msg": "添加成功"})
	return
}

func removeFromList(c *gin.Context) {
	var update removeListModel
	if err := c.ShouldBind(&update); err != nil {
		log.Println(err.Error())
		c.JSON(200, gin.H{"status": false, "msg": err.Error()})
		return
	}
	info, err := updateAllC("video", bson.M{"_id": bson.M{"$in": update.Videos}, "owner": bson.ObjectIdHex(c.MustGet("auth").(string))}, bson.M{"$unset": bson.M{"playlist": ""}})
	log.Println(info)
	if err != nil {
		c.JSON(200, gin.H{"status": false, "msg": err.Error()})
		return
	}
	c.JSON(200, gin.H{"status": true, "msg": "从列表移除成功"})
	return
}
