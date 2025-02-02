package openstack

import (
	"fmt"
	"testing"

	"github.com/gophercloud/gophercloud/openstack/imageservice/v2/images"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccImagesImageV2_basic(t *testing.T) {
	var image images.Image

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckImagesImageV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccImagesImageV2_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckImagesImageV2Exists("openstack_images_image_v2.image_1", &image),
					resource.TestCheckResourceAttr(
						"openstack_images_image_v2.image_1", "name", "Rancher TerraformAccTest"),
					resource.TestCheckResourceAttr(
						"openstack_images_image_v2.image_1", "container_format", "bare"),
					resource.TestCheckResourceAttr(
						"openstack_images_image_v2.image_1", "disk_format", "qcow2"),
					resource.TestCheckResourceAttr(
						"openstack_images_image_v2.image_1", "schema", "/v2/schemas/image"),
				),
			},
		},
	})
}

func TestAccImagesImageV2_name(t *testing.T) {
	var image images.Image

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckImagesImageV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccImagesImageV2_name_1,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckImagesImageV2Exists("openstack_images_image_v2.image_1", &image),
					resource.TestCheckResourceAttr(
						"openstack_images_image_v2.image_1", "name", "Rancher TerraformAccTest"),
				),
			},
			{
				Config: testAccImagesImageV2_name_2,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckImagesImageV2Exists("openstack_images_image_v2.image_1", &image),
					resource.TestCheckResourceAttr(
						"openstack_images_image_v2.image_1", "name", "TerraformAccTest Rancher"),
				),
			},
		},
	})
}

func TestAccImagesImageV2_tags(t *testing.T) {
	var image images.Image

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckImagesImageV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccImagesImageV2_tags_1,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckImagesImageV2Exists("openstack_images_image_v2.image_1", &image),
					testAccCheckImagesImageV2HasTag("openstack_images_image_v2.image_1", "foo"),
					testAccCheckImagesImageV2HasTag("openstack_images_image_v2.image_1", "bar"),
					testAccCheckImagesImageV2TagCount("openstack_images_image_v2.image_1", 2),
				),
			},
			{
				Config: testAccImagesImageV2_tags_2,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckImagesImageV2Exists("openstack_images_image_v2.image_1", &image),
					testAccCheckImagesImageV2HasTag("openstack_images_image_v2.image_1", "foo"),
					testAccCheckImagesImageV2HasTag("openstack_images_image_v2.image_1", "bar"),
					testAccCheckImagesImageV2HasTag("openstack_images_image_v2.image_1", "baz"),
					testAccCheckImagesImageV2TagCount("openstack_images_image_v2.image_1", 3),
				),
			},
			{
				Config: testAccImagesImageV2_tags_3,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckImagesImageV2Exists("openstack_images_image_v2.image_1", &image),
					testAccCheckImagesImageV2HasTag("openstack_images_image_v2.image_1", "foo"),
					testAccCheckImagesImageV2HasTag("openstack_images_image_v2.image_1", "baz"),
					testAccCheckImagesImageV2TagCount("openstack_images_image_v2.image_1", 2),
				),
			},
		},
	})
}

func TestAccImagesImageV2_visibility(t *testing.T) {
	var image images.Image

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckAdminOnly(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckImagesImageV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccImagesImageV2_visibility_1,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckImagesImageV2Exists("openstack_images_image_v2.image_1", &image),
					resource.TestCheckResourceAttr(
						"openstack_images_image_v2.image_1", "visibility", "private"),
				),
			},
			{
				Config: testAccImagesImageV2_visibility_2,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckImagesImageV2Exists("openstack_images_image_v2.image_1", &image),
					resource.TestCheckResourceAttr(
						"openstack_images_image_v2.image_1", "visibility", "public"),
				),
			},
		},
	})
}

