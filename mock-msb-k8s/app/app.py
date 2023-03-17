"""MSB k8s mock application"""
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

import http
import json
import uuid

from flask import Flask, request
from flask_restful import Api

app = Flask(__name__)
api = Api(app)

CONNECTIVITY_INFOS = []
DEFINITIONS = []
CONFIGURATIONS_TEMPLATES = []
PROFILES = []
INSTANCES = []
INSTANCE_EXAMPLE = {
  "id": "ID_GENERATED_BY_K8SPLUGIN",
  "namespace": "NAMESPACE_WHERE_INSTANCE_HAS_BEEN_DEPLOYED_AS_DERIVED_FROM_PROFILE",
  "release-name": "RELEASE_NAME_AS_COMPUTED_BASED_ON_INSTANTIATION_REQUEST_AND_PROFILE_DEFAULT",
  "request": {
    "rb-name": "test-rbdef",
    "rb-version": "v1",
    "profile-name": "p1",
    "cloud-region": "krd",
    "override-values": {
        "optionalDictOfParameters": "andTheirValues, like",
        "global.name": "dummy-name"
    },
    "labels": {
        "optionalLabelForInternalK8spluginInstancesMetadata": "dummy-value"
    },
    },
  "resources": [
        {
            "GVK": {
                "Group": "",
                "Kind": "ConfigMap",
                "Version": "v1"
            },
            "Name": "test-cm"
        },
        {
            "GVK": {
                "Group": "",
                "Kind": "Service",
                "Version": "v1"
            },
            "Name": "test-svc"
        },
        {
            "GVK": {
                "Group": "apps",
                "Kind": "Deployment",
                "Version": "v1"
            },
            "Name": "test-dep"
        }
  ]
}


@app.route('/v1/connectivity-info/<string:region_id>', methods=['GET', 'DELETE'])
def connectivity_info_get_delete(region_id):
    if request.method == 'GET':
        for conninfo in CONNECTIVITY_INFOS:
            if conninfo["cloud-region"] == region_id:
                return conninfo, http.HTTPStatus.OK
        else:
            return '', http.HTTPStatus.NOT_FOUND
    if request.method == 'DELETE':
        for conninfo in CONNECTIVITY_INFOS:
            if conninfo["cloud-region"] == region_id:
                CONNECTIVITY_INFOS.remove(conninfo)
                return '', http.HTTPStatus.OK
        else:
            return '', http.HTTPStatus.NOT_FOUND
    return '', http.HTTPStatus.METHOD_NOT_ALLOWED


@app.route('/v1/connectivity-info', methods=['POST'])
def connectivity_info_create():
    if request.method == 'POST':
        kubeconfig = request.files['file']
        metadata = json.loads(request.values['metadata'])
        CONNECTIVITY_INFOS.append({
            "cloud-region": metadata['cloud-region'],
            "cloud-owner": metadata['cloud-region'],
            "kubeconfig": kubeconfig.read().decode("utf-8")
        })
        return '', http.HTTPStatus.OK
    return '', http.HTTPStatus.METHOD_NOT_ALLOWED


@app.route('/v1/rb/definition', methods=['POST'])
def definition_create():
    if request.method == 'POST':
        data = json.loads(request.data)
        DEFINITIONS.append({
            "rb-name": data['rb-name'],
            "rb-version": data['rb-version'],
            "chart-name": data['chart-name'],
            "description": data['description'],
            "labels": data['labels']
        })
        return '', http.HTTPStatus.OK
    return '', http.HTTPStatus.METHOD_NOT_ALLOWED


@app.route('/v1/rb/definition/<string:rb_name>/<string:rb_version>/content', methods=['POST'])
def definition_upload_artifact(rb_name, rb_version):
    if request.method == 'POST':
        data = request.data
        return '', http.HTTPStatus.OK
    return '', http.HTTPStatus.METHOD_NOT_ALLOWED


@app.route('/v1/rb/definition/<string:rb_name>/<string:rb_version>', methods=['GET', 'DELETE'])
def definition_get_delete(rb_name, rb_version):
    if request.method == 'GET':
        for rb in DEFINITIONS:
            if rb['rb-name'] == rb_name and rb['rb-version'] == rb_version:
                return rb, http.HTTPStatus.OK
        else:
            return '', http.HTTPStatus.NOT_FOUND
    if request.method == 'DELETE':
        for rb in DEFINITIONS:
            if rb['rb-name'] == rb_name and rb['rb-version'] == rb_version:
                DEFINITIONS.remove(rb)
                return '', http.HTTPStatus.OK
        else:
            return '', http.HTTPStatus.NOT_FOUND
    return '', http.HTTPStatus.METHOD_NOT_ALLOWED


@app.route('/v1/rb/definition', methods=['GET'])
def definition_get_all():
    if request.method == 'GET':
        return json.dumps(DEFINITIONS), http.HTTPStatus.OK
    return '', http.HTTPStatus.METHOD_NOT_ALLOWED


