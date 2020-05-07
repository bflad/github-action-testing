package aws

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccAWSEbsVolumeDataSource_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAwsEbsVolumeDataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAwsEbsVolumeDataSourceID("data.aws_ebs_volume.ebs_volume"),
					resource.TestCheckResourceAttrSet("data.aws_ebs_volume.ebs_volume", "arn"),
					resource.TestCheckResourceAttr("data.aws_ebs_volume.ebs_volume", "size", "40"),
					resource.TestCheckResourceAttr("data.aws_ebs_volume.ebs_volume", "tags.%", "1"),
					resource.TestCheckResourceAttr("data.aws_ebs_volume.ebs_volume", "tags.Name", "External Volume"),
					resource.TestCheckResourceAttr("data.aws_ebs_volume.ebs_volume", "outpost_arn", ""),
				),
			},
		},
	})
}

func TestAccAWSEbsVolumeDataSource_multipleFilters(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAwsEbsVolumeDataSourceConfigWithMultipleFilters,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAwsEbsVolumeDataSourceID("data.aws_ebs_volume.ebs_volume"),
					resource.TestCheckResourceAttr("data.aws_ebs_volume.ebs_volume", "size", "10"),
					resource.TestCheckResourceAttr("data.aws_ebs_volume.ebs_volume", "volume_type", "gp2"),
					resource.TestCheckResourceAttr("data.aws_ebs_volume.ebs_volume", "tags.%", "1"),
					resource.TestCheckResourceAttr("data.aws_ebs_volume.ebs_volume", "tags.Name", "External Volume 1"),
				),
			},
		},
	})
}

func testAccCheckAwsEbsVolumeDataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find Volume data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Volume data source ID not set")
		}
		return nil
	}
}

const testAccCheckAwsEbsVolumeDataSourceConfig = `
data "aws_availability_zones" "available" {
  state = "available"

  filter {
    name   = "opt-in-status"
    values = ["opt-in-not-required"]
  }
}

resource "aws_ebs_volume" "example" {
  availability_zone = "${data.aws_availability_zones.available.names[0]}"
  type = "gp2"
  size = 40
  tags = {
    Name = "External Volume"
  }
}

data "aws_ebs_volume" "ebs_volume" {
  most_recent = true
  filter {
    name = "tag:Name"
    values = ["External Volume"]
  }
  filter {
    name = "volume-type"
    values = ["${aws_ebs_volume.example.type}"]
  }
}
`

const testAccCheckAwsEbsVolumeDataSourceConfigWithMultipleFilters = `
data "aws_availability_zones" "available" {
  state = "available"

  filter {
    name   = "opt-in-status"
    values = ["opt-in-not-required"]
  }
}

resource "aws_ebs_volume" "external1" {
  availability_zone = "${data.aws_availability_zones.available.names[0]}"
  type = "gp2"
  size = 10
  tags = {
    Name = "External Volume 1"
  }
}

data "aws_ebs_volume" "ebs_volume" {
  most_recent = true
  filter {
    name = "tag:Name"
    values = ["External Volume 1"]
  }
  filter {
    name = "size"
    values = ["${aws_ebs_volume.external1.size}"]
  }
  filter {
    name = "volume-type"
    values = ["${aws_ebs_volume.external1.type}"]
  }
}
`