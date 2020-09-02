package appsec

import (
	"fmt"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/client-v1"
	edge "github.com/akamai/AkamaiOPEN-edgegrid-golang/edgegrid"
)

// SecurityPolicy represents a collection of SecurityPolicy
//
// See: SecurityPolicy.GetSecurityPolicy()
// API Docs: // appsec v1
//
// https://developer.akamai.com/api/cloud_security/application_security/v1.html

type SecurityPolicyResponse struct {
	ConfigID int `json:"configId"`
	Version  int `json:"version"`
	Policies []struct {
		PolicyID                string `json:"policyId"`
		PolicyName              string `json:"policyName"`
		HasRatePolicyWithAPIKey bool   `json:"hasRatePolicyWithApiKey"`
		PolicySecurityControls  struct {
			ApplyApplicationLayerControls bool `json:"applyApplicationLayerControls"`
			ApplyNetworkLayerControls     bool `json:"applyNetworkLayerControls"`
			ApplyRateControls             bool `json:"applyRateControls"`
			ApplyReputationControls       bool `json:"applyReputationControls"`
			ApplyBotmanControls           bool `json:"applyBotmanControls"`
			ApplyAPIConstraints           bool `json:"applyApiConstraints"`
			ApplySlowPostControls         bool `json:"applySlowPostControls"`
		} `json:"policySecurityControls"`
	} `json:"policies"`
}

// NewSecurityPolicy creates a new *SecurityPolicy
func NewSecurityPolicyResponse() *SecurityPolicyResponse {
	return &SecurityPolicyResponse{}
}

// GetSecurityPolicy populates a *SecurityPolicy with it's related SecurityPolicy
//
// API Docs: // appsec v1
//
// https://developer.akamai.com/api/cloud_security/application_security/v1.html#getsecuritypolicy

func (securitypolicy *SecurityPolicyResponse) GetSecurityPolicy(correlationid string) error {

	req, err := client.NewRequest(
		Config,
		"GET",
		fmt.Sprintf(
			"/appsec/v1/configs/%d/versions/%d/security-policies?notMatched=false&detail=true",
			securitypolicy.ConfigID,
			securitypolicy.Version,
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

	if err = client.BodyJSON(res, securitypolicy); err != nil {
		return err
	}

	return nil

}
