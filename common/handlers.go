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

const HOST = "localhost:8081"
const DATABASE = "coffeeshop"
const COLLECTION = "drinks"

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
	// GET /byIngredients/:ingredients OPTIONAL

	//

	http.ListenAndServe(HOST, my_router)

}

//Handler Functions

//Static Homepage
func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the HomePage!")
	fmt.Println("Endpoint Hit: homePage")
}

//Function Wrapper for get all drinks
func getAllDrinks(s *mgo.Session) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		session := s.Copy()
		defer session.Close()

		c := session.DB(DATABASE).C(COLLECTION)

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
func drinkByName(s *mgo.Session) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		session := s.Copy()
		defer session.Close()

		c := session.DB(DATABASE).C(COLLECTION)
		//Get params from the call
		vars := mux.Vars(r)
		drink_name := vars["name"]

		var drink Drink

		err := c.Find(bson.M{"name": drink_name}).One(&drink)
		if err != nil {
			ErrorWithJSON(w, "Drink not Found", http.StatusNotFound)
			fmt.Print("Failed to find drink: ", err)
			fmt.Println("")
			return
		}
		if drink.Name == "" {
			ErrorWithJSON(w, "Drink not found", http.StatusNotFound)
			return
		}
		setAvailibility(&drink, time.Now())
		respBody, err := json.MarshalIndent(drink, "", "  ")
		if err != nil {
			log.Fatal(err)
		}
		ResponseWithJSON(w, respBody, http.StatusOK)

	}
}

//Create Drink
func createDrink(s *mgo.Session) func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {
		session := s.Copy()
		defer session.Close()

		var drink Drink

		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&drink)
		if err != nil {
			ErrorWithJSON(w, "Incorrect body", http.StatusBadRequest)
			return
		}
		drink.ID = bson.NewObjectId()
		setAvailibility(&drink, time.Now())
		c := session.DB(DATABASE).C(COLLECTION)
		err2 := c.Insert(drink)
		if err2 != nil {

			switch {
			case mgo.IsDup(err):
				ErrorWithJSON(w, "Drink already exists!", http.StatusBadRequest)
				return
			default:
				ErrorWithJSON(w, "Something went wrong", http.StatusInternalServerError)
				fmt.Print("Failed to insert drink", err2)
				return
			}
		}
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Location", r.URL.Path+"/"+string(drink.ID))
		w.WriteHeader(http.StatusCreated)
	}
}

//Delete drink

func removeDrink(s *mgo.Session) func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {
		session := s.Copy()
		defer session.Close()

		//Get params from the call
		vars := mux.Vars(r)
		drink_name := vars["name"]

		c := session.DB(DATABASE).C(COLLECTION)

		err := c.Remove(bson.M{"name": drink_name})
		if err != nil {
			switch err {
			case mgo.ErrNotFound:
				ErrorWithJSON(w, "Drink not found", http.StatusNotFound)
				return
			default:
				ErrorWithJSON(w, "Something went wrong", http.StatusInternalServerError)
				return
			}
		}
		w.WriteHeader(http.StatusOK)
	}
}

//--------------------------------------------------------------------------//
//Helper Functions
//--------------------------------------------------------------------------//

//Set availability of drink
func setAvailibility(drink *Drink, date time.Time) {

	if date.After(drink.StartDate) && date.Before(drink.EndDate) {
		drink.changeAvailability(true)
	}

}

//http status response writer
func ResponseWithJSON(w http.ResponseWriter, json []byte, code int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	w.Write(json)
}

//http status error writer
func ErrorWithJSON(w http.ResponseWriter, message string, code int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	fmt.Fprintf(w, "{message: %q}", message)
}
