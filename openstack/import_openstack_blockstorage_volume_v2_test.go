package openstack

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccBlockStorageV2Volume_importBasic(t *testing.T) {
	resourceName := "openstack_blockstorage_volume_v2.volume_1"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBlockStorageV2VolumeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBlockStorageV2Volume_basic,
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
