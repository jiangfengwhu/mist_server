package main

import (
	"crypto/md5"
	"fmt"
	"github.com/globalsign/mgo/bson"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
)

var lookowner = bson.M{"$lookup": bson.M{"from": "user", "localField": "owner", "foreignField": "_id", "as": "owner_doc"}}
var unwind = bson.M{"$unwind": "$owner_doc"}

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
func mktorrent(tpid string) (string, error) {
	id := globalConf.ResDir + "/video/" + tpid
	cmd := strings.Fields("ffmpeg -i " + id + " -movflags faststart -acodec copy -vcodec copy -y " + id + ".mp4")
	prc := exec.Command(cmd[0], cmd[1:]...)
	err := prc.Run()
	if err != nil {
		log.Println(err)
		return "", err
	}
	// cmd1 := strings.Fields("create-torrent " + id + ".mp4 --pieceLength=734003 --announce=" + globalConf.Announce + " --urlList=" + globalConf.Host + globalConf.ResRef + "/video/" + tpid + ".mp4 -o " + id + ".torrent")
	cmd1 := strings.Fields("ffmpeg -i " + id + ".mp4 -c copy -f dash -window_size 0 -seg_duration 10 -init_seg_name " + tpid + "init$RepresentationID$.m4s -media_seg_name " + tpid + "$RepresentationID$-$Number%05d$.m4s -use_template 0 -bsf:a aac_adtstoasc " + id + ".mpd")
	prc1 := exec.Command(cmd1[0], cmd1[1:]...)
	err = prc1.Run()
	if err != nil {
		log.Println(err)
		return "", err
	}
	// return globalConf.Host + globalConf.ResRef + "/video/" + tpid + ".mpd", nil
	return globalConf.ResRef + "/video/" + tpid + ".mpd", nil
}
func capCover(id string, sec string) (string, error) {
	id = globalConf.ResDir + "/video/" + id
	cmd := strings.Fields("ffmpeg -ss " + sec + " -i " + id + " -vframes 1 -r 1 -f image2 -y " + id + ".jpg")
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

func splitPath(inp string) string {
	filepaths := strings.Split(inp, "/")
	return filepaths[len(filepaths)-1]
}
