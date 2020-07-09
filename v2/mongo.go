package main

import (
	"context"
	"os"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type geo struct {
	Lng float32 `bson:"lng,omitempty"`
	Lat float32 `bson:"lat,omitempty"`
}

type mongoe struct {
	Event       `bson:",inline"`
	MongoUserID primitive.ObjectID `bson:"userid_mgo,omitempty"`
	GeoCoord    geo                `bson:"geo_coord,omitempty"`
}

func (standarde *Event) mongosave() {

	var e mongoe

	if oid, err := primitive.ObjectIDFromHex(standarde.UserID); err == nil {
		e = mongoe{*standarde, oid, geo{0, 0}}
	} else {
		e = mongoe{*standarde, primitive.NewObjectID(), geo{0, 0}}
	}
	e.GeoCoord.Lat = standarde.GeoLatitude
	e.GeoCoord.Lng = standarde.GeoLongitude
	e.GeoLatitude = 0
	e.GeoLongitude = 0

	if _, err := client.
		Database(os.Getenv("MONGO_DB")).
		Collection(os.Getenv("MONGO_COLLECTION")).
		InsertOne(context.Background(), e); err != nil {
		panic(err)
	}
}
