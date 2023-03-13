"""Mock SO orchestration request resource."""
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
from dataclasses import dataclass, field
from datetime import datetime
from typing import Dict

from flask_restful import Resource


@dataclass
class OrchestrationRequestData:
    """Orchestration request dataclass."""

    request_id: str
    created_at: datetime = field(default_factory=datetime.now)


ORCHESTRATION_REQUESTS = {}


def time_diff(dt: datetime, diff: int = 1) -> bool:
    """Check if given datetime has older (in seconds) than current datetime.

    Args:
        dt (datetime): Datetime to check
        diff (int, optional): Number of seconds to check. Defaults to 1.

    Returns:
        bool: True if datetime is older, False otherwise

    """
    return (datetime.now() - dt).seconds > diff


class OrchestrationRequest(Resource):
    """Orchestration request resource."""

    @staticmethod
    def reset() -> None:
        """Reset orchestration request resource.

        Clean ORCHESTRATION_REQUESTS dictionary

        """
        global ORCHESTRATION_REQUESTS
        ORCHESTRATION_REQUESTS = {}

    def get(self, orchestration_request_id: str) -> Dict[str, Dict[str, str]]:
        """Get orchestration request data.

        Return orchestration request data from ORCHESTRATION_REQUESTS dictionary.
            If it doesn't exist it creates that.

        Args:
            orchestration_request_id (str): Orchestration request id key value

        Returns:
            Dict[str, Dict[str, str]]: Orchestration request data
        """
        try:
            orchestration_request_data = ORCHESTRATION_REQUESTS[
                orchestration_request_id
            ]
        except KeyError:
            orchestration_request_data = OrchestrationRequestData(
                request_id=orchestration_request_id
            )
            ORCHESTRATION_REQUESTS[
                orchestration_request_id
            ] = orchestration_request_data
        return {
            "request": {
                "requestStatus": {
                    "requestState": "COMPLETE"
                    if time_diff(orchestration_request_data.created_at)
                    else "IN_PROGRESS"
                }
            }
        }
