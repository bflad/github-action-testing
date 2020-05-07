package aws

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-aws/aws/internal/keyvaluetags"
)

// add sweeper to delete known test subnets
func init() {
	resource.AddTestSweepers("aws_subnet", &resource.Sweeper{
		Name: "aws_subnet",
		F:    testSweepSubnets,
		Dependencies: []string{
			"aws_autoscaling_group",
			"aws_batch_compute_environment",
			"aws_elastic_beanstalk_environment",
			"aws_cloudhsm_v2_cluster",
			"aws_db_subnet_group",
			"aws_directory_service_directory",
			"aws_ec2_client_vpn_endpoint",
			"aws_ec2_transit_gateway_vpc_attachment",
			"aws_efs_file_system",
			"aws_eks_cluster",
			"aws_elasticache_cluster",
			"aws_elasticache_replication_group",
			"aws_elasticsearch_domain",
			"aws_elb",
			"aws_emr_cluster",
			"aws_fsx_lustre_file_system",
			"aws_fsx_windows_file_system",
			"aws_lambda_function",
			"aws_lb",
			"aws_mq_broker",
			"aws_msk_cluster",
			"aws_network_interface",
			"aws_redshift_cluster",
			"aws_route53_resolver_endpoint",
			"aws_sagemaker_notebook_instance",
			"aws_spot_fleet_request",
			"aws_vpc_endpoint",
		},
	})
}

func testSweepSubnets(region string) error {
	client, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting client: %s", err)
	}
	conn := client.(*AWSClient).ec2conn
	input := &ec2.DescribeSubnetsInput{}
	var sweeperErrs *multierror.Error

	err = conn.DescribeSubnetsPages(input, func(page *ec2.DescribeSubnetsOutput, lastPage bool) bool {
		for _, subnet := range page.Subnets {
			if subnet == nil {
				continue
			}

			id := aws.StringValue(subnet.SubnetId)
			input := &ec2.DeleteSubnetInput{
				SubnetId: subnet.SubnetId,
			}

			if aws.BoolValue(subnet.DefaultForAz) {
				log.Printf("[DEBUG] Skipping default EC2 Subnet: %s", id)
				continue
			}

			log.Printf("[INFO] Deleting EC2 Subnet: %s", id)

			// Handle eventual consistency, especially with lingering ENIs from Load Balancers and Lambda
			err := resource.Retry(5*time.Minute, func() *resource.RetryError {
				_, err := conn.DeleteSubnet(input)

				if isAWSErr(err, "DependencyViolation", "") {
					return resource.RetryableError(err)
				}

				if err != nil {
					return resource.NonRetryableError(err)
				}

				return nil
			})

			if isResourceTimeoutError(err) {
				_, err = conn.DeleteSubnet(input)
			}

			if err != nil {
				sweeperErr := fmt.Errorf("error deleting EC2 Subnet (%s): %w", id, err)
				log.Printf("[ERROR] %s", sweeperErr)
				sweeperErrs = multierror.Append(sweeperErrs, sweeperErr)
			}
		}

		return !lastPage
	})

	if testSweepSkipSweepError(err) {
		log.Printf("[WARN] Skipping EC2 Subnet sweep for %s: %s", region, err)
		return nil
	}

	if err != nil {
		return fmt.Errorf("Error describing subnets: %s", err)
	}

	return sweeperErrs.ErrorOrNil()
}

