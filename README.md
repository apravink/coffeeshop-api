# coffeeshop-api
---

## Dependencies  

* "gopkg.in/mgo.v2"  

## TODO  

- [x] Create a server  
- [x] Create DB connection
- [X] Create a route to display all  
- [x] CRUD routes and functions  
- [x] Add dateTime constraint    
- [ ] Filter by ingredient  

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

### Wishlist  
- Validation for create inputs
