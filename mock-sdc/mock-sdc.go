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
	e.GET("/sdc1/feProxy/onboarding-api/v1.0/items/:itemID/versions", getItemVersions)
	e.GET("/sdc1/feProxy/onboarding-api/v1.0/items/:itemID/versions/:versionID", getItemVersion)
	e.PUT("/sdc1/feProxy/onboarding-api/v1.0/items/:itemID/versions/:versionID/actions", updateItemVersion)
	e.GET("/sdc1/feProxy/onboarding-api/v1.0/vendor-license-models", getVendorServiceModels)
	e.POST("/sdc1/feProxy/onboarding-api/v1.0/vendor-license-models", postVendorServiceModels)
	e.PUT("/sdc1/feProxy/onboarding-api/v1.0/vendor-license-models/:vendorID/versions/:versionID/actions", updateVendorVersion)
	e.GET("/sdc1/feProxy/onboarding-api/v1.0/vendor-software-products", getVendorSoftwareProducts)
	e.POST("/sdc1/feProxy/onboarding-api/v1.0/vendor-software-products", postVendorSoftwareProducts)
	e.GET("/sdc1/feProxy/onboarding-api/v1.0/vendor-software-products/:vspID/versions/:versionID", getVspVersion)
	e.PUT("/sdc1/feProxy/onboarding-api/v1.0/vendor-software-products/:vspID/versions/:versionID/actions", updateVspVersion)
	e.POST("/sdc1/feProxy/onboarding-api/v1.0/vendor-software-products/:vspID/versions/:versionID/orchestration-template-candidate", uploadArtifacts)
	e.PUT("/sdc1/feProxy/onboarding-api/v1.0/vendor-software-products/:vspID/versions/:versionID/orchestration-template-candidate/process", validateArtifacts)
	e.GET("/sdc1/feProxy/rest/v1/followed", getAllResources)
	e.GET("/sdc1/feProxy/rest/v1/screen", getAllResources)
	e.GET("/sdc/v1/catalog/resources", getResources)
	e.GET("/sdc/v1/catalog/services", getServices)
	e.GET("/sdc/v1/artifactTypes", getArtifactTypes)
	e.GET("/sdc/v1/distributionKafkaData", distributionKafkaData)
	e.POST("/sdc/v1/registerForDistribution", registerForDistribution)
	e.POST("/sdc/v1/unRegisterForDistribution", unRegisterForDistribution)
	e.POST("/sdc1/feProxy/rest/v1/catalog/resources", postResources)
	e.POST("/sdc1/feProxy/rest/v1/catalog/resources/:resourceID/lifecycleState/:action", postResourceAction)
	e.POST("/sdc1/feProxy/rest/v1/catalog/services", postResources)
	e.POST("/sdc1/feProxy/rest/v1/catalog/services/:resourceID/resourceInstance", postAddResourceToService)
	e.POST("/sdc1/feProxy/rest/v1/catalog/services/:resourceID/lifecycleState/:action", postResourceAction)
	e.POST("/sdc1/feProxy/rest/v1/catalog/services/:resourceID/distribution-state/:action", postResourceAction)
	e.POST("/sdc1/feProxy/rest/v1/catalog/services/:resourceID/distribution/PROD/:action", postResourceAction)
	e.GET("/sdc1/feProxy/rest/v1/catalog/services/:resourceID/distribution", getDistribution)
	e.GET("/sdc1/feProxy/rest/v1/catalog/services/distribution/:distributionID", getDistributionList)
	e.GET("/sdc1/feProxy/rest/v1/catalog/services/:resourceID", getServiceUniqueIdentifier)
	e.POST("/sdc1/feProxy/rest/v1/catalog/services/:resourceID/resourceInstance/:vfID/artifacts", uploadTcaArtifact)
	e.POST("/sdc1/feProxy/rest/v1/catalog/services/:resourceID/properties", postResourceProperties)
	e.POST("/sdc1/feProxy/rest/v1/catalog/services/:resourceID/create/inputs", postResourceInputs)
	e.POST("/sdc1/feProxy/rest/v1/catalog/resources/:resourceID/create/inputs", postResourceInputs)
	e.GET("/sdc1/feProxy/rest/v1/catalog/resources/:resourceID/filteredDataByParams", getResourcefilteredData)
	e.GET("/sdc1/feProxy/rest/v1/catalog/services/:resourceID/filteredDataByParams", getResourcefilteredData)
	e.GET("/sdc1/feProxy/rest/v1/setup/ui", getCategories)
	e.POST("/reset", reset)
	generateInitialVendorList()
	generateInitialVspList()
	generateInitialResourceList()
	generateDistributionStatusList()
	e.Logger.Fatal(e.Start(":30206"))
}
