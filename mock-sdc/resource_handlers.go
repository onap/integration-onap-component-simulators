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
	"container/list"
	"encoding/json"
	"net/http"
	"errors"
	"io/ioutil"

	"github.com/labstack/echo"
	uuid "github.com/satori/go.uuid"
)

// ResourceLight describes license model in SDC
type ResourceLight struct {
	ID                 string `json:"uuid"`
	InvariantID        string `json:"invariantUUID"`
	ResourceType       string `json:"resourceType"`
	Name               string `json:"name"`
	Category           string `json:"category"`
	SubCategory        string `json:"subCategory"`
	LastUpdaterUserID  string `json:"lastUpdaterUserId"`
	LifecycleState     string `json:"lifecycleState"`
	Version            string `json:"version"`
	ToscaModelURL      string `json:"toscaModelURL"`
	DistributionStatus string `json:"distributionStatus"`
}

// SubCategory describes SubCategory model in SDC
type SubCategory struct {
	Name           string   `json:"name"`
	NormalizedName string   `json:"normalizedName"`
	UniqueID       string   `json:"uniqueId"`
	Icons          []string `json:"icons"`
	Groupings      string   `json:"groupings"`
	OwnerID        string   `json:"ownerId"`
	Empty          bool     `json:"empty"`
	Version        string   `json:"version"`
	Type           string   `json:"type"`
}

// Category describes Category model in SDC
type Category struct {
	Name           string        `json:"name"`
	NormalizedName string        `json:"normalizedName"`
	UniqueID       string        `json:"uniqueId"`
	Icons          []string      `json:"icons"`
	Subcategories  []SubCategory `json:"subcategories"`
	OwnerID        string        `json:"ownerId"`
	Empty          bool          `json:"empty"`
	Type           string        `json:"type"`
	Version        string        `json:"version"`
}

//ArtifactAdd Describes ressource component Instances artifacts in SDC
type ArtifactAdd struct {
	ArtifactName  string `json:"artifactName"`
	ArtifactLabel string `json:"artifactLabel"`
	ArtifactType  string `json:"artifactType"`
	Description   string `json:"description"`
}

//ComponentInstance Describes ressource component Instances in SDC
type ComponentInstance struct {
	UniqueID            string        `json:"uniqueId"`
	Name                string        `json:"name"`
	ComponentName       string        `json:"componentName"`
	OriginType          string        `json:"originType"`
	ComponentVersion    string        `json:"componentVersion"`
	DeploymentArtifacts []ArtifactAdd `json:"deploymentArtifacts"`
}

// Resource describes Resource model in SDC
type Resource struct {
	ID                           string              `json:"uuid"`
	InvariantID                  string              `json:"invariantUUID"`
	UniqueID                     string              `json:"uniqueId"`
	ResourceType                 string              `json:"resourceType"`
	Name                         string              `json:"name"`
	Category                     string              `json:"category"`
	SubCategory                  string              `json:"subCategory"`
	LastUpdaterUserID            string              `json:"lastUpdaterUserId"`
	LifecycleState               string              `json:"lifecycleState"`
	Version                      string              `json:"version"`
	ToscaModelURL                string              `json:"toscaModelURL"`
	Artifacts                    struct{}            `json:"artifacts"`
	Attributes                   []string            `json:"attributes"`
	Capabilities                 struct{}            `json:"capabilities"`
	Categories                   []Category          `json:"categories"`
	ComponentInstances           []ComponentInstance `json:"componentInstances"`
	ComponentInstancesAttributes struct{}            `json:"componentInstancesAttributes"`
	ComponentInstancesProperties struct{}            `json:"componentInstancesProperties"`
	ComponentType                string              `json:"componentType"`
	ContactID                    string              `json:"contactId"`
	CsarUUID                     string              `json:"csarUUID"`
	CsarVersion                  string              `json:"csarVersion"`
	DeploymentArtifacts          struct{}            `json:"deploymentArtifacts"`
	Description                  string              `json:"description"`
	Icon                         string              `json:"icon"`
	Properties                   []Property          `json:"properties"`
	Requirements                 struct{}            `json:"requirements"`
	Tags                         []string            `json:"tags"`
	ToscaArtifacts               struct{}            `json:"toscaArtifacts"`
	VendorName                   string              `json:"vendorName"`
	VendorRelease                string              `json:"vendorRelease"`
	DistributionStatus           string              `json:"distributionStatus"`
	DistributionID               string              `json:"distributionID"`
	Inputs                       []Input
}

// ResourceList is the way to return Resources in SDC via DeepLoad
type ResourceList struct {
	Resources []Resource `json:"resources"`
	Services  []Resource `json:"services"`
}

// ActionBody yolo
type ActionBody struct {
	UserRemarks string `json:"userRemarks"`
}

// ResourceAdd to a Service
type ResourceAdd struct {
	Name             string `json:"name"`
	ComponentVersion string `json:"componentVersion"`
	PosY             int    `json:"posY"`
	PosX             int    `json:"posX"`
	UniqueID         string `json:"uniqueId"`
	OriginType       string `json:"originType"`
	ComponentUID     string `json:"componentUid"`
	Icon             string `json:"icon"`
}

// DistributionIDResult format
type DistributionIDResult struct {
	DistributionID    string `json:"distributionID"`
	UserID            string `json:"userId"`
	DeployementStatus string `json:"deployementStatus"`
}

// DistributionIDList format
type DistributionIDList struct {
	DistributionStatusOfServiceList []DistributionIDResult `json:"distributionStatusOfServiceList"`
}

//DistributionStatus format
type DistributionStatus struct {
	OmfComponentID string `json:"omfComponentID"`
	Timestamp      string `json:"timestamp"`
	URL            string `json:"url"`
	Status         string `json:"status"`
	ErrorReason    string `json:"errorReason"`
}

// DistributionStatusList format
type DistributionStatusList struct {
	DistributionStatusList []DistributionStatus `json:"distributionStatusList"`
}

// NewUploadResult format
type NewUploadResult struct {
	Description  string `json:"description"`
	ArtifactType string `json:"artifactType"`
	ArtifactName string `json:"artifactName"`
}

//Property format
type Property struct {
	Name           string `json:"name"`
	Value          string `json:"value"`
	Type           string `json:"type"`
	UniqueID       string `json:"uniqueId"`
	ParentUniqueID string `json:"parentUniqueId"`
}

//Input format
type Input struct {
	Name     string `json:"name"`
	Value    string `json:"value"`
	Type     string `json:"type"`
	UniqueID string `json:"uniqueId"`
}

var resourceList []Resource
var distributionList []DistributionStatus

