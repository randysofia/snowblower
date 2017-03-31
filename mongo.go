package main

import (
	"fmt"
	"os"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type geo struct {
	Lng float32 `bson:"lng,omitempty"`
	Lat float32 `bson:"lat,omitempty"`
}

type mongoe struct {
	Event       `bson:",inline"`
	MongoUserID bson.ObjectId `bson:"userid_mgo,omitempty"`
	GeoCoord    geo           `bson:"geo_coord,omitempty"`
}

func (standarde *Event) mongosave() {

	var e mongoe

	if !bson.IsObjectIdHex(standarde.UserID) {
		e = mongoe{*standarde, bson.ObjectId(""), geo{0, 0}}
	} else {
		e = mongoe{*standarde, bson.ObjectIdHex(standarde.UserID), geo{0, 0}}
	}
	e.GeoCoord.Lat = standarde.GeoLatitude
	e.GeoCoord.Lng = standarde.GeoLongitude
	e.GeoLatitude = 0
	e.GeoLongitude = 0
	session, err := mgo.Dial(os.Getenv("MONGO_URI"))

	if err != nil {
		panic(err)
	}
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)

	c := session.DB(os.Getenv("MONGO_DB")).C(os.Getenv("MONGO_COLLECTION"))
	//	e.ID = bson.NewObjectId()

	err = c.Insert(e)
	if err != nil {
		panic(err)
	} else {

	}
	fmt.Println("saved")
}
