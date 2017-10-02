package common

import (
	"net/http"

	"github.com/gorilla/mux"
	mgo "gopkg.in/mgo.v2"
)

//Handle Requests
func HandleRequests(s *mgo.Session) {

	session := s.Copy()
	defer session.Close()
	my_router := mux.NewRouter().StrictSlash(true)

	my_router.HandleFunc("/", homePage).Methods("GET")
	my_router.HandleFunc("/drinks", getAllDrinks(session)).Methods("GET")
	// POST /drinks/
	my_router.HandleFunc("/drinks/", createDrink(session)).Methods("POST")
	// GET /drinks/{name}
	my_router.HandleFunc("/drinks/{name}", drinkByName(session)).Methods("GET")
	// DELETE /drinks/{name}
	my_router.HandleFunc("/drinks/{name}", removeDrink(session)).Methods("DELETE")
	// GET /byDate/{date}
	my_router.HandleFunc("/byDate/{date}", drinksByDate(session)).Methods("GET")
	// GET /byIngredients/:ingredients OPTIONAL

	http.ListenAndServe(HOST, my_router)

}
