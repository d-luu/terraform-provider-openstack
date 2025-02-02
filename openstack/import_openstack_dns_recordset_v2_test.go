package openstack

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDNSV2RecordSet_importBasic(t *testing.T) {
	zoneName := randomZoneName()
	resourceName := "openstack_dns_recordset_v2.recordset_1"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckDNS(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDNSV2RecordSetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDNSV2RecordSet_basic(zoneName),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
