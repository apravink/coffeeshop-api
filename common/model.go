//----------------------------------------------------------//
//-----------------------Model------------------------------//
//----------------------------------------------------------//

package common

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

//Drink object
type Drink struct {
	ID           bson.ObjectId `bson:"_id"`
	Name         string        `bson:"name"`
	Price        string        `bson:"price"`
	StartDate    time.Time     `bson:"startDate"`
	EndDate      time.Time     `bson:"endDate"`
	Ingredients  []string      `bson:"ingredients"`
	Availibility bool
}

func (d *Drink) changeAvailability(c bool) {
	d.Availibility = c
}

type Drinks []Drink
