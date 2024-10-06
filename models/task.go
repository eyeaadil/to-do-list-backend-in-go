package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Task struct {
	ID primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name  string `json:"name,omitempty" bson:"name omitempty"`
	Description string `json:"description,omitempty" bson:"description,omitempty"`
	Completed bool `json:"completed,omitempty" bson:"completed,omitempty"`
	Status string `json:"status,omitempty" bson:"status,omitempty"`
}