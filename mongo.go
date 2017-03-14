package main

import (
	"os"

	"gopkg.in/mgo.v2"
)

func (e *Event) mongosave() {
	session, err := mgo.Dial(os.Getenv("MONGO_URI"))

	if err != nil {
		panic(err)
	}
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)

	c := session.DB("luminous-analytics").C("snowplow")
	//	e.ID = bson.NewObjectId()

	err = c.Insert(e)
	if err != nil {
		panic(err)
	} else {

	}

}
