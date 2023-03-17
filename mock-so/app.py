"""SO mock application."""
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
import json

from flask import Flask, request
from flask_restful import Api

from resources.orchestration_request import OrchestrationRequest
from resources.service_instance import (
    NetworkInstance,
    NetworkInstanceList,
    ServiceInstance,
    ServiceInstanceList,
    VnfInstance,
    VnfInstanceList,
    VfModuleInstance,
    VfModuleInstanceList,
)


app = Flask(__name__)
api = Api(app)


@app.route("/reset")
def reset():
    """Reset endpoint.

    Reset all resources.

    Returns:
        str: Empty string, it has to returns anything

    """
    ServiceInstanceList.reset()
    OrchestrationRequest.reset()
    return ""


@app.route("/set_aai_mock", methods=["POST"])
def set_aai_mock():
    """Set A&AI mock url address.

    Set it for all resources which connects with A&AI mock.

    Returns:
        str: Empty string, it has to returns anything

    """
    aai_mock_url = json.loads(request.data)["AAI_MOCK"]
    ServiceInstance.set_aai_mock(aai_mock_url)
    ServiceInstanceList.set_aai_mock(aai_mock_url)
    VnfInstance.set_aai_mock(aai_mock_url)
    VfModuleInstance.set_aai_mock(aai_mock_url)
    VnfInstanceList.set_aai_mock(aai_mock_url)
    VfModuleInstanceList.set_aai_mock(aai_mock_url)
    NetworkInstance.set_aai_mock(aai_mock_url)
    NetworkInstanceList.set_aai_mock(aai_mock_url)
    return ""


api.add_resource(
    ServiceInstance,
    "/onap/so/infra/serviceInstantiation/v7/serviceInstances/<service_instance_id>",
)
api.add_resource(
    ServiceInstanceList, "/onap/so/infra/serviceInstantiation/v7/serviceInstances"
)
api.add_resource(
    VnfInstance,
    (
        "/onap/so/infra/serviceInstantiation/v7/serviceInstances/<service_instance_id>/"
        "vnfs/<vnf_instance_id>"
    ),
)
api.add_resource(
    VnfInstanceList,
    "/onap/so/infra/serviceInstantiation/v7/serviceInstances/<service_instance_id>/vnfs",
)
api.add_resource(
    VfModuleInstance,
    (
        "/onap/so/infra/serviceInstantiation/v7/serviceInstances/<service_instance_id>/vnfs/"
        "<vnf_instance_id>/vfModules/<vf_module_instance_id>"
    ),
)
api.add_resource(
    VfModuleInstanceList,
    (
        "/onap/so/infra/serviceInstantiation/v7/serviceInstances/<service_instance_id>/vnfs/"
        "<vnf_instance_id>/vfModules"
    ),
)
api.add_resource(
    OrchestrationRequest,
    "/onap/so/infra/orchestrationRequests/v7/<orchestration_request_id>",
)
api.add_resource(
    NetworkInstanceList,
    "/onap/so/infra/serviceInstantiation/v7/serviceInstances/<service_instance_id>/networks",
)
api.add_resource(
    NetworkInstance,
    (
        "/onap/so/infra/serviceInstantiation/v7/serviceInstances/<service_instance_id>/"
        "networks/<network_instance_id>"
    ),
)


if __name__ == "__main__":
    app.run(host="0.0.0.0", debug=True, port=5001)
