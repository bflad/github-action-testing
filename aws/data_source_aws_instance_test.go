package aws

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAWSInstanceDataSource_basic(t *testing.T) {
	resourceName := "aws_instance.test"
	datasourceName := "data.aws_instance.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccInstanceDataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(datasourceName, "ami", resourceName, "ami"),
					resource.TestCheckResourceAttrPair(datasourceName, "tags.%", resourceName, "tags.%"),
					resource.TestCheckResourceAttrPair(datasourceName, "instance_type", resourceName, "instance_type"),
					resource.TestCheckResourceAttrPair(datasourceName, "arn", resourceName, "arn"),
					resource.TestCheckNoResourceAttr(datasourceName, "user_data_base64"),
					resource.TestCheckResourceAttr(datasourceName, "outpost_arn", ""),
				),
			},
		},
	})
}

func TestAccAWSInstanceDataSource_tags(t *testing.T) {
	rInt := acctest.RandInt()
	resourceName := "aws_instance.test"
	datasourceName := "data.aws_instance.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccInstanceDataSourceConfig_Tags(rInt),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(datasourceName, "ami", resourceName, "ami"),
					resource.TestCheckResourceAttrPair(datasourceName, "tags.%", resourceName, "tags.%"),
					resource.TestCheckResourceAttrPair(datasourceName, "instance_type", resourceName, "instance_type"),
				),
			},
		},
	})
}

func TestAccAWSInstanceDataSource_AzUserData(t *testing.T) {
	resourceName := "aws_instance.test"
	datasourceName := "data.aws_instance.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccInstanceDataSourceConfig_AzUserData,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(datasourceName, "ami", resourceName, "ami"),
					resource.TestCheckResourceAttrPair(datasourceName, "tags.%", resourceName, "tags.%"),
					resource.TestCheckResourceAttrPair(datasourceName, "instance_type", resourceName, "instance_type"),
					resource.TestCheckResourceAttrPair(datasourceName, "availability_zone", resourceName, "availability_zone"),
					resource.TestCheckResourceAttrPair(datasourceName, "user_data", resourceName, "user_data"),
				),
			},
		},
	})
}

func TestAccAWSInstanceDataSource_gp2IopsDevice(t *testing.T) {
	resourceName := "aws_instance.test"
	datasourceName := "data.aws_instance.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccInstanceDataSourceConfig_gp2IopsDevice,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(datasourceName, "ami", resourceName, "ami"),
					resource.TestCheckResourceAttrPair(datasourceName, "instance_type", resourceName, "instance_type"),
					resource.TestCheckResourceAttrPair(datasourceName, "root_block_device.#", resourceName, "root_block_device.#"),
					resource.TestCheckResourceAttrPair(datasourceName, "root_block_device.0.volume_size", resourceName, "root_block_device.0.volume_size"),
					resource.TestCheckResourceAttrPair(datasourceName, "root_block_device.0.volume_type", resourceName, "root_block_device.0.volume_type"),
					resource.TestCheckResourceAttrPair(datasourceName, "root_block_device.0.device_name", resourceName, "root_block_device.0.device_name"),
					resource.TestCheckResourceAttrPair(datasourceName, "root_block_device.0.iops", resourceName, "root_block_device.0.iops"),
				),
			},
		},
	})
}

func TestAccAWSInstanceDataSource_blockDevices(t *testing.T) {
	resourceName := "aws_instance.test"
	datasourceName := "data.aws_instance.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccInstanceDataSourceConfig_blockDevices,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(datasourceName, "ami", resourceName, "ami"),
					resource.TestCheckResourceAttrPair(datasourceName, "instance_type", resourceName, "instance_type"),
					resource.TestCheckResourceAttrPair(datasourceName, "root_block_device.#", resourceName, "root_block_device.#"),
					resource.TestCheckResourceAttrPair(datasourceName, "root_block_device.0.volume_size", resourceName, "root_block_device.0.volume_size"),
					resource.TestCheckResourceAttrPair(datasourceName, "root_block_device.0.volume_type", resourceName, "root_block_device.0.volume_type"),
					resource.TestCheckResourceAttrPair(datasourceName, "root_block_device.0.device_name", resourceName, "root_block_device.0.device_name"),
					resource.TestCheckResourceAttrPair(datasourceName, "ebs_block_device.#", resourceName, "ebs_block_device.#"),
					//resource.TestCheckResourceAttrPair(datasourceName, "ephemeral_block_device.#", resourceName, "ephemeral_block_device.#"),
					// ephemeral block devices don't get saved properly due to API limitations, so this can't actually be tested right now
				),
			},
		},
	})
}