func generateInitialResourceList() {
	resourceList = nil
	resourceList = append(resourceList, Resource{
		ID:                "6c4952d2-0ecc-4697-a039-d9766565feae",
		InvariantID:       "803cbaf5-deea-4022-a731-709d285435d6",
		UniqueID:          "1e6e90ec-632a-492f-9511-f2787a2befaf",
		ResourceType:      "Configuration",
		Name:              "VLAN Network Receptor Configuration",
		Category:          "Configuration",
		SubCategory:       "Configuration",
		LastUpdaterUserID: "jh0003",
		LifecycleState:    "CERTIFIED",
		Version:           "1.0",
		ToscaModelURL:     "/sdc/v1/catalog/resources/6c4952d2-0ecc-4697-a039-d9766565feae/toscaModel",
	})
	resourceList = append(resourceList, Resource{
		ID:                "85a9a912-b0ca-4cc9-9dc4-a480546ef93b",
		InvariantID:       "2df7615c-38f5-45e2-ac40-f9a8f97baec2",
		UniqueID:          "1e6e90ec-632a-492f-9511-f2787a2befaf",
		ResourceType:      "CP",
		Name:              "contrailV2VLANSubInterfaceV2",
		Category:          "Generic",
		SubCategory:       "Network Elements",
		LastUpdaterUserID: "jh0003",
		LifecycleState:    "CERTIFIED",
		Version:           "1.0",
		ToscaModelURL:     "/sdc/v1/catalog/resources/85a9a912-b0ca-4cc9-9dc4-a480546ef93b/toscaModel",
	})
	resourceList = append(resourceList, Resource{
		ID:                "7c6b6644-590d-4e60-84d7-0dfba3ad4694",
		InvariantID:       "1e6e90ec-632a-492f-9511-f2787a2bef9f",
		UniqueID:          "1e6e90ec-632a-492f-9511-f2787a2befaf",
		ResourceType:      "VFC",
		Name:              "VDU Compute",
		Category:          "Generic",
		SubCategory:       "Infrastructure",
		LastUpdaterUserID: "jh0003",
		LifecycleState:    "CERTIFIED",
		Version:           "1.0",
		ToscaModelURL:     "/sdc/v1/catalog/resources/7c6b6644-590d-4e60-84d7-0dfba3ad4694/toscaModel",
	})
	resourceList = append(resourceList, Resource{
		ID:                "9391354f-8f25-462d-b331-841e6cc5c851",
		InvariantID:       "85cd3f14-cb9c-4a28-811b-d076e9a48303",
		UniqueID:          "1e6e90ec-632a-492f-9511-f2787a2befaf",
		ResourceType:      "VFC",
		Name:              "Cp",
		Category:          "Generic",
		SubCategory:       "Infrastructure",
		LastUpdaterUserID: "jh0003",
		LifecycleState:    "CERTIFIED",
		Version:           "1.0",
		ToscaModelURL:     "/sdc/v1/catalog/resources/9391354f-8f25-462d-b331-841e6cc5c851/toscaModel",
	})
}

func getResources(c echo.Context) error {
	resourceType := c.QueryParam("resourceType")
	resources := []ResourceLight{}
	for _, r := range resourceList {
		if (r.ComponentType != "SERVICE") &&
			((resourceType == "") || (r.ResourceType == resourceType)) {
			resources = append(resources, ResourceLight{
				ID:                r.ID,
				InvariantID:       r.InvariantID,
				ResourceType:      r.ResourceType,
				Name:              r.Name,
				Category:          r.Category,
				SubCategory:       r.SubCategory,
				LastUpdaterUserID: r.LastUpdaterUserID,
				LifecycleState:    r.LifecycleState,
				Version:           r.Version,
				ToscaModelURL:     r.ToscaModelURL,
			})
		}
	}
	if len(resources) != 0 {
		return c.JSON(http.StatusOK, resources)
	}
	return c.JSON(http.StatusNotFound, SdcError{
		Message:   "No Resources found",
		ErrorCode: "SVC4642",
		Status:    "Not Found"})
}

func getServices(c echo.Context) error {
	resources := []ResourceLight{}
	for _, r := range resourceList {
		if r.ComponentType == "SERVICE" {
			resources = append(resources, ResourceLight{
				ID:                 r.ID,
				InvariantID:        r.InvariantID,
				ResourceType:       r.ResourceType,
				Name:               r.Name,
				Category:           r.Category,
				SubCategory:        r.SubCategory,
				LastUpdaterUserID:  r.LastUpdaterUserID,
				LifecycleState:     r.LifecycleState,
				Version:            r.Version,
				ToscaModelURL:      r.ToscaModelURL,
				DistributionStatus: r.DistributionStatus,
			})
		}
	}
	if len(resources) != 0 {
		return c.JSON(http.StatusOK, resources)
	}
	return c.JSON(http.StatusNotFound, SdcError{
		Message:   "No Resources found",
		ErrorCode: "SVC4642",
		Status:    "Not Found"})
}

func postResources(c echo.Context) error {
	resource := new(Resource)
	if err := c.Bind(resource); err != nil {
		return err
	}
	for _, r := range resourceList {
		if resource.Name == r.Name && r.ResourceType == resource.ResourceType {
			return c.JSON(http.StatusBadRequest, SdcError{
				Message:   "Resource of same Name and ResourceType exists",
				ErrorCode: "SVC3642",
				Status:    "Exists"})
		}
	}
	resource.ID = uuid.Must(uuid.NewV4()).String()
	resource.InvariantID = uuid.Must(uuid.NewV4()).String()
	resource.UniqueID = uuid.Must(uuid.NewV4()).String()
	resource.Version = "0.1"
	resource.LifecycleState = "NOT_CERTIFIED_CHECKOUT"
	resource.DistributionStatus = "DISTRIBUTION_NOT_APPROVED"

	resourceList = append(resourceList, *resource)

	return c.JSON(http.StatusCreated, resource)
}

func postResourceAction(c echo.Context) error {
	resourceID := c.Param("resourceID")
	action := c.Param("action")
	actionBody := new(ActionBody)
	if err := c.Bind(actionBody); err != nil {
		return err
	}
	for i, r := range resourceList {
		if r.UniqueID == resourceID {
			if r.LifecycleState == "NOT_CERTIFIED_CHECKOUT" && action == "Certify" {
				resourceList[i].Version = "1.0"
				resourceList[i].LifecycleState = "CERTIFIED"
				return c.JSON(http.StatusCreated, resourceList[i])
			}
			if r.LifecycleState == "NOT_CERTIFIED_CHECKOUT" && action == "checkin" {
				resourceList[i].LifecycleState = "NOT_CERTIFIED_CHECKIN"
				return c.JSON(http.StatusOK, resourceList[i])
			}
			if r.LifecycleState == "NOT_CERTIFIED_CHECKIN" && action == "Certify" {
				resourceList[i].LifecycleState = "CERTIFIED"
				resourceList[i].Version = "1.0"
				resourceList[i].DistributionStatus = "DISTRIBUTION_APPROVED"
				return c.JSON(http.StatusOK, resourceList[i])
			}
			if r.LifecycleState == "CERTIFIED" &&
				r.DistributionStatus == "DISTRIBUTION_APPROVED" &&
				action == "activate" {
				resourceList[i].DistributionStatus = "DISTRIBUTED"
				return c.JSON(http.StatusOK, resourceList[i])
			}
			return c.JSON(http.StatusBadRequest, SdcError{
				Message:   "Cannot perform this action",
				ErrorCode: "SVC3642",
				Status:    "Bad Action"})
		}
	}
	return c.JSON(http.StatusNotFound, SdcError{
		Message:   "Resource not found",
		ErrorCode: "SVC4642",
		Status:    "Not Found"})
}

func postAddResourceToService(c echo.Context) error {
	resourceID := c.Param("resourceID")
	resourceAdd := new(ResourceAdd)
	if err := c.Bind(resourceAdd); err != nil {
		return err
	}
	for i, r := range resourceList {
		if r.UniqueID == resourceID {
			if r.LifecycleState == "NOT_CERTIFIED_CHECKOUT" {
				for _, rr := range resourceList {
					if rr.UniqueID == resourceAdd.UniqueID &&
						rr.UniqueID == resourceAdd.ComponentUID &&
						rr.Name == resourceAdd.Name &&
						rr.Version == resourceAdd.ComponentVersion &&
						rr.ResourceType == resourceAdd.OriginType {
						if rr.ResourceType == "VF" {
							ci := ComponentInstance{
								UniqueID:         uuid.Must(uuid.NewV4()).String(),
								Name:             resourceAdd.Name,
								ComponentName:    resourceAdd.Name,
								OriginType:       "VF",
								ComponentVersion: "1.0",
							}
							resourceList[i].ComponentInstances = append(r.ComponentInstances, ci)
						}
						return c.JSON(http.StatusCreated, r)
					}
				}
			}

			return c.JSON(http.StatusBadRequest, SdcError{
				Message:   "Cannot perform this action",
				ErrorCode: "SVC3642",
				Status:    "Bad Action"})
		}
	}

	return c.JSON(http.StatusNotFound, SdcError{
		Message:   "Resource not found",
		ErrorCode: "SVC4642",
		Status:    "Not Found"})
}

