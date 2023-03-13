// Copyright 2023 Deutsche Telekom AG, Orange
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import "time"

//describes operational policy configuration in CLAMP
type OperationalPolicy_payload struct {
	Name               string `json:"name"`
	JSONRepresentation struct {
		Title       string   `json:"title"`
		Type        string   `json:"type"`
		Description string   `json:"description"`
		Required    []string `json:"required"`
		Properties  struct {
			ID struct {
				Type        string `json:"type"`
				Description string `json:"description"`
			} `json:"id"`
			Actor struct {
				Type        string `json:"type"`
				Description string `json:"description"`
			} `json:"actor"`
			Operation struct {
				Type        string `json:"type"`
				Description string `json:"description"`
			} `json:"operation"`
			TimeRange struct {
				Title      string   `json:"title"`
				Type       string   `json:"type"`
				Required   []string `json:"required"`
				Properties struct {
					StartTime struct {
						Type   string `json:"type"`
						Format string `json:"format"`
					} `json:"start_time"`
					EndTime struct {
						Type   string `json:"type"`
						Format string `json:"format"`
					} `json:"end_time"`
				} `json:"properties"`
			} `json:"timeRange"`
			Limit struct {
				Type             string `json:"type"`
				Description      string `json:"description"`
				ExclusiveMinimum string `json:"exclusiveMinimum"`
			} `json:"limit"`
			TimeWindow struct {
				Type        string `json:"type"`
				Description string `json:"description"`
			} `json:"timeWindow"`
			TimeUnits struct {
				Type        string   `json:"type"`
				Description string   `json:"description"`
				Enum        []string `json:"enum"`
			} `json:"timeUnits"`
		} `json:"properties"`
	} `json:"jsonRepresentation"`
	ConfigurationsJSON struct {
		Actor      string `json:"actor"`
		Operation  string `json:"operation"`
		Limit      int    `json:"limit"`
		TimeWindow int    `json:"timeWindow"`
		TimeUnits  string `json:"timeUnits"`
	} `json:"configurationsJson"`
	PolicyModel struct {
		PolicyModelType string `json:"policyModelType"`
		Version         string `json:"version"`
		PolicyAcronym   string `json:"policyAcronym"`
		PolicyPdpGroup  struct {
			SupportedPdpGroups []struct {
				DefaultGroup []string `json:"defaultGroup"`
			} `json:"supportedPdpGroups"`
		} `json:"policyPdpGroup"`
		CreatedDate time.Time `json:"createdDate"`
		UpdatedDate time.Time `json:"updatedDate"`
		UpdatedBy   string    `json:"updatedBy"`
		CreatedBy   string    `json:"createdBy"`
	} `json:"policyModel"`
	CreatedDate time.Time `json:"createdDate"`
	UpdatedDate time.Time `json:"updatedDate"`
	UpdatedBy   string    `json:"updatedBy"`
	CreatedBy   string    `json:"createdBy"`
	PdpGroup    string    `json:"pdpGroup"`
	PdpSubgroup string    `json:"pdpSubgroup"`
}