func TestAccImagesImageV2_properties(t *testing.T) {
	var image_1 images.Image
	var image_2 images.Image
	var image_3 images.Image
	var image_4 images.Image
	var image_5 images.Image

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckImagesImageV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccImagesImageV2_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckImagesImageV2Exists("openstack_images_image_v2.image_1", &image_1),
					resource.TestCheckResourceAttrSet(
						"openstack_images_image_v2.image_1", "properties.os_hash_value"),
				),
			},
			{
				Config: testAccImagesImageV2_properties_1,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckImagesImageV2Exists("openstack_images_image_v2.image_1", &image_2),
					resource.TestCheckResourceAttr(
						"openstack_images_image_v2.image_1", "properties.foo", "bar"),
					resource.TestCheckResourceAttr(
						"openstack_images_image_v2.image_1", "properties.bar", "foo"),
					resource.TestCheckResourceAttrSet(
						"openstack_images_image_v2.image_1", "properties.os_hash_value"),
				),
			},
			{
				Config: testAccImagesImageV2_properties_2,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckImagesImageV2Exists("openstack_images_image_v2.image_1", &image_3),
					resource.TestCheckResourceAttr(
						"openstack_images_image_v2.image_1", "properties.foo", "bar"),
					resource.TestCheckResourceAttrSet(
						"openstack_images_image_v2.image_1", "properties.os_hash_value"),
				),
			},
			{
				Config: testAccImagesImageV2_properties_3,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckImagesImageV2Exists("openstack_images_image_v2.image_1", &image_4),
					resource.TestCheckResourceAttr(
						"openstack_images_image_v2.image_1", "properties.foo", "baz"),
					resource.TestCheckResourceAttrSet(
						"openstack_images_image_v2.image_1", "properties.os_hash_value"),
				),
			},
			{
				Config: testAccImagesImageV2_properties_4,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckImagesImageV2Exists("openstack_images_image_v2.image_1", &image_5),
					resource.TestCheckResourceAttr(
						"openstack_images_image_v2.image_1", "properties.foo", "baz"),
					resource.TestCheckResourceAttr(
						"openstack_images_image_v2.image_1", "properties.bar", "foo"),
					resource.TestCheckResourceAttrSet(
						"openstack_images_image_v2.image_1", "properties.os_hash_value"),
				),
			},
		},
	})
}

func TestAccImagesImageV2_webdownload(t *testing.T) {
	var image images.Image

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckGlanceImport(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckImagesImageV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccImagesImageV2_webdownload,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckImagesImageV2Exists("openstack_images_image_v2.image_1", &image),
					resource.TestCheckResourceAttr(
						"openstack_images_image_v2.image_1", "name", "Rancher TerraformAccTest"),
					resource.TestCheckResourceAttr(
						"openstack_images_image_v2.image_1", "container_format", "bare"),
					resource.TestCheckResourceAttr(
						"openstack_images_image_v2.image_1", "disk_format", "qcow2"),
					resource.TestCheckResourceAttr(
						"openstack_images_image_v2.image_1", "schema", "/v2/schemas/image"),
				),
			},
		},
	})
}

func testAccCheckImagesImageV2Destroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	imageClient, err := config.ImageV2Client(OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating OpenStack Image: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "openstack_images_image_v2" {
			continue
		}

		_, err := images.Get(imageClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("Image still exists")
		}
	}

	return nil
}

func testAccCheckImagesImageV2Exists(n string, image *images.Image) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)
		imageClient, err := config.ImageV2Client(OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating OpenStack Image: %s", err)
		}

		found, err := images.Get(imageClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("Image not found")
		}

		*image = *found

		return nil
	}
}

func testAccCheckImagesImageV2HasTag(n, tag string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)
		imageClient, err := config.ImageV2Client(OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating OpenStack Image: %s", err)
		}

		found, err := images.Get(imageClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("Image not found")
		}

		for _, v := range found.Tags {
			if tag == v {
				return nil
			}
		}

		return fmt.Errorf("Tag not found: %s", tag)
	}
}

func testAccCheckImagesImageV2TagCount(n string, expected int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)
		imageClient, err := config.ImageV2Client(OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating OpenStack Image: %s", err)
		}

		found, err := images.Get(imageClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("Image not found")
		}

		if len(found.Tags) != expected {
			return fmt.Errorf("Expecting %d tags, found %d", expected, len(found.Tags))
		}

		return nil
	}
}

const testAccImagesImageV2_basic = `
  resource "openstack_images_image_v2" "image_1" {
      name   = "Rancher TerraformAccTest"
      image_source_url = "https://releases.rancher.com/os/latest/rancheros-openstack.img"
      container_format = "bare"
      disk_format = "qcow2"

      timeouts {
        create = "10m"
      }
  }`

