package main

import (
	common "coffeeshop-api/common"

	"gopkg.in/mgo.v2"
)

func main() {

	session, err := mgo.Dial("localhost")
	if err != nil {
		panic(err)
	}

	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	ensureIndex(session)

	common.HandleRequests(session)
}

//Ensure index is configured correctly
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
}
