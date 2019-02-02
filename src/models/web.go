package models

import "github.com/globalsign/mgo/bson"

type Web struct {
	Id             bson.ObjectId `json:"id"        bson:"_id,omitempty"`
	WebId          string        `json:"webId"        bson:"webId,omitempty"`
	Url            string        `json:"url" binding:"required"`
	Web            string        `json:"web"`
	LastDateUpdate string        `json:"lastDateUpdate"`
}
