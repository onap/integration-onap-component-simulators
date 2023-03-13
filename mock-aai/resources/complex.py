"""A&AI Complex mock module."""
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

from flask_restful import reqparse, Resource, request

COMPLEXES = {}

parser = reqparse.RequestParser()


class Complex(Resource):
    """Complex resource class."""

    def put(self, physical_location_id: str):
        """Complex resource put method.

        Add complex data sent in JSON to COMPLEXES dictionary.

        Args:
            physical_location_id (str): Complex physical location id

        """
        COMPLEXES.update({physical_location_id: request.get_json()})

    @staticmethod
    def reset():
        """Reset Complex resource.

        Clean COMPLEXES dictionary

        """
        global COMPLEXES
        COMPLEXES = {}


class ComplexList(Resource):
    """List of complexes resource."""

    def get(self) -> Dict[str, List]:
        """Get the list of complexes.

        Return data from COMPLEXES dictionary.

        Returns:
            Dict[str, List]: Complexes dictionary

        """
        return {
            "complex": [complex_data for physical_location_id, complex_data in COMPLEXES.items()]
        }
