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
	my_router.HandleFunc("/drinks", getAllDrinks(session)).Methods("GET")
	// GET /drinks/{date}
	// GET /drinks/{name}
	// GET /drinks/:ingredients OPTIONAL
	// POST /drinks/
	// DELETE /all/drinks/{name}
	// DELETE /all/drinks/{id}
	//

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

		//Set availability of drinks according to current date
		for i, _ := range all_drinks {
			setAvailibility(&all_drinks[i], time.Now())
		}
		respBody, err := json.MarshalIndent(all_drinks, "", "  ")
		if err != nil {
			log.Fatal(err)
		}

		ResponseWithJSON(w, respBody, http.StatusOK)
	}
}

//Get a drink by name

//--------------------------------------------------------------------------//
//Helper Functions
//--------------------------------------------------------------------------//

//Set availability of drink
func setAvailibility(drink *Drink, date time.Time) {

	if date.After(drink.StartDate) && date.Before(drink.EndDate) {
		drink.changeAvailability(true)
	}

}

//Response Writer
func ResponseWithJSON(w http.ResponseWriter, json []byte, code int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	w.Write(json)
}
