"""A&AI Customer mock module."""
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

from flask_restful import reqparse, Resource, request


CUSTOMERS = {}


def random_resource_version() -> str:
    """Generate random resource version string value.

    Returns:
        str: UUID value

    """
    return str(uuid4())


class Customer(Resource):
    """Customer resource."""

    @staticmethod
    def reset():
        """Reset Cloud region resource.

        Clean CUSTOMERS dictionary

        """
        global CUSTOMERS
        CUSTOMERS = {}

    def get(self, global_customer_id: str):
        """Get customer.

        Get customer from CUSTOMERS dictionary.

        Args:
            global_customer_id (str): global customer id key value

        Returns:
            Dict[str, str]: Customer dictionary

        """
        return CUSTOMERS[global_customer_id]

    def put(self, global_customer_id: str):
        """Resource put method.

        Add customer data into CUSTOMERS dictionary.

        Args:
            global_customer_id (str): global customer id key value

        Returns:
            Tuple[str, int]: Response tuple. First element is a response body,
                the second one is HTTP response code.

        """
        CUSTOMERS[global_customer_id] = request.get_json()
        CUSTOMERS[global_customer_id]["resource-version"] = random_resource_version()
        return "", 201


class CustomerList(Resource):
    """List of customers resource."""

    def get(self) -> Dict[str, List]:
        """Get the list of customers.

        Return data from CUSTOMERS dictionary.

        Returns:
            Dict[str, List]: Customer dictionary

        """
        return {"customer": [data for global_customer_id, data in CUSTOMERS.items()]}


class ServiceSubscription(Resource):
    """Service subscription resource."""

    def put(self, global_customer_id: str, service_type: str) -> None:
        """Service subscription resource put method.

        Add service subscription data into CUSTOMERS dictionary.

        Args:
            global_customer_id (str): Global customer id key value
            service_type (str): Service type key value

        """
        try:
            CUSTOMERS[global_customer_id]["service_subscriptions"][
                service_type
            ] = {"service-type": service_type}
        except KeyError:
            CUSTOMERS[global_customer_id]["service_subscriptions"] = {
                service_type: {"service-type": service_type}
            }
        CUSTOMERS[global_customer_id]["service_subscriptions"][service_type].update(
            {"service-type": service_type}
        )


class ServiceSubscriptionList(Resource):
    """List of service subscriptions resource."""

    def get(self, global_customer_id: str) -> Dict[str, List]:
        """Get the list of service subscriptions.

        Return data from CUSTOMERS dictionary.

        Args:
            global_customer_id (str): Global customer id key value

        Returns:
            Dict[str, List]: Service subscriptions dictionary

        """
        service_type: str = request.args.get("service-type")
        if not service_type:
            return {
                "service-subscription": [
                    data
                    for service_type, data in CUSTOMERS[global_customer_id][
                        "service_subscriptions"
                    ].items()
                ]
            }
        try:
            return {
                "service-subscription": [
                    CUSTOMERS[global_customer_id]["service_subscriptions"][service_type]
                ]
            }
        except KeyError:
            return "", 404


class ServiceSubscriptionRelationship(Resource):
    """Service subscription relationship resource."""

    def put(self, global_customer_id: str, service_type: str) -> None:
        """Service subscription relationship resource put method.

        Add service subscription relationship data into CUSTOMERS dictionary.

        Args:
            global_customer_id (str): Global customer id key value
            service_type (str): Service type key value

        """
        try:
            CUSTOMERS[global_customer_id]["service_subscriptions"][service_type][
                "relationships"
            ].append(request.get_json())
        except KeyError:
            CUSTOMERS[global_customer_id]["service_subscriptions"][service_type][
                "relationships"
            ] = [request.get_json()]


