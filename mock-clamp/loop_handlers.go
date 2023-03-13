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

//describes the loop state in term of policy and DCAE
type State struct {
	ComponentState struct {
		StateName string `json:"stateName"`
	} `json:"componentState"`
}

//describes entity id in CLAMP got from SDC
type Resource struct {
	ResourceID struct {
		VfModuleModelName              string `json:"vfModuleModelName"`
		VfModuleModelInvariantUUID     string `json:"vfModuleModelInvariantUUID"`
		VfModuleModelUUID              string `json:"vfModuleModelUUID"`
		VfModuleModelVersion           string `json:"vfModuleModelVersion"`
		VfModuleModelCustomizationUUID string `json:"vfModuleModelCustomizationUUID"`
	} `json:"resourceID"`
}

//frequency limiter configuration
type Frequency_payload struct {
	ID         string `json:"id"`
	Actor      string `json:"actor"`
	Operation  string `json:"operation"`
	Limit      int    `json:"limit"`
	TimeWindow int    `json:"timeWindow"`
	TimeUnits  string `json:"timeUnits"`
}

/*
//describes operational policy configuration in CLAMP
type OperationalPolicy struct {
	Name               string                      `json:"name"`
	PolicyModel        Policy                      `json:"policyModel"`
	ConfigurationsJson struct{ Frequency_payload } `json:"configurationsJson"` //depends on operational policy model
}
*/
//describes TCA POLICY in CLAMP
type Tca_policy struct {
	Domain              string `json:"Domain"`
	MetricsPerEventName []struct {
		PolicyScope string `json:"policyScope"`
		Thresholds  []struct {
			Version               string `json:"version"`
			Severity              string `json:"severity"`
			ThresholdValue        int    `json:"thresholdValue"`
			ClosedLoopEventStatus string `json:"closedLoopEventStatus"`
			ClosedLoopControlName string `json:"closedLoopControlName"`
			Direction             string `json:"direction"`
			FieldPath             string `json:"fieldPath"`
		} `json:"thresholds"`
		EventName             string `json:"eventName"`
		PolicyVersion         string `json:"policyVersion"`
		ControlLoopSchemaType string `json:"controlLoopSchemaType"`
		PolicyName            string `json:"policyName"`
	} `json:"metricsPerEventName"`
}

//describes TCA POLICY payload in CLAMP
type Tca_policy_struct struct {
	TcaPolicy Tca_policy `json:"tca.policy"`
}

//describes microservice policy configuration in CLAMP
type MicroServicePolicy struct {
	Name               string            `json:"name"`
	ConfigurationsJson Tca_policy_struct `json:"configurationsJson"`
	PdpGroup           string            `json:"pdpGroup"`
	PdpSubgroup        string            `json:"pdpSubgroup"`
}

//LoopDetails describes loop innstance in CLAMP
type LoopDetails struct {
	Name                 string       `json:"name"`
	Template             LoopTemplate `json:"loopTemplate"`
	GlobalPropertiesJSON struct {
		DcaeDeployParameters struct {
			UniqueBlueprintParameters struct {
				PolicyID string `json:"policy_id"`
			} `json:"uniqueBlueprintParameters"`
		} `json:"dcaeDeployParameters"`
	} `json:"globalPropertiesJson"`
	LoopElementModelsUsed []string `json:"loopElementModelsUsed"` //microservices from sdc
	Components            struct {
		POLICY State `json:"POLICY"`
		DCAE   State `json:"DCAE"`
	} `json:"components"`
	ModelService struct {
		ResourceDetails struct {
			VFModule Resource `json:"VFModule"`
		} `json:"resourceDetails"`
	} `json:"modelService"`
	OperationalPolicies  []OperationalPolicy_payload `json:"operationalPolicies"`
	MicroServicePolicies []MicroServicePolicy        `json:"microServicePolicies"`
}

var loopInstanceList []LoopDetails

func generateInitialLoopInstances() {
	loopInstanceList = nil
	loop1 := new(LoopDetails)
	loop1.Template = templateList[0]
	loop1.Name = "intance_template01"
	loopInstanceList = append(loopInstanceList, *loop1)
	loop2 := new(LoopDetails)
	loop2.Name = "intance_template02"
	loop2.Template = templateList[1]
	loopInstanceList = append(loopInstanceList, *loop2)
}

func updateLoopDetails(c echo.Context) error {
	loopID := c.Param("loopID")
	for _, loop := range loopInstanceList {
		if loop.Name == loopID {
			return c.JSON(http.StatusOK, loop)
		}
	}
	return c.JSON(http.StatusNotFound, ClampError{
		Message: "No ClosedLoop found",
		Error:   "Not Found",
		Status:  "404"})
}

