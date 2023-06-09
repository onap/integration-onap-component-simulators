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
)

// VersionList describe the list return in SDC
type VersionList struct {
	ListCount int            `json:"listCount"`
	Results   []VersionLight `json:"results"`
}

// Version describe version model in SDC
type Version struct {
	ID               string `json:"id"`
	Name             string `json:"name"`
	Description      string `json:"description"`
	BaseID           string `json:"baseId"`
	Status           string `json:"status"`
	RealStatus       string
	CreationTime     int64 `json:"creationTime"`
	ModificationTime int64 `json:"modificationTime"`
	AdditionalInfo   struct {
		OptionalCreationMethods []string `json:"OptionalCreationMethods"`
	} `json:"additionalInfo"`
	State struct {
		SynchronizationState string `json:"synchronizationState"`
		Dirty                bool   `json:"dirty"`
	} `json:"state"`
}

// VersionLight describe version model in SDC
type VersionLight struct {
	ID               string `json:"id"`
	Name             string `json:"name"`
	Description      string `json:"description"`
	BaseID           string `json:"baseId"`
	Status           string `json:"status"`
	CreationTime     int64  `json:"creationTime"`
	ModificationTime int64  `json:"modificationTime"`
	AdditionalInfo   struct {
		OptionalCreationMethods []string `json:"OptionalCreationMethods"`
	} `json:"additionalInfo"`
}

// VersionDetails show details of a version
type VersionDetails struct {
	ID               string `json:"id"`
	Name             string `json:"name"`
	Description      string `json:"description"`
	Status           string `json:"status"`
	CreationTime     int64  `json:"creationTime"`
	ModificationTime int64  `json:"modificationTime"`
	State            struct {
		SynchronizationState string `json:"synchronizationState"`
		Dirty                bool   `json:"dirty"`
	} `json:"state"`
}

// CreatedVersion of a CreatedVendor in SDC
type CreatedVersion struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Status      string `json:"status"`
}

// CreatedItem model in SDC
type CreatedItem struct {
	ItemID  string         `json:"itemId"`
	Version CreatedVersion `json:"version"`
}

// Action describe the action on items in SDC
type Action struct {
	Action string `json:"action"`
}

// SdcError is the way to return Error in SDC
type SdcError struct {
	Status    string `json:"status"`
	ErrorCode string `json:"errorCode"`
	Message   string `json:"message"`
}

func getItemVersions(c echo.Context) error {
	itemID := c.Param("itemID")
	for _, v := range vendorList {
		if v.ID == itemID {
			versionList := generateVersionList(v.Versions)
			list := &VersionList{len(versionList), versionList}
			return c.JSON(http.StatusOK, list)
		}
	}
	for _, v := range vspList {
		if v.ID == itemID {
			versionList := generateVersionList(v.Versions)
			list := &VersionList{len(versionList), versionList}
			return c.JSON(http.StatusOK, list)
		}
	}
	return echo.NewHTTPError(http.StatusNotFound, "Item Not Found")
}

func generateVersionList(versions []Version) []VersionLight {
	versionList := []VersionLight{}
	for _, version := range versions {
		versionList = append(versionList, VersionLight{
			ID:               version.ID,
			Name:             version.Name,
			Description:      version.Description,
			BaseID:           version.BaseID,
			Status:           version.Status,
			CreationTime:     version.CreationTime,
			ModificationTime: version.ModificationTime,
			AdditionalInfo:   version.AdditionalInfo,
		})

	}
	return versionList
}

func getItemVersion(c echo.Context) error {
	itemID := c.Param("itemID")
	versionID := c.Param("versionID")
	for _, v := range vendorList {
		if v.ID == itemID {
			for _, version := range v.Versions {
				if version.ID == versionID {
					versionDetails := VersionDetails{
						ID:               version.ID,
						Name:             version.Name,
						Description:      version.Description,
						Status:           version.Status,
						CreationTime:     version.CreationTime,
						ModificationTime: version.ModificationTime,
						State:            version.State,
					}
					return c.JSON(http.StatusOK, versionDetails)
				}
			}
		}
	}
	for _, v := range vspList {
		if v.ID == itemID {
			for _, version := range v.Versions {
				if version.ID == versionID {
					versionDetails := VersionDetails{
						ID:               version.ID,
						Name:             version.Name,
						Description:      version.Description,
						Status:           version.Status,
						CreationTime:     version.CreationTime,
						ModificationTime: version.ModificationTime,
						State:            version.State,
					}
					return c.JSON(http.StatusOK, versionDetails)
				}
			}
		}
	}
	return echo.NewHTTPError(http.StatusNotFound, "Item Not Found")
}

