package common

import (
	"fmt"
	"net/http"
	"time"
)

//--------------------------------------------------------------------------//
//Helper Functions
//--------------------------------------------------------------------------//

//Set availability of drink
func setAvailibility(drink *Drink, date time.Time) {

	if date.After(drink.StartDate) && date.Before(drink.EndDate) {
		drink.changeAvailability(true)
	}

}

//Http status response writer
func ResponseWithJSON(w http.ResponseWriter, json []byte, code int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	w.Write(json)
}

//Http status error writer
func ErrorWithJSON(w http.ResponseWriter, message string, code int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	fmt.Fprintf(w, "{message: %q}", message)
}