// Test to verify that ebs_block_device kms_key_id does not elicit a panic
func TestAccAWSInstanceDataSource_EbsBlockDevice_KmsKeyId(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccInstanceDataSourceConfig_EbsBlockDevice_KmsKeyId,
			},
		},
	})
}

// Test to verify that root_block_device kms_key_id does not elicit a panic
func TestAccAWSInstanceDataSource_RootBlockDevice_KmsKeyId(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccInstanceDataSourceConfig_RootBlockDevice_KmsKeyId,
			},
		},
	})
}

func TestAccAWSInstanceDataSource_rootInstanceStore(t *testing.T) {
	resourceName := "aws_instance.test"
	datasourceName := "data.aws_instance.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccInstanceDataSourceConfig_rootInstanceStore,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(datasourceName, "ami", resourceName, "ami"),
					resource.TestCheckResourceAttrPair(datasourceName, "instance_type", resourceName, "instance_type"),
					resource.TestCheckResourceAttrPair(datasourceName, "ebs_block_device.#", resourceName, "ebs_block_device.#"),
					resource.TestCheckResourceAttrPair(datasourceName, "ebs_optimized", resourceName, "ebs_optimized"),
					resource.TestCheckResourceAttrPair(datasourceName, "root_block_device.#", resourceName, "root_block_device.#"),
				),
			},
		},
	})
}

func TestAccAWSInstanceDataSource_privateIP(t *testing.T) {
	resourceName := "aws_instance.test"
	datasourceName := "data.aws_instance.test"
	rName := fmt.Sprintf("tf-testacc-instance-%s", acctest.RandStringFromCharSet(12, acctest.CharSetAlphaNum))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccInstanceDataSourceConfig_privateIP(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(datasourceName, "ami", resourceName, "ami"),
					resource.TestCheckResourceAttrPair(datasourceName, "instance_type", resourceName, "instance_type"),
					resource.TestCheckResourceAttrPair(datasourceName, "private_ip", resourceName, "private_ip"),
				),
			},
		},
	})
}

func TestAccAWSInstanceDataSource_keyPair(t *testing.T) {
	resourceName := "aws_instance.test"
	datasourceName := "data.aws_instance.test"
	rName := fmt.Sprintf("tf-test-key-%d", acctest.RandInt())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccInstanceDataSourceConfig_keyPair(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(datasourceName, "ami", resourceName, "ami"),
					resource.TestCheckResourceAttrPair(datasourceName, "tags.%", resourceName, "tags.%"),
					resource.TestCheckResourceAttrPair(datasourceName, "instance_type", resourceName, "instance_type"),
					resource.TestCheckResourceAttrPair(datasourceName, "key_name", resourceName, "key_name"),
				),
			},
		},
	})
}

func TestAccAWSInstanceDataSource_VPC(t *testing.T) {
	resourceName := "aws_instance.test"
	datasourceName := "data.aws_instance.test"
	rName := fmt.Sprintf("tf-testacc-instance-%s", acctest.RandStringFromCharSet(12, acctest.CharSetAlphaNum))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccInstanceDataSourceConfig_VPC(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(datasourceName, "ami", resourceName, "ami"),
					resource.TestCheckResourceAttrPair(datasourceName, "instance_type", resourceName, "instance_type"),
					resource.TestCheckResourceAttrPair(datasourceName, "user_data", resourceName, "user_data"),
					resource.TestCheckResourceAttrPair(datasourceName, "associate_public_ip_address", resourceName, "associate_public_ip_address"),
					resource.TestCheckResourceAttrPair(datasourceName, "tenancy", resourceName, "tenancy"),
				),
			},
		},
	})
}

func TestAccAWSInstanceDataSource_PlacementGroup(t *testing.T) {
	resourceName := "aws_instance.test"
	datasourceName := "data.aws_instance.test"
	rName := fmt.Sprintf("tf-testacc-instance-%s", acctest.RandStringFromCharSet(12, acctest.CharSetAlphaNum))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccInstanceDataSourceConfig_PlacementGroup(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(datasourceName, "placement_group", resourceName, "placement_group"),
				),
			},
		},
	})
}

