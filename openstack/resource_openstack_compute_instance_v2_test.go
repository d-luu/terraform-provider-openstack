package openstack

import (
	"fmt"
	"reflect"
	"sort"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/blockstorage/v2/volumes"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/secgroups"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/volumeattach"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/servers"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/networks"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/ports"
	"github.com/gophercloud/gophercloud/pagination"
)

func TestAccComputeV2Instance_basic(t *testing.T) {
	var instance servers.Server

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeV2InstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeV2Instance_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2InstanceExists("openstack_compute_instance_v2.instance_1", &instance),
					testAccCheckComputeV2InstanceMetadata(&instance, "foo", "bar"),
					resource.TestCheckResourceAttr(
						"openstack_compute_instance_v2.instance_1", "all_metadata.foo", "bar"),
					resource.TestCheckResourceAttr(
						"openstack_compute_instance_v2.instance_1", "availability_zone", "nova"),
				),
			},
		},
	})
}

func TestAccComputeV2Instance_initialStateActive(t *testing.T) {
	var instance servers.Server

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeV2InstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeV2Instance_stateActive,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2InstanceExists("openstack_compute_instance_v2.instance_1", &instance),
					resource.TestCheckResourceAttr(
						"openstack_compute_instance_v2.instance_1", "power_state", "active"),
					testAccCheckComputeV2InstanceState(&instance, "active"),
				),
			},
			{
				Config: testAccComputeV2Instance_stateShutoff,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2InstanceExists("openstack_compute_instance_v2.instance_1", &instance),
					resource.TestCheckResourceAttr(
						"openstack_compute_instance_v2.instance_1", "power_state", "shutoff"),
					testAccCheckComputeV2InstanceState(&instance, "shutoff"),
				),
			},
			{
				Config: testAccComputeV2Instance_stateActive,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2InstanceExists("openstack_compute_instance_v2.instance_1", &instance),
					resource.TestCheckResourceAttr(
						"openstack_compute_instance_v2.instance_1", "power_state", "active"),
					testAccCheckComputeV2InstanceState(&instance, "active"),
				),
			},
		},
	})
}

func TestAccComputeV2Instance_initialStateShutoff(t *testing.T) {
	var instance servers.Server

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeV2InstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeV2Instance_stateShutoff,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2InstanceExists("openstack_compute_instance_v2.instance_1", &instance),
					resource.TestCheckResourceAttr(
						"openstack_compute_instance_v2.instance_1", "power_state", "shutoff"),
					testAccCheckComputeV2InstanceState(&instance, "shutoff"),
				),
			},
			{
				Config: testAccComputeV2Instance_stateActive,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2InstanceExists("openstack_compute_instance_v2.instance_1", &instance),
					resource.TestCheckResourceAttr(
						"openstack_compute_instance_v2.instance_1", "power_state", "active"),
					testAccCheckComputeV2InstanceState(&instance, "active"),
				),
			},
			{
				Config: testAccComputeV2Instance_stateShutoff,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2InstanceExists("openstack_compute_instance_v2.instance_1", &instance),
					resource.TestCheckResourceAttr(
						"openstack_compute_instance_v2.instance_1", "power_state", "shutoff"),
					testAccCheckComputeV2InstanceState(&instance, "shutoff"),
				),
			},
		},
	})
}

func TestAccComputeV2Instance_secgroupMulti(t *testing.T) {
	var instance_1 servers.Server
	var secgroup_1 secgroups.SecurityGroup

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeV2InstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeV2Instance_secgroupMulti,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2SecGroupExists(
						"openstack_compute_secgroup_v2.secgroup_1", &secgroup_1),
					testAccCheckComputeV2InstanceExists(
						"openstack_compute_instance_v2.instance_1", &instance_1),
				),
			},
		},
	})
}

func TestAccComputeV2Instance_secgroupMultiUpdate(t *testing.T) {
	var instance_1 servers.Server
	var secgroup_1, secgroup_2 secgroups.SecurityGroup

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeV2InstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeV2Instance_secgroupMultiUpdate_1,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2SecGroupExists(
						"openstack_compute_secgroup_v2.secgroup_1", &secgroup_1),
					testAccCheckComputeV2SecGroupExists(
						"openstack_compute_secgroup_v2.secgroup_2", &secgroup_2),
					testAccCheckComputeV2InstanceExists(
						"openstack_compute_instance_v2.instance_1", &instance_1),
				),
			},
			{
				Config: testAccComputeV2Instance_secgroupMultiUpdate_2,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2SecGroupExists(
						"openstack_compute_secgroup_v2.secgroup_1", &secgroup_1),
					testAccCheckComputeV2SecGroupExists(
						"openstack_compute_secgroup_v2.secgroup_2", &secgroup_2),
					testAccCheckComputeV2InstanceExists(
						"openstack_compute_instance_v2.instance_1", &instance_1),
				),
			},
		},
	})
}

func TestAccComputeV2Instance_bootFromVolumeImage(t *testing.T) {
	var instance servers.Server

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeV2InstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeV2Instance_bootFromVolumeImage,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2InstanceExists("openstack_compute_instance_v2.instance_1", &instance),
					testAccCheckComputeV2InstanceBootVolumeAttachment(&instance),
				),
			},
		},
	})
}

