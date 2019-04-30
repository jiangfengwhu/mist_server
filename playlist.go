package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/globalsign/mgo/bson"
)

type basicListModel struct {
	outListModel `bson:",inline"`
	Owner        bson.ObjectId   `bson:"owner" json:"-"`
	OwnerDoc     *basicUserModel `bson:"owner_doc,omitempty" json:"owner,omitempty"`
}
type outListModel struct {
	Title string        `bson:"title" json:"title" binding:"required"`
	Desc  string        `bson:"desc,omitempty" json:"desc,omitempty"`
	ID    bson.ObjectId `bson:"_id" json:"id"`
}
type updateListModel struct {
	ID     bson.ObjectId   `bson:"_id" json:"id" binding:"required"`
	Title string        `bson:"title" json:"title" binding:"required"`
	Desc  string        `bson:"desc,omitempty" json:"desc,omitempty"`
}
type addToListModel struct {
	ID     bson.ObjectId   `bson:"_id" json:"id" binding:"required"`
	Videos []bson.ObjectId `bson:"videos" json:"videos" binding:"required"`
}
type removeFromListModel struct {
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
	var update addToListModel
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
	var update removeFromListModel
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
func removeList(c *gin.Context) {
	id := c.Param("id")
	if len(id) != 24 {
		log.Println(id)
		c.JSON(200, gin.H{"status": false, "msg": "not correct id"})
		return
	}
	if err := delC("playlist", bson.M{"_id": bson.ObjectIdHex(id), "owner": bson.ObjectIdHex(c.MustGet("auth").(string))});err!=nil {
		c.JSON(200, gin.H{"status":false, "msg":"删除失败"})
		return
	}
	info, err := updateAllC("video", bson.M{"playlist": bson.ObjectIdHex(id)}, bson.M{"$unset": bson.M{"playlist": ""}})
	log.Println(info)
	if err != nil {
		c.JSON(200, gin.H{"status": false, "msg": err.Error()})
		return
	}
	c.JSON(200, gin.H{"status": true, "msg": "删除播放列表成功"})
	return
}
func updateList(c *gin.Context) {
	var list updateListModel
	if err := c.ShouldBind(&list); err != nil {
		log.Println(err)
		c.JSON(200, gin.H{"status": false, "msg": "信息不完整"})
		return
	}
	err := updateC("playlist", bson.M{"_id": list.ID, "owner": bson.ObjectIdHex(c.MustGet("auth").(string))}, bson.M{"$set": bson.M{"desc": list.Desc, "title": list.Title}})
	if err != nil {
		c.JSON(200, gin.H{"status": false, "msg": err.Error()})
		return
	}
	c.JSON(200, gin.H{"status": true, "msg": "更新播放列表信息成功", "list": gin.H{"id": list.ID, "desc": list.Desc, "title": list.Title}})
	return
}