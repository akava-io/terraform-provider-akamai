package akamai

import (
	"fmt"
	"strconv"

	appsec "github.com/akamai/AkamaiOPEN-edgegrid-golang/appsec-v1"
	edge "github.com/akamai/AkamaiOPEN-edgegrid-golang/edgegrid"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// appsec v1
//
// https://developer.akamai.com/api/cloud_security/application_security/v1.html
func resourceSelectedHostnames() *schema.Resource {
	return &schema.Resource{
		Create: resourceSelectedHostnamesUpdate,
		Read:   resourceSelectedHostnamesRead,
		Update: resourceSelectedHostnamesUpdate,
		Delete: resourceSelectedHostnamesDelete,
		Schema: map[string]*schema.Schema{
			"config_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"version": {
				Type:             schema.TypeInt,
				Required:         true,
				DiffSuppressFunc: suppressConfigurationCloneVersion,
			},
			"host_names": {
				Type:     schema.TypeList,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceSelectedHostnamesRead(d *schema.ResourceData, meta interface{}) error {
	CorrelationID := "[APPSEC][resourceSelectedHostnamesRead-" + CreateNonce() + "]"
	edge.PrintfCorrelation("[DEBUG]", CorrelationID, "  Read SelectedHostnames")

	selectedhostnames := appsec.NewSelectedHostnamesResponse()
	configid := d.Get("config_id").(int)
	version := d.Get("version").(int)
	err := selectedhostnames.GetSelectedHostnames(configid, version, CorrelationID)
	if err != nil {
		edge.PrintfCorrelation("[DEBUG]", CorrelationID, fmt.Sprintf("Error  %v\n", err))
		return nil
	}

	d.SetId(strconv.Itoa(configid))
	return nil
}

func resourceSelectedHostnamesDelete(d *schema.ResourceData, meta interface{}) error {
	CorrelationID := "[APPSEC][resourceSelectedHostnamesDelete-" + CreateNonce() + "]"
	edge.PrintfCorrelation("[DEBUG]", CorrelationID, "  Deleting SelectedHostnames")
	return schema.Noop(d, meta)
}

func resourceSelectedHostnamesUpdate(d *schema.ResourceData, meta interface{}) error {
	CorrelationID := "[APPSEC][resourceSelectedHostnamesUpdate-" + CreateNonce() + "]"
	edge.PrintfCorrelation("[DEBUG]", CorrelationID, "  Updating SelectedHostnames")
	selectedhostnames := appsec.NewSelectedHostnamesResponse()
	configid := d.Get("config_id").(int)
	version := d.Get("version").(int)
	hn := &appsec.SelectedHostnamesResponse{}

	hostnamelist := d.Get("host_names").([]interface{})

	for _, h := range hostnamelist {
		m := appsec.Hostname{}
		m.Hostname = h.(string)
		hn.HostnameList = append(hn.HostnameList, m)
	}

	selectedhostnames.HostnameList = hn.HostnameList
	err := selectedhostnames.UpdateSelectedHostnames(configid, version, CorrelationID)
	if err != nil {
		edge.PrintfCorrelation("[DEBUG]", CorrelationID, fmt.Sprintf("Error  %v\n", err))
		return err
	}
	return resourceSelectedHostnamesRead(d, meta)

}
