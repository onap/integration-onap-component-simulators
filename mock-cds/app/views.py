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

from aiohttp import web

BYTES = b'Response in bytes'
DICTIONARY = {
    "message": "Response in JSON",
    "success": True
    }
DICTIONARIES = {}

async def blueprint_enrich(request):
    """ Blueprint enrichment """
    return web.Response(body=BYTES, status=200)

async def blueprint_publish(request):
    """ Blueprint publishing """
    return web.Response(body=BYTES, status=200)

async def dictionary(request):
    """ Data dictionary """
    try:
        body = await request.json()
        DICTIONARIES[body["name"]] = body
    except ValueError:
        print("No JSON sent! Leave it because we used that endpoint during "
              "the availability and we won't break the integration tests")
    return web.json_response(data=DICTIONARY, status=200)

async def dictionary_get(request):
    """ Data dictionary get """
    name = request.match_info["name"]
    try:
        return web.json_response(data=DICTIONARIES[name], status=200)
    except KeyError:
        return web.Response(status=404)