package appsec

import (
	"fmt"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/client-v1"
	edge "github.com/akamai/AkamaiOPEN-edgegrid-golang/edgegrid"
)

// SecurityPolicyClone represents a collection of SecurityPolicyClone
//
// See: SecurityPolicyClone.GetSecurityPolicyClone()
// API Docs: // appsec v1
//
// https://developer.akamai.com/api/cloud_security/application_security/v1.html

type SecurityPolicyCloneResponse struct {
	ConfigID int        `json:"configId"`
	Policies []Policies `json:"policies"`
	Version  int        `json:"version"`
}

type Policies struct {
	HasRatePolicyWithAPIKey bool   `json:"hasRatePolicyWithApiKey"`
	PolicyID                string `json:"policyId"`
	PolicyName              string `json:"policyName"`
	PolicySecurityControls  struct {
		ApplyAPIConstraints           bool `json:"applyApiConstraints"`
		ApplyApplicationLayerControls bool `json:"applyApplicationLayerControls"`
		ApplyBotmanControls           bool `json:"applyBotmanControls"`
		ApplyNetworkLayerControls     bool `json:"applyNetworkLayerControls"`
		ApplyRateControls             bool `json:"applyRateControls"`
		ApplyReputationControls       bool `json:"applyReputationControls"`
		ApplySlowPostControls         bool `json:"applySlowPostControls"`
	}
}

type SecurityPolicyClonePost struct {
	CreateFromSecurityPolicy string `json:"createFromSecurityPolicy"`
	PolicyName               string `json:"policyName"`
	PolicyPrefix             string `json:"policyPrefix"`
}

type SecurityPolicyClonePostResponse struct {
	ConfigID               int    `json:"configId"`
	PolicyID               string `json:"policyId"`
	PolicyName             string `json:"policyName"`
	PolicySecurityControls struct {
		ApplyAPIConstraints           bool `json:"applyApiConstraints"`
		ApplyApplicationLayerControls bool `json:"applyApplicationLayerControls"`
		ApplyBotmanControls           bool `json:"applyBotmanControls"`
		ApplyNetworkLayerControls     bool `json:"applyNetworkLayerControls"`
		ApplyRateControls             bool `json:"applyRateControls"`
		ApplyReputationControls       bool `json:"applyReputationControls"`
		ApplySlowPostControls         bool `json:"applySlowPostControls"`
	} `json:"policySecurityControls"`
	Version int `json:"version"`
}

// New SecurityPoliciesClonePostResponse
func NewSecurityPolicyClonePostResponse() *SecurityPolicyClonePostResponse {
	SecurityPolicyClonePostResponse_new := &SecurityPolicyClonePostResponse{}
	return SecurityPolicyClonePostResponse_new
}

// NewSecurityPolicyClone creates a new *SecurityPolicyClone
func NewSecurityPolicyCloneResponse() *SecurityPolicyCloneResponse {
	SecurityPolicyClone_new := &SecurityPolicyCloneResponse{}
	return SecurityPolicyClone_new
}

// NewSecurityPolicyClone_post creates a new *SecurityPolicyClone_post
func NewSecurityPolicyClonePost() *SecurityPolicyClonePost {
	SecurityPolicyClone_new := &SecurityPolicyClonePost{}
	return SecurityPolicyClone_new
}

// GetSecurityPolicyClone populates a *SecurityPolicyClone with it's related SecurityPolicyClone
//
// API Docs: // appsec v1
//
// https://developer.akamai.com/api/cloud_security/application_security/v1.html#getsecuritypolicyclone

func (securitypolicyclone *SecurityPolicyCloneResponse) GetSecurityPolicyClone(correlationid string) (*SecurityPolicyCloneResponse, error) {

	req, err := client.NewRequest(
		Config,
		"GET",
		fmt.Sprintf(
			"/appsec/v1/configs/%d/versions/%d/security-policies?notMatched=false&detail=true",
			securitypolicyclone.ConfigID,
			securitypolicyclone.Version,
		),
		nil,
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
	securitypolicycloneresponse := NewSecurityPolicyCloneResponse()

	if err = client.BodyJSON(res, securitypolicyclone); err != nil {
		return nil, err
	}

	return securitypolicycloneresponse, nil

}

// Save will create a new SecurityPolicyClone. You cannot update a SecurityPolicyClone;
// trying to do so will result in an error.
//
// API Docs: // appsec v1
//
// https://developer.akamai.com/api/cloud_security/application_security/v1.html#postsecuritypolicyclone
func (securitypolicyclone *SecurityPolicyCloneResponse) Save(postpayload *SecurityPolicyClonePost, correlationid string) (*SecurityPolicyClonePostResponse, error) {
	req, err := client.NewJSONRequest(
		Config,
		"POST",
		fmt.Sprintf(
			"/appsec/v1/configs/%d/versions/%d/security-policies",
			securitypolicyclone.ConfigID,
			securitypolicyclone.Version,
		),
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
	spcr := &SecurityPolicyClonePostResponse{}
	if err = client.BodyJSON(res, spcr); err != nil {
		return nil, err
	}

	return spcr, nil
}
