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

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/satori/go.uuid"
)

//describes loop template in Clamp database
type LoopTemplate struct {
	Name            string `json:"name"`
	DcaeBlueprintId string `json:"dcaeBlueprintId"`
	ModelService    struct {
		ServiceDetails struct {
			Name string `json:"name"`
		} `json:"serviceDetails"`
	} `json:"modelService"`
}

//describes policy in Clamp database
type Policy struct {
	PolicyModelType string `json:"policyModelType"`
	Version         string `json:"version"`
	PolicyAcronym   string `json:"policyAcronym"`
	CreatedDate     string `json:"createdDate"`
	UpdatedDate     string `json:"updatedDate"`
	UpdatedBy       string `json:"updatedBy"`
	CreatedBy       string `json:"createdBy"`
}

//ClampError is the way to return Error in CLAMP
type ClampError struct {
	Message string `json:"message"`
	Error   string `json:"error"`
	Status  string `json:"status"`
}

var templateList []LoopTemplate

//must modify this function to generate template with service model name
func generateInitialTemplateList() {
	templateList = nil
	u1 := uuid.NewV4().String()
	loop1 := new(LoopTemplate)
	loop1.Name = "template_service01"
	loop1.DcaeBlueprintId = u1
	loop1.ModelService.ServiceDetails.Name = "service01"
	templateList = append(templateList, *loop1)
	u2 := uuid.NewV4().String()
	loop2 := new(LoopTemplate)
	loop2.Name = "template_service02"
	loop2.DcaeBlueprintId = u2
	loop2.ModelService.ServiceDetails.Name = "service02"
	templateList = append(templateList, *loop2)
}

var policyList []Policy

func generateInitialPolicyList() {
	policyList = nil
	policyList = append(policyList, Policy{
		PolicyModelType: "onap.policies.controlloop.MinMax",
		Version:         "1.0.0",
		PolicyAcronym:   "MinMax",
		CreatedDate:     "2020-04-30T09:03:30.362897Z",
		UpdatedDate:     "2020-04-30T09:03:30.362897Z",
		UpdatedBy:       "Not found",
		CreatedBy:       "Not found",
	})
	policyList = append(policyList, Policy{
		PolicyModelType: "onap.policies.controlloop.Guard",
		Version:         "1.0.0",
		PolicyAcronym:   "Guard",
		CreatedDate:     "2020-04-30T09:03:30.362897Z",
		UpdatedDate:     "2020-04-30T09:03:30.362897Z",
		UpdatedBy:       "Not found",
		CreatedBy:       "Not found",
	})
	policyList = append(policyList, Policy{
		PolicyModelType: "onap.policies.controlloop.guard.common.FrequencyLimiter",
		Version:         "1.0.0",
		PolicyAcronym:   "FrequencyLimiter",
		CreatedDate:     "2020-04-30T09:03:30.362897Z",
		UpdatedDate:     "2020-04-30T09:03:30.362897Z",
		UpdatedBy:       "Not found",
		CreatedBy:       "Not found",
	})
}

func getTemplates(c echo.Context) error {
	var templates []LoopTemplate
	for _, t := range templateList {
		//service must be distributed from sdc
		if t.DcaeBlueprintId != "" {
			templates = append(templates, t)
		}
	}
	if len(templates) != 0 {
		return c.JSON(http.StatusOK, templates)
	}
	return c.JSON(http.StatusNotFound, ClampError{
		Message: "No Templates found",
		Error:   "Not Found",
		Status:  "404"})
}

func getPolicies(c echo.Context) error {
	if len(policyList) != 0 {
		return c.JSON(http.StatusOK, policyList)
	}
	return c.JSON(http.StatusNotFound, ClampError{
		Message: "No Policies found",
		Error:   "Not Found",
		Status:  "404"})
}
