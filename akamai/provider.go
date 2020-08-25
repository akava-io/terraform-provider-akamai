package akamai

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	appsec "github.com/akamai/AkamaiOPEN-edgegrid-golang/appsec-v1"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/client-v1"
	dnsv2 "github.com/akamai/AkamaiOPEN-edgegrid-golang/configdns-v2"
	gtm "github.com/akamai/AkamaiOPEN-edgegrid-golang/configgtm-v1_4"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/edgegrid"
	"github.com/terraform-providers/terraform-provider-akamai/version"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/papi-v1"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

//const (
//	Version = "0.2.0"
//)

// Config contains the Akamai provider configuration (unused).
type Config struct {
	terraformVersion string
}

func getConfigOptions(section string) *schema.Resource {
	section = strings.ToUpper(section)

	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"host": {
				Type:     schema.TypeString,
				Optional: true,
				DefaultFunc: func() (interface{}, error) {
					if v := os.Getenv("AKAMAI_" + section + "_HOST"); v != "" {
						return v, nil
					} else if v := os.Getenv("AKAMAI_HOST"); v != "" {
						return v, nil
					}

					return nil, errors.New("host not set")
				},
			},
			"access_token": {
				Type:     schema.TypeString,
				Optional: true,
				DefaultFunc: func() (interface{}, error) {
					if v := os.Getenv("AKAMAI_" + section + "_ACCESS_TOKEN"); v != "" {
						return v, nil
					} else if v := os.Getenv("AKAMAI_ACCESS_TOKEN"); v != "" {
						return v, nil
					}

					return nil, errors.New("access_token not set")
				},
			},
			"client_token": {
				Type:     schema.TypeString,
				Optional: true,
				DefaultFunc: func() (interface{}, error) {
					if v := os.Getenv("AKAMAI_" + section + "_CLIENT_TOKEN"); v != "" {
						return v, nil
					} else if v := os.Getenv("AKAMAI_CLIENT_TOKEN"); v != "" {
						return v, nil
					}

					return nil, errors.New("client_token not set")
				},
			},
			"client_secret": {
				Type:     schema.TypeString,
				Optional: true,
				DefaultFunc: func() (interface{}, error) {
					if v := os.Getenv("AKAMAI_" + section + "_CLIENT_SECRET"); v != "" {
						return v, nil
					} else if v := os.Getenv("AKAMAI_CLIENT_SECRET"); v != "" {
						return v, nil
					}

					return nil, errors.New("client_secret not set")
				},
			},
			"max_body": {
				Type:     schema.TypeInt,
				Optional: true,
				DefaultFunc: func() (interface{}, error) {
					if v := os.Getenv("AKAMAI_" + section + "_MAX_SIZE"); v != "" {
						return v, nil
					} else if v := os.Getenv("AKAMAI_MAX_SIZE"); v != "" {
						return v, nil
					}

					return 131072, nil
				},
			},
		},
	}
}

