"""SO mock instance resources."""
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
from typing import Callable, Dict, List
from uuid import uuid4

import requests
from flask_restful import Resource, request


SERVICE_INSTANCES = {}


class SoResource(Resource):
    """Base SO resource class."""

    AAI_MOCK_URL = ""

    @classmethod
    def set_aai_mock(cls, aai_mock_url: str) -> None:
        """Set A&AI mock address.

        Instance resources needs to communicate with A&AI mock
            to update it's state.

        Args:
            aai_mock_url (str): A&AI mock url

        """
        cls.AAI_MOCK_URL = aai_mock_url

    @classmethod
    def reset(cls) -> None:
        """Reset SO resource.

        Clean SERVICE_INSTANCES dictionary
            and A&AI mock url address.

        """
        global SERVICE_INSTANCES
        SERVICE_INSTANCES = {}
        cls.AAI_MOCK_URL = ""

    def check_aai_mock_address(method: Callable) -> Callable:
        """Decorate method to check if A&AI address is set.

        If A&AI mock address is not set it returns 500 HTTP
            response for resource's method

        Args:
            method (Callable): method to decorate

        """

        def decorator(self, *args, **kwargs):
            if not self.AAI_MOCK_URL:
                return "A&AI mock address not set", 500
            return method(self, *args, **kwargs)

        return decorator


class ServiceInstance(SoResource):
    """Service instance resource class."""

    def delete(self, service_instance_id: str) -> Dict[str, Dict[str, str]]:
        """Delete service instance.

        Remove service instance data from SERVICE_INSTANCES
            dictionary.

        Args:
            service_instance_id (str): Service instance id key value

        Returns:
            Dict[str, Dict[str, str]]: Deletion request dictionary

        """
        service_instance = SERVICE_INSTANCES[service_instance_id]
        requests.delete(
            (
                f"{self.AAI_MOCK_URL}/aai/v16/business/customers/customer/"
                f"{service_instance['customerId']}/service-subscriptions/service-subscription/"
                f"{service_instance['serviceSubscriptionId']}/service-instances/service-instance/"
                f"{service_instance_id}"
            )
        )
        del SERVICE_INSTANCES[service_instance_id]
        return {"requestReferences": {"requestId": str(uuid4())}}


class ServiceInstanceList(SoResource):
    """Service instances list resource."""

    @SoResource.check_aai_mock_address
    def post(self) -> Dict[str, Dict[str, str]]:
        """Create service instance.

        Create service instance data dictionary.
            Call request to A&AI mock to create service instance there.

        Returns:
            Dict[str, Dict[str, str]]: Creation request dictionary

        """
        instance_id = str(uuid4())
        request_data = request.get_json()
        customer_id = request_data["requestDetails"]["subscriberInfo"][
            "globalSubscriberId"
        ]
        service_subscription_id = request_data["requestDetails"]["requestParameters"][
            "subscriptionServiceType"
        ]
        instance_name = request_data["requestDetails"]["requestInfo"]["instanceName"]
        service_instance = {
            "requestId": str(uuid4()),
            "instanceId": instance_id,
            "customerId": customer_id,
            "serviceSubscriptionId": service_subscription_id,
            "instanceName": instance_name,
        }
        requests.post(
            (
                f"{self.AAI_MOCK_URL}/aai/v16/business/customers/customer/{customer_id}/"
                f"service-subscriptions/service-subscription/{service_subscription_id}/"
                "service-instances"
            ),
            json={
                "service-instance-name": instance_name,
                "service-instance-id": instance_id,
            },
        )
        SERVICE_INSTANCES[service_instance["instanceId"]] = service_instance
        return {"requestReferences": service_instance}


class VnfInstance(SoResource):
    """Vnf instance resource."""

    @SoResource.check_aai_mock_address
    def delete(
        self, service_instance_id: str, vnf_instance_id: str
    ) -> Dict[str, Dict[str, str]]:
        """Delete vnf instance.

        Remove vnf instanca data from SERVICE_INSTANCES dictionary.
            Call DELETE request to A&AI mock.

        Args:
            service_instance_id (str): Service instance id key value
            vnf_instance_id (str): Vnf instance id key value

        Returns:
            Dict[str, Dict[str, str]]: Deletion request dictionary.

        """
        related_service = SERVICE_INSTANCES[service_instance_id]
        requests.delete(
            (
                f"{self.AAI_MOCK_URL}/aai/v16/business/customers/customer/"
                f"{related_service['customerId']}/service-subscriptions/service-subscription/"
                f"{related_service['serviceSubscriptionId']}/service-instances/service-instance/"
                f"{service_instance_id}/relationship-list"
            )
        )
        return {"requestReferences": {"requestId": str(uuid4())}}