func TestAccComputeV2Instance_bootFromVolumeVolume(t *testing.T) {
	var instance servers.Server

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeV2InstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeV2Instance_bootFromVolumeVolume,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2InstanceExists("openstack_compute_instance_v2.instance_1", &instance),
					testAccCheckComputeV2InstanceBootVolumeAttachment(&instance),
				),
			},
		},
	})
}

func TestAccComputeV2Instance_bootFromVolumeForceNew(t *testing.T) {
	var instance1_1 servers.Server
	var instance1_2 servers.Server

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeV2InstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeV2Instance_bootFromVolumeForceNew_1,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2InstanceExists(
						"openstack_compute_instance_v2.instance_1", &instance1_1),
				),
			},
			{
				Config: testAccComputeV2Instance_bootFromVolumeForceNew_2,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2InstanceExists(
						"openstack_compute_instance_v2.instance_1", &instance1_2),
					testAccCheckComputeV2InstanceInstanceIDsDoNotMatch(&instance1_1, &instance1_2),
				),
			},
		},
	})
}

func TestAccComputeV2Instance_blockDeviceNewVolume(t *testing.T) {
	var instance servers.Server

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeV2InstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeV2Instance_blockDeviceNewVolume,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2InstanceExists("openstack_compute_instance_v2.instance_1", &instance),
				),
			},
		},
	})
}

func TestAccComputeV2Instance_blockDeviceNewVolumeTypeAndBus(t *testing.T) {
	var instance servers.Server

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeV2InstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeV2Instance_blockDeviceNewVolumeTypeAndBus,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2InstanceExists("openstack_compute_instance_v2.instance_1", &instance),
				),
			},
		},
	})
}

func TestAccComputeV2Instance_blockDeviceExistingVolume(t *testing.T) {
	var instance servers.Server
	var volume volumes.Volume

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeV2InstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeV2Instance_blockDeviceExistingVolume,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2InstanceExists("openstack_compute_instance_v2.instance_1", &instance),
					testAccCheckBlockStorageV2VolumeExists(
						"openstack_blockstorage_volume_v2.volume_1", &volume),
				),
			},
		},
	})
}

// TODO: verify the personality really exists on the instance.
func TestAccComputeV2Instance_personality(t *testing.T) {
	var instance servers.Server

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeV2InstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeV2Instance_personality,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2InstanceExists("openstack_compute_instance_v2.instance_1", &instance),
				),
			},
		},
	})
}

func TestAccComputeV2Instance_multiEphemeral(t *testing.T) {
	var instance servers.Server

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeV2InstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeV2Instance_multiEphemeral,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2InstanceExists(
						"openstack_compute_instance_v2.instance_1", &instance),
				),
			},
		},
	})
}

func TestAccComputeV2Instance_accessIPv4(t *testing.T) {
	var instance servers.Server

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeV2InstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeV2Instance_accessIPv4,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2InstanceExists("openstack_compute_instance_v2.instance_1", &instance),
					resource.TestCheckResourceAttr(
						"openstack_compute_instance_v2.instance_1", "access_ip_v4", "192.168.1.100"),
				),
			},
		},
	})
}

func TestAccComputeV2Instance_changeFixedIP(t *testing.T) {
	var instance1_1 servers.Server
	var instance1_2 servers.Server

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeV2InstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeV2Instance_changeFixedIP_1,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2InstanceExists(
						"openstack_compute_instance_v2.instance_1", &instance1_1),
				),
			},
			{
				Config: testAccComputeV2Instance_changeFixedIP_2,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2InstanceExists(
						"openstack_compute_instance_v2.instance_1", &instance1_2),
					testAccCheckComputeV2InstanceInstanceIDsDoNotMatch(&instance1_1, &instance1_2),
				),
			},
		},
	})
}

func TestAccComputeV2Instance_stopBeforeDestroy(t *testing.T) {
	var instance servers.Server
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeV2InstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeV2Instance_stopBeforeDestroy,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2InstanceExists("openstack_compute_instance_v2.instance_1", &instance),
				),
			},
		},
	})
}

func TestAccComputeV2Instance_metadataRemove(t *testing.T) {
	var instance servers.Server

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeV2InstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeV2Instance_metadataRemove_1,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2InstanceExists("openstack_compute_instance_v2.instance_1", &instance),
					testAccCheckComputeV2InstanceMetadata(&instance, "foo", "bar"),
					testAccCheckComputeV2InstanceMetadata(&instance, "abc", "def"),
					resource.TestCheckResourceAttr(
						"openstack_compute_instance_v2.instance_1", "all_metadata.foo", "bar"),
					resource.TestCheckResourceAttr(
						"openstack_compute_instance_v2.instance_1", "all_metadata.abc", "def"),
				),
			},
			{
				Config: testAccComputeV2Instance_metadataRemove_2,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2InstanceExists("openstack_compute_instance_v2.instance_1", &instance),
					testAccCheckComputeV2InstanceMetadata(&instance, "foo", "bar"),
					testAccCheckComputeV2InstanceMetadata(&instance, "ghi", "jkl"),
					testAccCheckComputeV2InstanceNoMetadataKey(&instance, "abc"),
					resource.TestCheckResourceAttr(
						"openstack_compute_instance_v2.instance_1", "all_metadata.foo", "bar"),
					resource.TestCheckResourceAttr(
						"openstack_compute_instance_v2.instance_1", "all_metadata.ghi", "jkl"),
				),
			},
		},
	})
}

