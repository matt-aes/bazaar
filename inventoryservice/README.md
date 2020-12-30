# Inventory Server

This is a Flask-based service that returns a JSON dictionary providing
a list of all aircraft in the inventory.  

`inventoryservice/all`

returns all aircraft in the inventory as an array of JSON structures.

`inventoryservice/one/<registration>`

returns the `json` for the aircraft with that registration.

**Not yet implemented:** a query language allows selecting specific aircraft by single/twin,
number of seats, or a price range.

`inventoryservice/max_price?40000`

will return the JSON data:

```
[
 {  "registration": "N8204T",
    "model":        "Cessna 180",
    "price":        38000
  },
  { "registration": "CF-NLX",
    "model":        "Alon Ercoupe",
    "price":        14000
  },
  { "registration": "CF-LUN",
    "model":        "Alon Ercoupe",
    "price":        18000
  }
]
```

For this simple demo application, the aircraft inventory is simply stored in a `json` file
in `inventoryservice/data/inventory.json`.  This is read in at startup and saved in the `inventory`
variable in the main program in `inventoryservice/service.py`.

To test out the Service Preview functionality, you will need to stop and restart the program if
you change, add or remove any of the aircraft in the inventory file.  If you add a new aircraft
that does not have a corresponding specification already defined in the `specsservice`, you will
need to add an entry in the `specsservice/data/specs.json` file.  The model name in the
specification `json` entry must match exactly the model in the inventory entry.