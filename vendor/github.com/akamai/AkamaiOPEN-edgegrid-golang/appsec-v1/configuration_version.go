package appsec

import (
	"fmt"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/client-v1"
	edge "github.com/akamai/AkamaiOPEN-edgegrid-golang/edgegrid"
)

// ConfigurationVersion represents a collection of ConfigurationVersion
//
// See: ConfigurationVersion.GetConfigurationVersion()
// API Docs: // appsec v1
//
// https://developer.akamai.com/api/cloud_security/application_security/v1.html

type ConfigurationVersionResponse struct {
	ConfigID           int    `json:"configId"`
	ConfigName         string `json:"configName"`
	LastCreatedVersion int    `json:"lastCreatedVersion"`
	Page               int    `json:"page"`
	PageSize           int    `json:"pageSize"`
	TotalSize          int    `json:"totalSize"`
	VersionList        []struct {
		ConfigID   int `json:"configId"`
		Production struct {
			Status string `json:"status"`
		} `json:"production"`
		Staging struct {
			Status string `json:"status"`
		} `json:"staging"`
		Version int `json:"version"`
		BasedOn int `json:"basedOn,omitempty"`
	} `json:"versionList"`
}

// NewConfigurationVersion creates a new *ConfigurationVersion
func NewConfigurationVersionResponse() *ConfigurationVersionResponse {
	return &ConfigurationVersionResponse{}
}

// GetConfigurationVersion populates a *ConfigurationVersion with it's related ConfigurationVersion
//
// API Docs: // appsec v1
//
// https://developer.akamai.com/api/cloud_security/application_security/v1.html#getconfigurationversion

func (configurationversion *ConfigurationVersionResponse) GetConfigurationVersion(correlationid string) error {

	req, err := client.NewRequest(
		Config,
		"GET",
		fmt.Sprintf(
			"/appsec/v1/configs/%d/versions?page=-1&detail=false",
			configurationversion.ConfigID,
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

	if err = client.BodyJSON(res, configurationversion); err != nil {
		return err
	}

	return nil

}