func TestAccComputeV2Instance_forceDelete(t *testing.T) {
	var instance servers.Server
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeV2InstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeV2Instance_forceDelete,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2InstanceExists("openstack_compute_instance_v2.instance_1", &instance),
				),
			},
		},
	})
}

func TestAccComputeV2Instance_timeout(t *testing.T) {
	var instance servers.Server
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeV2InstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeV2Instance_timeout,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2InstanceExists("openstack_compute_instance_v2.instance_1", &instance),
				),
			},
		},
	})
}

func TestAccComputeV2Instance_networkModeAuto(t *testing.T) {
	var instance servers.Server
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeV2InstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeV2Instance_networkModeAuto,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2InstanceExists("openstack_compute_instance_v2.instance_1", &instance),
					testAccCheckComputeV2InstanceNetworkExists("openstack_compute_instance_v2.instance_1", &instance),
				),
			},
		},
	})
}

func TestAccComputeV2Instance_networkModeNone(t *testing.T) {
	var instance servers.Server
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeV2InstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeV2Instance_networkModeNone,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2InstanceExists("openstack_compute_instance_v2.instance_1", &instance),
					testAccCheckComputeV2InstanceNetworkDoesNotExist("openstack_compute_instance_v2.instance_1", &instance),
				),
			},
		},
	})
}

func TestAccComputeV2Instance_networkNameToID(t *testing.T) {
	var instance servers.Server
	var network networks.Network
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeV2InstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeV2Instance_networkNameToID,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2InstanceExists("openstack_compute_instance_v2.instance_1", &instance),
					testAccCheckNetworkingV2NetworkExists("openstack_networking_network_v2.network_1", &network),
					resource.TestCheckResourceAttrPtr(
						"openstack_compute_instance_v2.instance_1", "network.1.uuid", &network.ID),
				),
			},
		},
	})
}

func TestAccComputeV2Instance_crazyNICs(t *testing.T) {
	var instance servers.Server
	var network_1 networks.Network
	var network_2 networks.Network
	var port_1 ports.Port
	var port_2 ports.Port
	var port_3 ports.Port
	var port_4 ports.Port

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeV2InstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeV2Instance_crazyNICs,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2InstanceExists("openstack_compute_instance_v2.instance_1", &instance),
					testAccCheckNetworkingV2NetworkExists(
						"openstack_networking_network_v2.network_1", &network_1),
					testAccCheckNetworkingV2NetworkExists(
						"openstack_networking_network_v2.network_2", &network_2),
					testAccCheckNetworkingV2PortExists(
						"openstack_networking_port_v2.port_1", &port_1),
					testAccCheckNetworkingV2PortExists(
						"openstack_networking_port_v2.port_2", &port_2),
					testAccCheckNetworkingV2PortExists(
						"openstack_networking_port_v2.port_3", &port_3),
					testAccCheckNetworkingV2PortExists(
						"openstack_networking_port_v2.port_4", &port_4),
					resource.TestCheckResourceAttrPtr(
						"openstack_compute_instance_v2.instance_1", "network.1.uuid", &network_1.ID),
					resource.TestCheckResourceAttrPtr(
						"openstack_compute_instance_v2.instance_1", "network.2.uuid", &network_2.ID),
					resource.TestCheckResourceAttrPtr(
						"openstack_compute_instance_v2.instance_1", "network.3.uuid", &network_1.ID),
					resource.TestCheckResourceAttrPtr(
						"openstack_compute_instance_v2.instance_1", "network.4.uuid", &network_2.ID),
					resource.TestCheckResourceAttrPtr(
						"openstack_compute_instance_v2.instance_1", "network.5.uuid", &network_1.ID),
					resource.TestCheckResourceAttrPtr(
						"openstack_compute_instance_v2.instance_1", "network.6.uuid", &network_2.ID),
					resource.TestCheckResourceAttr(
						"openstack_compute_instance_v2.instance_1", "network.1.name", "network_1"),
					resource.TestCheckResourceAttr(
						"openstack_compute_instance_v2.instance_1", "network.2.name", "network_2"),
					resource.TestCheckResourceAttr(
						"openstack_compute_instance_v2.instance_1", "network.3.name", "network_1"),
					resource.TestCheckResourceAttr(
						"openstack_compute_instance_v2.instance_1", "network.4.name", "network_2"),
					resource.TestCheckResourceAttr(
						"openstack_compute_instance_v2.instance_1", "network.5.name", "network_1"),
					resource.TestCheckResourceAttr(
						"openstack_compute_instance_v2.instance_1", "network.6.name", "network_2"),
					resource.TestCheckResourceAttr(
						"openstack_compute_instance_v2.instance_1", "network.7.name", "network_1"),
					resource.TestCheckResourceAttr(
						"openstack_compute_instance_v2.instance_1", "network.8.name", "network_2"),
					resource.TestCheckResourceAttr(
						"openstack_compute_instance_v2.instance_1", "network.1.fixed_ip_v4", "192.168.1.100"),
					resource.TestCheckResourceAttr(
						"openstack_compute_instance_v2.instance_1", "network.2.fixed_ip_v4", "192.168.2.100"),
					resource.TestCheckResourceAttr(
						"openstack_compute_instance_v2.instance_1", "network.3.fixed_ip_v4", "192.168.1.101"),
					resource.TestCheckResourceAttr(
						"openstack_compute_instance_v2.instance_1", "network.4.fixed_ip_v4", "192.168.2.101"),
					resource.TestCheckResourceAttrPtr(
						"openstack_compute_instance_v2.instance_1", "network.5.port", &port_1.ID),
					resource.TestCheckResourceAttrPtr(
						"openstack_compute_instance_v2.instance_1", "network.6.port", &port_2.ID),
					resource.TestCheckResourceAttrPtr(
						"openstack_compute_instance_v2.instance_1", "network.7.port", &port_3.ID),
					resource.TestCheckResourceAttrPtr(
						"openstack_compute_instance_v2.instance_1", "network.8.port", &port_4.ID),
				),
			},
		},
	})
}

