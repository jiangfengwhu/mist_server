package main

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"golang.org/x/crypto/bcrypt"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

// Auth is auth middleware
func Auth() func(c *gin.Context) {
	return func(c *gin.Context) {
		re, err := getSession(c)
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		if val, ok := re["uid"]; !(ok && val != nil) {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		c.Set("auth", re["uid"])
	}
}

// CheckGuest is get info m
func CheckGuest() func(c *gin.Context) {
	return func(c *gin.Context) {
		re, err := getSession(c)
		if err != nil {
			c.AbortWithStatus(http.StatusBadGateway)
			return
		}
		if val, ok := re["uid"]; !(ok && val != nil) {
			c.Set("auth", "guest")
			return
		}
		c.Set("auth", re["uid"])
	}
}
func verifyCap(tk string) (*reCapResponseModel, error) {
	resp, err := http.PostForm("https://www.google.com/recaptcha/api/siteverify", url.Values{
		"secret":   {globalConf.RecapSecure},
		"response": {tk},
	})
	if err != nil {
		defer resp.Body.Close()
		return nil, err
	}
	reCapResp := reCapResponseModel{}
	err = json.NewDecoder(resp.Body).Decode(&reCapResp)
	if err != nil {
		return nil, err
	}
	log.Println(reCapResp)
	return &reCapResp, nil
}

func regist(c *gin.Context) {
	var user registModel
	if err := c.ShouldBind(&user); err != nil {
		c.JSON(200, gin.H{"msg": "信息不完整", "status": false})
		return
	}
	reCapResp, err := verifyCap(user.Token)
	if err != nil {
		c.JSON(200, gin.H{"status": false, "msg": "连接reCAPTCHA出错，请稍后再试"})
		return
	}
	if reCapResp.Success == true && reCapResp.Action == "mist_register" {
		handleUser(func(col *mgo.Collection) {
			re, _ := col.Find(bson.M{"email": user.Email}).Count()
			if re > 0 {
				c.JSON(200, gin.H{"status": false, "msg": "邮箱已注册，请登录"})
				return
			}
			activeCode := fmt.Sprintf("%d", rand.Int())
			expireDate := time.Now().Unix() + 1800
			err = sendActiveMail(user.Email, `<p>点击激活：<a>`+globalConf.Host+"/api/activeAccount?email="+user.Email+"&activeCode="+activeCode+"&expireDate="+strconv.FormatInt(expireDate, 10)+`</a></p>`)
			if err != nil {
				c.JSON(200, gin.H{"msg": "邮件服务异常", "status": false})
				return
			}
			encryptPw, err := bcrypt.GenerateFromPassword([]byte(user.Passwd), 10)
			if err != nil {
				c.JSON(200, gin.H{"msg": "密码加密失败", "status": false})
				return
			}
			err = col.Insert(bson.M{"email": user.Email, "passwd": string(encryptPw), "nickName": user.Name, "expireDate": expireDate, "activeCode": activeCode})
			if err != nil {
				c.JSON(200, gin.H{"msg": "数据库错误", "status": false})
				return
			}
			c.JSON(200, gin.H{"status": true, "msg": "注册成功，请激活"})
		})
		return
	}
	if reCapResp.Err[0] == "timeout-or-duplicate" {
		c.JSON(200, gin.H{"msg": "验证码失效", "status": false})
		return
	}
	c.JSON(200, gin.H{"msg": "你可能是个机器人", "status": false})
}

func activeAccount(c *gin.Context) {
	var user activeUserModel
	if err := c.ShouldBindQuery(&user); err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"msg": err.Error(), "status": false})
		return
	}
	handleUser(func(col *mgo.Collection) {
		var re activeUserModel
		err := col.Find(bson.M{"email": user.Email}).One(&re)
		if err != nil {
			c.JSON(200, gin.H{"msg": err.Error(), "status": false})
			return
		}
		if re.Expire-time.Now().Unix() < 0 {
			col.Remove(bson.M{"email": user.Email, "expireDate": re.Expire})
			c.JSON(200, gin.H{"msg": "链接失效", "status": false})
			return
		} else if re.ActiveCode == user.ActiveCode {
			err := col.Update(bson.M{"email": user.Email}, bson.M{"$unset": bson.M{"expireDate": "", "activeCode": ""}, "$set": bson.M{"jd": time.Now().Unix()}})
			if err != nil {
				c.JSON(200, gin.H{"msg": err.Error(), "status": false})
				return
			}
		}
		c.JSON(200, gin.H{"msg": "激活成功"})
		return
	})
}

