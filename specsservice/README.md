# Specifications Server

This is a Flask-based service that returns a JSON dictionary providing
the specifications for an aircraft model in inventory.

`specsservice/Cessna Skymaster`

will return the JSON data:

```
{ "model": "Cessna Skymaster",
    "type": "twin",
    "hp": 420,
    "seats": 6,
    "speed": 180,
    "range": 840,
    "load": 1700 }
```

For this simple demo application, the aircraft specifications are simply stored in a JSON file
in `specsservice/data/inventory.json`.  This is read in at startup and saved in the `specs`
variable in the main program in `specsservice/service.py`.

To test out the Service Preview functionality, you will need to stop and restart the program if
you change, add or remove any of the aircraft specifications.