func TestAccComputeV2Instance_tags(t *testing.T) {
	var instance servers.Server

	resourceName := "openstack_compute_instance_v2.instance_1"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNetworkingV2NetworkDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeV2Instance_tags_create,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2InstanceExists(resourceName, &instance),
					testAccCheckComputeV2InstanceTags(resourceName, []string{"tag1", "tag2", "tag3"}),
				),
			},
			{
				Config: testAccComputeV2Instance_tags_add,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2InstanceExists(resourceName, &instance),
					testAccCheckComputeV2InstanceTags(resourceName, []string{"tag1", "tag2", "tag3", "tag4"}),
				),
			},
			{
				Config: testAccComputeV2Instance_tags_delete,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2InstanceExists(resourceName, &instance),
					testAccCheckComputeV2InstanceTags(resourceName, []string{"tag2", "tag3"}),
				),
			},
			{
				Config: testAccComputeV2Instance_tags_clear,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2InstanceExists(resourceName, &instance),
					testAccCheckComputeV2InstanceTags(resourceName, nil),
				),
			},
		},
	})
}

func testAccCheckComputeV2InstanceDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	computeClient, err := config.ComputeV2Client(OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating OpenStack compute client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "openstack_compute_instance_v2" {
			continue
		}

		server, err := servers.Get(computeClient, rs.Primary.ID).Extract()
		if err == nil {
			if server.Status != "SOFT_DELETED" {
				return fmt.Errorf("Instance still exists")
			}
		}
	}

	return nil
}

func testAccCheckComputeV2InstanceExists(n string, instance *servers.Server) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)
		computeClient, err := config.ComputeV2Client(OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating OpenStack compute client: %s", err)
		}

		found, err := servers.Get(computeClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("Instance not found")
		}

		*instance = *found

		return nil
	}
}

func testAccCheckComputeV2InstanceDoesNotExist(n string, instance *servers.Server) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := testAccProvider.Meta().(*Config)
		computeClient, err := config.ComputeV2Client(OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating OpenStack compute client: %s", err)
		}

		_, err = servers.Get(computeClient, instance.ID).Extract()
		if err != nil {
			if _, ok := err.(gophercloud.ErrDefault404); ok {
				return nil
			}
			return err
		}

		return fmt.Errorf("Instance still exists")
	}
}

func testAccCheckComputeV2InstanceMetadata(
	instance *servers.Server, k string, v string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if instance.Metadata == nil {
			return fmt.Errorf("No metadata")
		}

		for key, value := range instance.Metadata {
			if k != key {
				continue
			}

			if v == value {
				return nil
			}

			return fmt.Errorf("Bad value for %s: %s", k, value)
		}

		return fmt.Errorf("Metadata not found: %s", k)
	}
}

func testAccCheckComputeV2InstanceNoMetadataKey(
	instance *servers.Server, k string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if instance.Metadata == nil {
			return nil
		}

		for key := range instance.Metadata {
			if k == key {
				return fmt.Errorf("Metadata found: %s", k)
			}
		}

		return nil
	}
}

func testAccCheckComputeV2InstanceBootVolumeAttachment(
	instance *servers.Server) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		var attachments []volumeattach.VolumeAttachment

		config := testAccProvider.Meta().(*Config)
		computeClient, err := config.ComputeV2Client(OS_REGION_NAME)
		if err != nil {
			return err
		}

		err = volumeattach.List(computeClient, instance.ID).EachPage(
			func(page pagination.Page) (bool, error) {

				actual, err := volumeattach.ExtractVolumeAttachments(page)
				if err != nil {
					return false, fmt.Errorf("Unable to lookup attachment: %s", err)
				}

				attachments = actual
				return true, nil
			})

		if len(attachments) == 1 {
			return nil
		}

		return fmt.Errorf("No attached volume found.")
	}
}

