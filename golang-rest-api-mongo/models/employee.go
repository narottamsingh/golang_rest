package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Employee struct {
	Id   primitive.ObjectID `json:"id,omitempty`
	Name string             `json:"name,omitempty"  validate:"required"`
	Age  int                `json:"age,omitempty" validare:"required"`
}
