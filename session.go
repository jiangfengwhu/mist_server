package main

import (
	"github.com/gin-gonic/gin"
	"github.com/kidstuff/mongostore"
	"log"
)

func getSession(c *gin.Context) (map[interface{}]interface{}, error) {
	dbsess := initSession.Copy()

	defer dbsess.Close()

	store := mongostore.NewMongoStore(dbsess.DB("mist").C("session"), 3600*12, true, []byte("mist"))
	store.Options.HttpOnly = true

	session, err := store.Get(c.Request, "mist")
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	return session.Values, nil
}

func setSession(c *gin.Context, cont map[string]interface{}) error {
	dbsess := initSession.Copy()
	defer dbsess.Close()

	store := mongostore.NewMongoStore(dbsess.DB("mist").C("session"), 3600*12, true, []byte("mist"))
	store.Options.HttpOnly = true
	// Get a session.
	session, err := store.Get(c.Request, "mist")
	if err != nil {
		log.Println(err.Error())
		return err
	}

	// Add a value.
	for key, val := range cont {
		session.Values[key] = val
	}
	// Save.
	if err = session.Save(c.Request, c.Writer); err != nil {
		log.Printf("Error saving session: %v", err)
		return err
	}
	return nil
}

func deleteSession(c *gin.Context) error {
	dbsess := initSession.Copy()
	defer dbsess.Close()
	store := mongostore.NewMongoStore(dbsess.DB("mist").C("session"), 3600*12, true, []byte("mist"))
	session, err := store.Get(c.Request, "mist")
	if err != nil {
		log.Println(err.Error())
		return err
	}
	session.Options.MaxAge = -1
	if session.Save(c.Request, c.Writer); err != nil {
		log.Println("save err", err.Error())
		return err
	}
	return nil
}
