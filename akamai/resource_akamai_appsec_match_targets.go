package akamai

import (
	"encoding/json"
	"fmt"
	"strconv"

	appsec "github.com/akamai/AkamaiOPEN-edgegrid-golang/appsec-v1"
	edge "github.com/akamai/AkamaiOPEN-edgegrid-golang/edgegrid"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
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
			"json": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"type", "sequence", "is_negative_path_match", "is_negative_file_extension_match", "default_file", "host_names", "file_paths", "file_extensions", "security_policy", "bypass_network_lists"},
			},

			"target_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"sequence": {
				Type:          schema.TypeInt,
				Optional:      true,
				ConflictsWith: []string{"json"},
			},
			"is_negative_path_match": {
				Type:          schema.TypeBool,
				Optional:      true,
				ConflictsWith: []string{"json"},
			},
			"is_negative_file_extension_match": {
				Type:          schema.TypeBool,
				Optional:      true,
				ConflictsWith: []string{"json"},
			},
			"default_file": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"json"},
			},
			"host_names": &schema.Schema{
				Type:          schema.TypeSet,
				Optional:      true,
				Elem:          &schema.Schema{Type: schema.TypeString},
				ConflictsWith: []string{"json"},
			},
			"file_paths": &schema.Schema{
				Type:          schema.TypeSet,
				Optional:      true,
				Elem:          &schema.Schema{Type: schema.TypeString},
				ConflictsWith: []string{"json"},
			},
			"file_extensions": &schema.Schema{
				Type:          schema.TypeSet,
				Optional:      true,
				Elem:          &schema.Schema{Type: schema.TypeString},
				ConflictsWith: []string{"json"},
			},
			"security_policy": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"json"},
			},
			"bypass_network_lists": &schema.Schema{
				Type:          schema.TypeSet,
				Optional:      true,
				Elem:          &schema.Schema{Type: schema.TypeString},
				ConflictsWith: []string{"json"},
			},
		},
	}
}

func resourceMatchTargetsCreate(d *schema.ResourceData, meta interface{}) error {
	CorrelationID := "[APPSEC][resourceMatchTargetsCreate-" + CreateNonce() + "]"
	edge.PrintfCorrelation("[DEBUG]", CorrelationID, " Creating MatchTargets")

	matchtargets := appsec.NewMatchTargetsResponse()

	jsonpostpayload, ok := d.GetOk("json")
	if ok {

		json.Unmarshal([]byte(jsonpostpayload.(string)), &matchtargets)
	} else {
		matchtargets.ConfigID = d.Get("config_id").(int)
		matchtargets.ConfigVersion = d.Get("version").(int)
		matchtargets.Type = d.Get("type").(string)
		matchtargets.Sequence = d.Get("sequence").(int)
		matchtargets.IsNegativePathMatch = d.Get("is_negative_path_match").(bool)
		matchtargets.IsNegativeFileExtensionMatch = d.Get("is_negative_file_extension_match").(bool)
		matchtargets.DefaultFile = d.Get("default_file").(string)
		matchtargets.Hostnames = setToStringSlice(d.Get("host_names").(*schema.Set))
		matchtargets.FilePaths = setToStringSlice(d.Get("file_paths").(*schema.Set))
		matchtargets.FileExtensions = setToStringSlice(d.Get("file_extensions").(*schema.Set))
		matchtargets.SecurityPolicy.PolicyID = d.Get("security_policy").(string)
		bypassnetworklists := d.Get("bypass_network_lists").(*schema.Set).List()

		for _, b := range bypassnetworklists {
			bl := appsec.BypassNetworkList{}
			bl.ID = b.(string)
			matchtargets.BypassNetworkLists = append(matchtargets.BypassNetworkLists, bl)
		}
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
	CorrelationID := "[APPSEC][resourceMatchTargetsUpdate-" + CreateNonce() + "]"
	edge.PrintfCorrelation("[DEBUG]", CorrelationID, " Updating MatchTargets")

	matchtargets := appsec.NewMatchTargetsResponse()

	matchtargets.ConfigID = d.Get("config_id").(int)
	matchtargets.ConfigVersion = d.Get("version").(int)
	matchtargets.Type = d.Get("type").(string)
	matchtargets.Sequence = d.Get("sequence").(int)
	matchtargets.IsNegativePathMatch = d.Get("is_negative_path_match").(bool)
	matchtargets.IsNegativeFileExtensionMatch = d.Get("is_negative_file_extension_match").(bool)
	matchtargets.DefaultFile = d.Get("default_file").(string)
	matchtargets.Hostnames = setToStringSlice(d.Get("host_names").(*schema.Set))
	matchtargets.FilePaths = setToStringSlice(d.Get("file_paths").(*schema.Set))
	matchtargets.FileExtensions = setToStringSlice(d.Get("file_extensions").(*schema.Set))
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
	CorrelationID := "[APPSEC][resourceMatchTargetsDelete-" + CreateNonce() + "]"
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
	CorrelationID := "[APPSEC][resourceMatchTargetsRead-" + CreateNonce() + "]"
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
	_, ok := d.GetOk("json")
	if ok {

		j, err := json.Marshal(matchtargets)
		if err == nil {
			d.Set("json", string(j))
		}

	} else {
		d.Set("type", matchtargets.Type)
		d.Set("sequence", matchtargets.Sequence)
		d.Set("is_negative_path_match", matchtargets.IsNegativePathMatch)
		d.Set("is_negative_file_extension_match", matchtargets.IsNegativeFileExtensionMatch)
		d.Set("default_file", matchtargets.DefaultFile)
		d.Set("host_names", matchtargets.Hostnames)
		d.Set("file_paths", matchtargets.FilePaths)
		d.Set("file_extensions", matchtargets.FileExtensions)
		d.Set("security_policy", matchtargets.SecurityPolicy.PolicyID)
		d.Set("target_id", matchtargets.TargetID)
	}
	d.SetId(strconv.Itoa(matchtargets.TargetID))

	return nil
}