func TestAccAWSInstanceDataSource_SecurityGroups(t *testing.T) {
	rInt := acctest.RandInt()
	resourceName := "aws_instance.test"
	datasourceName := "data.aws_instance.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccInstanceDataSourceConfig_SecurityGroups(rInt),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(datasourceName, "ami", resourceName, "ami"),
					resource.TestCheckResourceAttrPair(datasourceName, "instance_type", resourceName, "instance_type"),
					resource.TestCheckResourceAttrPair(datasourceName, "user_data", resourceName, "user_data"),
					resource.TestCheckResourceAttrPair(datasourceName, "vpc_security_group_ids.#", resourceName, "vpc_security_group_ids.#"),
					resource.TestCheckResourceAttrPair(datasourceName, "security_groups.#", resourceName, "security_groups.#"),
				),
			},
		},
	})
}

func TestAccAWSInstanceDataSource_VPCSecurityGroups(t *testing.T) {
	resourceName := "aws_instance.test"
	datasourceName := "data.aws_instance.test"
	rName := fmt.Sprintf("tf-testacc-instance-%s", acctest.RandStringFromCharSet(12, acctest.CharSetAlphaNum))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccInstanceDataSourceConfig_VPCSecurityGroups(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(datasourceName, "ami", resourceName, "ami"),
					resource.TestCheckResourceAttrPair(datasourceName, "instance_type", resourceName, "instance_type"),
					resource.TestCheckResourceAttrPair(datasourceName, "vpc_security_group_ids.#", resourceName, "vpc_security_group_ids.#"),
					resource.TestCheckResourceAttrPair(datasourceName, "security_groups.#", resourceName, "security_groups.#"),
				),
			},
		},
	})
}

func TestAccAWSInstanceDataSource_getPasswordData_trueToFalse(t *testing.T) {
	datasourceName := "data.aws_instance.test"
	rName := fmt.Sprintf("tf-testacc-instance-%s", acctest.RandStringFromCharSet(12, acctest.CharSetAlphaNum))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccInstanceDataSourceConfig_getPasswordData(rName, true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceName, "get_password_data", "true"),
					resource.TestCheckResourceAttrSet(datasourceName, "password_data"),
				),
			},
			{
				Config: testAccInstanceDataSourceConfig_getPasswordData(rName, false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceName, "get_password_data", "false"),
					resource.TestCheckNoResourceAttr(datasourceName, "password_data"),
				),
			},
		},
	})
}

func TestAccAWSInstanceDataSource_getPasswordData_falseToTrue(t *testing.T) {
	datasourceName := "data.aws_instance.test"
	rName := fmt.Sprintf("tf-testacc-instance-%s", acctest.RandStringFromCharSet(12, acctest.CharSetAlphaNum))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccInstanceDataSourceConfig_getPasswordData(rName, false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceName, "get_password_data", "false"),
					resource.TestCheckNoResourceAttr(datasourceName, "password_data"),
				),
			},
			{
				Config: testAccInstanceDataSourceConfig_getPasswordData(rName, true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceName, "get_password_data", "true"),
					resource.TestCheckResourceAttrSet(datasourceName, "password_data"),
				),
			},
		},
	})
}

func TestAccAWSInstanceDataSource_GetUserData(t *testing.T) {
	datasourceName := "data.aws_instance.test"
	rName := fmt.Sprintf("tf-testacc-instance-%s", acctest.RandStringFromCharSet(12, acctest.CharSetAlphaNum))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccInstanceDataSourceConfigGetUserData(rName, true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceName, "get_user_data", "true"),
					resource.TestCheckResourceAttr(datasourceName, "user_data_base64", "IyEvYmluL2Jhc2gKCmVjaG8gImhlbGxvIHdvcmxkIgo="),
				),
			},
			{
				Config: testAccInstanceDataSourceConfigGetUserData(rName, false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceName, "get_user_data", "false"),
					resource.TestCheckNoResourceAttr(datasourceName, "user_data_base64"),
				),
			},
			{
				Config: testAccInstanceDataSourceConfigGetUserData(rName, true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceName, "get_user_data", "true"),
					resource.TestCheckResourceAttr(datasourceName, "user_data_base64", "IyEvYmluL2Jhc2gKCmVjaG8gImhlbGxvIHdvcmxkIgo="),
				),
			},
		},
	})
}

