package main

import (
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"log"
)

var initSession *mgo.Session

func dial() {
	log.Println("call database")
	var err error
	initSession, err = mgo.Dial("mongodb://127.0.0.1:27017")
	if err != nil {
		log.Println("database connection failed")
	}
	initSession.SetMode(mgo.Monotonic, true)
}
func updateC(col string, selector interface{}, update interface{}) error {
	s := initSession.Copy()
	defer s.Close()
	return s.DB("mist").C(col).Update(selector, update)
}
func insertC(col string, docs interface{}) error {
	s := initSession.Copy()
	defer s.Close()
	return s.DB("mist").C(col).Insert(docs)
}
func findC(col string, query interface{}, isall bool, result interface{}) error {
	s := initSession.Copy()
	defer s.Close()
	if isall {
		return s.DB("mist").C(col).Find(query).Sort("-_id").All(result)
	}
	return s.DB("mist").C(col).Find(query).One(result)
}
func getCount(col string, query interface{}) (int, error) {
	s := initSession.Copy()
	defer s.Close()
	return s.DB("mist").C(col).Find(query).Count()
}
func pipiC(col string, pipeline interface{}, result interface{}, all bool) error {
	s := initSession.Copy()
	defer s.Close()
	if all {
		return s.DB("mist").C(col).Pipe(pipeline).All(result)
	}
	return s.DB("mist").C(col).Pipe(pipeline).One(result)
}
func latestC(col string, pipeline []bson.M, skip int, limit int, result interface{}) error {
	s := initSession.Copy()
	defer s.Close()
	pipeline = append(pipeline, lookowner, unwind, bson.M{"$sort": bson.M{"_id": -1}}, bson.M{"$limit": limit}, bson.M{"$skip": skip})
	return s.DB("mist").C(col).Pipe(pipeline).All(result)
}
func findAndModify(col string, findQuery interface{}, update interface{}, upsert bool, remove bool, renew bool, re interface{}) error {
	s := initSession.Copy()
	defer s.Close()
	change := mgo.Change{
		Update:    update,
		Upsert:    upsert,
		Remove:    remove,
		ReturnNew: renew,
	}
	_, err := s.DB("mist").C(col).Find(findQuery).Apply(change, re)
	return err
}
func handleUser(dbhandle func(col *mgo.Collection)) {
	s := initSession.Copy()
	defer s.Close()
	dbhandle(s.DB("mist").C("user"))
}
func handleVideo(dbhandle func(col *mgo.Collection)) {
	s := initSession.Copy()
	defer s.Close()
	dbhandle(s.DB("mist").C("video"))
}
