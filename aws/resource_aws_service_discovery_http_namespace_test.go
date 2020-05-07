package aws

import (
	"fmt"
	"log"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/servicediscovery"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-aws/aws/internal/service/servicediscovery/waiter"
)

func init() {
	resource.AddTestSweepers("aws_service_discovery_http_namespace", &resource.Sweeper{
		Name: "aws_service_discovery_http_namespace",
		F:    testSweepServiceDiscoveryHttpNamespaces,
		Dependencies: []string{
			"aws_service_discovery_service",
		},
	})
}

func testSweepServiceDiscoveryHttpNamespaces(region string) error {
	client, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting client: %s", err)
	}
	conn := client.(*AWSClient).sdconn
	var sweeperErrs *multierror.Error

	input := &servicediscovery.ListNamespacesInput{
		Filters: []*servicediscovery.NamespaceFilter{
			{
				Condition: aws.String(servicediscovery.FilterConditionEq),
				Name:      aws.String(servicediscovery.NamespaceFilterNameType),
				Values:    aws.StringSlice([]string{servicediscovery.NamespaceTypeHttp}),
			},
		},
	}

	err = conn.ListNamespacesPages(input, func(page *servicediscovery.ListNamespacesOutput, isLast bool) bool {
		if page == nil {
			return !isLast
		}

		for _, namespace := range page.Namespaces {
			if namespace == nil {
				continue
			}

			id := aws.StringValue(namespace.Id)
			input := &servicediscovery.DeleteNamespaceInput{
				Id: namespace.Id,
			}

			log.Printf("[INFO] Deleting Service Discovery HTTP Namespace: %s", id)
			output, err := conn.DeleteNamespace(input)

			if err != nil {
				sweeperErr := fmt.Errorf("error deleting Service Discovery HTTP Namespace (%s): %w", id, err)
				log.Printf("[ERROR] %s", sweeperErr)
				sweeperErrs = multierror.Append(sweeperErrs, sweeperErr)
				continue
			}

			if output != nil && output.OperationId != nil {
				if _, err := waiter.OperationSuccess(conn, aws.StringValue(output.OperationId)); err != nil {
					sweeperErr := fmt.Errorf("error waiting for Service Discovery HTTP Namespace (%s) deletion: %w", id, err)
					log.Printf("[ERROR] %s", sweeperErr)
					sweeperErrs = multierror.Append(sweeperErrs, sweeperErr)
					continue
				}
			}
		}

		return !isLast
	})

	if testSweepSkipSweepError(err) {
		log.Printf("[WARN] Skipping Service Discovery HTTP Namespaces sweep for %s: %s", region, err)
		return sweeperErrs.ErrorOrNil() // In case we have completed some pages, but had errors
	}

	if err != nil {
		sweeperErrs = multierror.Append(sweeperErrs, fmt.Errorf("error retrieving Service Discovery HTTP Namespaces: %w", err))
	}

	return sweeperErrs.ErrorOrNil()
}

func TestAccAWSServiceDiscoveryHttpNamespace_basic(t *testing.T) {
	resourceName := "aws_service_discovery_http_namespace.test"
	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandStringFromCharSet(8, acctest.CharSetAlpha))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t); testAccPreCheckAWSServiceDiscovery(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAwsServiceDiscoveryHttpNamespaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccServiceDiscoveryHttpNamespaceConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAwsServiceDiscoveryHttpNamespaceExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "arn"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
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

func TestAccAWSServiceDiscoveryHttpNamespace_Description(t *testing.T) {
	resourceName := "aws_service_discovery_http_namespace.test"
	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandStringFromCharSet(8, acctest.CharSetAlpha))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t); testAccPreCheckAWSServiceDiscovery(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAwsServiceDiscoveryHttpNamespaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccServiceDiscoveryHttpNamespaceConfigDescription(rName, "test"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAwsServiceDiscoveryHttpNamespaceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "description", "test"),
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

func testAccCheckAwsServiceDiscoveryHttpNamespaceDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*AWSClient).sdconn

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aws_service_discovery_http_namespace" {
			continue
		}

		input := &servicediscovery.GetNamespaceInput{
			Id: aws.String(rs.Primary.ID),
		}

		_, err := conn.GetNamespace(input)
		if err != nil {
			if isAWSErr(err, servicediscovery.ErrCodeNamespaceNotFound, "") {
				continue
			}
			return err
		}
	}
	return nil
}

func testAccCheckAwsServiceDiscoveryHttpNamespaceExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		conn := testAccProvider.Meta().(*AWSClient).sdconn

		input := &servicediscovery.GetNamespaceInput{
			Id: aws.String(rs.Primary.ID),
		}

		_, err := conn.GetNamespace(input)
		return err
	}
}

func testAccServiceDiscoveryHttpNamespaceConfig(rName string) string {
	return fmt.Sprintf(`
resource "aws_service_discovery_http_namespace" "test" {
  name = %q
}
`, rName)
}

func testAccServiceDiscoveryHttpNamespaceConfigDescription(rName, description string) string {
	return fmt.Sprintf(`
resource "aws_service_discovery_http_namespace" "test" {
  description = %q
  name        = %q
}
`, description, rName)
}