func TestAccAWSInstanceDataSource_GetUserData_NoUserData(t *testing.T) {
	resourceName := "aws_instance.test"
	datasourceName := "data.aws_instance.test"
	rName := fmt.Sprintf("tf-testacc-instance-%s", acctest.RandStringFromCharSet(12, acctest.CharSetAlphaNum))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccInstanceDataSourceConfigGetUserDataNoUserData(rName, true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceName, "get_user_data", "true"),
					resource.TestCheckNoResourceAttr(datasourceName, "user_data_base64"),
					resource.TestCheckResourceAttrPair(datasourceName, "user_data_base64", resourceName, "user_data_base64"),
				),
			},
			{
				Config: testAccInstanceDataSourceConfigGetUserDataNoUserData(rName, false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceName, "get_user_data", "false"),
					resource.TestCheckNoResourceAttr(datasourceName, "user_data_base64"),
					resource.TestCheckResourceAttrPair(datasourceName, "user_data_base64", resourceName, "user_data_base64"),
				),
			},
			{
				Config: testAccInstanceDataSourceConfigGetUserDataNoUserData(rName, true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceName, "get_user_data", "true"),
					resource.TestCheckNoResourceAttr(datasourceName, "user_data_base64"),
					resource.TestCheckResourceAttrPair(datasourceName, "user_data_base64", resourceName, "user_data_base64"),
				),
			},
		},
	})
}

func TestAccAWSInstanceDataSource_creditSpecification(t *testing.T) {
	resourceName := "aws_instance.test"
	datasourceName := "data.aws_instance.test"
	rName := fmt.Sprintf("tf-testacc-instance-%s", acctest.RandStringFromCharSet(12, acctest.CharSetAlphaNum))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{

				Config: testAccInstanceDataSourceConfig_creditSpecification(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(datasourceName, "instance_type", resourceName, "instance_type"),
					resource.TestCheckResourceAttrPair(datasourceName, "credit_specification.#", resourceName, "credit_specification.#"),
					resource.TestCheckResourceAttrPair(datasourceName, "credit_specification.0.cpu_credits", resourceName, "credit_specification.0.cpu_credits"),
				),
			},
		},
	})
}

func TestAccAWSInstanceDataSource_metadataOptions(t *testing.T) {
	resourceName := "aws_instance.test"
	datasourceName := "data.aws_instance.test"
	rName := acctest.RandomWithPrefix("tf-acc-test")

	resource.ParallelTest(t, resource.TestCase{
		// No subnet_id specified requires default VPC or EC2-Classic.
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckHasDefaultVpcOrEc2Classic(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccInstanceDataSourceConfig_metadataOptions(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(datasourceName, "metadata_options.#", resourceName, "metadata_options.#"),
					resource.TestCheckResourceAttrPair(datasourceName, "metadata_options.0.http_endpoint", resourceName, "metadata_options.0.http_endpoint"),
					resource.TestCheckResourceAttrPair(datasourceName, "metadata_options.0.http_tokens", resourceName, "metadata_options.0.http_tokens"),
					resource.TestCheckResourceAttrPair(datasourceName, "metadata_options.0.http_put_response_hop_limit", resourceName, "metadata_options.0.http_put_response_hop_limit"),
				),
			},
		},
	})
}

// Lookup based on InstanceID
const testAccInstanceDataSourceConfig = `
resource "aws_instance" "test" {
  # us-west-2
  ami = "ami-4fccb37f"
  instance_type = "m1.small"
  tags = {
    Name = "HelloWorld"
  }
}

data "aws_instance" "test" {
  filter {
    name = "instance-id"
    values = ["${aws_instance.test.id}"]
  }
}
`

// Use the tags attribute to filter
func testAccInstanceDataSourceConfig_Tags(rInt int) string {
	return fmt.Sprintf(`
resource "aws_instance" "test" {
  # us-west-2
  ami           = "ami-4fccb37f"
  instance_type = "m1.small"

  tags = {
    Name     = "HelloWorld"
    TestSeed = "%[1]d"
  }
}

data "aws_instance" "test" {
  instance_tags = {
    Name     = "${aws_instance.test.tags["Name"]}"
    TestSeed = "%[1]d"
  }
}
`, rInt)
}

