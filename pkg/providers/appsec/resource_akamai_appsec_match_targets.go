package appsec

import (
	"fmt"
	"strconv"

	appsec "github.com/akamai/AkamaiOPEN-edgegrid-golang/appsec-v1"
	edge "github.com/akamai/AkamaiOPEN-edgegrid-golang/edgegrid"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/tools"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// appsec v1
//
// https://developer.akamai.com/api/cloud_security/application_security/v1.html
func resourceMatchTargets() *schema.Resource {
	return &schema.Resource{
		Create: resourceMatchTargetsCreate,
		Read:   resourceMatchTargetsRead,
		Update: resourceMatchTargetsUpdate,
		Delete: resourceMatchTargetsDelete,
		Schema: map[string]*schema.Schema{
			"config_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"version": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"target_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"sequence": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"is_negative_path_match": {
				Type:     schema.TypeBool,
				Required: true,
			},
			"is_negative_file_extension_match": {
				Type:     schema.TypeBool,
				Required: true,
			},
			"default_file": {
				Type:     schema.TypeString,
				Required: true,
			},
			"hostnames": {
				Type:     schema.TypeSet,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"file_paths": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"file_extensions": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"security_policy": {
				Type:     schema.TypeString,
				Required: true,
			},
			"bypass_network_lists": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceMatchTargetsCreate(d *schema.ResourceData, meta interface{}) error {
	CorrelationID := "[APPSEC][resourceMatchTargetsCreate-" + tools.CreateNonce() + "]"
	edge.PrintfCorrelation("[DEBUG]", CorrelationID, " Creating MatchTargets")

	matchtargets := appsec.NewMatchTargetsResponse()

	matchtargets.ConfigID = d.Get("config_id").(int)
	matchtargets.ConfigVersion = d.Get("version").(int)
	matchtargets.Type = d.Get("type").(string)
	matchtargets.Sequence = d.Get("sequence").(int)
	matchtargets.IsNegativePathMatch = d.Get("is_negative_path_match").(bool)
	matchtargets.IsNegativeFileExtensionMatch = d.Get("is_negative_file_extension_match").(bool)
	matchtargets.DefaultFile = d.Get("default_file").(string)
	matchtargets.Hostnames = tools.SetToStringSlice(d.Get("hostnames").(*schema.Set))
	matchtargets.FilePaths = tools.SetToStringSlice(d.Get("file_paths").(*schema.Set))
	matchtargets.FileExtensions = tools.SetToStringSlice(d.Get("file_extensions").(*schema.Set))
	matchtargets.SecurityPolicy.PolicyID = d.Get("security_policy").(string)
	bypassnetworklists := d.Get("bypass_network_lists").(*schema.Set).List()

	for _, b := range bypassnetworklists {
		bl := appsec.BypassNetworkList{}
		bl.ID = b.(string)
		matchtargets.BypassNetworkLists = append(matchtargets.BypassNetworkLists, bl)
	}

	postresp, err := matchtargets.SaveMatchTargets(CorrelationID)
	if err != nil {
		edge.PrintfCorrelation("[DEBUG]", CorrelationID, fmt.Sprintf("Error  %v\n", err))
		return err
	}

	d.SetId(strconv.Itoa(postresp.TargetID))

	return resourceMatchTargetsRead(d, meta)
}

func resourceMatchTargetsUpdate(d *schema.ResourceData, meta interface{}) error {
	CorrelationID := "[APPSEC][resourceMatchTargetsUpdate-" + tools.CreateNonce() + "]"
	edge.PrintfCorrelation("[DEBUG]", CorrelationID, " Updating MatchTargets")

	matchtargets := appsec.NewMatchTargetsResponse()

	matchtargets.ConfigID = d.Get("config_id").(int)
	matchtargets.ConfigVersion = d.Get("version").(int)
	matchtargets.Type = d.Get("type").(string)
	matchtargets.Sequence = d.Get("sequence").(int)
	matchtargets.IsNegativePathMatch = d.Get("is_negative_path_match").(bool)
	matchtargets.IsNegativeFileExtensionMatch = d.Get("is_negative_file_extension_match").(bool)
	matchtargets.DefaultFile = d.Get("default_file").(string)
	matchtargets.Hostnames = tools.SetToStringSlice(d.Get("hostnames").(*schema.Set))
	matchtargets.FilePaths = tools.SetToStringSlice(d.Get("file_paths").(*schema.Set))
	matchtargets.FileExtensions = tools.SetToStringSlice(d.Get("file_extensions").(*schema.Set))
	matchtargets.SecurityPolicy.PolicyID = d.Get("security_policy").(string)
	matchtargets.TargetID, _ = strconv.Atoi(d.Id())
	bypassnetworklists := d.Get("bypass_network_lists").(*schema.Set).List()

	for _, b := range bypassnetworklists {
		bl := appsec.BypassNetworkList{}
		bl.ID = b.(string)
		matchtargets.BypassNetworkLists = append(matchtargets.BypassNetworkLists, bl)
	}

	err := matchtargets.UpdateMatchTargets(CorrelationID)
	if err != nil {
		edge.PrintfCorrelation("[DEBUG]", CorrelationID, fmt.Sprintf("Error  %v\n", err))
		return nil
	}

	return resourceMatchTargetsRead(d, meta)
}

func resourceMatchTargetsDelete(d *schema.ResourceData, meta interface{}) error {
	CorrelationID := "[APPSEC][resourceMatchTargetsDelete-" + tools.CreateNonce() + "]"
	edge.PrintfCorrelation("[DEBUG]", CorrelationID, "  Deleting MatchTargets")

	matchtargets := appsec.NewMatchTargetsResponse()

	matchtargets.ConfigID = d.Get("config_id").(int)
	matchtargets.ConfigVersion = d.Get("version").(int)
	matchtargets.TargetID, _ = strconv.Atoi(d.Id())

	err := matchtargets.DeleteMatchTargets(CorrelationID)
	if err != nil {
		edge.PrintfCorrelation("[DEBUG]", CorrelationID, fmt.Sprintf("Error  %v\n", err))
		return nil
	}

	d.SetId("")

	return nil
}

func resourceMatchTargetsRead(d *schema.ResourceData, meta interface{}) error {
	CorrelationID := "[APPSEC][resourceMatchTargetsRead-" + tools.CreateNonce() + "]"
	edge.PrintfCorrelation("[DEBUG]", CorrelationID, "  Read MatchTargets")

	matchtargets := appsec.NewMatchTargetsResponse()

	matchtargets.ConfigID = d.Get("config_id").(int)
	matchtargets.ConfigVersion = d.Get("version").(int)
	matchtargets.TargetID, _ = strconv.Atoi(d.Id())

	err := matchtargets.GetMatchTargets(CorrelationID)
	if err != nil {
		edge.PrintfCorrelation("[DEBUG]", CorrelationID, fmt.Sprintf("Error  %v\n", err))
		return err
	}

	edge.PrintfCorrelation("[DEBUG]", CorrelationID, fmt.Sprintf("CONFIG value  %v\n", matchtargets.TargetID))
	d.Set("type", matchtargets.Type)
	d.Set("sequence", matchtargets.Sequence)
	d.Set("is_negative_path_match", matchtargets.IsNegativePathMatch)
	d.Set("is_negative_file_extension_match", matchtargets.IsNegativeFileExtensionMatch)
	d.Set("default_file", matchtargets.DefaultFile)
	d.Set("hostnames", matchtargets.Hostnames)
	d.Set("file_paths", matchtargets.FilePaths)
	d.Set("file_extensions", matchtargets.FileExtensions)
	d.Set("security_policy", matchtargets.SecurityPolicy.PolicyID)
	d.Set("target_id", matchtargets.TargetID)
	d.SetId(strconv.Itoa(matchtargets.TargetID))

	return nil
}
