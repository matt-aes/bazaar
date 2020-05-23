import os
import json
import logging

from flask import Flask, Response
from pathlib import Path

# Configuration
PORT = 8080
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


@flask_app.route('/')
def return_all_inventory():
    if True:
        result = []
        for registration, item in inventory.items():
            result.append(item)

        return Response(json.dumps(result),  mimetype='application/json')
    else:
        return "Inventory not available", 404

if __name__ == '__main__':
    # Set up logging and run the server.
    logging.basicConfig(level=logging.INFO)

    # Load the inventory.  This is scoped to all app routes and available in handlers.
    inventory = load_inventory()

    logging.info(f"\nStarting pyservice on {HOST}:{PORT}\n")
    flask_app.run(host=HOST, port=PORT)

