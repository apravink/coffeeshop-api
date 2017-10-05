package main

import (
	common "github.com/apravink/coffeeshop-api/common"

	"gopkg.in/mgo.v2"
)

const DBHOST = "localhost:27017"

func main() {

	//Connect to DB
	session, err := mgo.Dial(DBHOST)
	if err != nil {
		panic(err)
	}

	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	ensureIndex(session)

	common.HandleRequests(session)
}

//Check if db and collection exist. If not, create them.
func ensureIndex(s *mgo.Session) {
	session := s.Copy()
	defer session.Close()
	c := session.DB("coffeeshop").C("drinks")

	index := mgo.Index{
		Key:        []string{"id"},
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	}
	err := c.EnsureIndex(index)
	if err != nil {
		panic(err)
	}
	println("Successfully connected to the database")
}
