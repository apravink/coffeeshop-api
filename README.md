# coffeeshop-api
---

## Dependencies  

* "gopkg.in/mgo.v2"
* "github.com/gorilla/mux"   

## TODO  

- [x] Create a server  
- [x] Create DB connection
- [X] Create a route to display all  
- [x] CRUD routes and functions  
- [x] Add dateTime constraint     

## API Documentation    


### Data Model  
ID - Primary key id  
Name - Name of the drink
Price - Price of the drink in CAD  
StartDate - Date when drink is available  
EndDate - Date after which drink is unavailable  
ingredients - Array of ingredients   
Availibility - Whether or not the drink is available today    


### Endpoints   
`GET "/"`- Returns a simple printed message   
`GET "/drinks"`- Returns all drinks  
`GET "/drinks/name"`- Returns a drink with the matching name  
`GET "/byDate/date"`- Returns drinks that are available on the given date   
`POST /drinks/` - Adds a new drink to the database   
`DELETE /drinks/name`- Deletes the drink with the maching name from the db  


### Wishlist  
- Validation for create inputs
