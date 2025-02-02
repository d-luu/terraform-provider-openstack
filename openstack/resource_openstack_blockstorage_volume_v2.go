package openstack

import (
	"fmt"
	"log"
	"time"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/blockstorage/extensions/schedulerhints"
	"github.com/gophercloud/gophercloud/openstack/blockstorage/v2/volumes"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/volumeattach"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceBlockStorageVolumeV2() *schema.Resource {
	return &schema.Resource{
		Create: resourceBlockStorageVolumeV2Create,
		Read:   resourceBlockStorageVolumeV2Read,
		Update: resourceBlockStorageVolumeV2Update,
		Delete: resourceBlockStorageVolumeV2Delete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"size": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},

			"name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
			},

			"description": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
			},

			"availability_zone": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},

			"metadata": {
				Type:     schema.TypeMap,
				Optional: true,
				ForceNew: false,
				Computed: true,
			},

			"snapshot_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"source_vol_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"image_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"volume_type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},

			"consistency_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"source_replica": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"attachment": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"instance_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"device": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
				Set: blockStorageVolumeV2AttachmentHash,
			},

			"scheduler_hints": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"different_host": {
							Type:     schema.TypeList,
							Optional: true,
							ForceNew: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"same_host": {
							Type:     schema.TypeList,
							Optional: true,
							ForceNew: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"query": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"local_to_instance": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"additional_properties": {
							Type:     schema.TypeMap,
							Optional: true,
							ForceNew: true,
						},
					},
				},
				Set: blockStorageExtensionsSchedulerHintsHash,
			},
		},
	}
}

func resourceBlockStorageVolumeV2Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	blockStorageClient, err := config.BlockStorageV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenStack block storage client: %s", err)
	}

	metadata := d.Get("metadata").(map[string]interface{})
	volumeCreateOpts := &volumes.CreateOpts{
		AvailabilityZone:   d.Get("availability_zone").(string),
		ConsistencyGroupID: d.Get("consistency_group_id").(string),
		Description:        d.Get("description").(string),
		ImageID:            d.Get("image_id").(string),
		Metadata:           expandToMapStringString(metadata),
		Name:               d.Get("name").(string),
		Size:               d.Get("size").(int),
		SnapshotID:         d.Get("snapshot_id").(string),
		SourceReplica:      d.Get("source_replica").(string),
		SourceVolID:        d.Get("source_vol_id").(string),
		VolumeType:         d.Get("volume_type").(string),
	}

	var createOpts schedulerhints.CreateOptsExt
	var schedulerHints schedulerhints.SchedulerHints

	schedulerHintsRaw := d.Get("scheduler_hints").(*schema.Set).List()
	if len(schedulerHintsRaw) > 0 {
		log.Printf("[DEBUG] openstack_blockstorage_volume_v2 scheduler hints: %+v", schedulerHintsRaw[0])
		schedulerHints = resourceBlockStorageSchedulerHints(schedulerHintsRaw[0].(map[string]interface{}))
	}
	createOpts = schedulerhints.CreateOptsExt{
		VolumeCreateOptsBuilder: volumeCreateOpts,
		SchedulerHints:          schedulerHints,
	}

	log.Printf("[DEBUG] openstack_blockstorage_volume_v2 create options: %#v", createOpts)

	v, err := volumes.Create(blockStorageClient, createOpts).Extract()
	if err != nil {
		return fmt.Errorf("Error creating openstack_blockstorage_volume_v2: %s", err)
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"downloading", "creating"},
		Target:     []string{"available"},
		Refresh:    blockStorageVolumeV2StateRefreshFunc(blockStorageClient, v.ID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf(
			"Error waiting for openstack_blockstorage_volume_v2 %s to become ready: %s", v.ID, err)
	}

	d.SetId(v.ID)

	return resourceBlockStorageVolumeV2Read(d, meta)
}

