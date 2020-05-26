import os
import json
import logging

from flask import Flask, request, Response
from pathlib import Path

# Configuration
PORT = 8082
HOST =  "0.0.0.0"
LOG_LEVEL = os.getenv("LOG_LEVEL", "DEBUG")

# Globals
flask_app       = Flask(__name__)
script_path     = Path(os.getcwd())
inventory_path  = os.path.join(script_path, 'data', "inventory.json")

# Load an inventory file from json.  Filename is given, load_inventory fetches from data directory.
def load_inventory():
    with open(inventory_path) as f:
        inventorydata = json.load(f)

    # return a new dictionary, keyed by registration number
    return { item["registration"]: item for item in inventorydata }


@flask_app.route('/all')
def return_all_inventory():
    result = []
    for registration, aircraft in inventory.items():
        result.append(aircraft)

    # Generate URL back to each aircraft's image and add it to each
    # aircraft dictionary item.  Host_URL has a trailing slash.
    host = request.host_url
    for aircraft in result:
        aircraft["imageUrl"] = f"{request.host_url}image/{aircraft['registration']}"

    return Response(json.dumps(result),  mimetype='application/json')

@flask_app.route('/one/<registration>')
def return_specific_item(registration):
    if registration in inventory.keys():
        aircraft = inventory[registration]
        aircraft["imageUrl"] = f"{request.host_url}image/{aircraft['registration']}"

        return Response(json.dumps(aircraft),  mimetype='application/json')
    else:
        return "Aircraft not in inventory", 404

if __name__ == '__main__':
    # Set up logging and run the server.
    logging.basicConfig(level=logging.INFO)

    # Load the inventory.  This is scoped to all app routes and available in handlers.
    inventory = load_inventory()

    logging.info(f"\nStarting pyservice on {HOST}:{PORT}\n")
    flask_app.run(host=HOST, port=PORT)

