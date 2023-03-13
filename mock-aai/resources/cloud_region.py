"""A&AI CloudRegion mock module."""
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
from typing import Dict, List, Tuple

from flask_restful import Resource, request

CLOUD_REGIONS = {}


class CloudRegion(Resource):
    """Cloud region resource."""

    def get(self, cloud_owner: str, cloud_region_id: str) -> Dict[str, str]:
        """Get cloud region.

        Get cloud region from CLOUD_REGIONS dictionary.

        Args:
            cloud_owner (str): cloud owner key value
            cloud_region_id (str): cloud region id key value

        Returns:
            Dict[str, str]: Cloud region dictionary

        """
        return CLOUD_REGIONS[cloud_owner][cloud_region_id]

    def put(self, cloud_owner: str, cloud_region_id: str) -> Tuple[str, int]:
        """Cloud region resource put method.

        Add cloud region data into CLOUD_REGIONS dictionary.

        Args:
            cloud_owner (str): Cloud owner key value
            cloud_region_id (str): Cloud region id key value

        Returns:
            Tuple[str, int]: Response tuple. First element is a response body,
                the second one is HTTP response code.
        """
        CLOUD_REGIONS.update({cloud_owner: {cloud_region_id: request.get_json()}})
        return "", 201

    @staticmethod
    def reset() -> None:
        """Reset Cloud region resource.

        Clean CLOUD_REGIONS dictionary

        """
        global CLOUD_REGIONS
        CLOUD_REGIONS = {}


class CloudRegionList(Resource):
    """List of cloud regions resource."""

    def get(self):
        """Get the list of cloud regions.

        Return data from CLOUD_REGIONS dictionary.

        Returns:
            Dict[str, List]: Cloud regions dictionary

        """
        return {
            "cloud-region": [
                data
                for cloud_owner, cloud_owner_dict in CLOUD_REGIONS.items()
                for cloud_owner, data in cloud_owner_dict.items()
            ]
        }


class CloudRegionRelationship(Resource):
    """Cloud region relationship resource."""

    def put(self, cloud_owner: str, cloud_region_id: str):
        """Cloud region relationship resource put method.

        Add cloud region relationship data into CLOUD_REGIONS dictionary.

        Args:
            cloud_owner (str): Cloud owner key value
            cloud_region_id (str): Cloud region id key value

        """
        try:
            CLOUD_REGIONS[cloud_owner][cloud_region_id]["relationships"].apped(request.get_json())
        except KeyError:
            CLOUD_REGIONS[cloud_owner][cloud_region_id]["relationships"] = [request.get_json()]


class CloudRegionRelationshipList(Resource):
    """List of cloud region relationships resource."""

    def get(self, cloud_owner: str, cloud_region_id: str) -> Dict[str, List]:
        """Get the list of cloud region relationships.

        Return data from CLOUD_REGIONS dictionary.

        Args:
            cloud_owner (str): Cloud owner key value
            cloud_region_id (str): Cloud region id key value

        Returns:
            Dict[str, List]: Cloud region relationships dictionary

        """
        try:
            return {"relationship": CLOUD_REGIONS[cloud_owner][cloud_region_id]["relationships"]}
        except KeyError:
            return {"relationship": []}


class Tenant(Resource):
    """Cloud region tenant resource."""

    def put(self, cloud_owner: str, cloud_region_id: str, tenant_id: str) -> None:
        """Cloud region tenant resource put method.

        Add cloud region tenant data into CLOUD_REGIONS dictionary.

        Args:
            cloud_owner (str): Cloud owner key value
            cloud_region_id (str): Cloud region id key value

        """
        try:
            CLOUD_REGIONS[cloud_owner][cloud_region_id]["tenants"].update(
                {tenant_id: request.get_json()}
            )
        except KeyError:
            CLOUD_REGIONS[cloud_owner][cloud_region_id]["tenants"] = {tenant_id: request.get_json()}

    def get(self, cloud_owner: str, cloud_region_id: str, tenant_id: str) -> Dict[str, str]:
        """Get cloud region tenant.

        Get cloud region tenant from CLOUD_REGIONS dictionary.

        Args:
            cloud_owner (str): cloud owner key value
            cloud_region_id (str): cloud region id key value

        Returns:
            Dict[str, str]: Cloud region tenant dictionary

        """
        try:
            return CLOUD_REGIONS[cloud_owner][cloud_region_id]["tenants"][tenant_id]
        except KeyError:
            return "", 404


class TenantList(Resource):
    """List of tenants resource."""

    def get(self, cloud_owner: str, cloud_region_id: str) -> Dict[str, List]:
        """Get the list of cloud region tenants.

        Return data from CLOUD_REGIONS dictionary.

        Args:
            cloud_owner (str): Cloud owner key value
            cloud_region_id (str): Cloud region id key value

        Returns:
            Dict[str, List]: Cloud region tenants dictionary

        """
        return {
            "tenant": [
                data
                for tenant_id, data in CLOUD_REGIONS[cloud_owner][cloud_region_id]
                .get("tenants", {})
                .items()
            ]
        }