func resourceBlockStorageVolumeV2Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	blockStorageClient, err := config.BlockStorageV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenStack block storage client: %s", err)
	}

	v, err := volumes.Get(blockStorageClient, d.Id()).Extract()
	if err != nil {
		return CheckDeleted(d, err, "Error retrieving openstack_blockstorage_volume_v2")
	}

	log.Printf("[DEBUG] Retrieved openstack_blockstorage_volume_v2 %s: %#v", d.Id(), v)

	d.Set("size", v.Size)
	d.Set("description", v.Description)
	d.Set("availability_zone", v.AvailabilityZone)
	d.Set("name", v.Name)
	d.Set("snapshot_id", v.SnapshotID)
	d.Set("source_vol_id", v.SourceVolID)
	d.Set("volume_type", v.VolumeType)
	d.Set("metadata", v.Metadata)
	d.Set("region", GetRegion(d, config))

	attachments := flattenBlockStorageVolumeV2Attachments(v.Attachments)
	log.Printf("[DEBUG] openstack_blockstorage_volume_v2 %s attachments: %#v", d.Id(), attachments)
	if err := d.Set("attachment", attachments); err != nil {
		log.Printf(
			"[DEBUG] unable to set openstack_blockstorage_volume_v2 %s attachments: %s", d.Id(), err)
	}

	return nil
}

func resourceBlockStorageVolumeV2Update(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	blockStorageClient, err := config.BlockStorageV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenStack block storage client: %s", err)
	}

	name := d.Get("name").(string)
	description := d.Get("description").(string)
	updateOpts := volumes.UpdateOpts{
		Name:        &name,
		Description: &description,
	}

	if d.HasChange("metadata") {
		metadata := d.Get("metadata").(map[string]interface{})
		updateOpts.Metadata = expandToMapStringString(metadata)
	}

	_, err = volumes.Update(blockStorageClient, d.Id(), updateOpts).Extract()
	if err != nil {
		return fmt.Errorf("Error updating openstack_blockstorage_volume_v2 %s: %s", d.Id(), err)
	}

	return resourceBlockStorageVolumeV2Read(d, meta)
}

func resourceBlockStorageVolumeV2Delete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	blockStorageClient, err := config.BlockStorageV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenStack block storage client: %s", err)
	}

	v, err := volumes.Get(blockStorageClient, d.Id()).Extract()
	if err != nil {
		return CheckDeleted(d, err, "Error retrieving openstack_blockstorage_volume_v2")
	}

	// Make sure this volume is detached from all instances before deleting.
	if len(v.Attachments) > 0 {
		computeClient, err := config.ComputeV2Client(GetRegion(d, config))
		if err != nil {
			return fmt.Errorf("Error creating OpenStack compute client: %s", err)
		}

		for _, volumeAttachment := range v.Attachments {
			log.Printf("[DEBUG] openstack_blockstorage_volume_v2 %s attachment: %#v", d.Id(), volumeAttachment)

			serverID := volumeAttachment.ServerID
			attachmentID := volumeAttachment.ID
			if err := volumeattach.Delete(computeClient, serverID, attachmentID).ExtractErr(); err != nil {
				// It's possible the volume was already detached by
				// openstack_compute_volume_attach_v2, so consider
				// a 404 acceptable and continue.
				if _, ok := err.(gophercloud.ErrDefault404); ok {
					continue
				}

				// A 409 is also acceptable because there's another
				// concurrent action happening.
				if _, ok := err.(gophercloud.ErrDefault409); ok {
					continue
				}

				return fmt.Errorf(
					"Error detaching openstack_blockstorage_volume_v2 %s from %s: %s", d.Id(), serverID, err)
			}
		}

		stateConf := &resource.StateChangeConf{
			Pending:    []string{"in-use", "attaching", "detaching"},
			Target:     []string{"available", "deleted"},
			Refresh:    blockStorageVolumeV2StateRefreshFunc(blockStorageClient, d.Id()),
			Timeout:    10 * time.Minute,
			Delay:      10 * time.Second,
			MinTimeout: 3 * time.Second,
		}

		_, err = stateConf.WaitForState()
		if err != nil {
			return fmt.Errorf(
				"Error waiting for openstack_blockstorage_volume_v2 %s to become available: %s", d.Id(), err)
		}
	}

	// It's possible that this volume was used as a boot device and is currently
	// in a "deleting" state from when the instance was terminated.
	// If this is true, just move on. It'll eventually delete.
	if v.Status != "deleting" {
		if err := volumes.Delete(blockStorageClient, d.Id(), nil).ExtractErr(); err != nil {
			return CheckDeleted(d, err, "Error deleting openstack_blockstorage_volume_v2")
		}
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"deleting", "downloading", "available"},
		Target:     []string{"deleted"},
		Refresh:    blockStorageVolumeV2StateRefreshFunc(blockStorageClient, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("Error waiting for openstack_blockstorage_volume_v2 %s to delete: %s", d.Id(), err)
	}

	return nil
}
