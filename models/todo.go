package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Todo struct {
	Id         primitive.ObjectID `json:"_id" bson:"_id"`
	Title      string             `json:"title" bson:"title"`
	Timestamp  int64              `json:"timestamp" bson:"timestamp"`
	IsComplete bool               `json:"isComplete" bson:"isComplete"`
	UserId     int64              `json:"userId" bson:"userId"`
}