func createLoopInstance(c echo.Context) error {
	loopID := c.Param("loopID")
	templateName := c.QueryParam("templateName")
	loop := new(LoopDetails)

	//must add the constraint of limit number of instances for a template
	for _, l := range loopInstanceList {
		if l.Name == loopID &&
			l.Template.Name == templateName {
			//in reality it's overwritten
			return c.JSON(http.StatusBadRequest, ClampError{
				Message: "loop of same Name and for same template exists",
				Error:   "Exists",
				Status:  "500"})
		}
	}
	loop.Name = loopID
	loop.Template.Name = templateName
	//For drools configuration
	resource := new(Resource)
	resource.ResourceID.VfModuleModelUUID = uuid.NewV4().String()
	resource.ResourceID.VfModuleModelName = uuid.NewV4().String()
	resource.ResourceID.VfModuleModelInvariantUUID = uuid.NewV4().String()
	resource.ResourceID.VfModuleModelVersion = uuid.NewV4().String()
	resource.ResourceID.VfModuleModelCustomizationUUID = uuid.NewV4().String()
	loop.ModelService.ResourceDetails.VFModule = *resource
	//must generate as much microservices as tca blueprints count
	loop.LoopElementModelsUsed = append(loop.LoopElementModelsUsed, "microservice01")
	nb_microservices := len(loop.LoopElementModelsUsed)
	for i := 0; i < nb_microservices; i++ {
		loop.MicroServicePolicies = append(loop.MicroServicePolicies, MicroServicePolicy{
			Name: "Microservice" + uuid.NewV4().String(),
		})
	}
	loop.GlobalPropertiesJSON.DcaeDeployParameters.UniqueBlueprintParameters.PolicyID = "Microservice" + uuid.NewV4().String()
	state := new(State)
	state.ComponentState.StateName = "NOT_SENT"
	loop.Components.POLICY = *state
	state.ComponentState.StateName = "BLUEPRINT_DEPLOYED"
	loop.Components.DCAE = *state

	loopInstanceList = append(loopInstanceList, *loop)
	return c.JSON(http.StatusCreated, *loop)
}

func addOperationaPolicy(c echo.Context) error {
	loopID := c.Param("loopID")
	policyType := c.Param("policyType")
	policyVersion := c.Param("policyVersion")
	op_policy := new(struct {
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
	})
	for _, p := range policyList {
		if p.PolicyModelType == policyType &&
			p.Version == policyVersion {
			for j, l := range loopInstanceList {
				if l.Name == loopID {
					op_policy.PolicyAcronym = p.PolicyAcronym
					loopInstanceList[j].OperationalPolicies = append(loopInstanceList[j].OperationalPolicies,
						OperationalPolicy_payload{
							Name:        "OPERATIONAL" + uuid.NewV4().String(),
							PolicyModel: *op_policy,
						})
					return c.JSON(http.StatusOK, loopInstanceList[j])
				}
			}
			return c.JSON(http.StatusBadRequest, ClampError{
				Message: "loop not found",
				Error:   "Not Found",
				Status:  "404"})
		}
	}
	return c.JSON(http.StatusBadRequest, ClampError{
		Message: "Policy not found",
		Error:   "Not Found",
		Status:  "404"})
}

//remove operation is not working in the real CLAMP
func removeOperationaPolicy(c echo.Context) error {
	loopID := c.Param("loopID")
	policyType := c.Param("policyType")
	policyVersion := c.Param("policyVersion")

	for j, l := range loopInstanceList {
		if l.Name == loopID {
			for i, pp := range l.OperationalPolicies {
				if pp.PolicyModel.PolicyModelType == policyType &&
					pp.PolicyModel.Version == policyVersion {
					loopInstanceList[j].OperationalPolicies = append(loopInstanceList[j].OperationalPolicies[:i],
						loopInstanceList[j].OperationalPolicies[i+1:]...)
					return c.JSON(http.StatusOK, l)
				}
			}
			return c.JSON(http.StatusBadRequest, ClampError{
				Message: "Policy not found",
				Error:   "Not Found",
				Status:  "404"})
		}
	}
	return c.JSON(http.StatusBadRequest, ClampError{
		Message: "loop not found",
		Error:   "Not Found",
		Status:  "404"})
}

