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


NETWORKS = {}


class Network(Resource):
    """Network resource."""

    @staticmethod
    def reset():
        """Reset resource for tests.

        Create new, empty NETWORKS dictionary

        """
        global NETWORKS
        NETWORKS = {}

    def get(self, network_instance_id: str) -> Dict[str, List]:
        """Get network instance data.

        Get data from NETWORKS dictionary

        Args:
            network_instance_id (str): Network instance id key value

        Returns:
            Dict[str, List]: Network instance data dictionary

        """
        try:
            return NETWORKS[network_instance_id]
        except KeyError:
            NETWORKS[network_instance_id] = {
                "network-id": network_instance_id,
                "is-bound-to-vpn": False,
                "is-provider-network": False,
                "is-shared-network": False,
                "is-external-network": False,
            }
            return NETWORKS[network_instance_id]
