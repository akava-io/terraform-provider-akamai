package appsec

import (
	"fmt"

	"time"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/client-v1"
	edge "github.com/akamai/AkamaiOPEN-edgegrid-golang/edgegrid"
)

// ExportConfiguration represents a collection of ExportConfiguration
//
// See: ExportConfiguration.GetExportConfiguration()
// API Docs: // appsec v1
//
// https://developer.akamai.com/api/cloud_security/application_security/v1.html

type ExportConfigurationResponse struct {
	ConfigID   int    `json:"configId"`
	ConfigName string `json:"configName"`
	Version    int    `json:"version"`
	Staging    struct {
		Action string    `json:"action"`
		Status string    `json:"status"`
		Time   time.Time `json:"time"`
	} `json:"staging"`
	Production struct {
		Action string    `json:"action"`
		Status string    `json:"status"`
		Time   time.Time `json:"time"`
	} `json:"production"`
	CreateDate         time.Time     `json:"createDate"`
	CreatedBy          string        `json:"createdBy"`
	SelectedHosts      []interface{} `json:"selectedHosts"`
	SelectableHosts    []string      `json:"selectableHosts"`
	RatePolicies       []interface{} `json:"ratePolicies"`
	ReputationProfiles []struct {
		Context         string `json:"context"`
		ContextReadable string `json:"contextReadable"`
		Enabled         bool   `json:"enabled"`
		ID              int    `json:"id"`
		Name            string `json:"name"`
		Threshold       int    `json:"threshold"`
	} `json:"reputationProfiles"`
	CustomRules []interface{} `json:"customRules"`
	Rulesets    []struct {
		ID               int       `json:"id"`
		RulesetVersionID int       `json:"rulesetVersionId"`
		Type             string    `json:"type"`
		ReleaseDate      time.Time `json:"releaseDate"`
		Rules            []struct {
			ID                  int      `json:"id"`
			InspectRequestBody  bool     `json:"inspectRequestBody"`
			InspectResponseBody bool     `json:"inspectResponseBody"`
			Outdated            bool     `json:"outdated"`
			RuleVersion         int      `json:"ruleVersion"`
			Score               int      `json:"score"`
			Tag                 string   `json:"tag,omitempty"`
			Title               string   `json:"title"`
			AttackGroups        []string `json:"attackGroups,omitempty"`
		} `json:"rules"`
		AttackGroups []struct {
			Group     string `json:"group"`
			GroupName string `json:"groupName"`
			Threshold int    `json:"threshold"`
		} `json:"attackGroups"`
	} `json:"rulesets"`
	MatchTargets struct {
		WebsiteTargets []struct {
			Type                      string `json:"type"`
			DefaultFile               string `json:"defaultFile"`
			EffectiveSecurityControls struct {
				ApplyApplicationLayerControls bool `json:"applyApplicationLayerControls"`
				ApplyNetworkLayerControls     bool `json:"applyNetworkLayerControls"`
				ApplyRateControls             bool `json:"applyRateControls"`
				ApplyReputationControls       bool `json:"applyReputationControls"`
				ApplySlowPostControls         bool `json:"applySlowPostControls"`
			} `json:"effectiveSecurityControls"`
			FilePaths                    []string `json:"filePaths"`
			Hostnames                    []string `json:"hostNames"`
			ID                           int      `json:"id"`
			IsNegativeFileExtensionMatch bool     `json:"isNegativeFileExtensionMatch"`
			IsNegativePathMatch          bool     `json:"isNegativePathMatch"`
			SecurityPolicy               struct {
				PolicyID string `json:"policyId"`
			} `json:"securityPolicy"`
			Sequence int `json:"sequence"`
		} `json:"websiteTargets"`
	} `json:"matchTargets"`
	SecurityPolicies []struct {
		ID                      string `json:"id"`
		Name                    string `json:"name"`
		HasRatePolicyWithAPIKey bool   `json:"hasRatePolicyWithApiKey"`
		SecurityControls        struct {
			ApplyAPIConstraints           bool `json:"applyApiConstraints"`
			ApplyApplicationLayerControls bool `json:"applyApplicationLayerControls"`
			ApplyBotmanControls           bool `json:"applyBotmanControls"`
			ApplyNetworkLayerControls     bool `json:"applyNetworkLayerControls"`
			ApplyRateControls             bool `json:"applyRateControls"`
			ApplyReputationControls       bool `json:"applyReputationControls"`
			ApplySlowPostControls         bool `json:"applySlowPostControls"`
		} `json:"securityControls"`
		WebApplicationFirewall struct {
			RuleActions []struct {
				Action           string `json:"action"`
				ID               int    `json:"id"`
				RulesetVersionID int    `json:"rulesetVersionId"`
			} `json:"ruleActions"`
			AttackGroupActions []struct {
				Action           string `json:"action"`
				Group            string `json:"group"`
				RulesetVersionID int    `json:"rulesetVersionId"`
			} `json:"attackGroupActions"`
		} `json:"webApplicationFirewall"`
		IPGeoFirewall struct {
			Block      string `json:"block"`
			IPControls struct {
				BlockedIPNetworkLists struct {
					NetworkList []string `json:"networkList"`
				} `json:"blockedIPNetworkLists"`
			} `json:"ipControls"`
		} `json:"ipGeoFirewall"`
	} `json:"securityPolicies"`
	AdvancedOptions struct {
		Logging struct {
			AllowSampling bool `json:"allowSampling"`
			Cookies       struct {
				Type string `json:"type"`
			} `json:"cookies"`
			CustomHeaders struct {
				Type string `json:"type"`
			} `json:"customHeaders"`
			StandardHeaders struct {
				Type string `json:"type"`
			} `json:"standardHeaders"`
		} `json:"logging"`
		Prefetch struct {
			AllExtensions      bool `json:"allExtensions"`
			EnableAppLayer     bool `json:"enableAppLayer"`
			EnableRateControls bool `json:"enableRateControls"`
		} `json:"prefetch"`
	} `json:"advancedOptions"`
}

// NewExportConfiguration creates a new *ExportConfiguration
func NewExportConfigurationResponse() *ExportConfigurationResponse {
	return &ExportConfigurationResponse{}
}

// GetExportConfiguration populates a *ExportConfiguration with it's related ExportConfiguration
//
// API Docs: // appsec v1
//
// https://developer.akamai.com/api/cloud_security/application_security/v1.html#getexportconfiguration

func (exportconfiguration *ExportConfigurationResponse) GetExportConfiguration(correlationid string) error {

	req, err := client.NewRequest(
		Config,
		"GET",
		fmt.Sprintf(
			"/appsec/v1/export/configs/%d/versions/%d",
			exportconfiguration.ConfigID,
			exportconfiguration.Version,
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

	if err = client.BodyJSON(res, exportconfiguration); err != nil {
		return err
	}

	return nil

}
