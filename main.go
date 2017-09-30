package main

import (
	"fmt"
	"net/http"

	"gopkg.in/mgo.v2"
)

func main() {
	manageDatabase()
	handleRequests()
}

func handleRequests() {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		fmt.Fprintf(w, "Welcome to the homepage!")
		fmt.Println("Homepage Hit!")
	})
	http.ListenAndServe(":8081", nil)
	fmt.Println("Blah")
}

func manageDatabase() {
	session, err := mgo.Dial("localhost")
	if err != nil {
		panic(err)
	}
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("coffeeshop").C("drinks2")

	index := mgo.Index{
		Key:        []string{"id"},
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	}
	err2 := c.EnsureIndex(index)
	if err2 != nil {
		panic(err)
	}
}