@app.route('/v1/rb/definition/<string:rb_name>/<string:rb_version>/profile', methods=['POST'])
def profile_create(**kwargs):
    if request.method == 'POST':
        data = json.loads(request.data)
        PROFILES.append({
            "rb-name": data['rb-name'],
            "rb-version": data['rb-version'],
            "profile-name": data['profile-name'],
            "release-name": data['release-name'],
            "namespace": data['namespace'],
            "kubernetes-version": data['kubernetes-version']
        })
        return '', http.HTTPStatus.OK
    return '', http.HTTPStatus.METHOD_NOT_ALLOWED


@app.route('/v1/rb/definition/<string:rb_name>/<string:rb_version>/profile/<string:profile_name>/content', methods=['POST'])
def profile_upload_artifact(rb_name, rb_version, profile_name):
    if request.method == 'POST':
        data = request.data
        return '', http.HTTPStatus.OK
    return '', http.HTTPStatus.METHOD_NOT_ALLOWED


@app.route('/v1/rb/definition/<string:rb_name>/<string:rb_version>/profile/<string:profile_name>',
           methods=['GET', 'DELETE'])
def profile_get_delete(rb_name, rb_version, profile_name):
    if request.method == 'GET':
        for profile in PROFILES:
            if profile['rb-name'] == rb_name and profile['rb-version'] == rb_version \
                    and profile["profile-name"] == profile_name:
                return profile, http.HTTPStatus.OK
        else:
            return '', http.HTTPStatus.NOT_FOUND
    if request.method == 'DELETE':
        for profile in PROFILES:
            if profile['rb-name'] == rb_name and profile['rb-version'] == rb_version \
                    and profile["profile-name"] == profile_name:
                PROFILES.remove(profile)
                return '', http.HTTPStatus.OK
        else:
            return '', http.HTTPStatus.NOT_FOUND
    return '', http.HTTPStatus.METHOD_NOT_ALLOWED


@app.route('/v1/rb/definition/<string:rb_name>/<string:rb_version>/profile', methods=['GET'])
def profile_get_all(rb_name, rb_version):
    if request.method == 'GET':
        profiles = []
        for profile in PROFILES:
            if profile['rb-name'] == rb_name and profile['rb-version'] == rb_version:
                profiles.append(profile)
        return json.dumps(PROFILES), http.HTTPStatus.OK
    return '', http.HTTPStatus.METHOD_NOT_ALLOWED

@app.route('/v1/instance', methods=['POST'])
def instance_create():
    if request.method == 'POST':
        data = json.loads(request.data)
        instance_details = INSTANCE_EXAMPLE
        instance_details["id"] = str(uuid.uuid4())
        instance_details["request"] = data
        INSTANCES.append(instance_details)
        return instance_details, http.HTTPStatus.OK
    return '', http.HTTPStatus.METHOD_NOT_ALLOWED

@app.route('/v1/instance/<string:instance_id>', methods=['GET', 'DELETE'])
def instance_get_delete(instance_id):
    if request.method == 'GET':
        for instance in INSTANCES:
            if instance['id'] == instance_id:
                return instance, http.HTTPStatus.OK
        else:
            return '', http.HTTPStatus.NOT_FOUND
    if request.method == 'DELETE':
        for instance in INSTANCES:
            if instance['id'] == instance_id:
                INSTANCES.remove(instance)
                return '', http.HTTPStatus.OK
        else:
            return '', http.HTTPStatus.NOT_FOUND
    return '', http.HTTPStatus.METHOD_NOT_ALLOWED


@app.route('/v1/instance', methods=['GET'])
def instance_get_all():
    if request.method == 'GET':
        return json.dumps(INSTANCES), http.HTTPStatus.OK
    return '', http.HTTPStatus.METHOD_NOT_ALLOWED


@app.route('/v1/rb/definition/<string:rb_name>/<string:rb_version>/config-template',
           methods=["POST"])
def configuration_template_create(rb_name, rb_version):
    if request.method == "POST":
        data = json.loads(request.data)
        configuration_template = data
        CONFIGURATIONS_TEMPLATES.append(configuration_template)
        return '', http.HTTPStatus.OK
    return '', http.HTTPStatus.METHOD_NOT_ALLOWED


@app.route('/v1/rb/definition/<string:rb_name>/<string:rb_version>/config-template/<string:template_name>',
           methods=["GET"])
def configuration_template_get(rb_name, rb_version, template_name):
    if request.method == "GET":
        for template in CONFIGURATIONS_TEMPLATES:
            if template['template-name'] == template_name:
                return json.dumps(template), http.HTTPStatus.OK
        else:
            return '', http.HTTPStatus.NOT_FOUND
    return '', http.HTTPStatus.METHOD_NOT_ALLOWED


@app.route('/v1/rb/definition/<string:rb_name>/<string:rb_version>/config-template',
           methods=["GET"])
def configuration_template_get_all(rb_name, rb_version):
    if request.method == 'GET':
        return json.dumps(CONFIGURATIONS_TEMPLATES), http.HTTPStatus.OK
    return '', http.HTTPStatus.METHOD_NOT_ALLOWED


if __name__ == "__main__":
    app.run(debug=True, host='0.0.0.0', port=5003)
