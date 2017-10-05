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

/*--------------------------------------------------------------------*/
//Handler Functions
/*--------------------------------------------------------------------*/
//Static Homepage
/* Params: w <http.ResponseWriter>
			   r <*http.Request>
 Description: Writes a simple message to the response
*/
func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the HomePage!")
	fmt.Println("Endpoint Hit: homePage")
}

//Get All Drinks
/* Params: s *mgo.Session

Returns: func

Description: Wrapper function factory for a function that returns all the
drinks in the database, along with their current availability
*/
func getAllDrinks(s *mgo.Session) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		session := s.Copy()
		defer session.Close()

		c := session.DB(DATABASE).C(COLLECTION)

		var all_drinks Drinks

		//Query Mongo for all entries
		err := c.Find(bson.M{}).All(&all_drinks)
		if err != nil {
			ErrorWithJSON(w, "Something went wrong", http.StatusInternalServerError)
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
/* Params: s *mgo.Session

Returns: func

Description:Wrapper function factory for a function that gathers the name
variable from the url and returns a drink from the database if there is
a match, or a StatusNotFound error
*/
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
/* Params: s *mgo.Session

Returns: func

Description:Wrapper function factory for a function that creates a new drink
in the database according to the accompanying JSON payload in its body.
*/
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

//Delete drink from database
/* Params: s *mgo.Session

Returns: func

Description:Wrapper function factory for a function that gathers the name
variable from the url and returns a drink from the database if there is
a match, or a StatusNotFound error
*/

func removeDrink(s *mgo.Session) func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {
		session := s.Copy()
		defer session.Close()

		//Get name from uri
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

//Get drinks available on a certain date
/*Returns: func

Description:Wrapper function factory for a function that gathers the date from
the url and returns a list of drinks from the database that are available
for that period
*/

func drinksByDate(s *mgo.Session) func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {
		session := s.Copy()
		defer session.Close()

		//Get date from uri
		vars := mux.Vars(r)
		date := vars["date"]
		time_layout := "2006-01-02"
		date_utc, _ := time.Parse(time_layout, date)
		var available_drinks, all_drinks Drinks

		c := session.DB(DATABASE).C(COLLECTION)

		//Grabbing all and filtering. ToDo: Query the DB directly instead
		err := c.Find(bson.M{}).All(&all_drinks)

		if err != nil {
			ErrorWithJSON(w, "Something went wrong", http.StatusInternalServerError)
			fmt.Println("Failed get all drinks: ", err)
			return
		}

		//Set availability according to given date
		for i, _ := range all_drinks {

			setAvailibility(&all_drinks[i], date_utc)

			if all_drinks[i].Availibility == true {
				available_drinks = append(available_drinks, all_drinks[i])
			}
		}
		if len(available_drinks) == 0 {
			ErrorWithJSON(w, "Drink not found", http.StatusNotFound)
			return

		}

		respBody, err := json.MarshalIndent(available_drinks, "", "  ")
		if err != nil {
			log.Fatal(err)
		}

		ResponseWithJSON(w, respBody, http.StatusOK)
	}
}

//Get Drinks by ingredients
/* Params: s *mgo.Session

Returns: func

Description: Wrapper function factory for a function that returns the drinks
matching the ingredient list passed in the URL
*/
func getDrinksByIngredient(s *mgo.Session) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		session := s.Copy()
		defer session.Close()

		c := session.DB(DATABASE).C(COLLECTION)

		//Get query string array
		r.ParseForm()
		ingredients := r.Form["ingredients"]
		if len(ingredients) < 1 {
			fmt.Println("No ingredients specified")
			ErrorWithJSON(w, "No ingredients specified", http.StatusNoContent)
		}

		var all_drinks Drinks

		//Pass ingredients list and query drinks with ingredients
		err := c.Find(bson.M{
			"ingredients": bson.M{"$all": ingredients}}).All(&all_drinks)
		if err != nil {
			ErrorWithJSON(w, "Something went wrong", http.StatusInternalServerError)
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
