package openstack

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccFWFirewallV1_importBasic(t *testing.T) {
	resourceName := "openstack_fw_firewall_v1.fw_1"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckFW(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckFWFirewallV1Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccFWFirewallV1_basic_1,
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
