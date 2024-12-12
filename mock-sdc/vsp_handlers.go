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
	"strings"
	"time"

	"github.com/labstack/echo"
	uuid "github.com/satori/go.uuid"
)

// Vsp describes software product in SDC
type Vsp struct {
	ID                        string `json:"id"`
	Icon                      string `json:"icon"`
	OnboardingMethod          string `json:"onboardingMethod"`
	Name                      string `json:"name"`
	Description               string `json:"description"`
	Owner                     string `json:"owner"`
	Status                    string `json:"status"`
	VendorName                string `json:"vendorName"`
	VendorID                  string `json:"vendorId"`
	Category                  string `json:"category"`
	SubCategory               string `json:"subCategory"`
	CandidateOnboardingOrigin string `json:"candidateOnboardingOrigin"`
	OnboardingOrigin          string `json:"onboardingOrigin"`
	NetworkPackageName        string `json:"networkPackageName"`
	ValidationData            struct {
		ImportStructure struct {
			Heat string `json:"heat"`
		} `json:"importStructure"`
	} `json:"validationData"`
	Versions []Version `json:"-"`
}

// VspLight describes software product in SDC lists
type VspLight struct {
	ID               string `json:"id"`
	OnboardingMethod string `json:"onboardingMethod"`
	Name             string `json:"name"`
	Description      string `json:"description"`
	Owner            string `json:"owner"`
	Status           string `json:"status"`
	VendorName       string `json:"vendorName"`
	VendorID         string `json:"vendorId"`
}

// VspList is the way to return Vsps in SDC
type VspList struct {
	ListCount int        `json:"listCount"`
	Results   []VspLight `json:"results"`
}

// VspDetailsDraft describes software product in SDC
type VspDetailsDraft struct {
	ID               string `json:"id"`
	Icon             string `json:"icon"`
	OnboardingMethod string `json:"onboardingMethod"`
	Name             string `json:"name"`
	Description      string `json:"description"`
	VendorName       string `json:"vendorName"`
	VendorID         string `json:"vendorId"`
	Version          string `json:"version"`
	Category         string `json:"category"`
	SubCategory      string `json:"subCategory"`
}

// VspDetailsUploaded describes software product in SDC
type VspDetailsUploaded struct {
	ID                        string `json:"id"`
	Icon                      string `json:"icon"`
	OnboardingMethod          string `json:"onboardingMethod"`
	Name                      string `json:"name"`
	Description               string `json:"description"`
	VendorName                string `json:"vendorName"`
	VendorID                  string `json:"vendorId"`
	Version                   string `json:"version"`
	Category                  string `json:"category"`
	SubCategory               string `json:"subCategory"`
	CandidateOnboardingOrigin string `json:"candidateOnboardingOrigin"`
	NetworkPackageName        string `json:"networkPackageName"`
}

// VspDetailsValidated describes software product in SDC
type VspDetailsValidated struct {
	ID                 string `json:"id"`
	Icon               string `json:"icon"`
	OnboardingMethod   string `json:"onboardingMethod"`
	Name               string `json:"name"`
	Description        string `json:"description"`
	VendorName         string `json:"vendorName"`
	VendorID           string `json:"vendorId"`
	Version            string `json:"version"`
	Category           string `json:"category"`
	SubCategory        string `json:"subCategory"`
	OnboardingOrigin   string `json:"onboardingOrigin"`
	NetworkPackageName string `json:"networkPackageName"`
	ValidationData     struct {
		ImportStructure struct {
			Heat string `json:"heat"`
		} `json:"importStructure"`
	} `json:"validationData"`
}

// NewVsp describe the vsp creation model in SDC
type NewVsp struct {
	Icon             string   `json:"iconRef"`
	Name             string   `json:"name"`
	VendorName       string   `json:"vendorName"`
	VendorID         string   `json:"vendorId"`
	Description      string   `json:"description"`
	Category         string   `json:"category"`
	SubCategory      string   `json:"subCategory"`
	LicensingData    struct{} `json:"licensingData"`
	OnboardingMethod string   `json:"onboardingMethod"`
}

