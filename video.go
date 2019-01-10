package main

import (
	"github.com/gin-gonic/gin"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"io"
	"log"
	"os"
	"time"
)

func checkBP(c *gin.Context) {
	vid := c.Param("id")
	if vid == "" {
		c.JSON(200, gin.H{"msg": "参数不完整"})
		return
	}
	if _, err := os.Stat(globalConf.ResDir + "/video/" + vid + ".jpg"); err != nil {
		if os.IsNotExist(err) {
			info, err := os.Stat(globalConf.ResDir + "/video/" + vid)
			if err != nil {
				if os.IsNotExist(err) {
					c.JSON(200, gin.H{"index": 0})
					return
				}
				c.JSON(200, gin.H{"msg": err})
				return
			}
			c.JSON(200, gin.H{"index": info.Size()})
			return
		}
		c.JSON(200, gin.H{"msg": err})
		return
	}
	c.JSON(200, gin.H{"index": -1, "path": globalConf.ResDir + "/video/" + vid + ".jpg"})
	return
}

func uploadFile(c *gin.Context) {
	var upf uploadFileModel
	if err := c.ShouldBind(&upf); err != nil {
		log.Println(err.Error())
		c.JSON(200, gin.H{"msg": "信息不完整", "status": false})
		return
	}
	file, err := c.FormFile("blob")
	if err != nil {
		c.JSON(200, gin.H{"msg": err.Error(), "status": false})
		return
	}
	f, err := os.OpenFile(globalConf.ResDir+"/video/"+upf.ID, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		c.JSON(200, gin.H{"status": false, "msg": err.Error()})
		return
	}
	defer f.Close()
	src, err := file.Open()
	if err != nil {
		c.JSON(200, gin.H{"status": false, "msg": err.Error()})
		return
	}
	defer src.Close()
	_, err = io.Copy(f, src)
	if err != nil {
		c.JSON(200, gin.H{"status": false, "msg": err.Error()})
		return
	}
	c.JSON(200, gin.H{"status": true})
	return
}

func createCollection(c *gin.Context) {
	var collection basicVCModel
	if err := c.ShouldBind(&collection); err != nil {
		log.Println(err)
		c.JSON(200, gin.H{"status": false, "msg": "信息不完整"})
		return
	}
	collection.ID = bson.NewObjectId()
	collection.Owner = bson.ObjectIdHex(c.MustGet("auth").(string))
	collection.Date = time.Now().Unix()
	collection.View = 0
	err := insertC("video", collection)

	if err != nil {
		log.Println(err.Error())
		c.JSON(200, gin.H{"status": false, "msg": "数据库错误"})
		return
	}
	c.JSON(200, gin.H{"status": true, "cid": collection.ID, "msg": "创建专辑成功"})
	return
}

func addVideo(c *gin.Context) {
	var vid addVideoModel
	if err := c.ShouldBind(&vid); err != nil {
		log.Println(err)
		c.JSON(200, gin.H{"status": false, "msg": err.Error()})
		return
	}
	log.Println(vid)
	handleVideo(func(col *mgo.Collection) {
		var collec basicVCModel
		err := col.Find(bson.M{"_id": vid.Cid, "owner": bson.ObjectIdHex(c.MustGet("auth").(string))}).One(&collec)
		if err != nil {
			c.JSON(200, gin.H{"status": false, "msg": err.Error()})
			return
		}
		tp, err := mktorrent(vid.Vid)
		if err != nil {
			log.Println(err)
			c.JSON(200, gin.H{"status": false, "msg": err})
			return
		}
		path, err := capCover(vid.Vid, "5")
		if err != nil {
			log.Println(err)
			c.JSON(200, gin.H{"status": false, "msg": err})
			return
		}
		if collec.Cover != "" {
			err = updateC("video", bson.M{"_id": vid.Cid}, bson.M{"$push": bson.M{"videos": bson.M{"date": time.Now().Unix(), "title": vid.Title, "_id": vid.Vid, "desc": vid.Desc, "cover": path, "path": tp}}})
		} else {
			err = updateC("video", bson.M{"_id": vid.Cid}, bson.M{"$push": bson.M{"videos": bson.M{"date": time.Now().Unix(), "title": vid.Title, "_id": vid.Vid, "desc": vid.Desc, "cover": path, "path": tp}}, "$set": bson.M{"cover": path}})
		}
		if err != nil {
			log.Println(err)
			c.JSON(200, gin.H{"status": false, "msg": err.Error()})
			return
		}

		c.JSON(200, gin.H{"status": true, "msg": "添加成功", "path": path})
		return
	})
}

