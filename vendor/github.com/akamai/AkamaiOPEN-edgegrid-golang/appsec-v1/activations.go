package appsec

import (
	"fmt"
	"time"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/client-v1"
	edge "github.com/akamai/AkamaiOPEN-edgegrid-golang/edgegrid"
)

// Activations represents a collection of Activations
//
// See: Activations.GetActivations()
// API Docs: // appsec v1
//
// https://developer.akamai.com/api/cloud_security/application_security/v1.html

type ActivationsResponse struct {
	client.Resource
	DispatchCount     int          `json:"dispatchCount"`
	ActivationID      int          `json:"activationId"`
	Action            string       `json:"action"`
	Status            StatusValue  `json:"status"`
	StatusChange      chan bool    `json:"-"`
	Network           NetworkValue `json:"network"`
	Estimate          string       `json:"estimate"`
	CreatedBy         string       `json:"createdBy"`
	CreateDate        time.Time    `json:"createDate"`
	ActivationConfigs []struct {
		ConfigID              int    `json:"configId"`
		ConfigName            string `json:"configName"`
		ConfigVersion         int    `json:"configVersion"`
		PreviousConfigVersion int    `json:"previousConfigVersion"`
	} `json:"activationConfigs"`
}

type ActivationsPost struct {
	Action             string   `json:"action"`
	Network            string   `json:"network"`
	Note               string   `json:"note"`
	NotificationEmails []string `json:"notificationEmails"`
	ActivationConfigs  []struct {
		ConfigID      int `json:"configId"`
		ConfigVersion int `json:"configVersion"`
	} `json:"activationConfigs"`
}

type ActivationConfigs struct {
	ConfigID      int `json:"configId"`
	ConfigVersion int `json:"configVersion"`
}

func NewActivationsResponse() *ActivationsResponse {
	ActivationsResponse_new := &ActivationsResponse{}
	ActivationsResponse_new.Init()
	return ActivationsResponse_new

}

// NewActivations_post creates a new *Activations_post
func NewActivationsPost() *ActivationsPost {
	Activations_new := &ActivationsPost{}
	return Activations_new
}

// GetActivations populates  *Activations with it's related Activations
//
// API Docs: // appsec v1
//
// https://developer.akamai.com/api/cloud_security/application_security/v1.html

func (activations *ActivationsResponse) GetActivations(activationid int, correlationid string) (time.Duration, error) {
	req, err := client.NewRequest(
		Config,
		"GET",
		fmt.Sprintf("/appsec/v1/activations/%d",
			activationid,
		),
		nil,
	)

	if err != nil {
		return 0, err
	}

	edge.PrintHttpRequestCorrelation(req, true, correlationid)

	res, err := client.Do(Config, req)

	if err != nil {
		return 0, err
	}

	edge.PrintHttpResponseCorrelation(res, true, correlationid)

	if client.IsError(res) {
		return 0, client.NewAPIError(res)
	}

	if err = client.BodyJSON(res, activations); err != nil {
		return 0, err
	}

	return time.Duration(30 * time.Second), nil
}

// GetLatestProductionActivation retrieves the latest activation for the production network
//
// Pass in a status to check for, defaults to StatusActive
func (activations *ActivationsResponse) GetLatestProductionActivation(status StatusValue) (*ActivationsResponse, error) {
	return activations.GetLatestActivation(NetworkProduction, status)
}

// GetLatestStagingActivation retrieves the latest activation for the staging network
//
// Pass in a status to check for, defaults to StatusActive
func (activations *ActivationsResponse) GetLatestStagingActivation(status StatusValue) (*ActivationsResponse, error) {
	return activations.GetLatestActivation(NetworkStaging, status)
}

// GetLatestActivation gets the latest activation for the specified network
//
// Default to NetworkProduction. Pass in a status to check for, defaults to StatusActive
//
// This can return an activation OR a deactivation. Check activation.ActivationType and activation.Status for what you're looking for
func (activations *ActivationsResponse) GetLatestActivation(network NetworkValue, status StatusValue) (*ActivationsResponse, error) {
	if network == "" {
		network = NetworkProduction
	}

	if status == "" {
		status = StatusActive
	}

	var latest *ActivationsResponse

	if latest == nil {
		return nil, fmt.Errorf("No activation found (network: %s, status: %s)", network, status)
	}

	return latest, nil
}

func (activations *ActivationsResponse) Init() {
	activations.Complete = make(chan bool, 1)
	activations.StatusChange = make(chan bool, 1)
}

// Save activates a given Configuration
//
// If acknowledgeWarnings is true and warnings are returned on the first attempt,
// a second attempt is made, acknowledging the warnings.
//

func (activations *ActivationsResponse) SaveActivations(postpayload *ActivationsPost, acknowledgeWarnings bool, correlationid string) (*ActivationsResponse, error) {

	req, err := client.NewJSONRequest(
		Config,
		"POST",
		"/appsec/v1/activations",
		postpayload,
	)

	if err != nil {
		return nil, err
	}

	edge.PrintHttpRequestCorrelation(req, true, correlationid)

	res, err := client.Do(Config, req)

	if err != nil {
		return nil, err
	}

	edge.PrintHttpResponseCorrelation(res, true, correlationid)

	if client.IsError(res) && (!acknowledgeWarnings || (acknowledgeWarnings && res.StatusCode != 400)) {
		return nil, client.NewAPIError(res)
	}

	ar := &ActivationsResponse{}
	if err = client.BodyJSON(res, ar); err != nil {
		return nil, err
	}

	req, err = client.NewRequest(
		Config,
		"GET",
		fmt.Sprintf("/appsec/v1/activations/%d",
			ar.ActivationID,
		),
		nil,
	)

	if err != nil {
		return nil, err
	}

	edge.PrintHttpRequestCorrelation(req, true, correlationid)

	res, err = client.Do(Config, req)

	if err != nil {
		return nil, err
	}

	edge.PrintHttpResponseCorrelation(res, true, correlationid)

	if err := client.BodyJSON(res, activations); err != nil {
		return nil, err
	}

	return ar, nil
}

