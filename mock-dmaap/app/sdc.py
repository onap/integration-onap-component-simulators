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

class Sdc:
    __jsonArtifact = {}

    def __init__(self, logger, events):
        self.__logger = logger
        self.__events = events

    def post_msg_to_topic(self, topic, receive_msg):
        self.__logger.info("received events: " + str(receive_msg))
        return self.__jsonArtifact

    def add_msg_to_topic(self, topic, receive_msg):
        self.__jsonArtifact = receive_msg
        self.__events.add(topic, receive_msg)
        self.__logger.info("received events: " + str(receive_msg))
        return receive_msg

    def is_event_from_sdc(self, event):
        if event[:3] == '14.':
            return True
        else:
            return False
