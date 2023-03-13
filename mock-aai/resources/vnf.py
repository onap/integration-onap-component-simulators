"""Vnf resources module."""
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
from typing import Dict, List
from uuid import uuid4

from flask_restful import Resource, request


VNFS = {}


class Vnf(Resource):
    """Vnf resource."""

    @staticmethod
    def reset():
        """Reset resource for tests.

        Create new, empty VNFS dictionary

        """
        global VNFS
        VNFS = {}

    def get(self, vnf_instance_id: str) -> Dict[str, List]:
        """Get vnf instance data.

        Get data from VNFS dictionary

        Args:
            vnf_instance_id (str): Vnf instance id key value

        Returns:
            Dict[str, List]: Vnf instance data dictionary

        """
        try:
            return VNFS[vnf_instance_id]
        except KeyError:
            VNFS[vnf_instance_id] = {"vnf-id": vnf_instance_id, "model-version-id": str(uuid4())}
            return VNFS[vnf_instance_id]


class VfModule(Resource):
    """Vf module resource."""

    def delete(self, vnf_instance_id: str, vf_module_instance_id: str) -> None:
        """Delete vf module.

        Removes vf module data from VNFS dictionary.

        Args:
            vnf_instance_id (str): Vnf instance id key value
            vf_module_instance_id (str): Vf module instance id key value

        """
        del VNFS[vnf_instance_id]["vf_modules"][vf_module_instance_id]


class VfModuleList(Resource):
    """Vf module list resource."""

    def post(self, vnf_instance_id: str) -> None:
        """Create vf module.

        Add vf module data into VNFS dictionary.

        Args:
            vnf_instance_id (str): Vnf instance id key value

        """
        vf_module_data = request.get_json()
        vf_module_dict = {vf_module_data["vf-module-id"]: vf_module_data}
        try:
            VNFS[vnf_instance_id]["vf_modules"].update(vf_module_dict)
        except KeyError:
            VNFS[vnf_instance_id]["vf_modules"] = vf_module_dict

    def get(self, vnf_instance_id: str) -> Dict[str, List]:
        """Get Vnf instance Vf modules list.

        Get data from VNFS dictionary

        Args:
            vnf_instance_id (str): Vnf instance id key value

        Returns:
            Dict[str, List]: Vnf instance vf modules dictionary
        """
        return {
            "vf-module": [
                data for vf_module_id, data in VNFS[vnf_instance_id].get("vf_modules", {}).items()
            ]
        }
