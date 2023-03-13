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


class Ves:

    def __init__(self, logger, events):
        self.__logger = logger
        self.__events = events

    def handle_new_event(self, topic, http_request):
        receive_events = self.__decode_request_data(http_request)
        for event in receive_events:
            self.__events.add(topic, json.loads(event))
        return {}, 200

    def __decode_request_data(self, data):
        receive_events = data.split("\n")
        receive_events = receive_events[:-1]
        self.__logger.info("received events: " + str(receive_events))
        correct_events = []
        for event in receive_events:
            self.__logger.info("received event: " + str(event))
            correct_events.append(self.__get_correct_json(event))
        return correct_events

    def __get_correct_json(self, incorrect_json):
        json_start_position = incorrect_json.find("{")
        correct_json = incorrect_json[json_start_position:]
        correct_json = correct_json.replace("\r", "").replace("\t", "").replace(" ", "")
        return correct_json