// Provider returns the Akamai terraform.Resource provider.
func Provider() *schema.Provider {
	//client.UserAgent = client.UserAgent + " terraform/" + Version

	provider := &schema.Provider{
		Schema: map[string]*schema.Schema{
			"edgerc": &schema.Schema{
				Optional: true,
				Type:     schema.TypeString,
			},
			"dns_section": &schema.Schema{
				Optional: true,
				Type:     schema.TypeString,
				Default:  "default",
			},
			"gtm_section": &schema.Schema{
				Optional: true,
				Type:     schema.TypeString,
				Default:  "default",
			},
			"papi_section": &schema.Schema{
				Optional: true,
				Type:     schema.TypeString,
				Default:  "default",
			},
			"property_section": &schema.Schema{
				Optional: true,
				Type:     schema.TypeString,
				Default:  "default",
			},
			"appsec_section": &schema.Schema{
				Optional: true,
				Type:     schema.TypeString,
				Default:  "default",
			},
			"property": &schema.Schema{
				Optional: true,
				Type:     schema.TypeSet,
				Elem:     getConfigOptions("property"),
			},
			"dns": &schema.Schema{
				Optional: true,
				Type:     schema.TypeSet,
				Elem:     getConfigOptions("dns"),
			},
			"gtm": &schema.Schema{
				Optional: true,
				Type:     schema.TypeSet,
				Elem:     getConfigOptions("gtm"),
			},
			"appsec": &schema.Schema{
				Optional: true,
				Type:     schema.TypeSet,
				Elem:     getConfigOptions("appsec"),
			},
		},
		DataSourcesMap: map[string]*schema.Resource{
			"akamai_authorities_set":             dataSourceAuthoritiesSet(),
			"akamai_contract":                    dataSourcePropertyContract(),
			"akamai_cp_code":                     dataSourceCPCode(),
			"akamai_dns_record_set":              dataSourceDNSRecordSet(),
			"akamai_group":                       dataSourcePropertyGroups(),
			"akamai_property_rules":              dataPropertyRules(),
			"akamai_property":                    dataSourceAkamaiProperty(),
			"akamai_gtm_default_datacenter":      dataSourceGTMDefaultDatacenter(),
			"akamai_appsec_configuration":        dataSourceConfiguration(),
			"akamai_appsec_export_configuration": dataSourceExportConfiguration(),
			"akamai_appsec_selectable_hostnames": dataSourceSelectableHostnames(),
		},
		ResourcesMap: map[string]*schema.Resource{
			"akamai_cp_code":                      resourceCPCode(),
			"akamai_dns_zone":                     resourceDNSv2Zone(),
			"akamai_dns_record":                   resourceDNSv2Record(),
			"akamai_edge_hostname":                resourceSecureEdgeHostName(),
			"akamai_property":                     resourceProperty(),
			"akamai_property_rules":               resourcePropertyRules(),
			"akamai_property_variables":           resourcePropertyVariables(),
			"akamai_property_activation":          resourcePropertyActivation(),
			"akamai_gtm_domain":                   resourceGTMv1Domain(),
			"akamai_gtm_datacenter":               resourceGTMv1Datacenter(),
			"akamai_gtm_property":                 resourceGTMv1Property(),
			"akamai_gtm_resource":                 resourceGTMv1Resource(),
			"akamai_gtm_cidrmap":                  resourceGTMv1Cidrmap(),
			"akamai_gtm_geomap":                   resourceGTMv1Geomap(),
			"akamai_gtm_asmap":                    resourceGTMv1ASmap(),
			"akamai_appsec_configuration_clone":   resourceConfigurationClone(),
			"akamai_appsec_selected_hostnames":    resourceSelectedHostnames(),
			"akamai_appsec_security_policy_clone": resourceSecurityPolicyClone(),
			"akamai_appsec_match_targets":         resourceMatchTargets(),
			"akamai_appsec_custom_rule":           resourceCustomRule(),
		},
	}
	//ConfigureFunc: providerConfigure,
	provider.ConfigureFunc = func(d *schema.ResourceData) (interface{}, error) {
		terraformVersion := provider.TerraformVersion
		if terraformVersion == "" {
			// Terraform 0.12 introduced this field to the protocol
			// We can therefore assume that if it's missing it's 0.10 or 0.11
			terraformVersion = "0.11+compatible"
		}
		tfUserAgent := provider.UserAgent("terraform-provider-akamai", version.ProviderVersion)
		client.UserAgent = fmt.Sprintf("%s ", tfUserAgent)
		return providerConfigure(d, terraformVersion)
	}
	return provider
}

func providerConfigure(d *schema.ResourceData, terraformVersion string) (interface{}, error) {
	log.Printf("[DEBUG] START providerConfigure  %s\n", terraformVersion)
	dnsv2Config, dnsErr := getConfigDNSV2Service(d)
	papiConfig, papiErr := getPAPIV1Service(d)
	gtmConfig, gtmErr := getConfigGTMV1Service(d)
	appsecConfig, appsecErr := getAPPSECV1Service(d)

	if dnsErr != nil && papiErr != nil && gtmErr != nil && appsecErr != nil || dnsv2Config == nil && papiConfig == nil && gtmConfig == nil && appsecConfig == nil {
		return nil, fmt.Errorf("One or more Akamai Edgegrid provider configurations must be defined")
	}

	config := Config{
		terraformVersion: terraformVersion,
	}

	return &config, nil
}

type resourceData interface {
	GetOk(string) (interface{}, bool)
	Get(string) interface{}
}

type set interface {
	List() []interface{}
}