//must review tca_policy struct
func addTcaConfig(c echo.Context) error {
	loopID := c.Param("loopID")
	data := new(MicroServicePolicy)
	if err := c.Bind(data); err != nil {
		return err
	}
	for j, l := range loopInstanceList {
		if l.Name == loopID {
			if l.MicroServicePolicies != nil {
				for i, _ := range l.MicroServicePolicies {
					loopInstanceList[j].MicroServicePolicies[i] = *data
				}
				return c.JSON(http.StatusOK, loopInstanceList[j].MicroServicePolicies[0])
			}
			return c.JSON(http.StatusBadRequest, ClampError{
				Message: "Microservice policy not found",
				Error:   "Not Found",
				Status:  "404"})
		}
	}
	return c.JSON(http.StatusBadRequest, ClampError{
		Message: "Loop not found",
		Error:   "Not Found",
		Status:  "404"})
}

func addOperationalPolicyConfig(c echo.Context) error {
	loopID := c.Param("loopID")
	for j, l := range loopInstanceList {
		if l.Name == loopID {
			/*
				//cannot bind a list as said in labstack echo bind
				data := new([]OperationalPolicy_payload)
				if err := c.Bind(data); err != nil {
					return err
				}
			*/
			if l.OperationalPolicies != nil {
				//loopInstanceList[j].OperationalPolicies = *data
				loopInstanceList[j].OperationalPolicies[len(l.OperationalPolicies)-1].ConfigurationsJSON.Actor = "Test"
				return c.JSON(http.StatusOK, loopInstanceList[j])
			}
			return c.JSON(http.StatusBadRequest, ClampError{
				Message: "Operational Policy not found",
				Error:   "Not Found",
				Status:  "404"})
		}
	}
	return c.JSON(http.StatusBadRequest, ClampError{
		Message: "Loop not found",
		Error:   "Not Found",
		Status:  "404"})
}

//util function
func checkPoliciesConfiguration(l LoopDetails) bool {
	var empty struct {
		Actor      string `json:"actor"`
		Operation  string `json:"operation"`
		Limit      int    `json:"limit"`
		TimeWindow int    `json:"timeWindow"`
		TimeUnits  string `json:"timeUnits"`
	}

	for i, _ := range l.MicroServicePolicies {
		if l.MicroServicePolicies[i].ConfigurationsJson.TcaPolicy.Domain == "" {
			return false
		}
	}
	for i, _ := range l.OperationalPolicies {
		if l.OperationalPolicies[i].ConfigurationsJSON == empty {
			return false
		}
	}
	return true
}

func putLoopAction(c echo.Context) error {
	action := c.Param("action")
	loopID := c.Param("loopID")
	state := new(State)
	for j, l := range loopInstanceList {
		if l.Name == loopID {
			//POLICY actions
			if action == "submit" {
				if checkPoliciesConfiguration(loopInstanceList[j]) {
					state.ComponentState.StateName = "SENT_AND_DEPLOYED"
					loopInstanceList[j].Components.POLICY = *state
					return c.JSON(http.StatusOK, loopInstanceList[j].Components.POLICY.ComponentState)
				}
				return c.JSON(http.StatusBadRequest, ClampError{
					Message: "Policies are not well Configured",
					Error:   "Bad Action",
					Status:  "401"})
			}
			if action == "stop" {
				if l.Components.POLICY.ComponentState.StateName == "NOT_SENT" {
					return c.JSON(http.StatusBadRequest, ClampError{
						Message: "Cannot perform this action",
						Error:   "Bad Action",
						Status:  "400"})
				}
				state.ComponentState.StateName = "SENT"
				loopInstanceList[j].Components.POLICY = *state
				return c.JSON(http.StatusOK, "{}")
			}
			if action == "restart" {
				state.ComponentState.StateName = "SENT_AND_DEPLOYED"
				loopInstanceList[j].Components.POLICY = *state
				return c.JSON(http.StatusOK, "{}")
			}
			//DCAE actions
			if action == "deploy" {
				//must add deploy failure
				state.ComponentState.StateName = "MICROSERVICE_INSTALLED_SUCCESSFULLY"
				loopInstanceList[j].Components.DCAE = *state
				return c.JSON(http.StatusOK, "{}")
			}
			if action == "undeploy" {
				state.ComponentState.StateName = "MICROSERVICE_UNINSTALLED_SUCCESSFULLY"
				loopInstanceList[j].Components.DCAE = *state
				return c.JSON(http.StatusOK, "{}")
			}
			//LOOP action
			if action == "delete" {
				loopInstanceList = append(loopInstanceList[:j], loopInstanceList[j+1:]...)
				return c.JSON(http.StatusOK, "{}")
			}
			//action failure
			return c.JSON(http.StatusBadRequest, ClampError{
				Message: "Cannot perform this action",
				Error:   "Bad Action",
				Status:  "400"})
		}
	}
	return c.JSON(http.StatusBadRequest, ClampError{
		Message: "Loop not found",
		Error:   "Not Found",
		Status:  "404"})
}