func TestAccAWSSubnet_basic(t *testing.T) {
	var v ec2.Subnet
	resourceName := "aws_subnet.test"

	testCheck := func(*terraform.State) error {
		if aws.StringValue(v.CidrBlock) != "10.1.1.0/24" {
			return fmt.Errorf("bad cidr: %s", aws.StringValue(v.CidrBlock))
		}

		if !aws.BoolValue(v.MapPublicIpOnLaunch) {
			return fmt.Errorf("bad MapPublicIpOnLaunch: %t", aws.BoolValue(v.MapPublicIpOnLaunch))
		}

		return nil
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:      func() { testAccPreCheck(t) },
		IDRefreshName: resourceName,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckSubnetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSubnetConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSubnetExists(resourceName, &v),
					testCheck,
					// ipv6 should be empty if disabled so we can still use the property in conditionals
					resource.TestCheckResourceAttr(resourceName, "ipv6_cidr_block", ""),
					testAccMatchResourceAttrRegionalARN(resourceName, "arn", "ec2", regexp.MustCompile("subnet/subnet-.+$")),
					testAccCheckResourceAttrAccountID(resourceName, "owner_id"),
					resource.TestCheckResourceAttrSet(resourceName, "availability_zone"),
					resource.TestCheckResourceAttrSet(resourceName, "availability_zone_id"),
					resource.TestCheckResourceAttr(resourceName, "outpost_arn", ""),
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

func TestAccAWSSubnet_ignoreTags(t *testing.T) {
	var providers []*schema.Provider
	var subnet ec2.Subnet
	resourceName := "aws_subnet.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories(&providers),
		CheckDestroy:      testAccCheckVpcDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSubnetConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSubnetExists(resourceName, &subnet),
					testAccCheckSubnetUpdateTags(&subnet, nil, map[string]string{"ignorekey1": "ignorevalue1"}),
				),
				ExpectNonEmptyPlan: true,
			},
			{
				Config:   testAccProviderConfigIgnoreTagsKeyPrefixes1("ignorekey") + testAccSubnetConfig,
				PlanOnly: true,
			},
			{
				Config:   testAccProviderConfigIgnoreTagsKeys1("ignorekey1") + testAccSubnetConfig,
				PlanOnly: true,
			},
		},
	})
}

func TestAccAWSSubnet_ipv6(t *testing.T) {
	var before, after ec2.Subnet
	resourceName := "aws_subnet.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:      func() { testAccPreCheck(t) },
		IDRefreshName: resourceName,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckSubnetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSubnetConfigIpv6,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSubnetExists(resourceName, &before),
					testAccCheckAwsSubnetIpv6BeforeUpdate(t, &before),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccSubnetConfigIpv6UpdateAssignIpv6OnCreation,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSubnetExists(resourceName, &after),
					testAccCheckAwsSubnetIpv6AfterUpdate(t, &after),
				),
			},
			{
				Config: testAccSubnetConfigIpv6UpdateIpv6Cidr,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSubnetExists(resourceName, &after),

					testAccCheckAwsSubnetNotRecreated(t, &before, &after),
				),
			},
		},
	})
}

func TestAccAWSSubnet_enableIpv6(t *testing.T) {
	var subnet ec2.Subnet
	resourceName := "aws_subnet.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:      func() { testAccPreCheck(t) },
		IDRefreshName: resourceName,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckSubnetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSubnetConfigPreIpv6,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSubnetExists(resourceName, &subnet),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccSubnetConfigIpv6,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSubnetExists(resourceName, &subnet),
				),
			},
		},
	})
}