func login(c *gin.Context) {
	var user loginModel
	if err := c.ShouldBind(&user); err != nil {
		c.JSON(200, gin.H{"msg": "请输入必要的信息", "status": false})
		return
	}
	reCapResp, err := verifyCap(user.Token)
	if err != nil {
		c.JSON(200, gin.H{"status": false, "msg": "连接reCAPTCHA出错，请稍后再试"})
		return
	}
	if !(reCapResp.Success == true && reCapResp.Action == "mist_login") {
		c.JSON(200, gin.H{"status": false, "msg": reCapResp.Err})
		return
	}
	handleUser(func(col *mgo.Collection) {
		var re loginDBModel
		err := col.Find(bson.M{"email": user.Email}).One(&re)
		if err != nil {
			c.JSON(200, gin.H{"msg": "用户不存在", "status": false})
			return
		}
		if len(re.ActiveCode) != 0 {
			c.JSON(200, gin.H{"msg": "账号未激活", "status": false})
			return
		}
		err = bcrypt.CompareHashAndPassword([]byte(re.Passwd), []byte(user.Passwd))
		if err != nil {
			c.JSON(200, gin.H{"msg": "密码错误", "status": false})
			return
		}
		outre := map[string]interface{}{
			"uid": re.ID.Hex(),
		}
		err = setSession(c, outre)
		if err != nil {
			c.JSON(200, gin.H{"status": false, "msg": "session数据库错误"})
			return
		}
		c.JSON(200, gin.H{"status": true, "user": re, "msg": "欢迎" + re.Name + "回来"})
	})
}

func logstatus(c *gin.Context) {
	ses, err := getSession(c)
	if err != nil {
		c.JSON(200, gin.H{"status": false, "msg": "获取登录信息失败"})
		return
	}
	if val, ok := ses["uid"]; ok && val != nil {
		var re basicUserModel
		err := pipiC("user", []bson.M{{"$match": bson.M{"_id": bson.ObjectIdHex(val.(string))}}}, &re, false)
		if err != nil {
			log.Println(err.Error())
			c.JSON(200, gin.H{"status": false, "msg": "数据库查询失败"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"user": re, "status": true})
	} else {
		c.JSON(http.StatusOK, gin.H{"status": false, "msg": "未登录"})
	}
}

func logout(c *gin.Context) {
	if err := deleteSession(c); err != nil {
		c.JSON(200, gin.H{"status": false, "msg": err.Error()})
		return
	}
	c.JSON(200, gin.H{"msg": "成功退出登录", "status": true})
}

func changepic(c *gin.Context) {
	var up userPicModel
	if err := c.ShouldBind(&up); err != nil {
		c.JSON(200, gin.H{"status": false, "msg": "参数不够"})
		return
	}
	key := ""
	switch up.Type {
	case "1":
		key = "avatar"
	case "2":
		key = "profilePic"
	case "3":
		key = "golden"
	default:
		c.JSON(200, gin.H{"status": false, "msg": "类型错误"})
		return
	}

	handleUser(func(col *mgo.Collection) {
		var tmpuser map[string]interface{}
		err := col.Find(bson.M{"_id": bson.ObjectIdHex(c.MustGet("auth").(string))}).One(&tmpuser)
		if err != nil {
			c.JSON(200, gin.H{"status": false, "msg": err.Error()})
			return
		}
		if tmpuser[key] != nil {
			oldpath := filepath.Base(tmpuser[key].(string))
			err := os.Remove(globalConf.ResDir + "/pics/" + oldpath)
			if err != nil {
				c.JSON(200, gin.H{"status": false, "msg": "文件系统错误"})
				return
			}
		}
		file, err := c.FormFile("pic")
		if err != nil {
			c.JSON(200, gin.H{"status": false, "msg": err.Error()})
			return
		}
		if err := c.SaveUploadedFile(file, globalConf.ResDir+"/pics/"+file.Filename); err != nil {
			c.JSON(200, gin.H{"status": false, "msg": err.Error()})
			return
		}
		path := globalConf.ResRef + "/pics/" + file.Filename
		err = col.Update(bson.M{"_id": bson.ObjectIdHex(c.MustGet("auth").(string))}, bson.M{"$set": bson.M{key: path}})
		if err != nil {
			c.JSON(200, gin.H{"status": false, "msg": err.Error()})
			return
		}
		c.JSON(200, gin.H{"status": true, "msg": "更换照片成功", "path": path})
	})
}
func updateInfo(c *gin.Context) {
	var cinfo changeInfoModel
	if err := c.ShouldBind(&cinfo); err != nil {
		log.Println(err)
		c.JSON(200, gin.H{"status": false, "msg": "信息不完整"})
		return
	}
	log.Println(cinfo)
	err := updateC("user", bson.M{"_id": bson.ObjectIdHex(c.MustGet("auth").(string))}, bson.M{"$set": cinfo})
	if err != nil {
		c.JSON(200, gin.H{"status": false, "msg": "数据库错误"})
		return
	}
	c.JSON(200, gin.H{"status": true, "msg": "更新个人信息成功"})
	return
}
func getUser(c *gin.Context) {
	id := c.Param("id")
	if len(id) != 24 {
		c.JSON(404, "wrong id")
		return
	}
	var user detailUserModel
	err := findC("user", bson.M{"_id": bson.ObjectIdHex(id)}, false, &user)
	if err != nil {
		c.JSON(404, false)
		return
	}
	re, err := getSession(c)
	user.Authed = re["uid"] == id
	c.JSON(200, user)
	return
}
