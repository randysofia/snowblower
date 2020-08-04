package main

import (
	"context"
	"os"

	"github.com/spark451/snowman"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type mongoe struct {
	Event       `bson:",inline"`
	MongoUserID primitive.ObjectID `bson:"userid_mgo,omitempty"`
	GeoCoord    *snowman.Geo       `bson:"geo_coord,omitempty"`
}

func (standarde *Event) mongosave() {

	var e mongoe

	if oid, err := primitive.ObjectIDFromHex(standarde.UserID); err == nil {
		e = mongoe{*standarde, oid, nil}
	} else {
		e = mongoe{*standarde, primitive.NewObjectID(), nil}
	}
	if e.GeoLatitude != 0 || e.GeoLongitude != 0 {
		e.GeoCoord = &snowman.Geo{
			Lat: standarde.GeoLatitude,
			Lng: standarde.GeoLongitude,
		}
		e.GeoLatitude = 0
		e.GeoLongitude = 0
	}

	if _, err := client.
		Database(os.Getenv("MONGO_DB")).
		Collection(os.Getenv("MONGO_COLLECTION")).
		InsertOne(context.Background(), e); err != nil {
		panic(err)
	}
}