func TestAccAWSSubnet_availabilityZoneId(t *testing.T) {
	var v ec2.Subnet
	resourceName := "aws_subnet.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:      func() { testAccPreCheck(t) },
		IDRefreshName: resourceName,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckSubnetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSubnetConfigAvailabilityZoneId,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSubnetExists(resourceName, &v),
					resource.TestCheckResourceAttrSet(resourceName, "availability_zone"),
					resource.TestCheckResourceAttr(resourceName, "availability_zone_id", "usw2-az3"),
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

func TestAccAWSSubnet_outpost(t *testing.T) {
	var v ec2.Subnet
	resourceName := "aws_subnet.test"

	outpostArn := os.Getenv("AWS_OUTPOST_ARN")
	if outpostArn == "" {
		t.Skip(
			"Environment variable AWS_OUTPOST_ARN is not set. " +
				"This environment variable must be set to the ARN of " +
				"a deployed Outpost to enable this test.")
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:      func() { testAccPreCheck(t) },
		IDRefreshName: resourceName,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckSubnetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSubnetConfigOutpost(outpostArn),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSubnetExists(
						resourceName, &v),
					resource.TestCheckResourceAttr(
						resourceName, "outpost_arn", outpostArn),
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

func testAccCheckAwsSubnetIpv6BeforeUpdate(t *testing.T, subnet *ec2.Subnet) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if subnet.Ipv6CidrBlockAssociationSet == nil {
			return fmt.Errorf("Expected IPV6 CIDR Block Association")
		}

		if !aws.BoolValue(subnet.AssignIpv6AddressOnCreation) {
			return fmt.Errorf("bad AssignIpv6AddressOnCreation: %t", aws.BoolValue(subnet.AssignIpv6AddressOnCreation))
		}

		return nil
	}
}

func testAccCheckAwsSubnetIpv6AfterUpdate(t *testing.T, subnet *ec2.Subnet) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if aws.BoolValue(subnet.AssignIpv6AddressOnCreation) {
			return fmt.Errorf("bad AssignIpv6AddressOnCreation: %t", aws.BoolValue(subnet.AssignIpv6AddressOnCreation))
		}

		return nil
	}
}

func testAccCheckAwsSubnetNotRecreated(t *testing.T,
	before, after *ec2.Subnet) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if *before.SubnetId != *after.SubnetId {
			t.Fatalf("Expected SubnetIDs not to change, but both got before: %s and after: %s", *before.SubnetId, *after.SubnetId)
		}
		return nil
	}
}

func testAccCheckSubnetDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*AWSClient).ec2conn

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aws_subnet" {
			continue
		}

		// Try to find the resource
		resp, err := conn.DescribeSubnets(&ec2.DescribeSubnetsInput{
			SubnetIds: []*string{aws.String(rs.Primary.ID)},
		})
		if err == nil {
			if len(resp.Subnets) > 0 {
				return fmt.Errorf("still exist.")
			}

			return nil
		}

		// Verify the error is what we want
		ec2err, ok := err.(awserr.Error)
		if !ok {
			return err
		}
		if ec2err.Code() != "InvalidSubnetID.NotFound" {
			return err
		}
	}

	return nil
}

func testAccCheckSubnetExists(n string, v *ec2.Subnet) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		conn := testAccProvider.Meta().(*AWSClient).ec2conn
		resp, err := conn.DescribeSubnets(&ec2.DescribeSubnetsInput{
			SubnetIds: []*string{aws.String(rs.Primary.ID)},
		})
		if err != nil {
			return err
		}
		if len(resp.Subnets) == 0 {
			return fmt.Errorf("Subnet not found")
		}

		*v = *resp.Subnets[0]

		return nil
	}
}

func testAccCheckSubnetUpdateTags(subnet *ec2.Subnet, oldTags, newTags map[string]string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := testAccProvider.Meta().(*AWSClient).ec2conn

		return keyvaluetags.Ec2UpdateTags(conn, aws.StringValue(subnet.SubnetId), oldTags, newTags)
	}
}

const testAccSubnetConfig = `
resource "aws_vpc" "test" {
	cidr_block = "10.1.0.0/16"
	tags = {
		Name = "terraform-testacc-subnet"
	}
}

resource "aws_subnet" "test" {
	cidr_block = "10.1.1.0/24"
	vpc_id = "${aws_vpc.test.id}"
	map_public_ip_on_launch = true
	tags = {
		Name = "tf-acc-subnet"
	}
}
`