func changeCover(c *gin.Context) {
	var cc changeCoverModel
	if err := c.ShouldBind(&cc); err != nil {
		c.JSON(200, gin.H{"status": false, "msg": err.Error()})
		return
	}
	if cc.Path == "" {
		file, err := c.FormFile("file")
		if err != nil {
			c.JSON(200, gin.H{"status": false, "msg": err.Error()})
			return
		}
		if err := c.SaveUploadedFile(file, globalConf.ResDir+"/cover/"+file.Filename); err != nil {
			c.JSON(200, gin.H{"status": false, "msg": err.Error()})
			return
		}
		cc.Path = globalConf.ResRef + "/cover/" + file.Filename
	}
	err := updateC("video", bson.M{"_id": cc.Cid, "owner": bson.ObjectIdHex(c.MustGet("auth").(string))}, bson.M{"$set": bson.M{"cover": cc.Path}})
	if err != nil {
		log.Println(cc.Cid, err.Error())
		c.JSON(200, gin.H{"status": false, "msg": err.Error()})
		return
	}
	c.JSON(200, gin.H{"status": true, "msg": "设置封面成功"})
	return
}

func latestVideo(c *gin.Context) {
	var params getVModel
	if err := c.ShouldBind(&params); err != nil {
		c.JSON(200, gin.H{"status": false, "msg": err.Error()})
		return
	}
	empty := []string{}
	var re []basicVCModel
	log.Println(*params.Start*params.Size, *params.Start*(params.Size+1))
	err := latestC("video", []bson.M{{"$match": bson.M{"videos": bson.M{"$exists": true, "$ne": empty}}}}, *params.Start*params.Size, params.Size*(*params.Start+1), &re)
	if err != nil {
		c.JSON(200, false)
		return
	}
	c.JSON(200, re)
	return
}

func getVideo(c *gin.Context) {
	err := bson.IsObjectIdHex(c.Param("id"))
	if !err {
		c.JSON(200, false)
		return
	}
	var re outVideoModel
	handleVideo(func(col *mgo.Collection) {
		err := col.Pipe([]bson.M{
			{"$match": bson.M{"_id": bson.ObjectIdHex(c.Param("id"))}},
			lookowner,
			unwind,
		}).One(&re)
		if err != nil {
			c.JSON(200, false)
			return
		}
		c.JSON(200, re)
		col.Update(bson.M{"_id": bson.ObjectIdHex(c.Param("id"))}, bson.M{"$inc": bson.M{"view": 1}})
	})
}
func myVideoAll(c *gin.Context) {
	var re []myVCModel
	err := findC("video", bson.M{"owner": bson.ObjectIdHex(c.MustGet("auth").(string))}, true, &re)
	if err != nil {
		c.JSON(200, false)
		return
	}
	c.JSON(200, re)
	return
}
func myVideo(c *gin.Context) {
	var params getVModel
	if err := c.ShouldBind(&params); err != nil {
		c.JSON(200, gin.H{"status": false, "msg": err.Error()})
		return
	}
	var re []myVCModel
	log.Println(*params.Start*params.Size, *params.Start*(params.Size+1))
	err := latestC("video", []bson.M{{"$match": bson.M{"owner": bson.ObjectIdHex(c.MustGet("auth").(string))}}}, *params.Start*params.Size, params.Size*(*params.Start+1), &re)
	log.Println(re)
	if err != nil {
		c.JSON(200, false)
		return
	}
	count, err := getCount("video", bson.M{"owner": bson.ObjectIdHex(c.MustGet("auth").(string))})
	if err != nil {
		c.JSON(200, false)
		return
	}
	c.JSON(200, gin.H{"docs": re, "counts": count})
	return
}