class VnfInstanceList(SoResource):
    """Vnf instances list resource."""

    @SoResource.check_aai_mock_address
    def post(self, service_instance_id: str) -> Dict[str, Dict[str, str]]:
        """Create vnf instance.

        Create vnf instance data dictionary.
            Call request to A&AI mock to create vnf instance there.

        Returns:
            Dict[str, Dict[str, str]]: Creation request dictionary

        """
        instance_id = str(uuid4())
        request_data = request.get_json()
        related_instance_id = request_data["requestDetails"]["relatedInstanceList"][0][
            "relatedInstance"
        ]["instanceId"]
        related_service = SERVICE_INSTANCES[related_instance_id]
        requests.post(
            (
                f"{self.AAI_MOCK_URL}/aai/v16/business/customers/customer/{related_service['customerId']}/"
                f"service-subscriptions/service-subscription/{related_service['serviceSubscriptionId']}/"
                f"service-instances/service-instance/{related_instance_id}/relationship-list"
            ),
            json={
                "related-to": "generic-vnf",
                "related-link": f"/aai/v16/network/generic-vnfs/generic-vnf/{instance_id}",
            },
        )
        return {
            "requestReferences": {"requestId": str(uuid4()), "instanceId": instance_id}
        }


class VfModuleInstance(SoResource):
    """Vf module instance resource class."""

    @SoResource.check_aai_mock_address
    def delete(
        self, service_instance_id: str, vnf_instance_id: str, vf_module_instance_id: str
    ) -> Dict[str, Dict[str, str]]:
        """Delete vf module instance.

        Call DELETE request to A&AI mock to delete vf module instance.

        Args:
            service_instance_id (str): Service instance id key value.
            vnf_instance_id (str): Vnf instance id key value.
            vf_module_instance_id (str): Vf module instance id key value.

        Returns:
            Dict[str, Dict[str, str]]: Deletion request dictionary

        """
        requests.delete(
            (
                f"{self.AAI_MOCK_URL}/aai/v16/network/generic-vnfs/generic-vnf/"
                f"{vnf_instance_id}/vf-modules/{vf_module_instance_id}"
            )
        )
        return {"requestReferences": {"requestId": str(uuid4())}}


class VfModuleInstanceList(SoResource):
    """Vf module instances list resource."""

    @SoResource.check_aai_mock_address
    def post(
        self, service_instance_id: str, vnf_instance_id: str
    ) -> Dict[str, Dict[str, str]]:
        """Create vf module instance.

        Call POST request to A&AI mock to create vf module instance.

        Args:
            service_instance_id (str): Service instance id key value
            vnf_instance_id (str): Vnf instance id key value

        Returns:
            Dict[str, Dict[str, str]]: Creation request dictionary

        """
        instance_id = str(uuid4())
        requests.post(
            (
                f"{self.AAI_MOCK_URL}/aai/v16/network/generic-vnfs/generic-vnf/"
                f"{vnf_instance_id}/vf-modules"
            ),
            json={"vf-module-id": instance_id},
        )
        return {
            "requestReferences": {"requestId": str(uuid4()), "instanceId": instance_id}
        }


class NetworkInstance(SoResource):
    """Network instance resource."""

    @SoResource.check_aai_mock_address
    def delete(
        self, service_instance_id: str, network_instance_id: str
    ) -> Dict[str, Dict[str, str]]:
        """Delete network instance.

        Remove network instanca data from SERVICE_INSTANCES dictionary.
            Call DELETE request to A&AI mock.

        Args:
            service_instance_id (str): Service instance id key value
            network_instance_id (str): Network instance id key value

        Returns:
            Dict[str, Dict[str, str]]: Deletion request dictionary.

        """
        related_service = SERVICE_INSTANCES[service_instance_id]
        requests.delete(
            (
                f"{self.AAI_MOCK_URL}/aai/v16/business/customers/customer/"
                f"{related_service['customerId']}/service-subscriptions/service-subscription/"
                f"{related_service['serviceSubscriptionId']}/service-instances/service-instance/"
                f"{service_instance_id}/relationship-list"
            )
        )
        return {"requestReferences": {"requestId": str(uuid4())}}


class NetworkInstanceList(SoResource):
    """Network instances list resource."""

    @SoResource.check_aai_mock_address
    def post(self, service_instance_id: str) -> Dict[str, Dict[str, str]]:
        """Create network instance.

        Create network instance data dictionary.
            Call request to A&AI mock to create network instance there.

        Returns:
            Dict[str, Dict[str, str]]: Creation request dictionary

        """
        instance_id = str(uuid4())
        request_data = request.get_json()
        related_instance_id = request_data["requestDetails"]["relatedInstanceList"][0][
            "relatedInstance"
        ]["instanceId"]
        related_service = SERVICE_INSTANCES[related_instance_id]
        requests.post(
            (
                f"{self.AAI_MOCK_URL}/aai/v16/business/customers/customer/{related_service['customerId']}/"
                f"service-subscriptions/service-subscription/{related_service['serviceSubscriptionId']}/"
                f"service-instances/service-instance/{related_instance_id}/relationship-list"
            ),
            json={
                "related-to": "l3-network",
                "related-link": f"/aai/v16/network/l3-networks/l3-network/{instance_id}",
            },
        )
        return {
            "requestReferences": {"requestId": str(uuid4()), "instanceId": instance_id}
        }
