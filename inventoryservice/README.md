# Inventory Server

This is a Flask-based service that returns a JSON dictionary providing
a list of all aircraft in the inventory.  A query language allows selecting
specific aircraft by single/twin, number of seats, or a price range.

inventory_service/inventory
    => Returns all aircraft in the inventory
    
inventory_service/inventory/max_price?40000

will return the JSON data:

[
 {  "registration": "N8204T",
    "model": "Cessna 180",
    "minimum bid": 38000
  },
  { "registration": "CF-NLX",
    "model": "Alon Ercoupe",
    "minimum bid": 14000
  },
  { "registration": "CF-LUN",
    "model": "Alon Ercoupe",
    "minimum bid": 18000
  }
]