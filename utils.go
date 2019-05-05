package main

import (
	"crypto/md5"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/globalsign/mgo/bson"
)

var lookowner = bson.M{"$lookup": bson.M{"from": "user", "localField": "owner", "foreignField": "_id", "as": "owner_doc"}}
var unwind = bson.M{"$unwind": "$owner_doc"}
var looklist = bson.M{"$lookup": bson.M{"from": "playlist", "localField": "playlist", "foreignField": "_id", "as": "listdoc"}}
var unwindlist = bson.M{"$unwind": bson.M{"path": "$listdoc", "preserveNullAndEmptyArrays": true}}
var lookcomments = bson.M{"$lookup": bson.M{"from": "comment", "localField": "comments", "foreignField": "_id", "as": "comments_doc"}}
var lookcomowner = bson.M{"$lookup": bson.M{"from": "user", "localField": "comments_doc.owner", "foreignField": "_id", "as": "owners_doc"}}

func getMd5(filename string) (string, error) {
	f, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer f.Close()

	hs := md5.New()
	if _, err := io.Copy(hs, f); err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", hs.Sum(nil)), nil
}
func getsubs(tpid string) []string {
	id := globalConf.ResDir + "/video/" + tpid
	re := make([]string, 0, 5)
	var counter uint8
	var loop func()
	loop = func() {
		cmd := strings.Fields("ffmpeg -i " + id + " -y -map 0:s:" + fmt.Sprint(counter) + " -c:s webvtt " + id + ".vtt")
		err := exec.Command(cmd[0], cmd[1:]...).Run()
		if err != nil {
			return
		}
		hash, _ := getMd5(id + ".vtt")
		os.Rename(id+".vtt", globalConf.ResDir+"/video/"+hash+".vtt")
		re = append(re, globalConf.ResRef+"/video/"+hash+".vtt")
		counter++
		loop()
	}
	loop()
	return re
}
func saveSubs(loc string) string {
	cmd := strings.Fields("ffmpeg -i " + loc + " " + loc + ".vtt")
	err := exec.Command(cmd[0], cmd[1:]...).Run()
	if err != nil {
		log.Println(err.Error())
		return ""
	}
	hash, _ := getMd5(loc + ".vtt")
	if err := os.Rename(loc+".vtt", globalConf.ResDir+"/video/"+hash+".vtt"); err != nil {
		return ""
	}
	os.Remove(loc)
	return globalConf.ResRef + "/video/" + hash + ".vtt"
}
func mktorrent(tpid string) (string, error) {
	id := globalConf.ResDir + "/video/" + tpid
	// cmd1 := strings.Fields("create-torrent " + id + ".mp4 --pieceLength=734003 --announce=" + globalConf.Announce + " --urlList=" + globalConf.Host + globalConf.ResRef + "/video/" + tpid + ".mp4 -o " + id + ".torrent")
	cmd1 := strings.Fields("ffmpeg -i " + id + " -codec copy -sn -g 48 -keyint_min 48 -start_number 0 -hls_time 10 -hls_playlist_type vod -hls_allow_cache 1 -f hls " + id + ".m3u8")

	prc1 := exec.Command(cmd1[0], cmd1[1:]...)
	err := prc1.Run()
	if err != nil {
		log.Println(err)
		return "", err
	}
	// return globalConf.Host + globalConf.ResRef + "/video/" + tpid + ".mpd", nil
	return globalConf.ResRef + "/video/" + tpid + ".m3u8", nil
}
func capCover(id string, sec string, origin bool) (string, error) {
	id = globalConf.ResDir + "/video/" + id
	cmdstr := "ffmpeg -ss " + sec + " -i " + id + " -vframes 1 -r 1 -vf scale=320:180:force_original_aspect_ratio=increase,crop=320:180 -f image2 -y " + id + ".jpg"
	if origin {
		cmdstr = "ffmpeg -ss " + sec + " -i " + id + " -vframes 1 -r 1 -f image2 -y " + id + ".jpg"
	}
	cmd := strings.Fields(cmdstr)
	prc := exec.Command(cmd[0], cmd[1:]...)
	err := prc.Run()
	if err != nil {
		log.Println(err)
		return "", err
	}
	hash, err := getMd5(id + ".jpg")
	if err != nil {
		log.Println(err)
		return "", err
	}
	if err := os.Rename(id+".jpg", globalConf.ResDir+"/video/"+hash+".jpg"); err != nil {
		return "", err
	}
	if err := os.Remove(id); err != nil {
		log.Println(err)
		return "", err
	}
	return globalConf.ResRef + "/video/" + hash + ".jpg", nil
}
func strtoobj(inp []string) []bson.ObjectId {
	re := make([]bson.ObjectId, 0)
	for _, val := range inp {
		re = append(re, bson.ObjectIdHex(val))
	}
	return re
}
