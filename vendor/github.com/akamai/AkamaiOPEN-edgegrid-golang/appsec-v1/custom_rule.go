package appsec

import (
	"fmt"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/client-v1"
	edge "github.com/akamai/AkamaiOPEN-edgegrid-golang/edgegrid"
)

// CustomRule represents a collection of CustomRule
//
// See: CustomRule.GetCustomRule()
// API Docs: // appsec v1
//
// https://developer.akamai.com/api/cloud_security/application_security/v1.html

type CustomRuleResponse struct {
	ID            int      `json:"id,omitempty"`
	Name          string   `json:"name"`
	Description   string   `json:"description"`
	Version       int      `json:"version,omitempty"`
	RuleActivated bool     `json:"ruleActivated"`
	Tag           []string `json:"tag"`
	Conditions    []struct {
		Type          string   `json:"type"`
		PositiveMatch bool     `json:"positiveMatch"`
		Value         []string `json:"value,omitempty"`
		ValueWildcard bool     `json:"valueWildcard,omitempty"`
		ValueCase     bool     `json:"valueCase,omitempty"`
		NameWildcard  bool     `json:"nameWildcard,omitempty"`
		Name          []string `json:"name,omitempty"`
		NameCase      bool     `json:"nameCase,omitempty"`
	} `json:"conditions"`
}

// NewCpCodes creates a new *CpCodes
func NewCustomRuleResponse() *CustomRuleResponse {
	return &CustomRuleResponse{}
}

// GetCustomRule populates a *CustomRule with it's related CustomRule
//
// API Docs: // appsec v1
//
// https://developer.akamai.com/api/cloud_security/application_security/v1.html#getcustomrule

func (customrule *CustomRuleResponse) GetCustomRule(configid int, ruleid int, correlationid string) error {

	req, err := client.NewRequest(
		Config,
		"GET",
		fmt.Sprintf(
			"/appsec/v1/configs/%d/custom-rules/%d",
			configid,
			ruleid,
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

	if err = client.BodyJSON(res, customrule); err != nil {
		return err
	}

	return nil

}

// Update will update a CustomRule.
//
// API Docs: // appsec v1
//
// https://developer.akamai.com/api/cloud_security/application_security/v1.html#putcustomrule
func (customrule *CustomRuleResponse) UpdateCustomRule(configid int, ruleid int, correlationid string) error {
	req, err := client.NewJSONRequest(
		Config,
		"PUT",
		fmt.Sprintf(
			"/appsec/v1/configs/%d/custom-rules/%d",
			configid,
			ruleid,
		),
		customrule,
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

// Save will create a new customrule.
//
//
// API Docs: // appsec v1
//
// https://developer.akamai.com/api/cloud_security/application_security/v1.html#postcustomrule
func (customrule *CustomRuleResponse) SaveCustomRule(configid int, correlationid string) (*CustomRuleResponse, error) {
	req, err := client.NewJSONRequest(
		Config,
		"POST",
		fmt.Sprintf(
			"/appsec/v1/configs/%d/custom-rules",
			configid,
		),
		customrule,
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

	if err = client.BodyJSON(res, customrule); err != nil {
		return nil, err
	}

	return customrule, nil
}

// Delete will delete a CustomRule
//
//
// API Docs: // appsec v1
//
// https://developer.akamai.com/api/cloud_security/application_security/v1.html#deletecustomrule
func (customrule *CustomRuleResponse) DeleteCustomRule(configid int, ruleid int, correlationid string) error {
	req, err := client.NewJSONRequest(
		Config,
		"DELETE",
		fmt.Sprintf(
			"/appsec/v1/configs/%d/custom-rules/%d",
			configid,
			ruleid,
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
