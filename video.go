package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
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
	// c.JSON(200, gin.H{"index": -1, "path": globalConf.ResDir + "/video/" + vid + ".jpg"})
	c.JSON(200, gin.H{"index": -1})
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
	// var collection basicVCModel
	// if err := c.ShouldBind(&collection); err != nil {
	// 	log.Println(err)
	// 	c.JSON(200, gin.H{"status": false, "msg": "信息不完整"})
	// 	return
	// }
	// collection.ID = bson.NewObjectId()
	// collection.Owner = bson.ObjectIdHex(c.MustGet("auth").(string))
	// collection.Date = time.Now().Unix()
	// collection.View = 0
	// err := insertC("video", collection)

	// if err != nil {
	// 	log.Println(err.Error())
	// 	c.JSON(200, gin.H{"status": false, "msg": "数据库错误"})
	// 	return
	// }
	// c.JSON(200, gin.H{"status": true, "cid": collection.ID, "msg": "创建专辑成功"})
	// return
}
func addVideo(c *gin.Context) {
	var video newVideoModel
	if err := c.ShouldBind(&video); err != nil {
		log.Println(err)
		c.JSON(200, gin.H{"status": false, "msg": err.Error()})
		return
	}
	log.Println(video)
	video.Date = time.Now().Unix()
	video.ID = bson.NewObjectId()
	video.Owner = bson.ObjectIdHex(c.MustGet("auth").(string))
	video.View = 0
	tsp := getsubs(video.Hash)
	tp, err := mktorrent(video.Hash)
	if err != nil {
		log.Println(err.Error())
		c.JSON(200, gin.H{"status": false, "msg": err.Error()})
		return
	}
	path, err := capCover(video.Hash, fmt.Sprint(video.CoverPos), false)
	if err != nil {
		log.Println(err)
		c.JSON(200, gin.H{"status": false, "msg": err.Error()})
		return
	}
	video.Path = tp
	video.Cover = path
	video.Subtitle = tsp
	err = insertC("video", video)
	if err != nil {
		log.Println(err.Error())
		c.JSON(200, gin.H{"status": false, "msg": err.Error()})
		return
	}
	c.JSON(200, gin.H{"status": true, "msg": "添加成功", "path": path})
	return
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
func searchVideo(c *gin.Context) {
	var par searchModel
	if err := c.ShouldBind(&par); err != nil {
		c.JSON(200, false)
		return
	}
	var re []faceVideoModel
	commentslength := bson.M{"$addFields": bson.M{"comments_length": bson.M{"$size": bson.M{"$ifNull": []interface{}{"$comments", []string{}}}}, "likes_length": bson.M{"$size": bson.M{"$ifNull": []interface{}{"$likes", []string{}}}}}}
	tagmatch := bson.M{"$match": bson.M{"$or": []bson.M{bson.M{"title": bson.M{"$regex": bson.RegEx{Pattern: par.Key, Options: "i"}}}, bson.M{"desc": bson.M{"$regex": bson.RegEx{Pattern: par.Key, Options: "i"}}}, bson.M{"listdoc.title": bson.M{"$regex": bson.RegEx{Pattern: par.Key, Options: "i"}}}, bson.M{"listdoc.desc": bson.M{"$regex": bson.RegEx{Pattern: par.Key, Options: "i"}}}}}}
	err := latestC("video", []bson.M{looklist, unwindlist, tagmatch, commentslength}, *par.Start*par.Size, par.Size*(*par.Start+1), &re)
	if err != nil {
		c.JSON(200, false)
		return
	}
	// notempty := bson.M{"$match": bson.M{"videos": bson.M{"$exists": true, "$ne": empty}}}
	c.JSON(200, re)
	return
}
func latestVideo(c *gin.Context) {
	var params getVModel
	if err := c.ShouldBind(&params); err != nil {
		c.JSON(200, gin.H{"status": false, "msg": err.Error()})
		return
	}
	// empty := []string{}
	tagsize := 11
	re := make([][]faceVideoModel, 1, tagsize)
	log.Println(*params.Start*params.Size, *params.Start*(params.Size+1))
	commentslength := bson.M{"$addFields": bson.M{"comments_length": bson.M{"$size": bson.M{"$ifNull": []interface{}{"$comments", []string{}}}}, "likes_length": bson.M{"$size": bson.M{"$ifNull": []interface{}{"$likes", []string{}}}}}}
	removeListDump := bson.M{"$group": bson.M{"_id": bson.M{"$ifNull": []string{"$playlist", "$_id"}}, "counterview": bson.M{"$sum": "$view"}, "counterlike": bson.M{"$sum": "$likes_length"}, "countercomments": bson.M{"$sum": "$comments_length"}, "counters": bson.M{"$sum": 1}, "doc": bson.M{"$last": "$$ROOT"}}}
	replaceroot := bson.M{"$replaceRoot": bson.M{"newRoot": bson.M{"$mergeObjects": []interface{}{"$doc", bson.M{"view": "$counterview", "likes_length": "$counterlike", "comments_length": "$countercomments", "nvideos": "$counters"}}}}}
	var err error
	if params.Tag == -1 {
		for i := 0; i < tagsize; i++ {
			tagmatch := bson.M{"$match": bson.M{"tag": i + 1}}
			err = latestC("video", []bson.M{tagmatch, commentslength, removeListDump, replaceroot, looklist, unwindlist}, *params.Start*params.Size, params.Size*(*params.Start+1), &re[i])
			if i != tagsize-1 {
				re = append(re, []faceVideoModel{})
			}
		}
	} else {
		tagmatch := bson.M{"$match": bson.M{"tag": params.Tag}}
		err = latestC("video", []bson.M{tagmatch, commentslength, removeListDump, replaceroot, looklist, unwindlist}, *params.Start*params.Size, params.Size*(*params.Start+1), &re[0])
	}
	if err != nil {
		c.JSON(200, false)
		return
	}
	// notempty := bson.M{"$match": bson.M{"videos": bson.M{"$exists": true, "$ne": empty}}}
	c.JSON(200, re)
	return
}

func getVideo(c *gin.Context) {
	err := bson.IsObjectIdHex(c.Param("id"))
	if !err {
		c.JSON(200, false)
		return
	}
	var re detailVideoModel
	handleVideo(func(col *mgo.Collection) {
		likeslength := bson.M{"$addFields": bson.M{"likes_length": bson.M{"$size": bson.M{"$ifNull": []interface{}{"$likes", []string{}}}}, "isLiked": bson.M{"$in": []interface{}{c.MustGet("auth").(string), bson.M{"$ifNull": []interface{}{"$likes", []string{}}}}}}}
		err := col.Pipe([]bson.M{
			{"$match": bson.M{"_id": bson.ObjectIdHex(c.Param("id"))}},
			lookowner,
			unwind,
			likeslength,
		}).One(&re)
		var recom []basicVideoModel
		tagLine := bson.M{"$match": bson.M{"tag": re.Tag}}
		if len(re.PlayList) != 0 {
			tagLine = bson.M{"$match": bson.M{"tag": re.Tag, "playlist": bson.M{"$ne": re.PlayList}}}
			if err := findC("playlist", bson.M{"_id": re.PlayList}, false, &re.ListDoc); err != nil {
				c.JSON(200, false)
				return
			}
			var plist []basicVideoModel
			listLine := bson.M{"$match": bson.M{"playlist": re.PlayList}}
			err := pipiC("video", []bson.M{listLine, bson.M{"$sort": bson.M{"_id": 1}}}, &plist, true)
			if err != nil {
				log.Println(err.Error())
				c.JSON(200, false)
				return
			}
			re.ListVideos = plist
		}
		err = latestC("video", []bson.M{tagLine}, 0, 12, &recom)
		if err != nil {
			c.JSON(200, false)
			return
		}
		re.Recommend = recom
		c.JSON(200, re)
		col.Update(bson.M{"_id": bson.ObjectIdHex(c.Param("id"))}, bson.M{"$inc": bson.M{"view": 1}})
	})
}
func videoAll(c *gin.Context) {
	id := c.Param("id")
	if len(id) != 24 {
		c.JSON(200, "wrong id")
		return
	}
	var re []userVideoModel
	err := pipiC("video", []bson.M{
		{"$match": bson.M{"owner": bson.ObjectIdHex(id)}},
		looklist,
		unwindlist,
	}, &re, true)
	if err != nil {
		c.JSON(200, false)
		return
	}
	c.JSON(200, re)
	return
}
func myVideo(c *gin.Context) {
	// var params getVModel
	// if err := c.ShouldBind(&params); err != nil {
	// 	c.JSON(200, gin.H{"status": false, "msg": err.Error()})
	// 	return
	// }
	// var re []myVCModel
	// log.Println(*params.Start*params.Size, *params.Start*(params.Size+1))
	// err := latestC("video", []bson.M{{"$match": bson.M{"owner": bson.ObjectIdHex(c.MustGet("auth").(string))}}}, *params.Start*params.Size, params.Size*(*params.Start+1), &re)
	// log.Println(re)
	// if err != nil {
	// 	c.JSON(200, false)
	// 	return
	// }
	// count, err := getCount("video", bson.M{"owner": bson.ObjectIdHex(c.MustGet("auth").(string))})
	// if err != nil {
	// 	c.JSON(200, false)
	// 	return
	// }
	// c.JSON(200, gin.H{"docs": re, "counts": count})
	// return
}

func delvideoc(c *gin.Context) {
	ids := c.PostFormArray("ids")
	log.Println(c.PostFormArray("ids"))
	handleVideo(func(col *mgo.Collection) {
		var dels []delVideoModel
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
			os.Remove(globalConf.ResDir + "/video/" + filepath.Base(vc.Cover))
			vs, err := filepath.Glob(globalConf.ResDir + "/video/" + strings.TrimSuffix(filepath.Base(vc.Path), ".m3u8") + "*")
			if err != nil {
				log.Println(err.Error())
				c.JSON(200, gin.H{"msg": "删除文件错误", "status": false})
				return
			}
			for _, v := range vs {
				os.Remove(v)
			}
			for _, v := range vc.Subs {
				os.Remove(globalConf.ResDir + "/video/" + filepath.Base(v))
			}
		}
		c.JSON(200, gin.H{"msg": "删除专辑成功", "status": true})
		return
	})
}
func updateVC(c *gin.Context) {
	var vc updateVideoModel
	if err := c.ShouldBind(&vc); err != nil {
		log.Println(err)
		c.JSON(200, gin.H{"status": false, "msg": "信息不完整"})
		return
	}
	err := updateC("video", bson.M{"_id": vc.ID, "owner": bson.ObjectIdHex(c.MustGet("auth").(string))}, bson.M{"$set": bson.M{"tag": vc.Tag, "desc": vc.Desc, "title": vc.Title}})
	if err != nil {
		c.JSON(200, gin.H{"status": false, "msg": err.Error()})
		return
	}
	c.JSON(200, gin.H{"status": true, "msg": "更新专辑成功"})
	return
}