const testAccImagesImageV2_name_1 = `
  resource "openstack_images_image_v2" "image_1" {
      name   = "Rancher TerraformAccTest"
      image_source_url = "https://releases.rancher.com/os/latest/rancheros-openstack.img"
      container_format = "bare"
      disk_format = "qcow2"
  }`

const testAccImagesImageV2_name_2 = `
  resource "openstack_images_image_v2" "image_1" {
      name   = "TerraformAccTest Rancher"
      image_source_url = "https://releases.rancher.com/os/latest/rancheros-openstack.img"
      container_format = "bare"
      disk_format = "qcow2"
  }`

const testAccImagesImageV2_tags_1 = `
  resource "openstack_images_image_v2" "image_1" {
      name   = "Rancher TerraformAccTest"
      image_source_url = "https://releases.rancher.com/os/latest/rancheros-openstack.img"
      container_format = "bare"
      disk_format = "qcow2"
      tags = ["foo","bar"]
  }`

const testAccImagesImageV2_tags_2 = `
  resource "openstack_images_image_v2" "image_1" {
      name   = "Rancher TerraformAccTest"
      image_source_url = "https://releases.rancher.com/os/latest/rancheros-openstack.img"
      container_format = "bare"
      disk_format = "qcow2"
      tags = ["foo","bar","baz"]
  }`

const testAccImagesImageV2_tags_3 = `
  resource "openstack_images_image_v2" "image_1" {
      name   = "Rancher TerraformAccTest"
      image_source_url = "https://releases.rancher.com/os/latest/rancheros-openstack.img"
      container_format = "bare"
      disk_format = "qcow2"
      tags = ["foo","baz"]
  }`

const testAccImagesImageV2_visibility_1 = `
  resource "openstack_images_image_v2" "image_1" {
      name   = "Rancher TerraformAccTest"
      image_source_url = "https://releases.rancher.com/os/latest/rancheros-openstack.img"
      container_format = "bare"
      disk_format = "qcow2"
      visibility = "private"
  }`

const testAccImagesImageV2_visibility_2 = `
  resource "openstack_images_image_v2" "image_1" {
      name   = "Rancher TerraformAccTest"
      image_source_url = "https://releases.rancher.com/os/latest/rancheros-openstack.img"
      container_format = "bare"
      disk_format = "qcow2"
      visibility = "public"
  }`

const testAccImagesImageV2_properties_1 = `
  resource "openstack_images_image_v2" "image_1" {
      name   = "Rancher TerraformAccTest"
      image_source_url = "https://releases.rancher.com/os/latest/rancheros-openstack.img"
      container_format = "bare"
      disk_format = "qcow2"

      properties = {
        foo = "bar"
        bar = "foo"
      }
  }`

const testAccImagesImageV2_properties_2 = `
  resource "openstack_images_image_v2" "image_1" {
      name   = "Rancher TerraformAccTest"
      image_source_url = "https://releases.rancher.com/os/latest/rancheros-openstack.img"
      container_format = "bare"
      disk_format = "qcow2"

      properties = {
        foo = "bar"
      }
  }`

const testAccImagesImageV2_properties_3 = `
  resource "openstack_images_image_v2" "image_1" {
      name   = "Rancher TerraformAccTest"
      image_source_url = "https://releases.rancher.com/os/latest/rancheros-openstack.img"
      container_format = "bare"
      disk_format = "qcow2"

      properties = {
        foo = "baz"
      }
  }`

const testAccImagesImageV2_properties_4 = `
  resource "openstack_images_image_v2" "image_1" {
      name   = "Rancher TerraformAccTest"
      image_source_url = "https://releases.rancher.com/os/latest/rancheros-openstack.img"
      container_format = "bare"
      disk_format = "qcow2"

      properties = {
        foo = "baz"
        bar = "foo"
      }
  }`

const testAccImagesImageV2_webdownload = `
  resource "openstack_images_image_v2" "image_1" {
      name   = "Rancher TerraformAccTest"
      image_source_url = "https://releases.rancher.com/os/latest/rancheros-openstack.img"
      container_format = "bare"
      disk_format = "qcow2"
      web_download = true

      timeouts {
        create = "10m"
      }
  }`
