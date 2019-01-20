package main

import (
	"github.com/globalsign/mgo/bson"
	"mime/multipart"
)

type registModel struct {
	Name   string `bson:"nickName" json:"nickName" binding:"required"`
	Email  string `bson:"email" json:"email" binding:"required"`
	Passwd string `bson:"passwd" json:"passwd" binding:"required"`
	Token  string `json:"token" binding:"required"`
}
type basicUserModel struct {
	Name   string        `bson:"nickName" json:"nickName"`
	Avatar string        `bson:"avatar" json:"avatar,omitempty"`
	ID     bson.ObjectId `bson:"_id" json:"uid"`
}
type changeInfoModel struct {
	Name     string `bson:"nickName" json:"nickName" binding:"required"`
	Sign     string `bson:"sign,omitempty" json:"sign,omitempty"`
	Dirth    string `json:"birth,omitempty" bson:"birth,omitempty"`
	Homepage string `json:"homepage,omitempty" bson:"homepage,omitempty"`
}

type reCapResponseModel struct {
	Success bool     `json:"success"`
	Err     []string `json:"error-codes,omitempty"`
	Time    string   `json:"challenge_ts,omitempty"`
	Host    string   `json:"hostname,omitempty"`
	Score   float32  `json:"score,omitempty"`
	Action  string   `json:"action,omitempty"`
}

type activeUserModel struct {
	Email      string `json:"email" bson:"email" form:"email" binding:"required"`
	ActiveCode string `json:"activeCode" bson:"activeCode" form:"activeCode" binding:"required"`
	Expire     int64  `json:"expireDate" bson:"expireDate" form:"expireDate" binding:"required"`
}
type changeUserModel struct {
	basicUserModel `bson:",inline"`
	ProPic         string `json:"profilePic,omitempty" bson:"profilePic"`
}
type loginModel struct {
	Email  string `json:"email" bson:"email" binding:"required"`
	Passwd string `json:"passwd" bson:"passwd" binding:"required"`
	Token  string `json:"token" binding:"required"`
}
type detailUserModel struct {
	ID       bson.ObjectId `json:"uid" bson:"_id"`
	Name     string        `json:"nickName" bson:"nickName"`
	Avatar   string        `json:"avatar,omitempty" bson:"avatar"`
	ProPic   string        `json:"profilePic,omitempty" bson:"profilePic"`
	Sign     string        `json:"sign,omitempty" bson:"sign,omitempty"`
	Gloden   string        `json:"golden,omitempty" bson:"golden,omitempty"`
	Jointime int64         `json:"jd" bson:"jd"`
	Authed   bool          `json:"authed,omitempty"`
	Dirth    string        `json:"birth,omitempty" bson:"birth,omitempty"`
	Homepage string        `json:"homepage,omitempty" bson:"homepage,omitempty"`
}
type loginDBModel struct {
	basicUserModel `bson:",inline"`
	Email          string `json:"-" bson:"email"`
	Passwd         string `json:"-" bson:"passwd"`
	ActiveCode     string `json:"-" bson:"activeCode"`
}
type vcScalfold struct {
	Title    string   `bson:"title" json:"title" binding:"required"`
	Price    *int     `bson:"price" json:"price" binding:"exists"`
	Desc     string   `bson:"desc" json:"desc"`
	Tags     []string `bson:"tags" json:"tags" binding:"required"`
	Cover    string   `bson:"cover" json:"cover"`
	Date     int64    `bson:"date" json:"date"`
	View     int64    `bson:"view" json:"view"`
	Comments int64    `bson:"comments_length,omitempty" json:"comments"`
}
type myVCModel struct {
	vcScalfold `bson:",inline"`
	ID         bson.ObjectId `bson:"_id" json:"id" binding:"required"`
}
type basicVCModel struct {
	vcScalfold `bson:",inline"`
	ID         bson.ObjectId  `bson:"_id" json:"id"`
	Owner      bson.ObjectId  `bson:"owner" json:"-"`
	OwnerDoc   basicUserModel `bson:"owner_doc,omitempty" json:"owner"`
	Likes      int64          `bson:"likes_length,omitempty" json:"likes,omitempty"`
	IsLiked    bool           `bson:"isLiked,omitempty" json:"isliked,omitempty"`
}
type basicVideoModel struct {
	Title string `bson:"title" json:"title"`
	Desc  string `bson:"desc" json:"desc"`
	Date  int64  `bson:"date" json:"date"`
	ID    string `bson:"_id" json:"id"`
	Cover string `bson:"cover" json:"cover"`
	Path  string `bson:"path" json:"path"`
}
type uploadFileModel struct {
	Blob *multipart.FileHeader `json:"blob" binding:"required"`
	ID   string                `json:"vid" form:"vid" binding:"required"`
}

