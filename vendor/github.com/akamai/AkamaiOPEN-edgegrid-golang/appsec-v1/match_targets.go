package appsec

import (
	"fmt"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/client-v1"
	edge "github.com/akamai/AkamaiOPEN-edgegrid-golang/edgegrid"
)

// MatchTargets represents a collection of MatchTargets
//
// See: MatchTargets.GetMatchTargets()
// API Docs: // appsec v1
//
// https://developer.akamai.com/api/cloud_security/application_security/v1.html

type MatchTargetsResponse struct {
	Type                      string `json:"type"`
	ConfigID                  int    `json:"configId"`
	ConfigVersion             int    `json:"configVersion"`
	DefaultFile               string `json:"defaultFile"`
	EffectiveSecurityControls struct {
		ApplyApplicationLayerControls bool `json:"applyApplicationLayerControls"`
		ApplyBotmanControls           bool `json:"applyBotmanControls"`
		ApplyNetworkLayerControls     bool `json:"applyNetworkLayerControls"`
		ApplyRateControls             bool `json:"applyRateControls"`
		ApplyReputationControls       bool `json:"applyReputationControls"`
		ApplySlowPostControls         bool `json:"applySlowPostControls"`
	} `json:"effectiveSecurityControls"`
	Hostnames                    []string `json:"hostnames"`
	IsNegativeFileExtensionMatch bool     `json:"isNegativeFileExtensionMatch"`
	IsNegativePathMatch          bool     `json:"isNegativePathMatch"`
	FilePaths                    []string `json:"filePaths"`
	FileExtensions               []string `json:"fileExtensions"`
	SecurityPolicy               struct {
		PolicyID string `json:"policyId"`
	} `json:"securityPolicy"`
	Sequence           int `json:"sequence"`
	TargetID           int `json:"targetId"`
	BypassNetworkLists []struct {
		Name string `json:"name"`
		ID   string `json:"id"`
	} `json:"bypassNetworkLists"`
}

type BypassNetworkList struct {
	Name string `json:"name"`
	ID   string `json:"id"`
}

// NewCpCodes creates a new *CpCodes
func NewMatchTargetsResponse() *MatchTargetsResponse {
	return &MatchTargetsResponse{}
}

// GetMatchTargets populates a *MatchTargets with it's related MatchTargets
//
// API Docs: // appsec v1
//
// https://developer.akamai.com/api/cloud_security/application_security/v1.html#getmatchtargets

func (matchtargets *MatchTargetsResponse) GetMatchTargets(correlationid string) error {

	req, err := client.NewRequest(
		Config,
		"GET",
		fmt.Sprintf(
			"/appsec/v1/configs/%d/versions/%d/match-targets/%d?includeChildObjectName=true",
			matchtargets.ConfigID,
			matchtargets.ConfigVersion,
			matchtargets.TargetID,
		),
		nil,
	)
	if err != nil {
		return err
	}

	edge.PrintHttpRequestCorrelation(req, true, correlationid)

	res, err := client.Do(Config, req)
	if err != nil {
		return err
	}

	edge.PrintHttpResponseCorrelation(res, true, correlationid)

	if client.IsError(res) {
		return client.NewAPIError(res)
	}

	if err = client.BodyJSON(res, matchtargets); err != nil {
		return err
	}

	return nil

}

// Update will update a MatchTargets.
//
// API Docs: // appsec v1
//
// https://developer.akamai.com/api/cloud_security/application_security/v1.html#putmatchtargets
func (matchtargets *MatchTargetsResponse) UpdateMatchTargets(correlationid string) error {
	req, err := client.NewJSONRequest(
		Config,
		"PUT",
		fmt.Sprintf(
			"/appsec/v1/configs/%d/versions/%d/match-targets/%d",
			matchtargets.ConfigID,
			matchtargets.ConfigVersion,
			matchtargets.TargetID,
		),
		matchtargets,
	)
	if err != nil {
		return err
	}

	edge.PrintHttpRequestCorrelation(req, true, correlationid)

	res, err := client.Do(Config, req)
	if err != nil {
		return err
	}

	edge.PrintHttpResponseCorrelation(res, true, correlationid)

	if client.IsError(res) {
		return client.NewAPIError(res)
	}

	return nil
}

// Save will create a new matchtargets.
//
//
// API Docs: // appsec v1
//
// https://developer.akamai.com/api/cloud_security/application_security/v1.html#postmatchtargets
func (matchtargets *MatchTargetsResponse) SaveMatchTargets(correlationid string) (*MatchTargetsResponse, error) {
	req, err := client.NewJSONRequest(
		Config,
		"POST",
		fmt.Sprintf(
			"/appsec/v1/configs/%d/versions/%d/match-targets",
			matchtargets.ConfigID,
			matchtargets.ConfigVersion,
		),
		matchtargets,
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

	if err = client.BodyJSON(res, matchtargets); err != nil {
		return nil, err
	}

	return matchtargets, nil
}

// Delete will delete a MatchTargets
//
//
// API Docs: // appsec v1
//
// https://developer.akamai.com/api/cloud_security/application_security/v1.html#deletematchtargets
func (matchtargets *MatchTargetsResponse) DeleteMatchTargets(correlationid string) error {
	req, err := client.NewJSONRequest(
		Config,
		"DELETE",
		fmt.Sprintf(
			"/appsec/v1/configs/%d/versions/%d/match-targets/%d",
			matchtargets.ConfigID,
			matchtargets.ConfigVersion,
			matchtargets.TargetID,
		),
		nil,
	)
	if err != nil {
		return err
	}

	edge.PrintHttpRequestCorrelation(req, true, correlationid)

	res, err := client.Do(Config, req)
	if err != nil {
		return err
	}

	edge.PrintHttpResponseCorrelation(res, true, correlationid)

	if client.IsError(res) {
		return client.NewAPIError(res)
	}

	return nil
}