func testAccCheckComputeV2InstanceInstanceIDsDoNotMatch(
	instance1, instance2 *servers.Server) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if instance1.ID == instance2.ID {
			return fmt.Errorf("Instance was not recreated.")
		}

		return nil
	}
}

func testAccCheckComputeV2InstanceState(
	instance *servers.Server, state string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if strings.ToLower(instance.Status) != state {
			return fmt.Errorf("Instance state is not match.")
		}

		return nil
	}
}

func testAccCheckComputeV2InstanceTags(name string, tags []string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("resource not found: %s", name)
		}

		if _, ok := rs.Primary.Attributes["tags.#"]; !ok {
			return fmt.Errorf("resource tags not found: %s.tags", name)
		}

		var rtags []string
		for key, val := range rs.Primary.Attributes {
			if !strings.HasPrefix(key, "tags.") {
				continue
			}

			if key == "tags.#" {
				continue
			}

			rtags = append(rtags, val)
		}

		sort.Strings(rtags)
		sort.Strings(tags)
		if !reflect.DeepEqual(rtags, tags) {
			return fmt.Errorf(
				"%s.tags: expected: %#v, got %#v", name, tags, rtags)
		}
		return nil
	}
}

func testAccCheckComputeV2InstanceNetworkExists(
	n string, instance *servers.Server) resource.TestCheckFunc {

	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("resource not found: %s", n)
		}

		networkCount, ok := rs.Primary.Attributes["network.#"]

		if !ok {
			return fmt.Errorf("network attributes not found: %s", n)
		}

		if networkCount != "1" {
			return fmt.Errorf("network should be exists when network mode 'auto': %s", n)
		}

		return nil
	}
}

func testAccCheckComputeV2InstanceNetworkDoesNotExist(
	n string, instance *servers.Server) resource.TestCheckFunc {

	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("resource not found: %s", n)
		}

		networkCount, ok := rs.Primary.Attributes["network.#"]

		if !ok {
			return fmt.Errorf("network attributes not found: %s", n)
		}

		if networkCount != "0" {
			return fmt.Errorf("network should not exists when network mode 'none': %s", n)
		}

		return nil
	}
}

var testAccComputeV2Instance_basic = fmt.Sprintf(`
resource "openstack_compute_instance_v2" "instance_1" {
  name = "instance_1"
  security_groups = ["default"]
  metadata = {
    foo = "bar"
  }
  network {
    uuid = "%s"
  }
}
`, OS_NETWORK_ID)

var testAccComputeV2Instance_secgroupMulti = fmt.Sprintf(`
resource "openstack_compute_secgroup_v2" "secgroup_1" {
  name = "secgroup_1"
  description = "a security group"
  rule {
    from_port = 22
    to_port = 22
    ip_protocol = "tcp"
    cidr = "0.0.0.0/0"
  }
}

resource "openstack_compute_instance_v2" "instance_1" {
  name = "instance_1"
  security_groups = ["default", "${openstack_compute_secgroup_v2.secgroup_1.name}"]
  network {
    uuid = "%s"
  }
}
`, OS_NETWORK_ID)

var testAccComputeV2Instance_secgroupMultiUpdate_1 = fmt.Sprintf(`
resource "openstack_compute_secgroup_v2" "secgroup_1" {
  name = "secgroup_1"
  description = "a security group"
  rule {
    from_port = 22
    to_port = 22
    ip_protocol = "tcp"
    cidr = "0.0.0.0/0"
  }
}

resource "openstack_compute_secgroup_v2" "secgroup_2" {
  name = "secgroup_2"
  description = "another security group"
  rule {
    from_port = 80
    to_port = 80
    ip_protocol = "tcp"
    cidr = "0.0.0.0/0"
  }
}

resource "openstack_compute_instance_v2" "instance_1" {
  name = "instance_1"
  security_groups = ["default"]
  network {
    uuid = "%s"
  }
}
`, OS_NETWORK_ID)

var testAccComputeV2Instance_secgroupMultiUpdate_2 = fmt.Sprintf(`
resource "openstack_compute_secgroup_v2" "secgroup_1" {
  name = "secgroup_1"
  description = "a security group"
  rule {
    from_port = 22
    to_port = 22
    ip_protocol = "tcp"
    cidr = "0.0.0.0/0"
  }
}

resource "openstack_compute_secgroup_v2" "secgroup_2" {
  name = "secgroup_2"
  description = "another security group"
  rule {
    from_port = 80
    to_port = 80
    ip_protocol = "tcp"
    cidr = "0.0.0.0/0"
  }
}

resource "openstack_compute_instance_v2" "instance_1" {
  name = "instance_1"
  security_groups = ["default", "${openstack_compute_secgroup_v2.secgroup_1.name}", "${openstack_compute_secgroup_v2.secgroup_2.name}"]
  network {
    uuid = "%s"
  }
}
`, OS_NETWORK_ID)

var testAccComputeV2Instance_bootFromVolumeImage = fmt.Sprintf(`
resource "openstack_compute_instance_v2" "instance_1" {
  name = "instance_1"
  security_groups = ["default"]
  block_device {
    uuid = "%s"
    source_type = "image"
    volume_size = 5
    boot_index = 0
    destination_type = "volume"
    delete_on_termination = true
  }
  network {
    uuid = "%s"
  }
}
`, OS_IMAGE_ID, OS_NETWORK_ID)