type addVideoModel struct {
	Vid   string        `json:"vid" binding:"required"`
	Cid   bson.ObjectId `json:"cid" binding:"required"`
	Title string        `json:"title" binding:"required"`
	Desc  string        `json:"desc"`
}

type changeCoverModel struct {
	Cid  bson.ObjectId `json:"cid" binding:"required"`
	Path string        `json:"path" form:"path"`
}

type outVideoModel struct {
	basicVCModel `bson:",inline"`
	Videos       []basicVideoModel `bson:"videos" json:"videos"`
}
type getVModel struct {
	Size  int  `form:"size" binding:"required"`
	Start *int `form:"fi" binding:"exists"`
}
type userPicModel struct {
	// Blob *multipart.FileHeader `form:"pic" binding:"required"`
	Type string `form:"type" binding:"required"`
}
type delVideoQ struct {
	Cid   string `form:"cid" binding:"required"`
	Vid   string `form:"vid" binding:"required"`
	Cover string `form:"cov" binding:"required"`
}
type circleModel struct {
	ID      bson.ObjectId `bson:"_id" json:"id"`
	Content string        `json:"cont,omitempty" bson:"cont,omitempty"`
	Pics    []string      `bson:"pics,omitempty" json:"pics,omitempty"`
	Date    int64         `bson:"date" json:"date"`
}
type outCircleModel struct {
	circleModel `bson:",inline"`
	Owner       bson.ObjectId  `bson:"owner" json:"-"`
	OwnerDoc    basicUserModel `bson:"owner_doc,omitempty" json:"owner"`
}
type commentModel struct {
	Owner   bson.ObjectId `bson:"owner" json:"owner"`
	Content string        `bson:"text" json:"text"`
	Date    int64         `bson:"date" json:"date"`
	Reply   int           `bson:"comments,omitempty" json:"reply,omitempty"`
	At      string        `bson:"at,omitempty" json:"at,omitempty"`
	ID      bson.ObjectId `bson:"_id" json:"id"`
	Likes   int64         `bson:"likes_length,omitempty" json:"likes,omitempty"`
	IsLiked bool          `bson:"isLiked,omitempty" json:"isliked,omitempty"`
}
type outCommentModel struct {
	CommentsDoc []commentModel   `bson:"comments_doc,omitempty" json:"comments,omitempty"`
	Owners      []basicUserModel `bson:"owners_doc,omitempty" json:"owners,omitempty"`
}
type postCommentModel struct {
	ID      string `json:"id" binding:"required"`
	Content string `json:"text" binding:"required"`
	Type    string `json:"type" binding:"required"`
	At      string `json:"at"`
}
type getCommentsModel struct {
	ID   string `bson:"_id" form:"id" binding:"required"`
	Type string `form:"type" binding:"required"`
}
type likeModel struct {
	ID   string `bson:"_id" form:"id" binding:"required"`
	Type string `form:"type" binding:"required"`
	Inc  *int8  `form:"inc" binding:"exists"`
}
