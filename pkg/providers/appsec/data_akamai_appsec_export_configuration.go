package appsec

import (
	"fmt"
	"strconv"

	appsec "github.com/akamai/AkamaiOPEN-edgegrid-golang/appsec-v1"
	edge "github.com/akamai/AkamaiOPEN-edgegrid-golang/edgegrid"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/jsonhooks-v1"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/tools"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceExportConfiguration() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceExportConfigurationRead,
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
				Type:        schema.TypeString,
				Computed:    true,
				Description: "JSON Export representation",
			},
		},
	}
}

func dataSourceExportConfigurationRead(d *schema.ResourceData, meta interface{}) error {
	CorrelationID := "[APPSEC][dataSourceExportConfigurationRead-" + tools.CreateNonce() + "]"

	edge.PrintfCorrelation("[DEBUG]", CorrelationID, "  Read ExportConfiguration")

	exportconfiguration := appsec.NewExportConfigurationResponse()
	exportconfiguration.ConfigID = d.Get("config_id").(int)
	exportconfiguration.Version = d.Get("version").(int)

	err := exportconfiguration.GetExportConfiguration(CorrelationID)
	if err != nil {
		edge.PrintfCorrelation("[DEBUG]", CorrelationID, fmt.Sprintf("Error  %v\n", err))
		return nil
	}

	edge.PrintfCorrelation("[DEBUG]", CorrelationID, fmt.Sprintf("ExportConfiguration   %v\n", exportconfiguration))

	jsonBody, err := jsonhooks.Marshal(exportconfiguration)
	if err != nil {
		return err
	}

	d.Set("json", string(jsonBody))
	d.SetId(strconv.Itoa(exportconfiguration.ConfigID))

	return nil
}