func getDistribution(c echo.Context) error {
	resourceID := c.Param("resourceID")
	for i, r := range resourceList {
		if r.ID == resourceID {
			distributionIDResult := new(DistributionIDResult)
			if r.DistributionStatus == "DISTRIBUTED" {
				if len(r.DistributionID) < 1 {
					resourceList[i].DistributionID = uuid.Must(uuid.NewV4()).String()
				}
				distributionIDResult.DeployementStatus = "Distributed"
				distributionIDResult.UserID = "Oper P(op0001)"
				distributionIDResult.DistributionID = resourceList[i].DistributionID
			}
			distributionIDList := new(DistributionIDList)
			distributionIDList.DistributionStatusOfServiceList = append(distributionIDList.DistributionStatusOfServiceList, *distributionIDResult)
			return c.JSON(http.StatusOK, distributionIDList)
		}
	}
	return c.JSON(http.StatusNotFound, SdcError{
		Message:   "Resource not found",
		ErrorCode: "SVC4642",
		Status:    "Not Found"})
}

func getDistributionList(c echo.Context) error {
	distributionID := c.Param("distributionID")
	for _, r := range resourceList {
		if r.DistributionID == distributionID {
			d := new(DistributionStatusList)
			d.DistributionStatusList = distributionList
			return c.JSON(http.StatusOK, d)
		}
	}
	return c.JSON(http.StatusNotFound, SdcError{
		Message:   "Resource not found",
		ErrorCode: "SVC4642",
		Status:    "Not Found"})
}

func getAllResources(c echo.Context) error {

	var listResources []Resource
	var listServices []Resource

	for _, r := range resourceList {
		if r.ComponentType != "SERVICE" {
			listResources = append(listResources, r)
		} else {
			listServices = append(listServices, r)
		}
	}
	list := &ResourceList{listResources, listServices}
	return c.JSON(http.StatusOK, list)
}

func getServiceUniqueIdentifier(c echo.Context) error {
	resourceID := c.Param("resourceID")
	for _, r := range resourceList {
		if r.UniqueID == resourceID {
			return c.JSON(http.StatusOK, r)
		}
	}

	return c.JSON(http.StatusNotFound, SdcError{
		Message:   "Resource not found",
		ErrorCode: "SVC4642",
		Status:    "Not Found"})
}

func uploadTcaArtifact(c echo.Context) error {
	resourceID := c.Param("resourceID")
	vfID := c.Param("vfID")
	for _, r := range resourceList {
		if r.UniqueID == resourceID {
			if r.LifecycleState == "NOT_CERTIFIED_CHECKOUT" {
				for _, cc := range r.ComponentInstances {
					if cc.UniqueID == vfID {
						newArtifact := new(ArtifactAdd)
						if err := c.Bind(newArtifact); err != nil {
							return err
						}
						cc.DeploymentArtifacts = append(cc.DeploymentArtifacts, *newArtifact)
						NewUploadResult := NewUploadResult{
							Description:  newArtifact.Description,
							ArtifactType: newArtifact.ArtifactType,
							ArtifactName: newArtifact.ArtifactName,
						}
						return c.JSON(http.StatusCreated, NewUploadResult)
					}
				}
			}
			return c.JSON(http.StatusBadRequest, SdcError{
				Message:   "Cannot perform this action",
				ErrorCode: "SVC3642",
				Status:    "Bad Action"})
		}
	}
	return c.JSON(http.StatusNotFound, SdcError{
		Message:   "Resource not found",
		ErrorCode: "SVC4642",
		Status:    "Not Found"})
}

func postResourceProperties(c echo.Context) error {
	resourceID := c.Param("resourceID")
	var bodyBytes []byte
	if c.Request().Body != nil {
		bodyBytes, _ = ioutil.ReadAll(c.Request().Body)
	}
	var dat map[string]interface{}
	if err := json.Unmarshal(bodyBytes, &dat); err != nil {
		return err
	}
	for i, r := range resourceList {
		if r.UniqueID == resourceID {
			for key := range dat {
				propertyBody := new(Property)
				propertyBody.Name = dat[key].(map[string]interface{})["name"].(string)
				propertyBody.Type = dat[key].(map[string]interface{})["type"].(string)
				resourceList[i].Properties = append(r.Properties, *propertyBody)
			}
			return c.JSON(http.StatusOK, "")
		}
	}
	return c.JSON(http.StatusNotFound, "")
}

func getResourceProperties(c echo.Context) error {
	resourceID := c.Param("resourceID")
	for _, r := range resourceList {
		if r.UniqueID == resourceID {
			return c.JSON(http.StatusOK, map[string][]Property{
				"properties": r.Properties,
			})
		}
	}
	return c.JSON(http.StatusNotFound, "")
}

func postResourceInputs(c echo.Context) error {
	resourceID := c.Param("resourceID")
	inputBody := new(Input)
	if err := c.Bind(inputBody); err != nil {
		return err
	}
	for i, r := range resourceList {
		if r.UniqueID == resourceID {
			resourceList[i].Inputs = append(r.Inputs, *inputBody)
			return c.JSON(http.StatusOK, r.Inputs)
		}
	}
	return c.JSON(http.StatusNotFound, "")
}

func getResourcefilteredData(c echo.Context) error {
	paramType := c.QueryParam("include")
	switch paramType {
	case "inputs":
		return getResourceInputs(c)
	case "properties":
		return getResourceProperties(c)
	}
	return errors.New("Invalid query param")
}

func getResourceInputs(c echo.Context) error {
	resourceID := c.Param("resourceID")
	for _, r := range resourceList {
		if r.UniqueID == resourceID {
			return c.JSON(http.StatusOK, map[string][]Input{
				"inputs": r.Inputs,
			})
		}
	}
	return c.JSON(http.StatusNotFound, "")
}

