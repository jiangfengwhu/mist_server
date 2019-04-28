package main

import (
	"mime/multipart"

	"github.com/globalsign/mgo/bson"
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

type basicVideoModel struct {
	Title    string        `bson:"title" json:"title" binding:"required"`
	Date     int64         `bson:"date" json:"date"`
	ID       bson.ObjectId `bson:"_id" json:"id"`
	Cover    string        `bson:"cover" json:"cover"`
	View     int64         `bson:"view" json:"view"`
	Comments int64         `bson:"comments_length,omitempty" json:"comments,omitempty"`
	PlayList bson.ObjectId `bson:"playlist,omitempty" json:"playlist,omitempty"`
}
type userVideoModel struct {
	basicVideoModel `bson:",inline"`
	ListDoc         *outListModel `json:"listdoc,omitempty" bson:"listdoc,omitempty"`
}
type faceVideoModel struct {
	basicVideoModel `bson:",inline"`
	Owner           bson.ObjectId  `bson:"owner" json:"-"`
	OwnerDoc        basicUserModel `bson:"owner_doc,omitempty" json:"owner"`
	outLikeModel    `bson:",inline"`
}
type newVideoModel struct {
	faceVideoModel `bson:",inline"`
	Desc           string `bson:"desc,omitempty" json:"desc,omitempty"`
	Path           string `bson:"path" json:"path"`
	Hash           string `bson:"-" json:"vid,omitempty" binding:"required"`
	Tag            int8   `bson:"tag" json:"tag" binding:"required"`
	CoverPos       uint16 `bson:"-" json:"coverPos" binding:"required"`
}
type updateVideoModel struct {
	ID    bson.ObjectId `bson:"_id" json:"id" binding:"required"`
	Tag   int8          `bson:"tag" json:"tag" binding:"required"`
	Desc  string        `bson:"desc,omitempty" json:"desc,omitempty"`
	Title string        `bson:"title" json:"title" binding:"required"`
}
type detailVideoModel struct {
	faceVideoModel `bson:",inline"`
	Desc           string            `bson:"desc,omitempty" json:"desc,omitempty"`
	Path           string            `bson:"path" json:"path"`
	Tag            int8              `bson:"tag" json:"tag" binding:"required"`
	Recommend      []basicVideoModel `json:"recommend"`
	ListVideos     []basicVideoModel `json:"plists,omitempty"`
	ListDoc        *outListModel     `json:"listdoc,omitempty" bson:"listdoc,omitempty"`
}
type delVideoModel struct {
	Path  string `bson:"path" json:"path"`
	Cover string `bson:"cover" json:"cover"`
}
type uploadFileModel struct {
	Blob *multipart.FileHeader `json:"blob" form:"blob" binding:"required"`
	ID   string                `json:"vid" form:"vid" binding:"required"`
}

type changeCoverModel struct {
	Cid  bson.ObjectId `json:"cid" binding:"required"`
	Path string        `json:"path" form:"path"`
}

type getVModel struct {
	Tag   int8 `form:"tag" binding:"required"`
	Size  int  `form:"size" binding:"required"`
	Start *int `form:"fi" binding:"exists"`
}
type getOModel struct {
	Key   string `form:"key"`
	Size  int    `form:"size" binding:"required"`
	Start *int   `form:"fi" binding:"exists"`
}
type userPicModel struct {
	// Blob *multipart.FileHeader `form:"pic" binding:"required"`
	Type string `form:"type" binding:"required"`
}

type circleModel struct {
	ID      bson.ObjectId `bson:"_id" json:"id"`
	Content string        `json:"cont,omitempty" bson:"cont,omitempty"`
	Embed   string        `json:"embed,omitempty" bson:"embed,omitempty"`
	Pics    []string      `bson:"pics,omitempty" json:"pics,omitempty"`
	Date    int64         `bson:"date" json:"date"`
	Cover   string        `bson:"cover,omitempty" json:"cover,omitempty"`
}
type outCircleModel struct {
	circleModel  `bson:",inline"`
	Owner        bson.ObjectId  `bson:"owner" json:"-"`
	OwnerDoc     basicUserModel `bson:"owner_doc,omitempty" json:"owner,omitempty"`
	outLikeModel `bson:",inline"`
	Comments     int64 `bson:"comments_length,omitempty" json:"comments,omitempty"`
}
type ownCircleModel struct {
	circleModel  `bson:",inline"`
	outLikeModel `bson:",inline"`
	Comments     int64 `bson:"comments_length,omitempty" json:"comments,omitempty"`
}
type newGalleryModel struct {
	Pics     []string       `bson:"pics" json:"pics"`
	Content  string         `json:"cont,omitempty" bson:"cont,omitempty"`
	ID       bson.ObjectId  `bson:"_id" json:"id"`
	Date     int64          `bson:"date" json:"date"`
	Owner    bson.ObjectId  `bson:"owner" json:"-"`
	OwnerDoc basicUserModel `bson:"owner_doc,omitempty" json:"owner"`
}
type faceGalleryModel struct {
	Pics    []string      `bson:"pics" json:"pics"`
	Content string        `json:"cont,omitempty" bson:"cont,omitempty"`
	ID      bson.ObjectId `bson:"_id" json:"id"`
}
type outGalleryModel struct {
	newGalleryModel `bson:",inline"`
	outLikeModel    `bson:",inline"`
	Comments        int64 `bson:"comments_length,omitempty" json:"comments,omitempty"`
}
type commentModel struct {
	Owner        bson.ObjectId `bson:"owner" json:"owner"`
	Content      string        `bson:"text" json:"text"`
	Date         int64         `bson:"date" json:"date"`
	Reply        int           `bson:"comments,omitempty" json:"reply,omitempty"`
	At           string        `bson:"at,omitempty" json:"at,omitempty"`
	ID           bson.ObjectId `bson:"_id" json:"id"`
	outLikeModel `bson:",inline"`
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
type outLikeModel struct {
	Likes   int64 `bson:"likes_length,omitempty" json:"likes,omitempty"`
	IsLiked bool  `bson:"isLiked,omitempty" json:"isliked,omitempty"`
}
type likeModel struct {
	ID   string `bson:"_id" form:"id" binding:"required"`
	Type string `form:"type" binding:"required"`
	Inc  *int8  `form:"inc" binding:"exists"`
}
type searchModel struct {
	Key   string `form:"key" binding:"required"`
	Size  int    `form:"size" binding:"required"`
	Start *int   `form:"fi" binding:"exists"`
}