// filter on tag, populate more attributes
const testAccInstanceDataSourceConfig_AzUserData = `
resource "aws_instance" "test" {
  # us-west-2
  ami = "ami-4fccb37f"
  availability_zone = "us-west-2a"

  instance_type = "m1.small"
  user_data = "test:-with-character's"
  tags = {
    TFAccTest = "YesThisIsATest"
  }
}

data "aws_instance" "test" {
  instance_id = "${aws_instance.test.id}"
}
`

// GP2IopsDevice
const testAccInstanceDataSourceConfig_gp2IopsDevice = `
resource "aws_instance" "test" {
  # us-west-2
  ami = "ami-55a7ea65"
  instance_type = "m3.medium"
  root_block_device {
    volume_type = "gp2"
    volume_size = 11
  }
}

data "aws_instance" "test" {
  instance_id = "${aws_instance.test.id}"
}
`

// Block Device
const testAccInstanceDataSourceConfig_blockDevices = `
resource "aws_instance" "test" {
  # us-west-2
  ami = "ami-55a7ea65"
  instance_type = "m3.medium"

  root_block_device {
    volume_type = "gp2"
    volume_size = 11
  }
  ebs_block_device {
    device_name = "/dev/sdb"
    volume_size = 9
  }
  ebs_block_device {
    device_name = "/dev/sdc"
    volume_size = 10
    volume_type = "io1"
    iops = 100
  }

  # Encrypted ebs block device
  ebs_block_device {
    device_name = "/dev/sdd"
    volume_size = 12
    encrypted = true
  }

  ephemeral_block_device {
    device_name = "/dev/sde"
    virtual_name = "ephemeral0"
  }
}

data "aws_instance" "test" {
  instance_id = "${aws_instance.test.id}"
}
`

const testAccInstanceDataSourceConfig_EbsBlockDevice_KmsKeyId = `
resource "aws_kms_key" "test" {
  deletion_window_in_days = 7
}

resource "aws_instance" "test" {
  # us-west-2
  ami = "ami-55a7ea65"
  instance_type = "m3.medium"

  root_block_device {
    volume_type = "gp2"
    volume_size = 11
  }
  ebs_block_device {
    device_name = "/dev/sdb"
    encrypted   = true
    kms_key_id = "${aws_kms_key.test.arn}"
    volume_size = 9
  }
}

data "aws_instance" "test" {
  instance_id = "${aws_instance.test.id}"
}
`

const testAccInstanceDataSourceConfig_RootBlockDevice_KmsKeyId = `
resource "aws_kms_key" "test" {
  deletion_window_in_days = 7
}

resource "aws_instance" "test" {
  # us-west-2
  ami = "ami-55a7ea65"
  instance_type = "m3.medium"

  root_block_device {
    encrypted   = true
    kms_key_id = "${aws_kms_key.test.arn}"
    volume_type = "gp2"
    volume_size = 11
  }
}

data "aws_instance" "test" {
  instance_id = "${aws_instance.test.id}"
}
`

const testAccInstanceDataSourceConfig_rootInstanceStore = `
resource "aws_instance" "test" {
  ami = "ami-44c36524"
  instance_type = "m3.medium"
}
data "aws_instance" "test" {
  instance_id = "${aws_instance.test.id}"
}
`

func testAccInstanceDataSourceConfig_privateIP(rName string) string {
	return testAccLatestAmazonLinuxHvmEbsAmiConfig() + testAccAwsInstanceVpcConfig(rName, false) + fmt.Sprintf(`
resource "aws_instance" "test" {
  ami           = "${data.aws_ami.amzn-ami-minimal-hvm-ebs.id}"
  instance_type = "t2.micro"
  subnet_id     = "${aws_subnet.test.id}"
  private_ip    = "10.1.1.42"
}

data "aws_instance" "test" {
  instance_id = "${aws_instance.test.id}"
}
`)
}