const testAccSubnetConfigPreIpv6 = `
resource "aws_vpc" "test" {
	cidr_block = "10.10.0.0/16"
	assign_generated_ipv6_cidr_block = true
	tags = {
		Name = "terraform-testacc-subnet-ipv6"
	}
}

resource "aws_subnet" "test" {
	cidr_block = "10.10.1.0/24"
	vpc_id = "${aws_vpc.test.id}"
	map_public_ip_on_launch = true
	tags = {
		Name = "tf-acc-subnet-ipv6"
	}
}
`

const testAccSubnetConfigIpv6 = `
resource "aws_vpc" "test" {
	cidr_block = "10.10.0.0/16"
	assign_generated_ipv6_cidr_block = true
	tags = {
		Name = "terraform-testacc-subnet-ipv6"
	}
}

resource "aws_subnet" "test" {
	cidr_block = "10.10.1.0/24"
	vpc_id = "${aws_vpc.test.id}"
	ipv6_cidr_block = "${cidrsubnet(aws_vpc.test.ipv6_cidr_block, 8, 1)}"
	map_public_ip_on_launch = true
	assign_ipv6_address_on_creation = true
	tags = {
		Name = "tf-acc-subnet-ipv6"
	}
}
`

const testAccSubnetConfigIpv6UpdateAssignIpv6OnCreation = `
resource "aws_vpc" "test" {
	cidr_block = "10.10.0.0/16"
	assign_generated_ipv6_cidr_block = true
	tags = {
		Name = "terraform-testacc-subnet-assign-ipv6-on-creation"
	}
}

resource "aws_subnet" "test" {
	cidr_block = "10.10.1.0/24"
	vpc_id = "${aws_vpc.test.id}"
	ipv6_cidr_block = "${cidrsubnet(aws_vpc.test.ipv6_cidr_block, 8, 1)}"
	map_public_ip_on_launch = true
	assign_ipv6_address_on_creation = false
	tags = {
		Name = "tf-acc-subnet-assign-ipv6-on-creation"
	}
}
`

const testAccSubnetConfigIpv6UpdateIpv6Cidr = `
resource "aws_vpc" "test" {
	cidr_block = "10.10.0.0/16"
	assign_generated_ipv6_cidr_block = true
	tags = {
		Name = "terraform-testacc-subnet-ipv6-update-cidr"
	}
}

resource "aws_subnet" "test" {
	cidr_block = "10.10.1.0/24"
	vpc_id = "${aws_vpc.test.id}"
	ipv6_cidr_block = "${cidrsubnet(aws_vpc.test.ipv6_cidr_block, 8, 3)}"
	map_public_ip_on_launch = true
	assign_ipv6_address_on_creation = false
	tags = {
		Name = "tf-acc-subnet-ipv6-update-cidr"
	}
}
`

const testAccSubnetConfigAvailabilityZoneId = `
resource "aws_vpc" "test" {
  cidr_block = "10.1.0.0/16"
  tags = {
    Name = "terraform-testacc-subnet"
  }
}

resource "aws_subnet" "test" {
  cidr_block = "10.1.1.0/24"
  vpc_id = "${aws_vpc.test.id}"
  availability_zone_id = "usw2-az3"
  tags = {
    Name = "tf-acc-subnet"
  }
}
`

func testAccSubnetConfigOutpost(outpostArn string) string {
	return fmt.Sprintf(`
data "aws_availability_zones" "current" {
  # Exclude usw2-az4 (us-west-2d) as it has limited instance types.
  blacklisted_zone_ids = ["usw2-az4"]
}

resource "aws_vpc" "test" {
  cidr_block = "10.1.0.0/16"
  tags = {
    Name = "terraform-testacc-subnet-outpost"
  }
}

resource "aws_subnet" "test" {
  cidr_block = "10.1.1.0/24"
  vpc_id = "${aws_vpc.test.id}"
  availability_zone       = "${data.aws_availability_zones.current.names[0]}"
  outpost_arn = "%s"
  tags = {
    Name = "tf-acc-subnet-outpost"
  }
}
`, outpostArn)
}