var testAccComputeV2Instance_bootFromVolumeVolume = fmt.Sprintf(`
resource "openstack_blockstorage_volume_v2" "vol_1" {
  name = "vol_1"
  size = 5
  image_id = "%s"
}

resource "openstack_compute_instance_v2" "instance_1" {
  name = "instance_1"
  security_groups = ["default"]
  block_device {
    uuid = "${openstack_blockstorage_volume_v2.vol_1.id}"
    source_type = "volume"
    boot_index = 0
    destination_type = "volume"
    delete_on_termination = true
  }
  network {
    uuid = "%s"
  }
}
`, OS_IMAGE_ID, OS_NETWORK_ID)

var testAccComputeV2Instance_bootFromVolumeForceNew_1 = fmt.Sprintf(`
resource "openstack_compute_instance_v2" "instance_1" {
  name = "instance_1"
  security_groups = ["default"]
  block_device {
    uuid = "%s"
    source_type = "image"
    volume_size = 5
    boot_index = 0
    destination_type = "volume"
    delete_on_termination = true
  }
  network {
    uuid = "%s"
  }
}
`, OS_IMAGE_ID, OS_NETWORK_ID)

var testAccComputeV2Instance_bootFromVolumeForceNew_2 = fmt.Sprintf(`
resource "openstack_compute_instance_v2" "instance_1" {
  name = "instance_1"
  security_groups = ["default"]
  block_device {
    uuid = "%s"
    source_type = "image"
    volume_size = 4
    boot_index = 0
    destination_type = "volume"
    delete_on_termination = true
  }
  network {
    uuid = "%s"
  }
}
`, OS_IMAGE_ID, OS_NETWORK_ID)

var testAccComputeV2Instance_blockDeviceNewVolume = fmt.Sprintf(`
resource "openstack_compute_instance_v2" "instance_1" {
  name = "instance_1"
  security_groups = ["default"]
  block_device {
    uuid = "%s"
    source_type = "image"
    destination_type = "local"
    boot_index = 0
    delete_on_termination = true
  }
  block_device {
    source_type = "blank"
    destination_type = "volume"
    volume_size = 1
    boot_index = 1
    delete_on_termination = true
  }
  network {
    uuid = "%s"
  }
}
`, OS_IMAGE_ID, OS_NETWORK_ID)

var testAccComputeV2Instance_blockDeviceNewVolumeTypeAndBus = fmt.Sprintf(`
resource "openstack_compute_instance_v2" "instance_1" {
  name = "instance_1"
  security_groups = ["default"]
  block_device {
    uuid = "%s"
    source_type = "image"
    destination_type = "local"
    boot_index = 0
		delete_on_termination = true
		device_type = "disk"
		disk_bus = "virtio"
  }
  block_device {
    source_type = "blank"
    destination_type = "volume"
    volume_size = 1
    boot_index = 1
		delete_on_termination = true
		device_type = "disk"
		disk_bus = "virtio"
  }
  network {
    uuid = "%s"
  }
}
`, OS_IMAGE_ID, OS_NETWORK_ID)

var testAccComputeV2Instance_blockDeviceExistingVolume = fmt.Sprintf(`
resource "openstack_blockstorage_volume_v2" "volume_1" {
  name = "volume_1"
  size = 1
}

resource "openstack_compute_instance_v2" "instance_1" {
  name = "instance_1"
  security_groups = ["default"]
  block_device {
    uuid = "%s"
    source_type = "image"
    destination_type = "local"
    boot_index = 0
    delete_on_termination = true
  }
  block_device {
    uuid = "${openstack_blockstorage_volume_v2.volume_1.id}"
    source_type = "volume"
    destination_type = "volume"
    boot_index = 1
    delete_on_termination = true
  }
  network {
    uuid = "%s"
  }
}
`, OS_IMAGE_ID, OS_NETWORK_ID)

var testAccComputeV2Instance_personality = fmt.Sprintf(`
resource "openstack_compute_instance_v2" "instance_1" {
  name = "instance_1"
  security_groups = ["default"]
  personality {
    file = "/tmp/foobar.txt"
    content = "happy"
  }
  personality {
    file = "/tmp/barfoo.txt"
    content = "angry"
  }
  network {
    uuid = "%s"
  }
}
`, OS_NETWORK_ID)

var testAccComputeV2Instance_multiEphemeral = fmt.Sprintf(`
resource "openstack_compute_instance_v2" "instance_1" {
  name = "terraform-test"
  security_groups = ["default"]
  block_device {
    boot_index = 0
    delete_on_termination = true
    destination_type = "local"
    source_type = "image"
    uuid = "%s"
  }
  block_device {
    boot_index = -1
    delete_on_termination = true
    destination_type = "local"
    source_type = "blank"
    volume_size = 1
  }
  block_device {
    boot_index = -1
    delete_on_termination = true
    destination_type = "local"
    source_type = "blank"
    volume_size = 1
  }
  network {
    uuid = "%s"
  }
}
`, OS_IMAGE_ID, OS_NETWORK_ID)