class ServiceSubscriptionRelationshipList(Resource):
    """Service subscription relationships list resource."""

    def get(self, global_customer_id: str, service_type: str) -> Dict[str, List]:
        """Get the list of service subscription relationships.

        Return data from CUSTOMERS dictionary.

        Args:
            global_customer_id (str): Global customer id key value
            service_type (str): Service type key value

        Returns:
            Dict[str, List]: Service subscription relationships dictionary

        """
        return {
            "relationship": CUSTOMERS[global_customer_id]["service_subscriptions"][
                service_type
            ].get("relationships", [])
        }


class ServiceSubscriptionInstance(Resource):
    """Service subscription instance resource."""

    def delete(self, global_customer_id: str, service_type: str, service_instance_id: str) -> None:
        """Delete service subscription instance.

        Removes data from CUSTOMERS dictionary.

        Args:
            global_customer_id (str): Global customer id key value
            service_type (str): Service type key value
            service_instance_id (str): Service instance id key value

        """
        del CUSTOMERS[global_customer_id]["service_subscriptions"][service_type][
            "service_instances"
        ][service_instance_id]


class ServiceSubscriptionInstanceList(Resource):
    """Service subscription instances list resource."""

    def get(self, global_customer_id: str, service_type: str) -> Dict[str, List]:
        """Get service subscription's service instances.

        Returns data from CUSTOMERS dictionary

        Args:
            global_customer_id (str): Global customer id key value
            service_type (str): Service type key value

        Returns:
            Dict[str, List]: Service instances dictionary
        """
        return {
            "service-instance": [
                data
                for instance_id, data in CUSTOMERS[global_customer_id]["service_subscriptions"][
                    service_type
                ]
                .get("service_instances", dict())
                .items()
            ]
        }

    def post(self, global_customer_id: str, service_type: str) -> None:
        """Add service instance to service subscription.

        Add service instance data dictionary to service subscription's
            service instances dictionary.

        Args:
            global_customer_id (str): Global customer id key value
            service_type (str): Service type key value

        """
        request_data = request.get_json()
        instance_id = request_data["service-instance-id"]
        try:
            CUSTOMERS[global_customer_id]["service_subscriptions"][service_type][
                "service_instances"
            ][instance_id] = request_data
        except KeyError:
            CUSTOMERS[global_customer_id]["service_subscriptions"][service_type][
                "service_instances"
            ] = {instance_id: request_data}


class ServiceSubscriptionInstanceRelationshipList(Resource):
    """Service subscription instance relationships list resource."""

    def post(self, global_customer_id: str, service_type: str, service_instance_id) -> None:
        """Add relationship into service instance relationships list.

        Args:
            global_customer_id (str): Global customer id key value
            service_type (str): Service type key value
            service_instance_id (str): Service instance id key value

        """
        try:
            CUSTOMERS[global_customer_id]["service_subscriptions"][service_type][
                "service_instances"
            ][service_instance_id]["relationships"].append(request.get_json())
        except KeyError:
            CUSTOMERS[global_customer_id]["service_subscriptions"][service_type][
                "service_instances"
            ][service_instance_id]["relationships"] = [request.get_json()]

    def get(
        self, global_customer_id: str, service_type: str, service_instance_id: str
    ) -> Dict[str, List]:
        """Get the service instance relationships list.

        Args:
            global_customer_id (str): Global customer id key value
            service_type (str): Service type key value
            service_instance_id (str): Service instance id key value

        Returns:
            Dict[str, List]: Service instance relationships dictionary

        """
        return {
            "relationship": CUSTOMERS[global_customer_id]["service_subscriptions"][service_type][
                "service_instances"
            ][service_instance_id].get("relationships", [])
        }

    def delete(self, global_customer_id: str, service_type: str, service_instance_id: str) -> None:
        """Delete service subscription instance relationships.

        Make relationships list clea.

        Args:
            global_customer_id (str): Global customer id key value
            service_type (str): Service type key value
            service_instance_id (str): Service instance id key value

        """
        CUSTOMERS[global_customer_id]["service_subscriptions"][service_type]["service_instances"][
            service_instance_id
        ]["relationships"] = []
