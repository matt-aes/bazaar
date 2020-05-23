# Specifications Server

This is a Flask-based service that returns a JSON dictionary providing
the specifications for an aicraft model in inventory.

specs_service/spec/Cessna Skymaster

will return the JSON data:

  { "model": "Cessna Skymaster",
    "type": "twin",
    "hp": 420,
    "seats": 6,
    "speed": 180,
    "range": 840,
    "load": 1700 }