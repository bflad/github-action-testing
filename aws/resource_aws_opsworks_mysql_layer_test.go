package aws

import (
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/service/opsworks"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

// These tests assume the existence of predefined Opsworks IAM roles named `aws-opsworks-ec2-role`
// and `aws-opsworks-service-role`.

func TestAccAWSOpsworksMysqlLayer_basic(t *testing.T) {
	var opslayer opsworks.Layer
	stackName := acctest.RandomWithPrefix("tf-acc-test")
	resourceName := "aws_opsworks_mysql_layer.test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAwsOpsworksMysqlLayerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAwsOpsworksMysqlLayerConfigVpcCreate(stackName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSOpsworksLayerExists(resourceName, &opslayer),
					resource.TestCheckResourceAttr(resourceName, "name", stackName)),
			},
		},
	})
}

func TestAccAWSOpsworksMysqlLayer_tags(t *testing.T) {
	var opslayer opsworks.Layer
	stackName := acctest.RandomWithPrefix("tf-acc-test")
	resourceName := "aws_opsworks_mysql_layer.test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAwsOpsworksMysqlLayerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAwsOpsworksMysqlLayerConfigTags1(stackName, "key1", "value1"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSOpsworksLayerExists(resourceName, &opslayer),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.key1", "value1"),
				),
			},
			{
				Config: testAccAwsOpsworksMysqlLayerConfigTags2(stackName, "key1", "value1updated", "key2", "value2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSOpsworksLayerExists(resourceName, &opslayer),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.key1", "value1updated"),
					resource.TestCheckResourceAttr(resourceName, "tags.key2", "value2"),
				),
			},
			{
				Config: testAccAwsOpsworksMysqlLayerConfigTags1(stackName, "key2", "value2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSOpsworksLayerExists(resourceName, &opslayer),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.key2", "value2"),
				),
			},
		},
	})
}

func testAccCheckAwsOpsworksMysqlLayerDestroy(s *terraform.State) error {
	return testAccCheckAwsOpsworksLayerDestroy("aws_opsworks_mysql_layer", s)
}

func testAccAwsOpsworksMysqlLayerConfigVpcCreate(name string) string {
	return testAccAwsOpsworksStackConfigVpcCreate(name) +
		testAccAwsOpsworksCustomLayerSecurityGroups(name) +
		fmt.Sprintf(`
resource "aws_opsworks_mysql_layer" "test" {
  stack_id = "${aws_opsworks_stack.tf-acc.id}"
  name     = "%s"

  custom_security_group_ids = [
    "${aws_security_group.tf-ops-acc-layer1.id}",
    "${aws_security_group.tf-ops-acc-layer2.id}",
  ]
}
`, name)
}

func testAccAwsOpsworksMysqlLayerConfigTags1(name, tagKey1, tagValue1 string) string {
	return testAccAwsOpsworksStackConfigVpcCreate(name) +
		testAccAwsOpsworksCustomLayerSecurityGroups(name) +
		fmt.Sprintf(`
resource "aws_opsworks_mysql_layer" "test" {
  stack_id = "${aws_opsworks_stack.tf-acc.id}"
  name     = "%s"

  custom_security_group_ids = [
    "${aws_security_group.tf-ops-acc-layer1.id}",
    "${aws_security_group.tf-ops-acc-layer2.id}",
  ]

  tags = {
    %[2]q = %[3]q
  }
}
`, name, tagKey1, tagValue1)
}

func testAccAwsOpsworksMysqlLayerConfigTags2(name, tagKey1, tagValue1, tagKey2, tagValue2 string) string {
	return testAccAwsOpsworksStackConfigVpcCreate(name) +
		testAccAwsOpsworksCustomLayerSecurityGroups(name) +
		fmt.Sprintf(`
resource "aws_opsworks_mysql_layer" "test" {
  stack_id = "${aws_opsworks_stack.tf-acc.id}"
  name     = "%s"

  custom_security_group_ids = [
    "${aws_security_group.tf-ops-acc-layer1.id}",
    "${aws_security_group.tf-ops-acc-layer2.id}",
  ]

  tags = {
    %[2]q = %[3]q
    %[4]q = %[5]q
  }
}
`, name, tagKey1, tagValue1, tagKey2, tagValue2)
}