func updateVideo(c *gin.Context) {
	// var video addVideoModel
	// if err := c.ShouldBind(&video); err != nil {
	// 	log.Println(err)
	// 	c.JSON(200, gin.H{"status": false, "msg": "信息不完整"})
	// 	return
	// }
	// log.Println(video)
	// err := updateC("video", bson.M{"_id": video.Cid, "owner": bson.ObjectIdHex(c.MustGet("auth").(string)), "videos": bson.M{"$elemMatch": bson.M{"_id": video.Vid}}}, bson.M{"$set": bson.M{"videos.$.title": video.Title, "videos.$.desc": video.Desc}})
	// if err != nil {
	// 	c.JSON(200, gin.H{"status": false, "msg": "数据库错误"})
	// 	return
	// }
	// c.JSON(200, gin.H{"status": true, "msg": "更新视频信息成功"})
	// return
}
func checkOWN(c *gin.Context) {
	id := c.Param("id")
	num, err := getCount("video", bson.M{"_id": bson.ObjectIdHex(id), "owner": bson.ObjectIdHex(c.MustGet("auth").(string))})
	if err != nil {
		log.Println(err.Error())
		c.JSON(200, false)
		return
	}
	if num == 0 {
		c.JSON(200, false)
		return
	}
	c.JSON(200, true)
	return
}