func testAccInstanceDataSourceConfig_keyPair(rName string) string {
	return testAccLatestAmazonLinuxHvmEbsAmiConfig() + fmt.Sprintf(`
resource "aws_key_pair" "test" {
  key_name   = %[1]q
  public_key = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQD3F6tyPEFEzV0LX3X8BsXdMsQz1x2cEikKDEY0aIj41qgxMCP/iteneqXSIFZBp5vizPvaoIR3Um9xK7PGoW8giupGn+EPuxIA4cDM4vzOqOkiMPhz5XK0whEjkVzTo4+S0puvDZuwIsdiW9mxhJc7tgBNL0cYlWSYVkz4G/fslNfRPW5mYAM49f4fhtxPb5ok4Q2Lg9dPKVHO/Bgeu5woMc7RY0p1ej6D4CKFE6lymSDJpW0YHX/wqE9+cfEauh7xZcG0q9t2ta6F6fmX0agvpFyZo8aFbXeUBr7osSCJNgvavWbM/06niWrOvYX2xwWdhXmXSrbX8ZbabVohBK41 phodgson@thoughtworks.com"
}

resource "aws_instance" "test" {
  ami           = "${data.aws_ami.amzn-ami-minimal-hvm-ebs.id}"
  instance_type = "t1.micro"
  key_name      = "${aws_key_pair.test.key_name}"

  tags = {
    Name = %[1]q
  }
}

data "aws_instance" "test" {
  filter {
    name   = "tag:Name"
    values = [%[1]q]
  }

  filter {
    name   = "key-name"
    values = ["${aws_instance.test.key_name}"]
  }
}
`, rName)
}

func testAccInstanceDataSourceConfig_VPC(rName string) string {
	return testAccLatestAmazonLinuxHvmEbsAmiConfig() + testAccAwsInstanceVpcConfig(rName, false) + fmt.Sprintf(`
resource "aws_instance" "test" {
  ami                         = "${data.aws_ami.amzn-ami-minimal-hvm-ebs.id}"
  instance_type               = "m1.small"
  subnet_id                   = "${aws_subnet.test.id}"
  associate_public_ip_address = true
  tenancy                     = "dedicated"
  # pre-encoded base64 data
  user_data                   = "3dc39dda39be1205215e776bad998da361a5955d"
}

data "aws_instance" "test" {
  instance_id = "${aws_instance.test.id}"
}
`)
}

func testAccInstanceDataSourceConfig_PlacementGroup(rName string) string {
	return testAccLatestAmazonLinuxHvmEbsAmiConfig() + testAccAwsInstanceVpcConfig(rName, false) + fmt.Sprintf(`
resource "aws_placement_group" "test" {
  name     = %[1]q
  strategy = "cluster"
}

# Limitations: https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/placement-groups.html#concepts-placement-groups
resource "aws_instance" "test" {
  ami                         = "${data.aws_ami.amzn-ami-minimal-hvm-ebs.id}"
  instance_type               = "c3.large"
  subnet_id                   = "${aws_subnet.test.id}"
  associate_public_ip_address = true
  placement_group             = "${aws_placement_group.test.name}"

  # pre-encoded base64 data
  user_data = "3dc39dda39be1205215e776bad998da361a5955d"
}

data "aws_instance" "test" {
  instance_id = "${aws_instance.test.id}"
}
`, rName)
}

func testAccInstanceDataSourceConfig_SecurityGroups(rInt int) string {
	return fmt.Sprintf(`
provider "aws" {
  region = "us-east-1"
}

resource "aws_security_group" "tf_test_foo" {
  name        = "tf_test_foo-%d"
  description = "foo"

  ingress {
    protocol  = "icmp"
    from_port = -1
    to_port   = -1
    self      = true
  }
}

resource "aws_instance" "test" {
  ami             = "ami-408c7f28"
  instance_type   = "m1.small"
  security_groups = ["${aws_security_group.tf_test_foo.name}"]
  user_data       = "foo:-with-character's"
}

data "aws_instance" "test" {
  instance_id = "${aws_instance.test.id}"
}
`, rInt)
}

func testAccInstanceDataSourceConfig_VPCSecurityGroups(rName string) string {
	return testAccLatestAmazonLinuxHvmEbsAmiConfig() +
		testAccAwsInstanceVpcConfig(rName, false) +
		testAccAwsInstanceVpcSecurityGroupConfig(rName) +
		fmt.Sprintf(`
resource "aws_instance" "test" {
  ami                    = "${data.aws_ami.amzn-ami-minimal-hvm-ebs.id}"
  instance_type          = "t1.micro"
  vpc_security_group_ids = ["${aws_security_group.test.id}"]
  subnet_id              = "${aws_subnet.test.id}"
  depends_on             = ["aws_internet_gateway.test"]
}

data "aws_instance" "test" {
  instance_id = "${aws_instance.test.id}"
}
`)
}

