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
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/labstack/gommon/log"
)

func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Logger.SetLevel(log.DEBUG)
	e.GET("/", index)
	e.GET("/restservices/clds/v2/templates/", getTemplates)
	e.GET("/restservices/clds/v2/policyToscaModels/", getPolicies)

	e.GET("/restservices/clds/v2/loop/:loopID", updateLoopDetails)
	e.GET("/restservices/clds/v2/loop/getstatus/:loopID", updateLoopDetails) //refresh status
	e.POST("/restservices/clds/v2/loop/create/:loopID", createLoopInstance)
	e.PUT("/restservices/clds/v2/loop/addOperationaPolicy/:loopID/policyModel/:policyType/:policyVersion", addOperationaPolicy)
	e.PUT("/restservices/clds/v2/loop/removeOperationaPolicy/:loopID/policyModel/:policyType/:policyVersion", removeOperationaPolicy)
	e.POST("/restservices/clds/v2/loop/updateMicroservicePolicy/:loopID", addTcaConfig)                //modify
	e.POST("/restservices/clds/v2/loop/updateOperationalPolicies/:loopID", addOperationalPolicyConfig) //modify
	e.PUT("/restservices/clds/v2/loop/:action/:loopID", putLoopAction)
	e.POST("/reset", reset)
	generateInitialTemplateList()
	generateInitialPolicyList()
	generateInitialLoopInstances()
	e.Logger.Fatal(e.Start(":30258"))
}
