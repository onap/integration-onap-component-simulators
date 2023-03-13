"""
   Copyright 2023 Deutsche Telekom AG, Nokia, Orange

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
"""
import json
import requests
import logging as sys_logging
from flask import Flask, logging, request

FAULT_TOPIC = "fault"
dmaap_url = ""

app = Flask(__name__)
sys_logging.basicConfig(level=sys_logging.DEBUG)
logger = logging.create_logger(app)


@app.route("/set_dmaap_address", methods=['POST'])
def set_dmaap_address():
    logger.debug(request.json)
    global dmaap_url
    dmaap_url = get_dmaap_mock_url()
    return {}, 202


@app.route("/eventListener/<version>", methods=['POST'])
def event_sec_fault_output(version):
    logger.debug(request.json)
    event = json.dumps(request.json["event"]) \
        .replace('\n', ' ') \
        .__add__("\n")
    send_event_to_dmaap(dmaap_url, event, FAULT_TOPIC)
    return handle_new_event(version)


@app.route("/eventListener/<version>/eventBatch", methods=['POST'])
def event_sec_fault_output_batch(version):
    logger.debug(request.json)
    dmaap_mock_url = dmaap_url
    event = prepare_event_list_for_dmaap()
    send_event_to_dmaap(dmaap_mock_url, event, FAULT_TOPIC)
    return handle_new_event(version)


def send_event_to_dmaap(dmaap_mock_url, event, topic):
    byte_event = change_from_str_to_byte_array(event)
    requests.post("{}/events/{}".format(dmaap_mock_url, topic), data=byte_event)


def handle_new_event(version):
    return {}, 202


def change_from_str_to_byte_array(event):
    b = bytearray()
    b.extend(event.encode())
    return b


def prepare_event_list_for_dmaap():
    event_list = []
    for event in request.json["eventList"]:
        event_list.append(json.dumps(event).replace('\n', ' '))
    event = "\n".join(event_list).__add__("\n")
    return event


def get_dmaap_mock_url():
    return request.json["DMAAP_MOCK"]


if __name__ == "__main__":
    app.run(host='0.0.0.0', port=30417, debug=True)
