# Inventory Server

This is a Flask-based service that returns a JSON dictionary providing
a list of all aircraft in the inventory.  A query language allows selecting
specific aircraft by single/twin, number of seats, or a price range.

inventoryservice/all
    => Returns all aircraft in the inventory
    
inventoryservice/max_price?40000

will return the JSON data:

[
 {  "registration": "N8204T",
    "model": "Cessna 180",
    "price": 38000
  },
  { "registration": "CF-NLX",
    "model": "Alon Ercoupe",
    "price": 14000
  },
  { "registration": "CF-LUN",
    "model": "Alon Ercoupe",
    "price": 18000
  }
]