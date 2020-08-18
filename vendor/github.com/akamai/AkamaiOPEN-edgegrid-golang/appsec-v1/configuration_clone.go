package appsec

import (
	"fmt"

	"time"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/client-v1"
	edge "github.com/akamai/AkamaiOPEN-edgegrid-golang/edgegrid"
)

// ConfigurationClone represents a collection of ConfigurationClone
//
// See: ConfigurationClone.GetConfigurationClone()
// API Docs: // appsec v1
//
// https://developer.akamai.com/api/cloud_security/application_security/v1.html

type ConfigurationCloneResponse struct {
	ConfigID     int       `json:"configId"`
	ConfigName   string    `json:"configName"`
	Version      int       `json:"version"`
	VersionNotes string    `json:"versionNotes"`
	CreateDate   time.Time `json:"createDate"`
	CreatedBy    string    `json:"createdBy"`
	BasedOn      int       `json:"basedOn"`
	Production   struct {
		Status string    `json:"status"`
		Time   time.Time `json:"time"`
	} `json:"production"`
	Staging struct {
		Status string `json:"status"`
	} `json:"staging"`
}

type ConfigurationClonePost struct {
	CreateFromVersion int  `json:"createFromVersion"`
	RuleUpdate        bool `json:"ruleUpdate"`
}

// NewConfigurationClone creates a new *ConfigurationClone
func NewConfigurationCloneResponse() *ConfigurationCloneResponse {
	ConfigurationClone_new := &ConfigurationCloneResponse{}
	return ConfigurationClone_new
}

// NewConfigurationClone_post creates a new *ConfigurationClone_post
func NewConfigurationClonePost() *ConfigurationClonePost {
	ConfigurationClone_new := &ConfigurationClonePost{}
	return ConfigurationClone_new
}

// GetConfigurationClone populates a *ConfigurationClone with it's related ConfigurationClone
//
// API Docs: // appsec v1
//
// https://developer.akamai.com/api/cloud_security/application_security/v1.html#getconfigurationclone

func (configurationclone *ConfigurationCloneResponse) GetConfigurationClone(correlationid string) error {

	req, err := client.NewRequest(
		Config,
		"GET",
		fmt.Sprintf(
			"/appsec/v1/configs/%d/versions/%d",
			configurationclone.ConfigID,
			configurationclone.Version,
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

	if err = client.BodyJSON(res, configurationclone); err != nil {
		return err
	}

	return nil

}

// Save will create a new ConfigurationClone. You cannot update a ConfigurationClone;
// trying to do so will result in an error.
//
// API Docs: // appsec v1
//
// https://developer.akamai.com/api/cloud_security/application_security/v1.html#postconfigurationclone
func (configurationclone *ConfigurationCloneResponse) Save(postpayload *ConfigurationClonePost, correlationid string) (*ConfigurationCloneResponse, error) {
	req, err := client.NewJSONRequest(
		Config,
		"POST",
		fmt.Sprintf(
			"/appsec/v1/configs/%d/versions",
			configurationclone.ConfigID,
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
	ccr := &ConfigurationCloneResponse{}
	if err = client.BodyJSON(res, ccr); err != nil {
		return nil, err
	}

	return ccr, nil
}
