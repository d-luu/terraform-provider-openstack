package openstack

import (
	"fmt"
	"log"

	"github.com/gophercloud/gophercloud/openstack/blockstorage/v3/volumes"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceBlockStorageVolumeV3() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceBlockStorageVolumeV3Read,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"status": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"metadata": {
				Type:     schema.TypeMap,
				Optional: true,
				Computed: true,
			},

			// Computed values
			"bootable": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"volume_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"size": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"source_volume_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"multiattach": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func dataSourceBlockStorageVolumeV3Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	client, err := config.BlockStorageV3Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenStack block storage client: %s", err)
	}

	listOpts := volumes.ListOpts{
		Metadata: expandToMapStringString(d.Get("metadata").(map[string]interface{})),
		Name:     d.Get("name").(string),
		Status:   d.Get("status").(string),
	}

	allPages, err := volumes.List(client, listOpts).AllPages()
	if err != nil {
		return fmt.Errorf("Unable to query openstack_blockstorage_volume_v3: %s", err)
	}

	allVolumes, err := volumes.ExtractVolumes(allPages)
	if err != nil {
		return fmt.Errorf("Unable to retrieve openstack_blockstorage_volume_v3: %s", err)
	}

	if len(allVolumes) > 1 {
		return fmt.Errorf("Your openstack_blockstorage_volume_v3 query returned multiple results.")
	}

	if len(allVolumes) < 1 {
		return fmt.Errorf("Your openstack_blockstorage_volume_v3 query returned no results.")
	}

	return dataSourceBlockStorageVolumeV3Attributes(d, allVolumes[0])
}

func dataSourceBlockStorageVolumeV3Attributes(d *schema.ResourceData, volume volumes.Volume) error {
	d.SetId(volume.ID)
	d.Set("name", volume.Name)
	d.Set("status", volume.Status)
	d.Set("bootable", volume.Bootable)
	d.Set("volume_type", volume.VolumeType)
	d.Set("size", volume.Size)
	d.Set("source_volume_id", volume.SourceVolID)
	d.Set("multiattach", volume.Multiattach)

	if err := d.Set("metadata", volume.Metadata); err != nil {
		log.Printf("[DEBUG] Unable to set metadata for openstack_blockstorage_volume_v3 %s: %s", volume.ID, err)
	}

	return nil
}