// PollStatus will responsibly poll till the configuration is active or an error occurs
//
// The Activation.StatusChange is a channel that can be used to
// block on status changes. If a new valid status is returned, true will
// be sent to the channel, otherwise, false will be sent.
//
//	go activation.PollStatus(activationid)
//	for activation.Status != edgegrid.StatusActive {
//		select {
//		case statusChanged := <-activation.StatusChange:
//			if statusChanged == false {
//				break
//			}
//		case <-time.After(time.Minute * 30):
//			break
//		}
//	}
//
//	if activation.Status == edgegrid.StatusActive {
//		// Activation succeeded
//	}
func (activations *ActivationsResponse) PollStatus(activationid int, correlationid string) bool {
	currentStatus := activations.Status
	edge.PrintfCorrelation("[DEBUG]", correlationid, fmt.Sprintf("currentStatus  %v\n", currentStatus))
	var retry time.Duration = 0

	for currentStatus != StatusActive {
		time.Sleep(retry)

		var err error
		retry, err = activations.GetActivations(activationid, correlationid)
		edge.PrintfCorrelation("[DEBUG]", correlationid, fmt.Sprintf("Get currentStatus  %v Polled status %v \n", currentStatus, activations.Status))
		if err != nil {
			activations.StatusChange <- false
			return false
		}

		if activations.Network == NetworkStaging && retry > time.Minute {
			retry = time.Minute
		}

		if err != nil {
			activations.StatusChange <- false
			return false
		}

		if currentStatus != activations.Status {
			currentStatus = activations.Status
			edge.PrintfCorrelation("[DEBUG]", correlationid, fmt.Sprintf("SET NEW currentStatus  %v Polled status %v \n", currentStatus, activations.Status))
			activations.StatusChange <- true
		}
	}

	return true
}

// Delete will delete a Activations
//
//
// API Docs: // appsec v1
//
// https://developer.akamai.com/api/cloud_security/application_security/v1.html#deleteactivations
func (activations *ActivationsResponse) DeactivateActivations(postpayload *ActivationsPost, correlationid string) (*ActivationsResponse, error) {
	req, err := client.NewJSONRequest(
		Config,
		"POST",
		"/appsec/v1/activations",
		postpayload,
	)

	if err != nil {
		return nil, err
	}

	edge.PrintHttpRequestCorrelation(req, true, correlationid)

	res, err := client.Do(Config, req)

	if err != nil {
		return nil, err
	}

	edge.PrintHttpResponseCorrelation(res, true, correlationid)

	if client.IsError(res) {
		return nil, client.NewAPIError(res)
	}

	ar := &ActivationsResponse{}
	if err = client.BodyJSON(res, ar); err != nil {
		return nil, err
	}

	return ar, nil
}

// ActivationValue is used to create an "enum" of possible Activation.ActivationType values
type ActivationValue string

// NetworkValue is used to create an "enum" of possible Activation.Network values
type NetworkValue string

// StatusValue is used to create an "enum" of possible Activation.Status values
type StatusValue string

const (
	// ActivationTypeActivate Activation.ActivationType value ACTIVATE
	ActivationTypeActivate ActivationValue = "ACTIVATE"
	// ActivationTypeDeactivate Activation.ActivationType value DEACTIVATE
	ActivationTypeDeactivate ActivationValue = "DEACTIVATE"

	// NetworkProduction Activation.Network value PRODUCTION
	NetworkProduction NetworkValue = "PRODUCTION"
	// NetworkStaging Activation.Network value STAGING
	NetworkStaging NetworkValue = "STAGING"

	// StatusActive Activation.Status value ACTIVE
	StatusActive StatusValue = "ACTIVATED"
	// StatusInactive Activation.Status value INACTIVE
	StatusInactive StatusValue = "INACTIVE"
	// StatusPending Activation.Status value RECEIVED
	StatusPending StatusValue = "RECEIVED"
	// StatusZone1 Activation.Status value ZONE_1
	StatusZone1 StatusValue = "ZONE_1"
	// StatusZone2 Activation.Status value ZONE_2
	StatusZone2 StatusValue = "ZONE_2"
	// StatusZone3 Activation.Status value ZONE_3
	StatusZone3 StatusValue = "ZONE_3"
	// StatusAborted Activation.Status value ABORTED
	StatusAborted StatusValue = "ABORTED"
	// StatusFailed Activation.Status value FAILED
	StatusFailed StatusValue = "FAILED"
	// StatusDeactivated Activation.Status value DEACTIVATED
	StatusDeactivated StatusValue = "DEACTIVATED"
	// StatusPendingDeactivation Activation.Status value PENDING_DEACTIVATION
	StatusPendingDeactivation StatusValue = "PENDING_DEACTIVATION"
	// StatusNew Activation.Status value NEW
	StatusNew StatusValue = "NEW"
)
