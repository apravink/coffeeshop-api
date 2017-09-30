use coffeeshop;
var bulk = db.drinks.initializeUnorderedBulkOp();
bulk.insert( {id:1, name:"Iced Coffee", price:"1.99", startDate: new Date("January 1, 2017") , endDate: new Date("February 1, 2017") , ingredients: ["coffee", "sugar", "ice"] });
bulk.insert( {id:2, name:"Tea", price:"1.50", startDate: new Date("February 1, 2017") , endDate:new Date("March 1, 2017") , ingredients: ["tea", "sugar"] });
bulk.insert( {id:3, name:"Coffee", price:"1.50", startDate: new Date("January 1, 2017") , endDate:new Date("December 1, 2017") , ingredients: ["coffee", "sugar"] });
bulk.insert( {id:4, name:"Latte", price:"2.50", startDate: new Date("March 1, 2017") , endDate:new Date("April 1, 2017") , ingredients: ["coffee","milk", "sugar"] });
bulk.execute();
