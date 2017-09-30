package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	common "coffeeshop/common"

	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func main() {

	session, err := mgo.Dial("localhost")
	if err != nil {
		panic(err)
	}

	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	ensureIndex(session)

	handleRequests(session)
}

func ensureIndex(s *mgo.Session) {
	session := s.Copy()
	defer session.Close()
	c := session.DB("coffeeshop").C("drinks2")

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

//Handle Requests
func handleRequests(s *mgo.Session) {

	session := s.Copy()
	defer session.Close()
	my_router := mux.NewRouter().StrictSlash(true)

	my_router.HandleFunc("/", homePage).Methods("GET")
	my_router.HandleFunc("/all", getAllDrinks(session)).Methods("GET")

	http.ListenAndServe(":8081", my_router)

}

//Handler Functions
func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the HomePage!")
	fmt.Println("Endpoint Hit: homePage")
}

func getAllDrinks(s *mgo.Session) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		session := s.Copy()
		defer session.Close()

		c := session.DB("coffeeshop").C("drinks")

		var all_drinks common.Drinks
		err := c.Find(bson.M{}).All(&all_drinks)
		if err != nil {

			fmt.Println("Failed get all drinks: ", err)
			return
		}

		respBody, err := json.MarshalIndent(all_drinks, "", "  ")
		if err != nil {
			log.Fatal(err)
		}

		ResponseWithJSON(w, respBody, http.StatusOK)
	}
}

func ResponseWithJSON(w http.ResponseWriter, json []byte, code int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	w.Write(json)
}