func delvideoc(c *gin.Context) {
	ids := c.PostFormArray("ids")
	log.Println(c.PostFormArray("ids"))
	handleVideo(func(col *mgo.Collection) {
		var dels []outVideoModel
		line := bson.M{"_id": bson.M{"$in": strtoobj(ids)}, "owner": bson.ObjectIdHex(c.MustGet("auth").(string))}
		err := col.Find(line).All(&dels)
		if err != nil {
			c.JSON(200, gin.H{"msg": err.Error(), "status": false})
			return
		}
		_, err = col.RemoveAll(line)
		if err != nil {
			c.JSON(200, gin.H{"msg": err.Error(), "status": false})
			return
		}
		for _, vc := range dels {
			for _, video := range vc.Videos {
				os.Remove(globalConf.ResDir + "/video/" + splitPath(video.Cover))
				os.Remove(globalConf.ResDir + "/video/" + splitPath(video.Path))
				os.Remove(globalConf.ResDir + "/video/" + video.ID + ".mp4")
			}
		}
		c.JSON(200, gin.H{"msg": "删除专辑成功", "status": true})
		return
	})
}
func updateVC(c *gin.Context) {
	var vc myVCModel
	if err := c.ShouldBind(&vc); err != nil {
		log.Println(err)
		c.JSON(200, gin.H{"status": false, "msg": "信息不完整"})
		return
	}
	err := updateC("video", bson.M{"_id": vc.ID, "owner": bson.ObjectIdHex(c.MustGet("auth").(string))}, bson.M{"$set": bson.M{"tags": vc.Tags, "desc": vc.Desc, "price": vc.Price, "title": vc.Title}})
	if err != nil {
		c.JSON(200, gin.H{"status": false, "msg": "数据库错误"})
		return
	}
	c.JSON(200, gin.H{"status": true, "msg": "更新专辑成功"})
	return
}

func updateVideo(c *gin.Context) {
	var video addVideoModel
	if err := c.ShouldBind(&video); err != nil {
		log.Println(err)
		c.JSON(200, gin.H{"status": false, "msg": "信息不完整"})
		return
	}
	log.Println(video)
	err := updateC("video", bson.M{"_id": video.Cid, "owner": bson.ObjectIdHex(c.MustGet("auth").(string)), "videos": bson.M{"$elemMatch": bson.M{"_id": video.Vid}}}, bson.M{"$set": bson.M{"videos.$.title": video.Title, "videos.$.desc": video.Desc}})
	if err != nil {
		c.JSON(200, gin.H{"status": false, "msg": "数据库错误"})
		return
	}
	c.JSON(200, gin.H{"status": true, "msg": "更新视频信息成功"})
	return
}

func delvideo(c *gin.Context) {
	var params delVideoQ
	if err := c.ShouldBindQuery(&params); err != nil {
		log.Println(err.Error(), params)
		c.JSON(200, gin.H{"status": false, "msg": err.Error()})
		return
	}
	err := updateC("video", bson.M{"_id": bson.ObjectIdHex(params.Cid), "owner": bson.ObjectIdHex(c.MustGet("auth").(string)), "videos": bson.M{"$elemMatch": bson.M{"_id": params.Vid, "cover": globalConf.ResRef + "/video/" + params.Cover}}}, bson.M{"$pull": bson.M{"videos": bson.M{"_id": params.Vid}}})
	if err != nil {
		c.JSON(200, gin.H{"status": false, "msg": err.Error()})
		return
	}
	os.Remove(globalConf.ResDir + "/video/" + params.Vid + ".mp4")
	os.Remove(globalConf.ResDir + "/video/" + params.Vid + ".torrent")
	os.Remove(globalConf.ResDir + "/video/" + params.Cover)
	c.JSON(200, gin.H{"status": true, "msg": "删除视频成功"})
	return
}
