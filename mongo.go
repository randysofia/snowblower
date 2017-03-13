package main

import (
	"encoding/json"

	"gopkg.in/mgo.v2"
)

func (e *Event) mongosave() {
	session, err := mgo.Dial("")

	if err != nil {
		panic(err)
	}
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)

	c := session.DB("luminous-analytics").C("snowplow")
	//	e.ID = bson.NewObjectId()

	var x, y map[string]interface{}
	json.Unmarshal([]byte(e.TmpContexts), &x)
	e.Contexts = x
	json.Unmarshal([]byte(e.TmpUnstructuredEvent), &y)
	e.UnstructuredEvent = y
	err = c.Insert(e)
	if err != nil {
		panic(err)
	} else {

	}

}