type delSubModel struct {
	ID  bson.ObjectId `bson:"_id" json:"id" binding:"required"`
	Sub string        `json:"sub" binding:"required"`
}

func delSub(c *gin.Context) {
	var del delSubModel
	if err := c.ShouldBind(&del); err != nil {
		log.Println(err.Error())
		return
	}
	err := updateC("video", bson.M{"_id": del.ID, "owner": bson.ObjectIdHex(c.MustGet("auth").(string))}, bson.M{"$pull": bson.M{"subtitle": del.Sub}})
	if err != nil {
		c.JSON(200, gin.H{"status": false, "msg": err.Error()})
		return
	}
	if err := os.Remove(globalConf.ResDir + "/video/" + filepath.Base(del.Sub)); err != nil {
		log.Println(err.Error())
		c.JSON(200, gin.H{"msg": "删除文件错误", "status": false})
		return
	}
	c.JSON(200, gin.H{"status": true, "msg": "删除字幕成功"})
	return
}
func addSub(c *gin.Context) {
	var sub uploadFileModel
	if err := c.ShouldBind(&sub); err != nil {
		log.Println(err.Error())
		return
	}
	if err := c.SaveUploadedFile(sub.Blob, globalConf.ResDir+"/video/"+sub.Blob.Filename); err != nil {
		c.JSON(200, gin.H{"status": false, "msg": err.Error()})
		return
	}
	sp := saveSubs(globalConf.ResDir + "/video/" + sub.Blob.Filename)
	err := updateC("video", bson.M{"_id": bson.ObjectIdHex(sub.ID), "owner": bson.ObjectIdHex(c.MustGet("auth").(string))}, bson.M{"$addToSet": bson.M{"subtitle": sp}})
	if err != nil {
		c.JSON(200, gin.H{"status": false, "msg": err.Error()})
		return
	}
	c.JSON(200, gin.H{"status": true, "sub": sp, "msg": "添加字幕成功"})
	return
}
func delvideo(c *gin.Context) {
	// var params delVideoQ
	// if err := c.ShouldBindQuery(&params); err != nil {
	// 	log.Println(err.Error(), params)
	// 	c.JSON(200, gin.H{"status": false, "msg": err.Error()})
	// 	return
	// }
	// err := updateC("video", bson.M{"_id": bson.ObjectIdHex(params.Cid), "owner": bson.ObjectIdHex(c.MustGet("auth").(string)), "videos": bson.M{"$elemMatch": bson.M{"_id": params.Vid, "cover": globalConf.ResRef + "/video/" + params.Cover}}}, bson.M{"$pull": bson.M{"videos": bson.M{"_id": params.Vid}}})
	// if err != nil {
	// 	c.JSON(200, gin.H{"status": false, "msg": err.Error()})
	// 	return
	// }
	// os.Remove(globalConf.ResDir + "/video/" + params.Cover)
	// vs, err := filepath.Glob(globalConf.ResDir + "/video/" + params.Vid + "*")
	// if err != nil {
	// 	log.Println(err.Error())
	// 	c.JSON(200, gin.H{"msg": "删除文件错误", "status": false})
	// 	return
	// }
	// for _, v := range vs {
	// 	os.Remove(v)
	// }
	// c.JSON(200, gin.H{"status": true, "msg": "删除视频成功"})
	// return
}