// ArtifactUploadResult bla
type ArtifactUploadResult struct {
	Errors             struct{} `json:"errors"`
	Status             string   `json:"status"`
	OnboardingOrigin   string   `json:"onboardingOrigin"`
	NetworkPackageName string   `json:"networkPackageName"`
}

// ArtifactValidationResult bla
type ArtifactValidationResult struct {
	Errors    struct{} `json:"errors"`
	Status    string   `json:"status"`
	FileNames []string `json:"fileNames"`
}

// CsarCreateResult bla
type CsarCreateResult struct {
	Description   string `json:"description"`
	VspName       string `json:"vspName"`
	Version       string `json:"version"`
	PackageID     string `json:"packageId"`
	Category      string `json:"category"`
	SubCategory   string `json:"subCategory"`
	VendorName    string `json:"vendorName"`
	VendorRelease string `json:"vendorRelease"`
	PackageType   string `json:"packageType"`
	ResourceType  string `json:"resourceType"`
}

var vspList []Vsp

func generateInitialVspList() {
	vspList = []Vsp{}
}

func getVendorSoftwareProducts(c echo.Context) error {
	vspLights := []VspLight{}
	for _, v := range vspList {
		vspLights = append(vspLights, VspLight{
			ID:               v.ID,
			OnboardingMethod: v.OnboardingMethod,
			Name:             v.Name,
			Description:      v.Description,
			Owner:            v.Owner,
			Status:           v.Status,
			VendorName:       v.VendorName,
			VendorID:         v.VendorID,
		})
	}
	list := &VspList{len(vspLights), vspLights}
	return c.JSON(http.StatusOK, list)
}