var testAccComputeV2Instance_accessIPv4 = fmt.Sprintf(`
resource "openstack_networking_network_v2" "network_1" {
  name = "network_1"
}

resource "openstack_networking_subnet_v2" "subnet_1" {
  name = "subnet_1"
  network_id = "${openstack_networking_network_v2.network_1.id}"
  cidr = "192.168.1.0/24"
  ip_version = 4
  enable_dhcp = true
  no_gateway = true
}

resource "openstack_compute_instance_v2" "instance_1" {
  depends_on = ["openstack_networking_subnet_v2.subnet_1"]

  name = "instance_1"
  security_groups = ["default"]

  network {
    uuid = "%s"
  }

  network {
    uuid = "${openstack_networking_network_v2.network_1.id}"
    fixed_ip_v4 = "192.168.1.100"
    access_network = true
  }
}
`, OS_NETWORK_ID)

var testAccComputeV2Instance_changeFixedIP_1 = fmt.Sprintf(`
resource "openstack_compute_instance_v2" "instance_1" {
  name = "instance_1"
  security_groups = ["default"]
  network {
    uuid = "%s"
    fixed_ip_v4 = "10.0.0.24"
  }
}
`, OS_NETWORK_ID)

var testAccComputeV2Instance_changeFixedIP_2 = fmt.Sprintf(`
resource "openstack_compute_instance_v2" "instance_1" {
  name = "instance_1"
  security_groups = ["default"]
  network {
    uuid = "%s"
    fixed_ip_v4 = "10.0.0.25"
  }
}
`, OS_NETWORK_ID)

var testAccComputeV2Instance_stopBeforeDestroy = fmt.Sprintf(`
resource "openstack_compute_instance_v2" "instance_1" {
  name = "instance_1"
  security_groups = ["default"]
  stop_before_destroy = true
  network {
    uuid = "%s"
  }
}
`, OS_NETWORK_ID)

var testAccComputeV2Instance_detachPortsBeforeDestroy = fmt.Sprintf(`

resource "openstack_networking_port_v2" "port_1" {
  name = "port_1"
  network_id = "%s"
  admin_state_up = "true"
}


resource "openstack_compute_instance_v2" "instance_1" {
  name = "instance_1"
  security_groups = ["default"]
  vendor_options {
    detach_ports_before_destroy = true
  }
  network {
    port = "${openstack_networking_port_v2.port_1.id}"
  }
}
`, OS_NETWORK_ID)

var testAccComputeV2Instance_metadataRemove_1 = fmt.Sprintf(`
resource "openstack_compute_instance_v2" "instance_1" {
  name = "instance_1"
  security_groups = ["default"]
  metadata = {
    foo = "bar"
    abc = "def"
  }
  network {
    uuid = "%s"
  }
}
`, OS_NETWORK_ID)

var testAccComputeV2Instance_metadataRemove_2 = fmt.Sprintf(`
resource "openstack_compute_instance_v2" "instance_1" {
  name = "instance_1"
  security_groups = ["default"]
  metadata = {
    foo = "bar"
    ghi = "jkl"
  }
  network {
    uuid = "%s"
  }
}
`, OS_NETWORK_ID)

var testAccComputeV2Instance_forceDelete = fmt.Sprintf(`
resource "openstack_compute_instance_v2" "instance_1" {
  name = "instance_1"
  security_groups = ["default"]
  force_delete = true
  network {
    uuid = "%s"
  }
}
`, OS_NETWORK_ID)

var testAccComputeV2Instance_timeout = fmt.Sprintf(`
resource "openstack_compute_instance_v2" "instance_1" {
  name = "instance_1"
  security_groups = ["default"]

  timeouts {
    create = "10m"
  }
  network {
    uuid = "%s"
  }
}
`, OS_NETWORK_ID)

var testAccComputeV2Instance_networkModeAuto = fmt.Sprintf(`
resource "openstack_compute_instance_v2" "instance_1" {
  name = "instance_1"

  network_mode = "auto"
}
`)

var testAccComputeV2Instance_networkModeNone = fmt.Sprintf(`
resource "openstack_compute_instance_v2" "instance_1" {
  name = "test-instance-1"

  network_mode = "none"
}
`)

var testAccComputeV2Instance_networkNameToID = fmt.Sprintf(`
resource "openstack_networking_network_v2" "network_1" {
  name = "network_1"
}

resource "openstack_networking_subnet_v2" "subnet_1" {
  name = "subnet_1"
  network_id = "${openstack_networking_network_v2.network_1.id}"
  cidr = "192.168.1.0/24"
  ip_version = 4
  enable_dhcp = true
  no_gateway = true
}

resource "openstack_compute_instance_v2" "instance_1" {
  depends_on = ["openstack_networking_subnet_v2.subnet_1"]

  name = "instance_1"
  security_groups = ["default"]

  network {
    uuid = "%s"
  }

  network {
    name = "${openstack_networking_network_v2.network_1.name}"
  }

}
`, OS_NETWORK_ID)

