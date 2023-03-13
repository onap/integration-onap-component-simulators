"""A&AI mock application."""
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
from flask_restful import Resource, Api

from resources.cloud_region import (
    CloudRegion,
    CloudRegionList,
    CloudRegionRelationshipList,
    CloudRegionRelationship,
    Tenant,
    TenantList,
)
from resources.complex import Complex, ComplexList
from resources.customer import (
    Customer,
    CustomerList,
    ServiceSubscription,
    ServiceSubscriptionInstance,
    ServiceSubscriptionList,
    ServiceSubscriptionRelationship,
    ServiceSubscriptionRelationshipList,
    ServiceSubscriptionInstanceList,
    ServiceSubscriptionInstanceRelationshipList,
)
from resources.network import Network
from resources.vnf import Vnf, VfModule, VfModuleList

app = Flask(__name__)
api = Api(app)

API_VERSIONS_SUPPORTED = ["v16", "v20", "v23", "v26"]

def generate_urls_for_all_versions(endpoint_pattern: str) -> list:
    """Helper function.

    Generates list of endpoints by replacing <version> pattern for each supported API version.

    Args:
        endpoint_pattern (str): url pattern with <version> tag to be replaced

    Returns:
        list: List of string values reporesenting endpoints that are currently supported

    """
    return [endpoint_pattern.replace("<version>", api_version) for api_version in API_VERSIONS_SUPPORTED]

@app.route("/reset")
def reset() -> str:
    """Reset endpoint.

    Reset all resources.

    Returns:
        str: Empty string, it has to returns anything

    """
    CloudRegion.reset()
    Complex.reset()
    Customer.reset()
    Network.reset()
    Vnf.reset()
    return ""


# Cloud region resource
api.add_resource(CloudRegionList, *generate_urls_for_all_versions("/aai/<version>/cloud-infrastructure/cloud-regions"))
api.add_resource(
    CloudRegion,
    *generate_urls_for_all_versions("/aai/<version>/cloud-infrastructure/cloud-regions/cloud-region/<cloud_owner>/<cloud_region_id>"),
)
api.add_resource(
    CloudRegionRelationshipList,
    *generate_urls_for_all_versions("/aai/<version>/cloud-infrastructure/cloud-regions/cloud-region/<cloud_owner>/<cloud_region_id>/relationship-list"),
)
api.add_resource(
    CloudRegionRelationship,
    *generate_urls_for_all_versions("/aai/<version>/cloud-infrastructure/cloud-regions/cloud-region/<cloud_owner>/<cloud_region_id>/relationship-list/relationship"),
)
api.add_resource(
    TenantList,
    *generate_urls_for_all_versions("/aai/<version>/cloud-infrastructure/cloud-regions/cloud-region/<cloud_owner>/<cloud_region_id>/tenants"),
)
api.add_resource(
    Tenant,
    *generate_urls_for_all_versions("/aai/<version>/cloud-infrastructure/cloud-regions/cloud-region/<cloud_owner>/<cloud_region_id>/tenants/tenant/<tenant_id>"),
)
# Complex resource
api.add_resource(ComplexList, *generate_urls_for_all_versions("/aai/<version>/cloud-infrastructure/complexes"))
api.add_resource(
    Complex, *generate_urls_for_all_versions("/aai/<version>/cloud-infrastructure/complexes/complex/<physical_location_id>")
)
# Customer resource
api.add_resource(CustomerList, *generate_urls_for_all_versions("/aai/<version>/business/customers"))
api.add_resource(Customer, *generate_urls_for_all_versions("/aai/<version>/business/customers/customer/<global_customer_id>"))
api.add_resource(
    ServiceSubscriptionList,
    *generate_urls_for_all_versions("/aai/<version>/business/customers/customer/<global_customer_id>/service-subscriptions"),
)
api.add_resource(
    ServiceSubscription,
    *generate_urls_for_all_versions("/aai/<version>/business/customers/customer/<global_customer_id>/service-subscriptions/service-subscription/<service_type>"),
)
api.add_resource(
    ServiceSubscriptionRelationshipList,
    *generate_urls_for_all_versions("/aai/<version>/business/customers/customer/<global_customer_id>/service-subscriptions/service-subscription/<service_type>/relationship-list"),
)
api.add_resource(
    ServiceSubscriptionRelationship,
    *generate_urls_for_all_versions("/aai/<version>/business/customers/customer/<global_customer_id>/service-subscriptions/service-subscription/<service_type>/relationship-list/relationship"),
)
api.add_resource(
    ServiceSubscriptionInstance,
    *generate_urls_for_all_versions("/aai/<version>/business/customers/customer/<global_customer_id>/service-subscriptions/service-subscription/<service_type>/service-instances/service-instance/<service_instance_id>"),
)
api.add_resource(
    ServiceSubscriptionInstanceList,
    *generate_urls_for_all_versions("/aai/<version>/business/customers/customer/<global_customer_id>/service-subscriptions/service-subscription/<service_type>/service-instances"),
)
api.add_resource(
    ServiceSubscriptionInstanceRelationshipList,
    *generate_urls_for_all_versions("/aai/<version>/business/customers/customer/<global_customer_id>/service-subscriptions/service-subscription/<service_type>/service-instances/service-instance/<service_instance_id>/relationship-list"),
)
# VNF resource
api.add_resource(Vnf, *generate_urls_for_all_versions("/aai/<version>/network/generic-vnfs/generic-vnf/<vnf_instance_id>"))
api.add_resource(
    VfModule,
    *generate_urls_for_all_versions("/aai/<version>/network/generic-vnfs/generic-vnf/<vnf_instance_id>/vf-modules/<vf_module_instance_id>"),
)
api.add_resource(
    VfModuleList,
    *generate_urls_for_all_versions("/aai/<version>/network/generic-vnfs/generic-vnf/<vnf_instance_id>/vf-modules"),
)
api.add_resource(
    Network, *generate_urls_for_all_versions("/aai/<version>/network/l3-networks/l3-network/<network_instance_id>")
)


if __name__ == "__main__":
    app.run(host="0.0.0.0", debug=True)
