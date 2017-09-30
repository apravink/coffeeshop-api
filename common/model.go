package common

import (
	"time"
)

//Drink object
type Drink struct {
	ID           int       `json:"id"`
	Name         string    `json:"name"`
	Price        string    `json:"price"`
	StartDate    time.Time `bson:"startDate"`
	EndDate      time.Time `bson:"endDate"`
	Ingredients  []string  `bson:"ingredients"`
	Availibility bool
}

func (d *Drink) changeAvailability(c bool) {
	d.Availibility = c
}

type Drinks []Drink