func postVendorSoftwareProducts(c echo.Context) error {
	newVsp := new(NewVsp)
	if err := c.Bind(newVsp); err != nil {
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
	vspList = append(vspList, Vsp{
		ID:               u2,
		OnboardingMethod: "NetworkPackage",
		Name:             newVsp.Name,
		Description:      newVsp.Description,
		Owner:            "cs0008",
		Status:           "ACTIVE",
		VendorName:       newVsp.VendorName,
		VendorID:         newVsp.VendorID,
		Category:         "resourceNewCategory.generic",
		SubCategory:      "resourceNewCategory.generic.abstract",
		Icon:             "icon",
		Versions:         []Version{version}})

	createdVsp := CreatedItem{
		ItemID: u2,
		Version: CreatedVersion{
			ID:          u1,
			Name:        "1.0",
			Description: "Initial version",
			Status:      "Draft",
		},
	}
	return c.JSON(http.StatusCreated, createdVsp)
}

func uploadArtifacts(c echo.Context) error {
	vspID := c.Param("vspID")
	versionID := c.Param("versionID")
	for i, v := range vspList {
		if v.ID == vspID {
			for j, version := range v.Versions {
				if version.ID == versionID {
					if version.RealStatus == "Draft" {
						file, err := c.FormFile("upload")
						if err != nil {
							return err
						}
						fileName := strings.Split(file.Filename, ".")[0]
						fileExtension := strings.Split(file.Filename, ".")[1]
						vspList[i].NetworkPackageName = fileName
						vspList[i].CandidateOnboardingOrigin = fileExtension
						vspList[i].Versions[j].RealStatus = "Uploaded"
						var empty struct{}
						artifactUploadResult := ArtifactUploadResult{
							Errors:             empty,
							Status:             "Success",
							OnboardingOrigin:   fileExtension,
							NetworkPackageName: fileName,
						}
						return c.JSON(http.StatusCreated, artifactUploadResult)
					}
				}
			}
		}
	}
	return echo.NewHTTPError(http.StatusNotFound, "Item Not Found")
}

func validateArtifacts(c echo.Context) error {
	vspID := c.Param("vspID")
	versionID := c.Param("versionID")
	for i, v := range vspList {
		if v.ID == vspID {
			for j, version := range v.Versions {
				if version.ID == versionID {
					if version.RealStatus == "Uploaded" {
						vspList[i].OnboardingOrigin = vspList[i].CandidateOnboardingOrigin
						vspList[i].Versions[j].State.Dirty = true
						vspList[i].ValidationData = struct {
							ImportStructure struct {
								Heat string `json:"heat"`
							} `json:"importStructure"`
						}{
							ImportStructure: struct {
								Heat string `json:"heat"`
							}{
								Heat: "Yes",
							},
						}
						vspList[i].Versions[j].RealStatus = "Validated"
						var empty struct{}
						artifactValidationResult := ArtifactValidationResult{
							Errors: empty,
							FileNames: []string{
								"base_ubuntu16.env",
								"base_ubuntu16.yaml",
							},
							Status: "Success",
						}
						return c.JSON(http.StatusOK, artifactValidationResult)
					}
				}
			}
		}
	}
	return echo.NewHTTPError(http.StatusNotFound, "Item Not Found")
}

func updateVspVersion(c echo.Context) error {
	action := new(Action)
	if err := c.Bind(action); err != nil {
		return err
	}
	vspID := c.Param("vspID")
	versionID := c.Param("versionID")
	for i, v := range vspList {
		if v.ID == vspID {
			for j, version := range v.Versions {
				if version.ID == versionID {
					if action.Action == "Submit" {
						if version.RealStatus == "Commited" {
							vspList[i].Versions[j].RealStatus = "Certified"
							vspList[i].Versions[j].Status = "Certified"
							return c.String(http.StatusOK, "{}")
						}
						return echo.NewHTTPError(http.StatusNotFound, "Item not in good state")
					}
					if action.Action == "Create_Package" {
						if version.RealStatus == "Certified" {
							csarCreateResult := CsarCreateResult{
								Description:   v.Description,
								VspName:       v.Name,
								Version:       version.Name,
								PackageID:     v.ID,
								Category:      v.Category,
								SubCategory:   v.SubCategory,
								VendorName:    v.VendorName,
								VendorRelease: "1.0",
								PackageType:   "CSAR",
								ResourceType:  "VF",
							}
							return c.JSON(http.StatusOK, csarCreateResult)
						}
						return echo.NewHTTPError(http.StatusNotFound, "Item not in good state")
					}
					return echo.NewHTTPError(http.StatusNotFound, "Unknown Action")
				}
			}
		}
	}
	return echo.NewHTTPError(http.StatusNotFound, "Item Not Found")
}

func getVspVersion(c echo.Context) error {
	vspID := c.Param("vspID")
	versionID := c.Param("versionID")
	for _, v := range vspList {
		if v.ID == vspID {
			for _, version := range v.Versions {
				if version.ID == versionID {
					if version.RealStatus == "Draft" {
						vspDetails := VspDetailsDraft{
							ID:               v.ID,
							Icon:             v.Icon,
							OnboardingMethod: v.OnboardingMethod,
							Name:             v.Name,
							Description:      v.Description,
							VendorName:       v.VendorName,
							VendorID:         v.VendorID,
							Version:          version.ID,
							Category:         v.Category,
							SubCategory:      v.SubCategory,
						}
						return c.JSON(http.StatusOK, vspDetails)
					}
					if version.RealStatus == "Uploaded" {
						vspDetails := VspDetailsUploaded{
							ID:                        v.ID,
							Icon:                      v.Icon,
							OnboardingMethod:          v.OnboardingMethod,
							Name:                      v.Name,
							Description:               v.Description,
							VendorName:                v.VendorName,
							VendorID:                  v.VendorID,
							Version:                   version.ID,
							Category:                  v.Category,
							SubCategory:               v.SubCategory,
							CandidateOnboardingOrigin: v.CandidateOnboardingOrigin,
							NetworkPackageName:        v.NetworkPackageName,
						}
						return c.JSON(http.StatusOK, vspDetails)
					}
					vspDetails := VspDetailsValidated{
						ID:                 v.ID,
						Icon:               v.Icon,
						OnboardingMethod:   v.OnboardingMethod,
						Name:               v.Name,
						Description:        v.Description,
						VendorName:         v.VendorName,
						VendorID:           v.VendorID,
						Version:            version.ID,
						Category:           v.Category,
						SubCategory:        v.SubCategory,
						OnboardingOrigin:   v.OnboardingOrigin,
						NetworkPackageName: v.NetworkPackageName,
						ValidationData:     v.ValidationData,
					}
					return c.JSON(http.StatusOK, vspDetails)
				}
			}
		}
	}
	return echo.NewHTTPError(http.StatusNotFound, "Item Not Found")
}
