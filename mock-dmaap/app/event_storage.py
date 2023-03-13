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

class EventStorage:
    __events = {}

    def get_events(self):
        return self.__events

    def get_events_from_topic(self, topic):
        if self.__events.__contains__(topic):
            return self.__events[topic]
        else:
            return []

    def add(self, topic, event):
        if self.__events.__contains__(topic):
            self.__events[topic].append(event)
        else:
            self.__events[topic] = [event]

    def delete_events(self, topic):
        if self.__events.__contains__(topic):
            self.__events[topic].clear()

    def clear(self):
        self.__events.clear()
