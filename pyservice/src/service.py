import os
import logging
import requests

from flask import Flask

PORT = 5000
HOST =  "0.0.0.0"
LOG_LEVEL = os.getenv("LOG_LEVEL", "DEBUG")

app = Flask(__name__)

@app.route('/')
def hello_world():
    logging.info("handled /")
    return 'Hello Python World!'

if __name__ == '__main__':
    logging.basicConfig(level=logging.INFO)
    logging.info(f"\nStarting pyservice on {HOST}:{PORT}\n")
    app.run(host=HOST, port=PORT)

