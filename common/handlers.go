package common

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

//Handle Requests
func HandleRequests(s *mgo.Session) {

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

//Get all drinks
func getAllDrinks(s *mgo.Session) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		session := s.Copy()
		defer session.Close()

		c := session.DB("coffeeshop").C("drinks")

		var all_drinks Drinks

		err := c.Find(bson.M{}).All(&all_drinks)
		if err != nil {

			fmt.Println("Failed get all drinks: ", err)
			return
		}
		all_drinks = setAvailibility(all_drinks)
		respBody, err := json.MarshalIndent(all_drinks, "", "  ")
		if err != nil {
			log.Fatal(err)
		}

		ResponseWithJSON(w, respBody, http.StatusOK)
	}
}

//Helper Functions

//TODO
func setAvailibility(drinks Drinks) Drinks {
	drinks_set := drinks
	current_date := time.Now()

	for i, drink := range drinks_set {
		if current_date.After(drink.StartDate) && current_date.Before(drink.EndDate) {
			drinks_set[i].changeAvailability(true)
		} else {
			drinks_set[i].changeAvailability(false)
		}
	}
	fmt.Print(drinks_set)
	return drinks_set
}
func ResponseWithJSON(w http.ResponseWriter, json []byte, code int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	w.Write(json)
}
