package aws

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccDataSourceAwsDirectoryServiceDirectory_SimpleAD(t *testing.T) {
	alias := acctest.RandomWithPrefix("tf-acc-test")
	resourceName := "aws_directory_service_directory.test-simple-ad"
	dataSourceName := "data.aws_directory_service_directory.test-simple-ad"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAwsDirectoryServiceDirectoryConfig_SimpleAD(alias),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "type", "SimpleAD"),
					resource.TestCheckResourceAttr(dataSourceName, "size", "Small"),
					resource.TestCheckResourceAttr(dataSourceName, "name", "tf-testacc-corp.neverland.com"),
					resource.TestCheckResourceAttr(dataSourceName, "description", "tf-testacc SimpleAD"),
					resource.TestCheckResourceAttr(dataSourceName, "short_name", "corp"),
					resource.TestCheckResourceAttr(dataSourceName, "alias", alias),
					resource.TestCheckResourceAttr(dataSourceName, "enable_sso", "false"),
					resource.TestCheckResourceAttr(dataSourceName, "vpc_settings.#", "1"),
					resource.TestCheckResourceAttrPair(dataSourceName, "vpc_settings.0.vpc_id", resourceName, "vpc_settings.0.vpc_id"),
					resource.TestCheckResourceAttrPair(dataSourceName, "vpc_settings.0.subnet_ids", resourceName, "vpc_settings.0.subnet_ids"),
					resource.TestCheckResourceAttr(dataSourceName, "access_url", fmt.Sprintf("%s.awsapps.com", alias)),
					resource.TestCheckResourceAttrPair(dataSourceName, "dns_ip_addresses", resourceName, "dns_ip_addresses"),
					resource.TestCheckResourceAttrPair(dataSourceName, "security_group_id", resourceName, "security_group_id"),
				),
			},
		},
	})
}

func TestAccDataSourceAwsDirectoryServiceDirectory_MicrosoftAD(t *testing.T) {
	alias := acctest.RandomWithPrefix("tf-acc-test")
	resourceName := "aws_directory_service_directory.test-microsoft-ad"
	dataSourceName := "data.aws_directory_service_directory.test-microsoft-ad"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAwsDirectoryServiceDirectoryConfig_MicrosoftAD(alias),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "type", "MicrosoftAD"),
					resource.TestCheckResourceAttr(dataSourceName, "edition", "Standard"),
					resource.TestCheckResourceAttr(dataSourceName, "name", "tf-testacc-corp.neverland.com"),
					resource.TestCheckResourceAttr(dataSourceName, "description", "tf-testacc MicrosoftAD"),
					resource.TestCheckResourceAttr(dataSourceName, "short_name", "corp"),
					resource.TestCheckResourceAttr(dataSourceName, "alias", alias),
					resource.TestCheckResourceAttr(dataSourceName, "enable_sso", "false"),
					resource.TestCheckResourceAttr(dataSourceName, "vpc_settings.#", "1"),
					resource.TestCheckResourceAttrPair(dataSourceName, "vpc_settings.0.vpc_id", resourceName, "vpc_settings.0.vpc_id"),
					resource.TestCheckResourceAttrPair(dataSourceName, "vpc_settings.0.subnet_ids", resourceName, "vpc_settings.0.subnet_ids"),
					resource.TestCheckResourceAttr(dataSourceName, "access_url", fmt.Sprintf("%s.awsapps.com", alias)),
					resource.TestCheckResourceAttrPair(dataSourceName, "dns_ip_addresses", resourceName, "dns_ip_addresses"),
					resource.TestCheckResourceAttrPair(dataSourceName, "security_group_id", resourceName, "security_group_id"),
				),
			},
		},
	})
}

func testAccDataSourceAwsDirectoryServiceDirectoryConfig_Prerequisites(adType string) string {
	return fmt.Sprintf(`
data "aws_availability_zones" "available" {
  state = "available"

  filter {
    name   = "opt-in-status"
    values = ["opt-in-not-required"]
  }
}

resource "aws_vpc" "main" {
  cidr_block = "10.0.0.0/16"

  tags = {
    Name = "tf-testacc-%s"
  }
}

resource "aws_subnet" "primary" {
  vpc_id = "${aws_vpc.main.id}"
  availability_zone = "${data.aws_availability_zones.available.names[0]}"
  cidr_block = "10.0.1.0/24"

  tags = {
    Name = "tf-testacc-%s-primary"
  }
}
resource "aws_subnet" "secondary" {
  vpc_id = "${aws_vpc.main.id}"
  availability_zone = "${data.aws_availability_zones.available.names[1]}"
  cidr_block = "10.0.2.0/24"

  tags = {
    Name = "tf-testacc-%s-secondary"
  }
}
`, adType, adType, adType)
}

func testAccDataSourceAwsDirectoryServiceDirectoryConfig_SimpleAD(alias string) string {
	return testAccDataSourceAwsDirectoryServiceDirectoryConfig_Prerequisites("simple-ad") + fmt.Sprintf(`
resource "aws_directory_service_directory" "test-simple-ad" {
  type = "SimpleAD"
  size = "Small"
  name = "tf-testacc-corp.neverland.com"
  description = "tf-testacc SimpleAD"
  short_name = "corp"
  password = "#S1ncerely"
  
  alias = %q
  enable_sso = false

  vpc_settings {
    vpc_id = "${aws_vpc.main.id}"
    subnet_ids = ["${aws_subnet.primary.id}", "${aws_subnet.secondary.id}"]
  }
}

data "aws_directory_service_directory" "test-simple-ad" {
  directory_id = "${aws_directory_service_directory.test-simple-ad.id}"
}
`, alias)
}

func testAccDataSourceAwsDirectoryServiceDirectoryConfig_MicrosoftAD(alias string) string {
	return testAccDataSourceAwsDirectoryServiceDirectoryConfig_Prerequisites("microsoft-ad") + fmt.Sprintf(`
resource "aws_directory_service_directory" "test-microsoft-ad" {
  type = "MicrosoftAD"
  edition = "Standard"
  name = "tf-testacc-corp.neverland.com"
  description = "tf-testacc MicrosoftAD"
  short_name = "corp"
  password = "#S1ncerely"
  
  alias = %q
  enable_sso = false

  vpc_settings {
    vpc_id = "${aws_vpc.main.id}"
    subnet_ids = ["${aws_subnet.primary.id}", "${aws_subnet.secondary.id}"]
  }
}

data "aws_directory_service_directory" "test-microsoft-ad" {
  directory_id = "${aws_directory_service_directory.test-microsoft-ad.id}"
}
`, alias)
}
