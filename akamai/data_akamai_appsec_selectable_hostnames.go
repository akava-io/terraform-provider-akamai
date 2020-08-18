package akamai

import (
	"fmt"
	"strconv"

	appsec "github.com/akamai/AkamaiOPEN-edgegrid-golang/appsec-v1"
	edge "github.com/akamai/AkamaiOPEN-edgegrid-golang/edgegrid"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/jsonhooks-v1"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceSelectableHostnames() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceSelectableHostnamesRead,
		Schema: map[string]*schema.Schema{
			"config_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"version": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
			},
			"host_names": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "JSON List of hostnames",
			},
		},
	}
}

func dataSourceSelectableHostnamesRead(d *schema.ResourceData, meta interface{}) error {
	CorrelationID := "[APPSEC][dataSourceSelectableHostnamesRead-" + CreateNonce() + "]"

	edge.PrintfCorrelation("[DEBUG]", CorrelationID, "  Read SelectableHostnames")

	selectablehostnames := appsec.NewSelectableHostnamesResponse()
	selectablehostnames.ConfigID = d.Get("config_id").(int)
	selectablehostnames.ConfigVersion = d.Get("version").(int)

	err := selectablehostnames.GetSelectableHostnames(CorrelationID)
	if err != nil {
		edge.PrintfCorrelation("[DEBUG]", CorrelationID, fmt.Sprintf("Error  %v\n", err))
		return nil
	}

	edge.PrintfCorrelation("[DEBUG]", CorrelationID, fmt.Sprintf("SelectableHostnames   %v\n", selectablehostnames))

	jsonBody, err := jsonhooks.Marshal(selectablehostnames)
	if err != nil {
		return err
	}
	d.Set("host_names", string(jsonBody))
	d.SetId(strconv.Itoa(selectablehostnames.ConfigID))

	return nil
}