func testAccInstanceDataSourceConfig_getPasswordData(rName string, val bool) string {
	return testAccLatestWindowsServer2016CoreAmiConfig() + fmt.Sprintf(`
resource "aws_key_pair" "test" {
  key_name   = %[1]q
  public_key = "ssh-rsa AAAAB3NzaC1yc2EAAAABJQAAAQEAq6U3HQYC4g8WzU147gZZ7CKQH8TgYn3chZGRPxaGmHW1RUwsyEs0nmombmIhwxudhJ4ehjqXsDLoQpd6+c7BuLgTMvbv8LgE9LX53vnljFe1dsObsr/fYLvpU9LTlo8HgHAqO5ibNdrAUvV31ronzCZhms/Gyfdaue88Fd0/YnsZVGeOZPayRkdOHSpqme2CBrpa8myBeL1CWl0LkDG4+YCURjbaelfyZlIApLYKy3FcCan9XQFKaL32MJZwCgzfOvWIMtYcU8QtXMgnA3/I3gXk8YDUJv5P4lj0s/PJXuTM8DygVAUtebNwPuinS7wwonm5FXcWMuVGsVpG5K7FGQ== tf-acc-winpasswordtest"
}

resource "aws_instance" "test" {
  ami           = "${data.aws_ami.win2016core-ami.id}"
  instance_type = "t2.medium"
  key_name      = "${aws_key_pair.test.key_name}"
}

data "aws_instance" "test" {
  instance_id = "${aws_instance.test.id}"

  get_password_data = %[2]t
}
`, rName, val)
}

func testAccInstanceDataSourceConfigGetUserData(rName string, getUserData bool) string {
	return testAccLatestAmazonLinuxHvmEbsAmiConfig() + testAccAwsInstanceVpcConfig(rName, false) + fmt.Sprintf(`
resource "aws_instance" "test" {
  ami           = "${data.aws_ami.amzn-ami-minimal-hvm-ebs.id}"
  instance_type = "t2.micro"
  subnet_id     = "${aws_subnet.test.id}"

  user_data = <<EUD
#!/bin/bash

echo "hello world"
EUD
}

data "aws_instance" "test" {
  get_user_data = %[2]t
  instance_id   = "${aws_instance.test.id}"
}
`, rName, getUserData)
}

func testAccInstanceDataSourceConfigGetUserDataNoUserData(rName string, getUserData bool) string {
	return testAccLatestAmazonLinuxHvmEbsAmiConfig() + testAccAwsInstanceVpcConfig(rName, false) + fmt.Sprintf(`
resource "aws_instance" "test" {
  ami           = "${data.aws_ami.amzn-ami-minimal-hvm-ebs.id}"
  instance_type = "t2.micro"
  subnet_id     = "${aws_subnet.test.id}"
}

data "aws_instance" "test" {
  get_user_data = %[2]t
  instance_id   = "${aws_instance.test.id}"
}
`, rName, getUserData)
}

func testAccInstanceDataSourceConfig_creditSpecification(rName string) string {
	return testAccLatestAmazonLinuxHvmEbsAmiConfig() + testAccAwsInstanceVpcConfig(rName, false) + fmt.Sprintf(`
resource "aws_instance" "test" {
  ami           = "${data.aws_ami.amzn-ami-minimal-hvm-ebs.id}"
  instance_type = "t2.micro"
  subnet_id     = "${aws_subnet.test.id}"

  credit_specification {
    cpu_credits = "unlimited"
  }
}

data "aws_instance" "test" {
  instance_id = "${aws_instance.test.id}"
}
`)
}

func testAccInstanceDataSourceConfig_metadataOptions(rName string) string {
	return composeConfig(
		testAccLatestAmazonLinuxHvmEbsAmiConfig(),
		testAccAvailableEc2InstanceTypeForRegion("t3.micro", "t2.micro"),
		fmt.Sprintf(`
resource "aws_instance" "test" {
  ami           = data.aws_ami.amzn-ami-minimal-hvm-ebs.id
  instance_type = data.aws_ec2_instance_type_offering.available.instance_type

  tags = {
    Name = %[1]q
  }

  metadata_options {
    http_endpoint               = "enabled"
    http_tokens                 = "required"
    http_put_response_hop_limit = 2
  }
}

data "aws_instance" "test" {
  instance_id = aws_instance.test.id
}
`, rName))
}
