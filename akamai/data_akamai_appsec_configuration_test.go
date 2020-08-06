package akamai

import (
	"log"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccAkamaiConfiguration_basic(t *testing.T) {
	dataSourceName := "data.akamai_appsec_configuration.appsecconfiguration"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAkamaiConfigurationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAkamaiConfigurationConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "id"),
				),
			},
		},
	})
}

func testAccAkamaiConfigurationConfig() string {
	return `
provider "akamai" {
  appsec_section = "default"
}


data "akamai_appsec_configuration" "appsecconfiguration" {
    name = "Akamai Tools"
   }


`
}

func testAccCheckAkamaiConfigurationDestroy(s *terraform.State) error {
	log.Printf("[DEBUG] [Akamai Configuration] Configuration Destroy skipped ")
	return nil
}