var testAccComputeV2Instance_crazyNICs = fmt.Sprintf(`
resource "openstack_networking_network_v2" "network_1" {
  name = "network_1"
}

resource "openstack_networking_subnet_v2" "subnet_1" {
  name = "subnet_1"
  network_id = "${openstack_networking_network_v2.network_1.id}"
  cidr = "192.168.1.0/24"
  ip_version = 4
  enable_dhcp = true
  no_gateway = true
}

resource "openstack_networking_network_v2" "network_2" {
  name = "network_2"
}

resource "openstack_networking_subnet_v2" "subnet_2" {
  name = "subnet_2"
  network_id = "${openstack_networking_network_v2.network_2.id}"
  cidr = "192.168.2.0/24"
  ip_version = 4
  enable_dhcp = true
  no_gateway = true
}

resource "openstack_networking_port_v2" "port_1" {
  name = "port_1"
  network_id = "${openstack_networking_network_v2.network_1.id}"
  admin_state_up = "true"

  fixed_ip {
    subnet_id = "${openstack_networking_subnet_v2.subnet_1.id}"
    ip_address = "192.168.1.103"
  }
}

resource "openstack_networking_port_v2" "port_2" {
  name = "port_2"
  network_id = "${openstack_networking_network_v2.network_2.id}"
  admin_state_up = "true"

  fixed_ip {
    subnet_id = "${openstack_networking_subnet_v2.subnet_2.id}"
    ip_address = "192.168.2.103"
  }
}

resource "openstack_networking_port_v2" "port_3" {
  name = "port_3"
  network_id = "${openstack_networking_network_v2.network_1.id}"
  admin_state_up = "true"

  fixed_ip {
    subnet_id = "${openstack_networking_subnet_v2.subnet_1.id}"
    ip_address = "192.168.1.104"
  }
}

resource "openstack_networking_port_v2" "port_4" {
  name = "port_4"
  network_id = "${openstack_networking_network_v2.network_2.id}"
  admin_state_up = "true"

  fixed_ip {
    subnet_id = "${openstack_networking_subnet_v2.subnet_2.id}"
    ip_address = "192.168.2.104"
  }
}

resource "openstack_compute_instance_v2" "instance_1" {
  depends_on = [
    "openstack_networking_subnet_v2.subnet_1",
    "openstack_networking_subnet_v2.subnet_2",
    "openstack_networking_port_v2.port_1",
    "openstack_networking_port_v2.port_2",
  ]

  name = "instance_1"
  security_groups = ["default"]

  network {
    uuid = "%s"
  }

  network {
    uuid = "${openstack_networking_network_v2.network_1.id}"
    fixed_ip_v4 = "192.168.1.100"
  }

  network {
    uuid = "${openstack_networking_network_v2.network_2.id}"
    fixed_ip_v4 = "192.168.2.100"
  }

  network {
    uuid = "${openstack_networking_network_v2.network_1.id}"
    fixed_ip_v4 = "192.168.1.101"
  }

  network {
    uuid = "${openstack_networking_network_v2.network_2.id}"
    fixed_ip_v4 = "192.168.2.101"
  }

  network {
    port = "${openstack_networking_port_v2.port_1.id}"
  }

  network {
    port = "${openstack_networking_port_v2.port_2.id}"
  }

  network {
    port = "${openstack_networking_port_v2.port_3.id}"
  }

  network {
    port = "${openstack_networking_port_v2.port_4.id}"
  }
}
`, OS_NETWORK_ID)

var testAccComputeV2Instance_stateActive = fmt.Sprintf(`
resource "openstack_compute_instance_v2" "instance_1" {
  name = "instance_1"
  security_groups = ["default"]
  power_state = "active"
  network {
    uuid = "%s"
  }
}
`, OS_NETWORK_ID)

var testAccComputeV2Instance_stateShutoff = fmt.Sprintf(`
resource "openstack_compute_instance_v2" "instance_1" {
  name = "instance_1"
  security_groups = ["default"]
  power_state = "shutoff"
  network {
    uuid = "%s"
  }
}
`, OS_NETWORK_ID)

var testAccComputeV2Instance_tags_create = fmt.Sprintf(`
resource "openstack_compute_instance_v2" "instance_1" {
  name = "instance_1"
  security_groups = ["default"]
  network {
    uuid = "%s"
  }
  tags = ["tag1", "tag2", "tag3"]
}
`, OS_NETWORK_ID)

var testAccComputeV2Instance_tags_add = fmt.Sprintf(`
resource "openstack_compute_instance_v2" "instance_1" {
  name = "instance_1"
  security_groups = ["default"]
  network {
    uuid = "%s"
  }
  tags = ["tag1", "tag2", "tag3", "tag4"]
}
`, OS_NETWORK_ID)

var testAccComputeV2Instance_tags_delete = fmt.Sprintf(`
resource "openstack_compute_instance_v2" "instance_1" {
  name = "instance_1"
  security_groups = ["default"]
  network {
    uuid = "%s"
  }
  tags = ["tag2", "tag3"]
}
`, OS_NETWORK_ID)

var testAccComputeV2Instance_tags_clear = fmt.Sprintf(`
resource "openstack_compute_instance_v2" "instance_1" {
  name = "instance_1"
  security_groups = ["default"]
  network {
    uuid = "%s"
  }
}
`, OS_NETWORK_ID)