func getConfigDNSV2Service(d resourceData) (*edgegrid.Config, error) {
	var DNSv2Config edgegrid.Config
	var err error
	if _, ok := d.GetOk("dns"); ok {
		config := d.Get("dns").(set).List()[0].(map[string]interface{})

		DNSv2Config = edgegrid.Config{
			Host:         config["host"].(string),
			AccessToken:  config["access_token"].(string),
			ClientToken:  config["client_token"].(string),
			ClientSecret: config["client_secret"].(string),
			MaxBody:      config["max_body"].(int),
		}

		dnsv2.Init(DNSv2Config)
		return &DNSv2Config, nil
	}

	edgerc := d.Get("edgerc").(string)
	section := d.Get("dns_section").(string)
	DNSv2Config, err = edgegrid.Init(edgerc, section)
	if err != nil {
		return nil, err
	}

	dnsv2.Init(DNSv2Config)
	edgegrid.SetupLogging()
	return &DNSv2Config, nil
}

func getConfigGTMV1Service(d resourceData) (*edgegrid.Config, error) {
	var GTMv1Config edgegrid.Config
	var err error
	if _, ok := d.GetOk("gtm"); ok {
		config := d.Get("gtm").(set).List()[0].(map[string]interface{})

		GTMv1Config = edgegrid.Config{
			Host:         config["host"].(string),
			AccessToken:  config["access_token"].(string),
			ClientToken:  config["client_token"].(string),
			ClientSecret: config["client_secret"].(string),
			MaxBody:      config["max_body"].(int),
		}

		gtm.Init(GTMv1Config)
		edgegrid.SetupLogging()
		return &GTMv1Config, nil
	}

	edgerc := d.Get("edgerc").(string)
	section := d.Get("gtm_section").(string)
	GTMv1Config, err = edgegrid.Init(edgerc, section)
	if err != nil {
		return nil, err
	}

	gtm.Init(GTMv1Config)
	return &GTMv1Config, nil
}

func getPAPIV1Service(d resourceData) (*edgegrid.Config, error) {
	var papiConfig edgegrid.Config
	if _, ok := d.GetOk("property"); ok {
		log.Printf("[DEBUG] Setting property config via HCL")
		config := d.Get("property").(set).List()[0].(map[string]interface{})

		papiConfig = edgegrid.Config{
			Host:         config["host"].(string),
			AccessToken:  config["access_token"].(string),
			ClientToken:  config["client_token"].(string),
			ClientSecret: config["client_secret"].(string),
			MaxBody:      config["max_body"].(int),
		}

		papi.Init(papiConfig)
		return &papiConfig, nil
	}

	var err error
	edgerc := d.Get("edgerc").(string)
	if section, ok := d.GetOk("property_section"); ok && section != "default" {
		papiConfig, err = edgegrid.Init(edgerc, section.(string))
	} else if section, ok := d.GetOk("papi_section"); ok && section != "default" {
		papiConfig, err = edgegrid.Init(edgerc, section.(string))
	} else {
		papiConfig, err = edgegrid.Init(edgerc, "default")
	}

	if err != nil {
		return nil, err
	}

	papi.Init(papiConfig)
	return &papiConfig, nil
}

func getAPPSECV1Service(d resourceData) (*edgegrid.Config, error) {
	var appsecConfig edgegrid.Config
	if _, ok := d.GetOk("appsec"); ok {
		log.Printf("[DEBUG] Setting appsec config via HCL")
		config := d.Get("appsec").(set).List()[0].(map[string]interface{})

		appsecConfig = edgegrid.Config{
			Host:         config["host"].(string),
			AccessToken:  config["access_token"].(string),
			ClientToken:  config["client_token"].(string),
			ClientSecret: config["client_secret"].(string),
			MaxBody:      config["max_body"].(int),
		}

		appsec.Init(appsecConfig)
		return &appsecConfig, nil
	}

	var err error
	edgerc := d.Get("edgerc").(string)
	if section, ok := d.GetOk("appsec_section"); ok && section != "default" {
		appsecConfig, err = edgegrid.Init(edgerc, section.(string))
	} else if section, ok := d.GetOk("appsec_section"); ok && section != "default" {
		appsecConfig, err = edgegrid.Init(edgerc, section.(string))
	} else {
		appsecConfig, err = edgegrid.Init(edgerc, "default")
	}

	if err != nil {
		return nil, err
	}

	appsec.Init(appsecConfig)
	return &appsecConfig, nil
}