func updateItemVersion(c echo.Context) error {
	action := new(Action)
	if err := c.Bind(action); err != nil {
		return err
	}
	itemID := c.Param("itemID")
	versionID := c.Param("versionID")
	for i, v := range vspList {
		if v.ID == itemID {
			for j, version := range v.Versions {
				if version.ID == versionID {
					if action.Action == "Commit" {
						if version.RealStatus == "Validated" {
							vspList[i].Versions[j].RealStatus = "Commited"
							vspList[i].Versions[j].State.Dirty = false
							return c.String(http.StatusOK, "{}")
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

func getCategories(c echo.Context) error {
	return c.String(http.StatusOK, `{"categories":{"resourceCategories":[{"name":"Application L4+","normalizedName":"application l4+","uniqueId":"resourceNewCategory.application l4+","icons":null,"subcategories":[{"name":"Media Servers","normalizedName":"media servers","uniqueId":"resourceNewCategory.application l4+.media servers","icons":["applicationServer"],"groupings":null,"version":null,"ownerId":null,"empty":false,"type":null},{"name":"Database","normalizedName":"database","uniqueId":"resourceNewCategory.application l4+.database","icons":["database"],"groupings":null,"version":null,"ownerId":null,"empty":false,"type":null},{"name":"Border Element","normalizedName":"border element","uniqueId":"resourceNewCategory.application l4+.border element","icons":["borderElement"],"groupings":null,"version":null,"ownerId":null,"empty":false,"type":null},{"name":"Application Server","normalizedName":"application server","uniqueId":"resourceNewCategory.application l4+.application server","icons":["applicationServer"],"groupings":null,"version":null,"ownerId":null,"empty":false,"type":null},{"name":"Firewall","normalizedName":"firewall","uniqueId":"resourceNewCategory.application l4+.firewall","icons":["firewall"],"groupings":null,"version":null,"ownerId":null,"empty":false,"type":null},{"name":"Call Control","normalizedName":"call control","uniqueId":"resourceNewCategory.application l4+.call control","icons":["call_controll"],"groupings":null,"version":null,"ownerId":null,"empty":false,"type":null},{"name":"Web Server","normalizedName":"web server","uniqueId":"resourceNewCategory.application l4+.web server","icons":["applicationServer"],"groupings":null,"version":null,"ownerId":null,"empty":false,"type":null},{"name":"Load Balancer","normalizedName":"load balancer","uniqueId":"resourceNewCategory.application l4+.load balancer","icons":["loadBalancer"],"groupings":null,"version":null,"ownerId":null,"empty":false,"type":null}],"version":null,"ownerId":null,"empty":false,"type":null},{"name":"Network Connectivity","normalizedName":"network connectivity","uniqueId":"resourceNewCategory.network connectivity","icons":null,"subcategories":[{"name":"Connection Points","normalizedName":"connection points","uniqueId":"resourceNewCategory.network connectivity.connection points","icons":["cp"],"groupings":null,"version":null,"ownerId":null,"empty":false,"type":null},{"name":"Virtual Links","normalizedName":"virtual links","uniqueId":"resourceNewCategory.network connectivity.virtual links","icons":["vl"],"groupings":null,"version":null,"ownerId":null,"empty":false,"type":null}],"version":null,"ownerId":null,"empty":false,"type":null},{"name":"Allotted Resource","normalizedName":"allotted resource","uniqueId":"resourceNewCategory.allotted resource","icons":null,"subcategories":[{"name":"BRG","normalizedName":"brg","uniqueId":"resourceNewCategory.allotted resource.brg","icons":["brg"],"groupings":null,"version":null,"ownerId":null,"empty":false,"type":null},{"name":"TunnelXConn","normalizedName":"tunnelxconn","uniqueId":"resourceNewCategory.allotted resource.tunnelxconn","icons":["tunnel_x_connect"],"groupings":null,"version":null,"ownerId":null,"empty":false,"type":null},{"name":"IP Mux Demux","normalizedName":"ip mux demux","uniqueId":"resourceNewCategory.allotted resource.ip mux demux","icons":["ip_mux_demux"],"groupings":null,"version":null,"ownerId":null,"empty":false,"type":null},{"name":"Security Zone","normalizedName":"security zone","uniqueId":"resourceNewCategory.allotted resource.security zone","icons":["security_zone"],"groupings":null,"version":null,"ownerId":null,"empty":false,"type":null},{"name":"Service Admin","normalizedName":"service admin","uniqueId":"resourceNewCategory.allotted resource.service admin","icons":["service_admin"],"groupings":null,"version":null,"ownerId":null,"empty":false,"type":null},{"name":"Allotted Resource","normalizedName":"allotted resource","uniqueId":"resourceNewCategory.allotted resource.allotted resource","icons":["allotted_resource"],"groupings":null,"version":null,"ownerId":null,"empty":false,"type":null},{"name":"Contrail Route","normalizedName":"contrail route","uniqueId":"resourceNewCategory.allotted resource.contrail route","icons":["contrail_route"],"groupings":null,"version":null,"ownerId":null,"empty":false,"type":null}],"version":null,"ownerId":null,"empty":false,"type":null},{"name":"Configuration","normalizedName":"configuration","uniqueId":"resourceNewCategory.configuration","icons":null,"subcategories":[{"name":"Configuration","normalizedName":"configuration","uniqueId":"resourceNewCategory.configuration.configuration","icons":["pmc"],"groupings":null,"version":null,"ownerId":null,"empty":false,"type":null}],"version":null,"ownerId":null,"empty":false,"type":null},{"name":"Network L4+","normalizedName":"network l4+","uniqueId":"resourceNewCategory.network l4+","icons":null,"subcategories":[{"name":"Common Network Resources","normalizedName":"common network resources","uniqueId":"resourceNewCategory.network l4+.common network resources","icons":["network"],"groupings":null,"version":null,"ownerId":null,"empty":false,"type":null}],"version":null,"ownerId":null,"empty":false,"type":null},{"name":"Generic","normalizedName":"generic","uniqueId":"resourceNewCategory.generic","icons":null,"subcategories":[{"name":"Abstract","normalizedName":"abstract","uniqueId":"resourceNewCategory.generic.abstract","icons":["objectStorage","compute"],"groupings":null,"version":null,"ownerId":null,"empty":false,"type":null},{"name":"Network Service","normalizedName":"network service","uniqueId":"resourceNewCategory.generic.network service","icons":["network"],"groupings":null,"version":null,"ownerId":null,"empty":false,"type":null},{"name":"Rules","normalizedName":"rules","uniqueId":"resourceNewCategory.generic.rules","icons":["networkrules","securityrules"],"groupings":null,"version":null,"ownerId":null,"empty":false,"type":null},{"name":"Infrastructure","normalizedName":"infrastructure","uniqueId":"resourceNewCategory.generic.infrastructure","icons":["connector"],"groupings":null,"version":null,"ownerId":null,"empty":false,"type":null},{"name":"Network Elements","normalizedName":"network elements","uniqueId":"resourceNewCategory.generic.network elements","icons":["network","connector"],"groupings":null,"version":null,"ownerId":null,"empty":false,"type":null},{"name":"Database","normalizedName":"database","uniqueId":"resourceNewCategory.generic.database","icons":["database"],"groupings":null,"version":null,"ownerId":null,"empty":false,"type":null}],"version":null,"ownerId":null,"empty":false,"type":null},{"name":"DCAE Component","normalizedName":"dcae component","uniqueId":"resourceNewCategory.dcae component","icons":null,"subcategories":[{"name":"Analytics","normalizedName":"analytics","uniqueId":"resourceNewCategory.dcae component.analytics","icons":["dcae_analytics"],"groupings":null,"version":null,"ownerId":null,"empty":false,"type":null},{"name":"Database","normalizedName":"database","uniqueId":"resourceNewCategory.dcae component.database","icons":["dcae_database"],"groupings":null,"version":null,"ownerId":null,"empty":false,"type":null},{"name":"Policy","normalizedName":"policy","uniqueId":"resourceNewCategory.dcae component.policy","icons":["dcae_policy"],"groupings":null,"version":null,"ownerId":null,"empty":false,"type":null},{"name":"Machine Learning","normalizedName":"machine learning","uniqueId":"resourceNewCategory.dcae component.machine learning","icons":["dcae_machineLearning"],"groupings":null,"version":null,"ownerId":null,"empty":false,"type":null},{"name":"Microservice","normalizedName":"microservice","uniqueId":"resourceNewCategory.dcae component.microservice","icons":["dcae_microservice"],"groupings":null,"version":null,"ownerId":null,"empty":false,"type":null},{"name":"Source","normalizedName":"source","uniqueId":"resourceNewCategory.dcae component.source","icons":["dcae_source"],"groupings":null,"version":null,"ownerId":null,"empty":false,"type":null},{"name":"Collector","normalizedName":"collector","uniqueId":"resourceNewCategory.dcae component.collector","icons":["dcae_collector"],"groupings":null,"version":null,"ownerId":null,"empty":false,"type":null},{"name":"Utility","normalizedName":"utility","uniqueId":"resourceNewCategory.dcae component.utility","icons":["dcae_utilty"],"groupings":null,"version":null,"ownerId":null,"empty":false,"type":null}],"version":null,"ownerId":null,"empty":false,"type":null},{"name":"Network L2-3","normalizedName":"network l2-3","uniqueId":"resourceNewCategory.network l2-3","icons":null,"subcategories":[{"name":"Gateway","normalizedName":"gateway","uniqueId":"resourceNewCategory.network l2-3.gateway","icons":["gateway"],"groupings":null,"version":null,"ownerId":null,"empty":false,"type":null},{"name":"WAN Connectors","normalizedName":"wan connectors","uniqueId":"resourceNewCategory.network l2-3.wan connectors","icons":["network","connector","port"],"groupings":null,"version":null,"ownerId":null,"empty":false,"type":null},{"name":"Infrastructure","normalizedName":"infrastructure","uniqueId":"resourceNewCategory.network l2-3.infrastructure","icons":["ucpe"],"groupings":null,"version":null,"ownerId":null,"empty":false,"type":null},{"name":"Router","normalizedName":"router","uniqueId":"resourceNewCategory.network l2-3.router","icons":["router","vRouter"],"groupings":null,"version":null,"ownerId":null,"empty":false,"type":null},{"name":"LAN Connectors","normalizedName":"lan connectors","uniqueId":"resourceNewCategory.network l2-3.lan connectors","icons":["network","connector","port"],"groupings":null,"version":null,"ownerId":null,"empty":false,"type":null}],"version":null,"ownerId":null,"empty":false,"type":null},{"name":"Template","normalizedName":"template","uniqueId":"resourceNewCategory.template","icons":null,"subcategories":[{"name":"Base Monitoring Template","normalizedName":"base monitoring template","uniqueId":"resourceNewCategory.template.base monitoring template","icons":["monitoring_template"],"groupings":null,"version":null,"ownerId":null,"empty":false,"type":null},{"name":"Monitoring Template","normalizedName":"monitoring template","uniqueId":"resourceNewCategory.template.monitoring template","icons":["monitoring_template"],"groupings":null,"version":null,"ownerId":null,"empty":false,"type":null}],"version":null,"ownerId":null,"empty":false,"type":null}],"serviceCategories":[{"name":"Mobility","normalizedName":"mobility","uniqueId":"serviceNewCategory.mobility","icons":["mobility"],"subcategories":null,"version":null,"ownerId":null,"empty":false,"type":null},{"name":"Network L4+","normalizedName":"network l4+","uniqueId":"serviceNewCategory.network l4+","icons":["network_l_4"],"subcategories":null,"version":null,"ownerId":null,"empty":false,"type":null},{"name":"E2E Service","normalizedName":"e2e service","uniqueId":"serviceNewCategory.e2e service","icons":["network_l_1-3"],"subcategories":null,"version":null,"ownerId":null,"empty":false,"type":null},{"name":"VoIP Call Control","normalizedName":"voip call control","uniqueId":"serviceNewCategory.voip call control","icons":["call_controll"],"subcategories":null,"version":null,"ownerId":null,"empty":false,"type":null},{"name":"Network Service","normalizedName":"network service","uniqueId":"serviceNewCategory.network service","icons":["network_l_1-3"],"subcategories":null,"version":null,"ownerId":null,"empty":false,"type":null},{"name":"Network L1-3","normalizedName":"network l1-3","uniqueId":"serviceNewCategory.network l1-3","icons":["network_l_1-3"],"subcategories":null,"version":null,"ownerId":null,"empty":false,"type":null},{"name":"Partner Domain Service","normalizedName":"partner domain service","uniqueId":"serviceNewCategory.partner domain service","icons":["partner_domain_service"],"subcategories":null,"version":null,"ownerId":null,"empty":false,"type":null}],"productCategories":[]},"version":"1.6.7"}`)
}
