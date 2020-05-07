package aws

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccAWSIAMInstanceProfile_basic(t *testing.T) {
	var conf iam.GetInstanceProfileOutput
	resourceName := "aws_iam_instance_profile.test"
	rName := acctest.RandString(5)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAWSInstanceProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAwsIamInstanceProfileConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSInstanceProfileExists(resourceName, &conf),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAWSIAMInstanceProfile_withRoleNotRoles(t *testing.T) {
	var conf iam.GetInstanceProfileOutput
	resourceName := "aws_iam_instance_profile.test"

	rName := acctest.RandString(5)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAWSInstanceProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSInstanceProfileWithRoleSpecified(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSInstanceProfileExists(resourceName, &conf),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"name_prefix"},
			},
		},
	})
}

func TestAccAWSIAMInstanceProfile_withoutRole(t *testing.T) {
	var conf iam.GetInstanceProfileOutput
	resourceName := "aws_iam_instance_profile.test"
	rName := acctest.RandString(5)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAWSInstanceProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAwsIamInstanceProfileConfigWithoutRole(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSInstanceProfileExists(resourceName, &conf),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"name_prefix"},
			},
		},
	})
}

func TestAccAWSIAMInstanceProfile_namePrefix(t *testing.T) {
	var conf iam.GetInstanceProfileOutput
	rName := acctest.RandString(5)
	resourceName := "aws_iam_instance_profile.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:        func() { testAccPreCheck(t) },
		IDRefreshName:   resourceName,
		IDRefreshIgnore: []string{"name_prefix"},
		Providers:       testAccProviders,
		CheckDestroy:    testAccCheckAWSInstanceProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSInstanceProfilePrefixNameConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSInstanceProfileExists(resourceName, &conf),
					testAccCheckAWSInstanceProfileGeneratedNamePrefix(
						resourceName, "test-"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"name_prefix"},
			},
		},
	})
}

func testAccCheckAWSInstanceProfileGeneratedNamePrefix(resource, prefix string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		r, ok := s.RootModule().Resources[resource]
		if !ok {
			return fmt.Errorf("Resource not found")
		}
		name, ok := r.Primary.Attributes["name"]
		if !ok {
			return fmt.Errorf("Name attr not found: %#v", r.Primary.Attributes)
		}
		if !strings.HasPrefix(name, prefix) {
			return fmt.Errorf("Name: %q, does not have prefix: %q", name, prefix)
		}
		return nil
	}
}

func testAccCheckAWSInstanceProfileDestroy(s *terraform.State) error {
	iamconn := testAccProvider.Meta().(*AWSClient).iamconn

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aws_iam_instance_profile" {
			continue
		}

		// Try to get role
		_, err := iamconn.GetInstanceProfile(&iam.GetInstanceProfileInput{
			InstanceProfileName: aws.String(rs.Primary.ID),
		})
		if err == nil {
			return fmt.Errorf("still exist.")
		}

		if isAWSErr(err, "NoSuchEntity", "") {
			continue
		}

		return err
	}

	return nil
}

func testAccCheckAWSInstanceProfileExists(n string, res *iam.GetInstanceProfileOutput) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Instance Profile name is set")
		}

		iamconn := testAccProvider.Meta().(*AWSClient).iamconn

		resp, err := iamconn.GetInstanceProfile(&iam.GetInstanceProfileInput{
			InstanceProfileName: aws.String(rs.Primary.ID),
		})
		if err != nil {
			return err
		}

		*res = *resp

		return nil
	}
}

func testAccAwsIamInstanceProfileConfig(rName string) string {
	return fmt.Sprintf(`
resource "aws_iam_role" "test" {
  name               = "test-%s"
  assume_role_policy = "{\"Version\":\"2012-10-17\",\"Statement\":[{\"Effect\":\"Allow\",\"Principal\":{\"Service\":[\"ec2.amazonaws.com\"]},\"Action\":[\"sts:AssumeRole\"]}]}"
}

resource "aws_iam_instance_profile" "test" {
  name  = "test"
  roles = ["${aws_iam_role.test.name}"]
}
`, rName)
}

func testAccAwsIamInstanceProfileConfigWithoutRole(rName string) string {
	return fmt.Sprintf(`
resource "aws_iam_instance_profile" "test" {
  name = "test-%s"
}
`, rName)
}

func testAccAWSInstanceProfilePrefixNameConfig(rName string) string {
	return fmt.Sprintf(`
resource "aws_iam_role" "test" {
  name               = "test-%s"
  assume_role_policy = "{\"Version\":\"2012-10-17\",\"Statement\":[{\"Effect\":\"Allow\",\"Principal\":{\"Service\":[\"ec2.amazonaws.com\"]},\"Action\":[\"sts:AssumeRole\"]}]}"
}

resource "aws_iam_instance_profile" "test" {
  name_prefix = "test-"
  roles       = ["${aws_iam_role.test.name}"]
}
`, rName)
}

func testAccAWSInstanceProfileWithRoleSpecified(rName string) string {
	return fmt.Sprintf(`
resource "aws_iam_role" "test" {
  name               = "test-%s"
  assume_role_policy = "{\"Version\":\"2012-10-17\",\"Statement\":[{\"Effect\":\"Allow\",\"Principal\":{\"Service\":[\"ec2.amazonaws.com\"]},\"Action\":[\"sts:AssumeRole\"]}]}"
}

resource "aws_iam_instance_profile" "test" {
  name_prefix = "test-"
  role        = "${aws_iam_role.test.name}"
}
`, rName)
}
