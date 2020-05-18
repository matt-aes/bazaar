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
    return '/ handled: Hello Python World!'


@app.route('/hello')
def hello_world_2():
    logging.info("handled /hello")
    return '/hello handled: Hello Python World!'


@app.route('/fwd')
def forward_to_go_service():
    logging.info("handled /fwd")

    # Now make a request of the go service and return its value.
    r = requests.get("http://goservice/hello", verify=False)
    return f"/goservice/hello returned {r.status_code} => {r.text}"

if __name__ == '__main__':
    logging.basicConfig(level=logging.INFO)
    logging.info(f"\nStarting pyservice on {HOST}:{PORT}\n")
    app.run(host=HOST, port=PORT)

