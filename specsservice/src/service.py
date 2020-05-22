import os
import json
import logging

from flask import Flask, render_template
from pathlib import Path

# Configuration
PORT = 5000
HOST =  "0.0.0.0"
LOG_LEVEL = os.getenv("LOG_LEVEL", "DEBUG")

# Globals
flask_app     = Flask(__name__)
script_path   = Path(os.getcwd())
specs_path    = os.path.join(script_path.parent, 'data', "specs.json")

# Load a specifications file from json.  Filename is given, load_specs fetches from data directory.
def load_specs():
    with open(specs_path) as f:
        specsdata = json.load(f)

    # return a new dictionary, keyed by model
    return { spec["model"]: spec for spec in specsdata }


# The only GET route: given an aircraft model name, return the specifications.
@flask_app.route('/spec/<model>')
def return_specification(model):
    if model in specs.keys():
        return specs[model]
    else:
        return "Aircraft model not found", 404


if __name__ == '__main__':
    # Load the specifications -- This is scoped to all app routes and available in handlers.
    specs = load_specs()

    # Set up logging and run the server.
    logging.basicConfig(level=logging.INFO)
    logging.info(f"\nStarting pyservice on {HOST}:{PORT}\n")
    flask_app.run(host=HOST, port=PORT)

