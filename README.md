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

- Runs on `https://localhost:8081`  
- Connects to MongoDB running on default port `27017`
- Adding a `dummy.js` file which can be used to seed the database if needed
  using `mongo dummy.js`

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
- [x] Docker!  
- [x] Filtering using the ingredients array  
- [ ] Pagination
