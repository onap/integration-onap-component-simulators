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
from views import blueprint_enrich, blueprint_publish, dictionary, dictionary_get

def setup_routes(app):
    app.router.add_post('/api/v1/blueprint-model/enrich', blueprint_enrich)
    app.router.add_post('/api/v1/blueprint-model/publish', blueprint_publish)
    app.router.add_post('/api/v1/dictionary', dictionary)
    app.router.add_get('/api/v1/dictionary/{name}', dictionary_get)