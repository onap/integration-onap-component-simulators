"""SDNC mock application."""
"""
   Copyright 2023 Deutsche Telekom AG, Orange

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
from flask import Flask
from flask_restful import Api

from resources.preload import Preload


app = Flask(__name__)
api = Api(app)


api.add_resource(
    Preload,
    "/restconf/operations/GENERIC-RESOURCE-API:preload-vf-module-topology-operation",
    "/restconf/operations/GENERIC-RESOURCE-API:preload-network-topology-operation",
)


if __name__ == "__main__":
    app.run(host="0.0.0.0", debug=True, port=5002)