func generateDistributionStatusList() {
	distributionList = nil
	distributionList = append(distributionList, DistributionStatus{
		OmfComponentID: "SO-COpenSource-Env11",
		Timestamp:      "1574774740421",
		URL:            "/sdc/v1/catalog/services/Test12/1.0/resourceInstances/vsrx0/artifacts/vf-license-model.xml",
		Status:         "NOT_NOTIFIED",
		ErrorReason:    "null",
	})
	distributionList = append(distributionList, DistributionStatus{
		OmfComponentID: "aai-ml",
		Timestamp:      "1574774737842",
		URL:            "/sdc/v1/catalog/services/Test12/1.0/resourceInstances/vsrx0/artifacts/vf-license-model.xml",
		Status:         "NOT_NOTIFIED",
		ErrorReason:    "null",
	})
	distributionList = append(distributionList, DistributionStatus{
		OmfComponentID: "SO-COpenSource-Env11",
		Timestamp:      "1574774740421",
		URL:            "/sdc/v1/catalog/services/Test12/1.0/artifacts/service-Test12-template.yml",
		Status:         "NOT_NOTIFIED",
		ErrorReason:    "null",
	})
	distributionList = append(distributionList, DistributionStatus{
		OmfComponentID: "cds",
		Timestamp:      "1574774726254",
		URL:            "/sdc/v1/catalog/services/Test12/1.0/artifacts/service-Test12-csar.csar",
		Status:         "NOTIFIED",
		ErrorReason:    "null",
	})
	distributionList = append(distributionList, DistributionStatus{
		OmfComponentID: "multicloud-windriver-id",
		Timestamp:      "1574774731805",
		URL:            "/sdc/v1/catalog/services/Test12/1.0/artifacts/service-Test12-csar.csar",
		Status:         "NOTIFIED",
		ErrorReason:    "null",
	})
	distributionList = append(distributionList, DistributionStatus{
		OmfComponentID: "sdc-COpenSource-Env11-sdnc-dockero",
		Timestamp:      "1574774720318",
		URL:            "/sdc/v1/catalog/services/Test12/1.0/resourceInstances/vsrx0/artifacts/vsrx0_modules.json",
		Status:         "NOT_NOTIFIED",
		ErrorReason:    "null",
	})
	distributionList = append(distributionList, DistributionStatus{
		OmfComponentID: "sdc-COpenSource-Env11-sdnc-dockero",
		Timestamp:      "1574774737396",
		URL:            "/sdc/v1/catalog/services/Test12/1.0/artifacts/service-Test12-csar.csar",
		Status:         "DEPLOY_OK",
		ErrorReason:    "null",
	})
	distributionList = append(distributionList, DistributionStatus{
		OmfComponentID: "clamp",
		Timestamp:      "1574774737925",
		URL:            "/sdc/v1/catalog/services/Test12/1.0/resourceInstances/vsrx0/artifacts/vendor-license-model.xml",
		Status:         "NOT_NOTIFIED",
		ErrorReason:    "null",
	})
	distributionList = append(distributionList, DistributionStatus{
		OmfComponentID: "multicloud-k8s-id",
		Timestamp:      "1574774750490",
		URL:            "/sdc/v1/catalog/services/Test12/1.0/resourceInstances/vsrx0/artifacts/vsrx0_modules.json",
		Status:         "DEPLOY_OK",
		ErrorReason:    "null",
	})
	distributionList = append(distributionList, DistributionStatus{
		OmfComponentID: "multicloud-k8s-id",
		Timestamp:      "1574774736174",
		URL:            "/sdc/v1/catalog/services/Test12/1.0/resourceInstances/vsrx0/artifacts/vsrx0_modules.json",
		Status:         "NOTIFIED",
		ErrorReason:    "null",
	})
	distributionList = append(distributionList, DistributionStatus{
		OmfComponentID: "multicloud-windriver-id",
		Timestamp:      "1574774731805",
		URL:            "/sdc/v1/catalog/services/Test12/1.0/resourceInstances/vsrx0/artifacts/base_ubuntu16.env",
		Status:         "NOTIFIED",
		ErrorReason:    "null",
	})
	distributionList = append(distributionList, DistributionStatus{
		OmfComponentID: "SO-COpenSource-Env11",
		Timestamp:      "1574774740421",
		URL:            "/sdc/v1/catalog/services/Test12/1.0/artifacts/service-Test12-csar.csar",
		Status:         "NOTIFIED",
		ErrorReason:    "null",
	})
	distributionList = append(distributionList, DistributionStatus{
		OmfComponentID: "clamp",
		Timestamp:      "1574774737925",
		URL:            "/sdc/v1/catalog/services/Test12/1.0/resourceInstances/vsrx0/artifacts/vf-license-model.xml",
		Status:         "NOT_NOTIFIED",
		ErrorReason:    "null",
	})
	distributionList = append(distributionList, DistributionStatus{
		OmfComponentID: "policy-id",
		Timestamp:      "1574774728667",
		URL:            "/sdc/v1/catalog/services/Test12/1.0/resourceInstances/vsrx0/artifacts/base_ubuntu16.yaml",
		Status:         "NOTIFIED",
		ErrorReason:    "null",
	})
	distributionList = append(distributionList, DistributionStatus{
		OmfComponentID: "multicloud-starlingx-id",
		Timestamp:      "1574774737784",
		URL:            "/sdc/v1/catalog/services/Test12/1.0/resourceInstances/vsrx0/artifacts/base_ubuntu16.env",
		Status:         "NOTIFIED",
		ErrorReason:    "null",
	})
	distributionList = append(distributionList, DistributionStatus{
		OmfComponentID: "clamp",
		Timestamp:      "1574774737925",
		URL:            "/sdc/v1/catalog/services/Test12/1.0/artifacts/service-Test12-template.yml",
		Status:         "NOT_NOTIFIED",
		ErrorReason:    "null",
	})
	distributionList = append(distributionList, DistributionStatus{
		OmfComponentID: "multicloud-k8s-id",
		Timestamp:      "1574774744883",
		URL:            "/sdc/v1/catalog/services/Test12/1.0/resourceInstances/vsrx0/artifacts/vsrx0_modules.json",
		Status:         "DOWNLOAD_OK",
		ErrorReason:    "null",
	})
	distributionList = append(distributionList, DistributionStatus{
		OmfComponentID: "multicloud-windriver-id",
		Timestamp:      "1574774731805",
		URL:            "/sdc/v1/catalog/services/Test12/1.0/resourceInstances/vsrx0/artifacts/vsrx0_modules.json",
		Status:         "NOTIFIED",
		ErrorReason:    "null",
	})
	distributionList = append(distributionList, DistributionStatus{
		OmfComponentID: "multicloud-windriver-id",
		Timestamp:      "1574774731805",
		URL:            "/sdc/v1/catalog/services/Test12/1.0/resourceInstances/vsrx0/artifacts/base_ubuntu16.yaml",
		Status:         "NOTIFIED",
		ErrorReason:    "null",
	})
	distributionList = append(distributionList, DistributionStatus{
		OmfComponentID: "SO-COpenSource-Env11",
		Timestamp:      "1574774740421",
		URL:            "/sdc/v1/catalog/services/Test12/1.0/resourceInstances/vsrx0/artifacts/base_ubuntu16.yaml",
		Status:         "NOTIFIED",
		ErrorReason:    "null",
	})
	distributionList = append(distributionList, DistributionStatus{
		OmfComponentID: "SO-COpenSource-Env11",
		Timestamp:      "1574774756945",
		URL:            "/sdc/v1/catalog/services/Test12/1.0/resourceInstances/vsrx0/artifacts/base_ubuntu16.env",
		Status:         "DEPLOY_OK",
		ErrorReason:    "null",
	})
	distributionList = append(distributionList, DistributionStatus{
		OmfComponentID: "SO-COpenSource-Env11",
		Timestamp:      "1574774752508",
		URL:            "/sdc/v1/catalog/services/Test12/1.0/resourceInstances/vsrx0/artifacts/base_ubuntu16.yaml",
		Status:         "DOWNLOAD_OK",
		ErrorReason:    "null",
	})
	distributionList = append(distributionList, DistributionStatus{
		OmfComponentID: "cds",
		Timestamp:      "1574774726254",
		URL:            "/sdc/v1/catalog/services/Test12/1.0/resourceInstances/vsrx0/artifacts/vsrx0_modules.json",
		Status:         "NOT_NOTIFIED",
		ErrorReason:    "null",
	})
	distributionList = append(distributionList, DistributionStatus{
		OmfComponentID: "multicloud-k8s-id",
		Timestamp:      "1574774736174",
		URL:            "/sdc/v1/catalog/services/Test12/1.0/resourceInstances/vsrx0/artifacts/vendor-license-model.xml",
		Status:         "NOT_NOTIFIED",
		ErrorReason:    "null",
	})
	distributionList = append(distributionList, DistributionStatus{
		OmfComponentID: "aai-ml",
		Timestamp:      "1574774745892",
		URL:            "/sdc/v1/catalog/services/Test12/1.0/artifacts/service-Test12-csar.csar",
		Status:         "DOWNLOAD_OK",
		ErrorReason:    "null",
	})
	distributionList = append(distributionList, DistributionStatus{
		OmfComponentID: "cds",
		Timestamp:      "1574774726254",
		URL:            "/sdc/v1/catalog/services/Test12/1.0/resourceInstances/vsrx0/artifacts/vf-license-model.xml",
		Status:         "NOT_NOTIFIED",
		ErrorReason:    "null",
	})
	distributionList = append(distributionList, DistributionStatus{
		OmfComponentID: "aai-ml",
		Timestamp:      "1574774737842",
		URL:            "/sdc/v1/catalog/services/Test12/1.0/artifacts/service-Test12-csar.csar",
		Status:         "NOTIFIED",
		ErrorReason:    "null",
	})
	distributionList = append(distributionList, DistributionStatus{
		OmfComponentID: "SO-COpenSource-Env11",
		Timestamp:      "1574774757951",
		URL:            "/sdc/v1/catalog/services/Test12/1.0/artifacts/service-Test12-csar.csar",
		Status:         "DEPLOY_OK",
		ErrorReason:    "null",
	})
	distributionList = append(distributionList, DistributionStatus{
		OmfComponentID: "aai-ml",
		Timestamp:      "1574774737842",
		URL:            "/sdc/v1/catalog/services/Test12/1.0/resourceInstances/vsrx0/artifacts/base_ubuntu16.yaml",
		Status:         "NOT_NOTIFIED",
		ErrorReason:    "null",
	})
	distributionList = append(distributionList, DistributionStatus{
		OmfComponentID: "multicloud-starlingx-id",
		Timestamp:      "1574774737784",
		URL:            "/sdc/v1/catalog/services/Test12/1.0/resourceInstances/vsrx0/artifacts/vf-license-model.xml",
		Status:         "NOT_NOTIFIED",
		ErrorReason:    "null",
	})
	distributionList = append(distributionList, DistributionStatus{
		OmfComponentID: "clamp",
		Timestamp:      "1574774737925",
		URL:            "/sdc/v1/catalog/services/Test12/1.0/artifacts/service-Test12-csar.csar",
		Status:         "NOTIFIED",
		ErrorReason:    "null",
	})
	distributionList = append(distributionList, DistributionStatus{
		OmfComponentID: "cds",
		Timestamp:      "1574774726254",
		URL:            "/sdc/v1/catalog/services/Test12/1.0/artifacts/service-Test12-template.yml",
		Status:         "NOT_NOTIFIED",
		ErrorReason:    "null",
	})
	distributionList = append(distributionList, DistributionStatus{
		OmfComponentID: "dcae-sch",
		Timestamp:      "1574774724083",
		URL:            "/sdc/v1/catalog/services/Test12/1.0/resourceInstances/vsrx0/artifacts/vf-license-model.xml",
		Status:         "NOT_NOTIFIED",
		ErrorReason:    "null",
	})
	distributionList = append(distributionList, DistributionStatus{
		OmfComponentID: "multicloud-k8s-id",
		Timestamp:      "1574774749381",
		URL:            "",
		Status:         "COMPONENT_DONE_OK",
		ErrorReason:    "null",
	})
	distributionList = append(distributionList, DistributionStatus{
		OmfComponentID: "multicloud-k8s-id",
		Timestamp:      "1574774736174",
		URL:            "/sdc/v1/catalog/services/Test12/1.0/resourceInstances/vsrx0/artifacts/vf-license-model.xml",
		Status:         "NOT_NOTIFIED",
		ErrorReason:    "null",
	})
	distributionList = append(distributionList, DistributionStatus{
		OmfComponentID: "multicloud-k8s-id",
		Timestamp:      "1574774746542",
		URL:            "/sdc/v1/catalog/services/Test12/1.0/resourceInstances/vsrx0/artifacts/base_ubuntu16.yaml",
		Status:         "DOWNLOAD_OK",
		ErrorReason:    "null",
	})
	distributionList = append(distributionList, DistributionStatus{
		OmfComponentID: "policy-id",
		Timestamp:      "1574774728667",
		URL:            "/sdc/v1/catalog/services/Test12/1.0/artifacts/service-Test12-template.yml",
		Status:         "NOT_NOTIFIED",
		ErrorReason:    "null",
	})
	distributionList = append(distributionList, DistributionStatus{
		OmfComponentID: "multicloud-k8s-id",
		Timestamp:      "1574774736174",
		URL:            "/sdc/v1/catalog/services/Test12/1.0/resourceInstances/vsrx0/artifacts/base_ubuntu16.env",
		Status:         "NOTIFIED",
		ErrorReason:    "null",
	})
	distributionList = append(distributionList, DistributionStatus{
		OmfComponentID: "cds",
		Timestamp:      "1574774726254",
		URL:            "/sdc/v1/catalog/services/Test12/1.0/resourceInstances/vsrx0/artifacts/base_ubuntu16.yaml",
		Status:         "NOT_NOTIFIED",
		ErrorReason:    "null",
	})
	distributionList = append(distributionList, DistributionStatus{
		OmfComponentID: "cds",
		Timestamp:      "1574774735595",
		URL:            "/sdc/v1/catalog/services/Test12/1.0/artifacts/service-Test12-csar.csar",
		Status:         "DOWNLOAD_OK",
		ErrorReason:    "null",
	})
	distributionList = append(distributionList, DistributionStatus{
		OmfComponentID: "multicloud-starlingx-id",
		Timestamp:      "1574774737784",
		URL:            "/sdc/v1/catalog/services/Test12/1.0/resourceInstances/vsrx0/artifacts/vsrx0_modules.json",
		Status:         "NOTIFIED",
		ErrorReason:    "null",
	})
	distributionList = append(distributionList, DistributionStatus{
		OmfComponentID: "multicloud-k8s-id",
		Timestamp:      "1574774736174",
		URL:            "/sdc/v1/catalog/services/Test12/1.0/artifacts/service-Test12-template.yml",
		Status:         "NOT_NOTIFIED",
		ErrorReason:    "null",
	})
	distributionList = append(distributionList, DistributionStatus{
		OmfComponentID: "dcae-sch",
		Timestamp:      "1574774724083",
		URL:            "/sdc/v1/catalog/services/Test12/1.0/resourceInstances/vsrx0/artifacts/vsrx0_modules.json",
		Status:         "NOT_NOTIFIED",
		ErrorReason:    "null",
	})
	distributionList = append(distributionList, DistributionStatus{
		OmfComponentID: "SO-COpenSource-Env11",
		Timestamp:      "1574774748752",
		URL:            "/sdc/v1/catalog/services/Test12/1.0/artifacts/service-Test12-csar.csar",
		Status:         "DOWNLOAD_OK",
		ErrorReason:    "null",
	})
	distributionList = append(distributionList, DistributionStatus{
		OmfComponentID: "cds",
		Timestamp:      "1574774736609",
		URL:            "/sdc/v1/catalog/services/Test12/1.0/artifacts/service-Test12-csar.csar",
		Status:         "COMPONENT_DONE_OK",
		ErrorReason:    "null",
	})
	distributionList = append(distributionList, DistributionStatus{
		OmfComponentID: "SO-COpenSource-Env11",
		Timestamp:      "1574774969197",
		URL:            "",
		Status:         "DISTRIBUTION_COMPLETE_OK",
		ErrorReason:    "null",
	})
	distributionList = append(distributionList, DistributionStatus{
		OmfComponentID: "dcae-sch",
		Timestamp:      "1574774724083",
		URL:            "/sdc/v1/catalog/services/Test12/1.0/artifacts/service-Test12-csar.csar",
		Status:         "NOT_NOTIFIED",
		ErrorReason:    "null",
	})
	distributionList = append(distributionList, DistributionStatus{
		OmfComponentID: "aai-ml",
		Timestamp:      "1574774737842",
		URL:            "/sdc/v1/catalog/services/Test12/1.0/artifacts/service-Test12-template.yml",
		Status:         "NOT_NOTIFIED",
		ErrorReason:    "null",
	})
	distributionList = append(distributionList, DistributionStatus{
		OmfComponentID: "aai-ml",
		Timestamp:      "1574774750517",
		URL:            "/sdc/v1/catalog/services/Test12/1.0/artifacts/service-Test12-csar.csar",
		Status:         "DEPLOY_OK",
		ErrorReason:    "null",
	})
	distributionList = append(distributionList, DistributionStatus{
		OmfComponentID: "multicloud-windriver-id",
		Timestamp:      "1574774744623",
		URL:            "",
		Status:         "COMPONENT_DONE_OK",
		ErrorReason:    "null",
	})
	distributionList = append(distributionList, DistributionStatus{
		OmfComponentID: "dcae-sch",
		Timestamp:      "1574774724083",
		URL:            "/sdc/v1/catalog/services/Test12/1.0/resourceInstances/vsrx0/artifacts/base_ubuntu16.yaml",
		Status:         "NOT_NOTIFIED",
		ErrorReason:    "null",
	})
	distributionList = append(distributionList, DistributionStatus{
		OmfComponentID: "dcae-sch",
		Timestamp:      "1574774724083",
		URL:            "/sdc/v1/catalog/services/Test12/1.0/resourceInstances/vsrx0/artifacts/base_ubuntu16.env",
		Status:         "NOT_NOTIFIED",
		ErrorReason:    "null",
	})
	distributionList = append(distributionList, DistributionStatus{
		OmfComponentID: "multicloud-windriver-id",
		Timestamp:      "1574774731805",
		URL:            "/sdc/v1/catalog/services/Test12/1.0/artifacts/service-Test12-template.yml",
		Status:         "NOT_NOTIFIED",
		ErrorReason:    "null",
	})
	distributionList = append(distributionList, DistributionStatus{
		OmfComponentID: "multicloud-starlingx-id",
		Timestamp:      "1574774750947",
		URL:            "",
		Status:         "COMPONENT_DONE_OK",
		ErrorReason:    "null",
	})
	distributionList = append(distributionList, DistributionStatus{
		OmfComponentID: "sdc-COpenSource-Env11-sdnc-dockero",
		Timestamp:      "1574774735764",
		URL:            "/sdc/v1/catalog/services/Test12/1.0/artifacts/service-Test12-csar.csar",
		Status:         "DOWNLOAD_OK",
		ErrorReason:    "null",
	})
	distributionList = append(distributionList, DistributionStatus{
		OmfComponentID: "multicloud-starlingx-id",
		Timestamp:      "1574774737784",
		URL:            "/sdc/v1/catalog/services/Test12/1.0/artifacts/service-Test12-csar.csar",
		Status:         "NOTIFIED",
		ErrorReason:    "null",
	})
	distributionList = append(distributionList, DistributionStatus{
		OmfComponentID: "cds",
		Timestamp:      "1574774726254",
		URL:            "/sdc/v1/catalog/services/Test12/1.0/resourceInstances/vsrx0/artifacts/vendor-license-model.xml",
		Status:         "NOT_NOTIFIED",
		ErrorReason:    "null",
	})
	distributionList = append(distributionList, DistributionStatus{
		OmfComponentID: "clamp",
		Timestamp:      "1574774750026",
		URL:            "/sdc/v1/catalog/services/Test12/1.0/artifacts/service-Test12-csar.csar",
		Status:         "ALREADY_DOWNLOADED",
		ErrorReason:    "null",
	})
	distributionList = append(distributionList, DistributionStatus{
		OmfComponentID: "SO-COpenSource-Env11",
		Timestamp:      "1574774753902",
		URL:            "/sdc/v1/catalog/services/Test12/1.0/resourceInstances/vsrx0/artifacts/base_ubuntu16.env",
		Status:         "DOWNLOAD_OK",
		ErrorReason:    "null",
	})
	distributionList = append(distributionList, DistributionStatus{
		OmfComponentID: "multicloud-k8s-id",
		Timestamp:      "1574774736174",
		URL:            "/sdc/v1/catalog/services/Test12/1.0/artifacts/service-Test12-csar.csar",
		Status:         "NOTIFIED",
		ErrorReason:    "null",
	})
	distributionList = append(distributionList, DistributionStatus{
		OmfComponentID: "policy-id",
		Timestamp:      "1574774728667",
		URL:            "/sdc/v1/catalog/services/Test12/1.0/artifacts/service-Test12-csar.csar",
		Status:         "NOTIFIED",
		ErrorReason:    "null",
	})
	distributionList = append(distributionList, DistributionStatus{
		OmfComponentID: "policy-id",
		Timestamp:      "1574774737376",
		URL:            "/sdc/v1/catalog/services/Test12/1.0/artifacts/service-Test12-csar.csar",
		Status:         "DOWNLOAD_OK",
		ErrorReason:    "null",
	})
	distributionList = append(distributionList, DistributionStatus{
		OmfComponentID: "policy-id",
		Timestamp:      "1574774728667",
		URL:            "/sdc/v1/catalog/services/Test12/1.0/resourceInstances/vsrx0/artifacts/vsrx0_modules.json",
		Status:         "NOT_NOTIFIED",
		ErrorReason:    "null",
	})
	distributionList = append(distributionList, DistributionStatus{
		OmfComponentID: "sdc-COpenSource-Env11-sdnc-dockero",
		Timestamp:      "1574774720318",
		URL:            "/sdc/v1/catalog/services/Test12/1.0/artifacts/service-Test12-template.yml",
		Status:         "NOT_NOTIFIED",
		ErrorReason:    "null",
	})
	distributionList = append(distributionList, DistributionStatus{
		OmfComponentID: "SO-COpenSource-Env11",
		Timestamp:      "1574774754939",
		URL:            "/sdc/v1/catalog/services/Test12/1.0/resourceInstances/vsrx0/artifacts/vsrx0_modules.json",
		Status:         "DEPLOY_OK",
		ErrorReason:    "null",
	})
	distributionList = append(distributionList, DistributionStatus{
		OmfComponentID: "multicloud-windriver-id",
		Timestamp:      "1574774731805",
		URL:            "/sdc/v1/catalog/services/Test12/1.0/resourceInstances/vsrx0/artifacts/vf-license-model.xml",
		Status:         "NOT_NOTIFIED",
		ErrorReason:    "null",
	})
	distributionList = append(distributionList, DistributionStatus{
		OmfComponentID: "sdc-COpenSource-Env11-sdnc-dockero",
		Timestamp:      "1574774720318",
		URL:            "/sdc/v1/catalog/services/Test12/1.0/artifacts/service-Test12-csar.csar",
		Status:         "NOTIFIED",
		ErrorReason:    "null",
	})
	distributionList = append(distributionList, DistributionStatus{
		OmfComponentID: "sdc-COpenSource-Env11-sdnc-dockero",
		Timestamp:      "1574774720318",
		URL:            "/sdc/v1/catalog/services/Test12/1.0/resourceInstances/vsrx0/artifacts/base_ubuntu16.yaml",
		Status:         "NOT_NOTIFIED",
		ErrorReason:    "null",
	})
	distributionList = append(distributionList, DistributionStatus{
		OmfComponentID: "multicloud-k8s-id",
		Timestamp:      "1574774748190",
		URL:            "/sdc/v1/catalog/services/Test12/1.0/resourceInstances/vsrx0/artifacts/base_ubuntu16.env",
		Status:         "DOWNLOAD_OK",
		ErrorReason:    "null",
	})
	distributionList = append(distributionList, DistributionStatus{
		OmfComponentID: "policy-id",
		Timestamp:      "1574774728667",
		URL:            "/sdc/v1/catalog/services/Test12/1.0/resourceInstances/vsrx0/artifacts/vendor-license-model.xml",
		Status:         "NOT_NOTIFIED",
		ErrorReason:    "null",
	})
	distributionList = append(distributionList, DistributionStatus{
		OmfComponentID: "multicloud-windriver-id",
		Timestamp:      "1574774742014",
		URL:            "/sdc/v1/catalog/services/Test12/1.0/resourceInstances/vsrx0/artifacts/base_ubuntu16.yaml",
		Status:         "DOWNLOAD_OK",
		ErrorReason:    "null",
	})
	distributionList = append(distributionList, DistributionStatus{
		OmfComponentID: "sdc-COpenSource-Env11-sdnc-dockero",
		Timestamp:      "1574774720318",
		URL:            "/sdc/v1/catalog/services/Test12/1.0/resourceInstances/vsrx0/artifacts/base_ubuntu16.env",
		Status:         "NOT_NOTIFIED",
		ErrorReason:    "null",
	})
	distributionList = append(distributionList, DistributionStatus{
		OmfComponentID: "policy-id",
		Timestamp:      "1574774728667",
		URL:            "/sdc/v1/catalog/services/Test12/1.0/resourceInstances/vsrx0/artifacts/base_ubuntu16.env",
		Status:         "NOT_NOTIFIED",
		ErrorReason:    "null",
	})
	distributionList = append(distributionList, DistributionStatus{
		OmfComponentID: "multicloud-windriver-id",
		Timestamp:      "1574774745715",
		URL:            "/sdc/v1/catalog/services/Test12/1.0/resourceInstances/vsrx0/artifacts/vsrx0_modules.json",
		Status:         "DEPLOY_OK",
		ErrorReason:    "null",
	})
	distributionList = append(distributionList, DistributionStatus{
		OmfComponentID: "SO-COpenSource-Env11",
		Timestamp:      "1574774740421",
		URL:            "/sdc/v1/catalog/services/Test12/1.0/resourceInstances/vsrx0/artifacts/vsrx0_modules.json",
		Status:         "NOTIFIED",
		ErrorReason:    "null",
	})
	distributionList = append(distributionList, DistributionStatus{
		OmfComponentID: "dcae-sch",
		Timestamp:      "1574774724083",
		URL:            "/sdc/v1/catalog/services/Test12/1.0/resourceInstances/vsrx0/artifacts/vendor-license-model.xml",
		Status:         "NOT_NOTIFIED",
		ErrorReason:    "null",
	})
	distributionList = append(distributionList, DistributionStatus{
		OmfComponentID: "multicloud-k8s-id",
		Timestamp:      "1574774736174",
		URL:            "/sdc/v1/catalog/services/Test12/1.0/resourceInstances/vsrx0/artifacts/base_ubuntu16.yaml",
		Status:         "NOTIFIED",
		ErrorReason:    "null",
	})
	distributionList = append(distributionList, DistributionStatus{
		OmfComponentID: "sdc-COpenSource-Env11-sdnc-dockero",
		Timestamp:      "1574774738399",
		URL:            "",
		Status:         "COMPONENT_DONE_OK",
		ErrorReason:    "null",
	})
	distributionList = append(distributionList, DistributionStatus{
		OmfComponentID: "multicloud-windriver-id",
		Timestamp:      "1574774740487",
		URL:            "/sdc/v1/catalog/services/Test12/1.0/resourceInstances/vsrx0/artifacts/vsrx0_modules.json",
		Status:         "DOWNLOAD_OK",
		ErrorReason:    "null",
	})
	distributionList = append(distributionList, DistributionStatus{
		OmfComponentID: "multicloud-windriver-id",
		Timestamp:      "1574774731805",
		URL:            "/sdc/v1/catalog/services/Test12/1.0/resourceInstances/vsrx0/artifacts/vendor-license-model.xml",
		Status:         "NOT_NOTIFIED",
		ErrorReason:    "null",
	})
	distributionList = append(distributionList, DistributionStatus{
		OmfComponentID: "policy-id",
		Timestamp:      "1574774740598",
		URL:            "",
		Status:         "COMPONENT_DONE_OK",
		ErrorReason:    "null",
	})
	distributionList = append(distributionList, DistributionStatus{
		OmfComponentID: "clamp",
		Timestamp:      "1574774737925",
		URL:            "/sdc/v1/catalog/services/Test12/1.0/resourceInstances/vsrx0/artifacts/vsrx0_modules.json",
		Status:         "NOTIFIED",
		ErrorReason:    "null",
	})
	distributionList = append(distributionList, DistributionStatus{
		OmfComponentID: "aai-ml",
		Timestamp:      "1574774737842",
		URL:            "/sdc/v1/catalog/services/Test12/1.0/resourceInstances/vsrx0/artifacts/base_ubuntu16.env",
		Status:         "NOT_NOTIFIED",
		ErrorReason:    "null",
	})
	distributionList = append(distributionList, DistributionStatus{
		OmfComponentID: "clamp",
		Timestamp:      "1574774737925",
		URL:            "/sdc/v1/catalog/services/Test12/1.0/resourceInstances/vsrx0/artifacts/base_ubuntu16.env",
		Status:         "NOT_NOTIFIED",
		ErrorReason:    "null",
	})
	distributionList = append(distributionList, DistributionStatus{
		OmfComponentID: "multicloud-starlingx-id",
		Timestamp:      "1574774737784",
		URL:            "/sdc/v1/catalog/services/Test12/1.0/artifacts/service-Test12-template.yml",
		Status:         "NOT_NOTIFIED",
		ErrorReason:    "null",
	})
	distributionList = append(distributionList, DistributionStatus{
		OmfComponentID: "multicloud-starlingx-id",
		Timestamp:      "1574774749858",
		URL:            "/sdc/v1/catalog/services/Test12/1.0/resourceInstances/vsrx0/artifacts/base_ubuntu16.env",
		Status:         "DOWNLOAD_OK",
		ErrorReason:    "null",
	})
	distributionList = append(distributionList, DistributionStatus{
		OmfComponentID: "multicloud-starlingx-id",
		Timestamp:      "1574774748343",
		URL:            "/sdc/v1/catalog/services/Test12/1.0/resourceInstances/vsrx0/artifacts/base_ubuntu16.yaml",
		Status:         "DOWNLOAD_OK",
		ErrorReason:    "null",
	})
	distributionList = append(distributionList, DistributionStatus{
		OmfComponentID: "aai-ml",
		Timestamp:      "1574774751522",
		URL:            "",
		Status:         "COMPONENT_DONE_OK",
		ErrorReason:    "null",
	})
	distributionList = append(distributionList, DistributionStatus{
		OmfComponentID: "cds",
		Timestamp:      "1574774726254",
		URL:            "/sdc/v1/catalog/services/Test12/1.0/resourceInstances/vsrx0/artifacts/base_ubuntu16.env",
		Status:         "NOT_NOTIFIED",
		ErrorReason:    "null",
	})
	distributionList = append(distributionList, DistributionStatus{
		OmfComponentID: "sdc-COpenSource-Env11-sdnc-dockero",
		Timestamp:      "1574774720318",
		URL:            "/sdc/v1/catalog/services/Test12/1.0/resourceInstances/vsrx0/artifacts/vendor-license-model.xml",
		Status:         "NOT_NOTIFIED",
		ErrorReason:    "null",
	})
	distributionList = append(distributionList, DistributionStatus{
		OmfComponentID: "SO-COpenSource-Env11",
		Timestamp:      "1574774740421",
		URL:            "/sdc/v1/catalog/services/Test12/1.0/resourceInstances/vsrx0/artifacts/base_ubuntu16.env",
		Status:         "NOTIFIED",
		ErrorReason:    "null",
	})
	distributionList = append(distributionList, DistributionStatus{
		OmfComponentID: "policy-id",
		Timestamp:      "1574774739512",
		URL:            "/sdc/v1/catalog/services/Test12/1.0/artifacts/service-Test12-csar.csar",
		Status:         "DEPLOY_OK",
		ErrorReason:    "null",
	})
	distributionList = append(distributionList, DistributionStatus{
		OmfComponentID: "multicloud-starlingx-id",
		Timestamp:      "1574774737784",
		URL:            "/sdc/v1/catalog/services/Test12/1.0/resourceInstances/vsrx0/artifacts/vendor-license-model.xml",
		Status:         "NOT_NOTIFIED",
		ErrorReason:    "null",
	})
	distributionList = append(distributionList, DistributionStatus{
		OmfComponentID: "multicloud-starlingx-id",
		Timestamp:      "1574774746621",
		URL:            "/sdc/v1/catalog/services/Test12/1.0/resourceInstances/vsrx0/artifacts/vsrx0_modules.json",
		Status:         "DOWNLOAD_OK",
		ErrorReason:    "null",
	})
	distributionList = append(distributionList, DistributionStatus{
		OmfComponentID: "clamp",
		Timestamp:      "1574774751028",
		URL:            "/sdc/v1/catalog/services/Test12/1.0/artifacts/service-Test12-csar.csar",
		Status:         "ALREADY_DEPLOYED",
		ErrorReason:    "null",
	})
	distributionList = append(distributionList, DistributionStatus{
		OmfComponentID: "SO-COpenSource-Env11",
		Timestamp:      "1574774740421",
		URL:            "/sdc/v1/catalog/services/Test12/1.0/resourceInstances/vsrx0/artifacts/vendor-license-model.xml",
		Status:         "NOT_NOTIFIED",
		ErrorReason:    "null",
	})
	distributionList = append(distributionList, DistributionStatus{
		OmfComponentID: "multicloud-starlingx-id",
		Timestamp:      "1574774737784",
		URL:            "/sdc/v1/catalog/services/Test12/1.0/resourceInstances/vsrx0/artifacts/base_ubuntu16.yaml",
		Status:         "NOTIFIED",
		ErrorReason:    "null",
	})
	distributionList = append(distributionList, DistributionStatus{
		OmfComponentID: "sdc-COpenSource-Env11-sdnc-dockero",
		Timestamp:      "1574774720318",
		URL:            "/sdc/v1/catalog/services/Test12/1.0/resourceInstances/vsrx0/artifacts/vf-license-model.xml",
		Status:         "NOT_NOTIFIED",
		ErrorReason:    "null",
	})
	distributionList = append(distributionList, DistributionStatus{
		OmfComponentID: "policy-id",
		Timestamp:      "1574774728667",
		URL:            "/sdc/v1/catalog/services/Test12/1.0/resourceInstances/vsrx0/artifacts/vf-license-model.xml",
		Status:         "NOT_NOTIFIED",
		ErrorReason:    "null",
	})
	distributionList = append(distributionList, DistributionStatus{
		OmfComponentID: "aai-ml",
		Timestamp:      "1574774737842",
		URL:            "/sdc/v1/catalog/services/Test12/1.0/resourceInstances/vsrx0/artifacts/vsrx0_modules.json",
		Status:         "NOT_NOTIFIED",
		ErrorReason:    "null",
	})
	distributionList = append(distributionList, DistributionStatus{
		OmfComponentID: "SO-COpenSource-Env11",
		Timestamp:      "1574774751123",
		URL:            "/sdc/v1/catalog/services/Test12/1.0/resourceInstances/vsrx0/artifacts/vsrx0_modules.json",
		Status:         "DOWNLOAD_OK",
		ErrorReason:    "null",
	})
	distributionList = append(distributionList, DistributionStatus{
		OmfComponentID: "aai-ml",
		Timestamp:      "1574774737842",
		URL:            "/sdc/v1/catalog/services/Test12/1.0/resourceInstances/vsrx0/artifacts/vendor-license-model.xml",
		Status:         "NOT_NOTIFIED",
		ErrorReason:    "null",
	})
	distributionList = append(distributionList, DistributionStatus{
		OmfComponentID: "dcae-sch",
		Timestamp:      "1574774724083",
		URL:            "/sdc/v1/catalog/services/Test12/1.0/artifacts/service-Test12-template.yml",
		Status:         "NOT_NOTIFIED",
		ErrorReason:    "null",
	})
	distributionList = append(distributionList, DistributionStatus{
		OmfComponentID: "multicloud-starlingx-id",
		Timestamp:      "1574774752037",
		URL:            "/sdc/v1/catalog/services/Test12/1.0/resourceInstances/vsrx0/artifacts/vsrx0_modules.json",
		Status:         "DEPLOY_OK",
		ErrorReason:    "null",
	})
	distributionList = append(distributionList, DistributionStatus{
		OmfComponentID: "SO-COpenSource-Env11",
		Timestamp:      "1574774755942",
		URL:            "/sdc/v1/catalog/services/Test12/1.0/resourceInstances/vsrx0/artifacts/base_ubuntu16.yaml",
		Status:         "DEPLOY_OK",
		ErrorReason:    "null",
	})
	distributionList = append(distributionList, DistributionStatus{
		OmfComponentID: "clamp",
		Timestamp:      "1574774737925",
		URL:            "/sdc/v1/catalog/services/Test12/1.0/resourceInstances/vsrx0/artifacts/base_ubuntu16.yaml",
		Status:         "NOT_NOTIFIED",
		ErrorReason:    "null",
	})
	distributionList = append(distributionList, DistributionStatus{
		OmfComponentID: "multicloud-windriver-id",
		Timestamp:      "1574774743466",
		URL:            "/sdc/v1/catalog/services/Test12/1.0/resourceInstances/vsrx0/artifacts/base_ubuntu16.env",
		Status:         "DOWNLOAD_OK",
		ErrorReason:    "null",
	})
}

func getArtifactTypes(c echo.Context) error {
	list := []string{"HEAT"}
	return c.JSON(http.StatusOK, list)
}

func registerForDistribution(c echo.Context) error {
	distributionRegistration := map[string]string{
		"distrNotificationTopicName":"testName",
		"distrStatusTopicName":"testTopic",
	}
	return c.JSON(http.StatusOK, distributionRegistration)
}

func unRegisterForDistribution(c echo.Context) error {
	return c.JSON(http.StatusOK, list.New())
}

func distributionKafkaData(c echo.Context) error {
	kafkaData := map[string]string{
	    "kafkaBootStrapServer":"localhost:43219",
		"distrNotificationTopicName":"SDC-DIST-NOTIF-TOPIC",
		"distrStatusTopicName":"SDC-DIST-STATUS-TOPIC",
	}
	return c.JSON(http.StatusOK, kafkaData)
}

