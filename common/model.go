package common

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

//Drink object
type Drink struct {
	ID          bson.ObjectId `bson:"id"`
	Name        string        `json:"name"`
	Price       string        `json:"price"`
	StartDate   time.Time     `bson:"startDate"`
	EndDate     time.Time     `bson:"endDate"`
	Ingredients []string      `bson:"ingredients"`
}

type Drinks []Drink
