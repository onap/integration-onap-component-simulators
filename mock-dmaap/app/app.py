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
import logging as sys_logging
from event_storage import EventStorage
from sdc import Sdc
from ves import Ves

from flask import Flask, request, logging, Response

app = Flask(__name__)

sys_logging.basicConfig(level=sys_logging.DEBUG)
logger = logging.create_logger(app)
apiKey = {"key": "test_key", "secret": "test_secret"}
event_storage = EventStorage()
sdc = Sdc(logger, event_storage)
ves = Ves(logger, event_storage)


@app.route("/apiKeys/create", methods=['POST'])
def create_api_key():
    resp = Response(json.dumps(apiKey))
    return resp


@app.route("/apiKeys/<path:key>", methods=['DELETE'])
def delete_api_key(key):
    if key == apiKey["key"]:
        return {}, 200
    return {"no such key"}, 404


@app.route("/reset", methods=['GET'])
def reset_events():
    event_storage.clear()
    return event_storage


@app.route("/events", methods=['GET'])
def get_events():
    resp = Response(json.dumps(event_storage))
    resp.headers['Content-Type'] = 'application/json'
    return resp


@app.route("/events/<path:topic>/<path:consumer_group>/<path:consumer_id>", methods=['GET'])
def get_events_from_topic(topic, consumer_group, consumer_id):
    events = event_storage.get_events_from_topic(topic)
    resp = Response(json.dumps(events))
    event_storage.delete_events(topic)
    resp.headers['Content-Type'] = 'application/json'
    return resp


@app.route("/events/<path:topic>", methods=['POST'])
def post_msg_to_topic(topic):
    receive_msg = request.data.decode("utf-8")
    if sdc.is_event_from_sdc(receive_msg):
        return sdc.post_msg_to_topic(topic, receive_msg)
    else:
        return ves.handle_new_event(topic, receive_msg)


@app.route("/events/<path:topic>/add", methods=['POST'])
def add_msg_to_topic(topic):
    receive_msg = request.data.decode("utf-8")
    return sdc.add_msg_to_topic(topic, receive_msg)


@app.route("/topics", methods=['GET'])
def get_topics():
    topics = {
        'topics': ['org.onap.dmaap.mr.PNF_REGISTRATION', 'SDC-DISTR-STATUS-TOPIC-AUTO', 'SDC-DISTR-NOTIF-TOPIC-AUTO',
                   'org.onap.dmaap.mr.PNF_READY', 'POLICY-PDP-PAP', 'POLICY-NOTIFICATION',
                   'unauthenticated.SEC_3GPP_FAULTSUPERVISION_OUTPUT', '__consumer_offsets',
                   'org.onap.dmaap.mr.mirrormakeragent']}
    resp = Response(json.dumps(topics))
    resp.headers['Content-Type'] = 'application/json'
    return resp


if __name__ == "__main__":
    app.run(host='0.0.0.0', port=3904)
