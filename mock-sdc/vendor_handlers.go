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
	"time"

	"github.com/labstack/echo"
	uuid "github.com/satori/go.uuid"
)

// Vendor describes license model in SDC
type Vendor struct {
	ID          string `json:"id"`
	Type        string `json:"type"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Owner       string `json:"owner"`
	Status      string `json:"status"`
	Properties  struct {
	} `json:"properties"`
	Versions []Version `json:"-"`
}

// NewVendor describe the vendor creation model in SDC
type NewVendor struct {
	IconRef     string `json:"iconRef"`
	VendorName  string `json:"vendorName"`
	Description string `json:"description"`
}

// VendorList is the way to return Vendors in SDC
type VendorList struct {
	ListCount int      `json:"listCount"`
	Results   []Vendor `json:"results"`
}

var vendorList []Vendor

func generateInitialVendorList() {
	vendorList = nil
	var empty struct{}
	version1 := Version{
		ID:               "61c134e128f54119934b3960c77a3f33",
		Name:             "1.0",
		Description:      "Initial version",
		BaseID:           "",
		Status:           "Certified",
		RealStatus:       "Certified",
		CreationTime:     1559565688604,
		ModificationTime: 1559565787436,
		AdditionalInfo: struct {
			OptionalCreationMethods []string `json:"OptionalCreationMethods"`
		}{
			OptionalCreationMethods: []string{"major"}},
	}
	version2 := Version{
		ID:               "2e3ba48c748d47e3bd4afdd8348bdfb9",
		Name:             "1.0",
		Description:      "Initial version",
		BaseID:           "",
		Status:           "Certified",
		RealStatus:       "Certified",
		CreationTime:     1559562354868,
		ModificationTime: 1559562421476,
		AdditionalInfo: struct {
			OptionalCreationMethods []string `json:"OptionalCreationMethods"`
		}{
			OptionalCreationMethods: []string{"major"}},
	}
	vendorList = append(vendorList,
		Vendor{
			ID:          "212a52b2630749388a7693086ac1467e",
			Type:        "vlm",
			Name:        "wvfw",
			Description: "wvfw",
			Owner:       "cs0008",
			Status:      "ACTIVE",
			Properties:  empty,
			Versions:    []Version{version1}})
	vendorList = append(vendorList, Vendor{
		ID:          "e78eb0b1c73e43138f705cd92c0e4ace",
		Type:        "vlm",
		Name:        "vfw_test",
		Description: "test vfw",
		Owner:       "cs0008",
		Status:      "ACTIVE",
		Properties:  empty,
		Versions:    []Version{version2}})
}

func getVendorServiceModels(c echo.Context) error {
	list := &VendorList{len(vendorList), vendorList}
	return c.JSON(http.StatusOK, list)
}

func postVendorServiceModels(c echo.Context) error {
	newVendor := new(NewVendor)
	if err := c.Bind(newVendor); err != nil {
		return err
	}

	u1 := uuid.NewV4().String()
	version := Version{
		ID:               u1,
		Name:             "1.0",
		Description:      "Initial version",
		BaseID:           "",
		Status:           "Draft",
		RealStatus:       "Draft",
		CreationTime:     (time.Now().UnixNano() / 1000000),
		ModificationTime: (time.Now().UnixNano() / 1000000),
		AdditionalInfo: struct {
			OptionalCreationMethods []string `json:"OptionalCreationMethods"`
		}{
			OptionalCreationMethods: []string{"major"}},
		State: struct {
			SynchronizationState string `json:"synchronizationState"`
			Dirty                bool   `json:"dirty"`
		}{
			SynchronizationState: "UpToDate",
			Dirty:                false,
		},
	}
	u2 := uuid.NewV4().String()
	var empty struct{}
	vendorList = append(vendorList, Vendor{
		ID:          u2,
		Type:        "vlm",
		Name:        newVendor.VendorName,
		Description: newVendor.Description,
		Owner:       "cs0008",
		Status:      "ACTIVE",
		Properties:  empty,
		Versions:    []Version{version}})

	createdVendor := CreatedItem{
		ItemID: u2,
		Version: CreatedVersion{
			ID:          u1,
			Name:        "1.0",
			Description: "Initial version",
			Status:      "Draft",
		},
	}
	return c.JSON(http.StatusCreated, createdVendor)
}

func updateVendorVersion(c echo.Context) error {
	action := new(Action)
	if err := c.Bind(action); err != nil {
		return err
	}
	if action.Action == "Submit" {
		vendorID := c.Param("vendorID")
		versionID := c.Param("versionID")
		for i, vendor := range vendorList {
			if vendor.ID == vendorID {
				for j, version := range vendor.Versions {
					if version.ID == versionID {
						vendorList[i].Versions[j].Status = "Certified"
						vendorList[i].Versions[j].RealStatus = "Certified"
					}
					return c.String(http.StatusOK, "{}")
				}
				return echo.NewHTTPError(http.StatusNotFound, "Version Not Found")
			}
		}
		return echo.NewHTTPError(http.StatusNotFound, "Vendor Not Found")
	}
	return echo.NewHTTPError(http.StatusNotFound, "Unknown Action")
}
