package aws

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/batch"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccAWSBatchJobDefinition_basic(t *testing.T) {
	var jd batch.JobDefinition
	rName := acctest.RandomWithPrefix("tf-acc-test")
	resourceName := "aws_batch_job_definition.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t); testAccPreCheckAWSBatch(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBatchJobDefinitionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBatchJobDefinitionConfigName(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBatchJobDefinitionExists(resourceName, &jd),
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

func TestAccAWSBatchJobDefinition_ContainerProperties_Advanced(t *testing.T) {
	var jd batch.JobDefinition
	compare := batch.JobDefinition{
		Parameters: map[string]*string{
			"param1": aws.String("val1"),
			"param2": aws.String("val2"),
		},
		RetryStrategy: &batch.RetryStrategy{
			Attempts: aws.Int64(int64(1)),
		},
		Timeout: &batch.JobTimeout{
			AttemptDurationSeconds: aws.Int64(int64(60)),
		},
		ContainerProperties: &batch.ContainerProperties{
			Command: []*string{aws.String("ls"), aws.String("-la")},
			Environment: []*batch.KeyValuePair{
				{Name: aws.String("VARNAME"), Value: aws.String("VARVAL")},
			},
			Image:  aws.String("busybox"),
			Memory: aws.Int64(int64(512)),
			MountPoints: []*batch.MountPoint{
				{ContainerPath: aws.String("/tmp"), ReadOnly: aws.Bool(false), SourceVolume: aws.String("tmp")},
			},
			ResourceRequirements: []*batch.ResourceRequirement{},
			Ulimits: []*batch.Ulimit{
				{HardLimit: aws.Int64(int64(1024)), Name: aws.String("nofile"), SoftLimit: aws.Int64(int64(1024))},
			},
			Vcpus: aws.Int64(int64(1)),
			Volumes: []*batch.Volume{
				{
					Host: &batch.Host{SourcePath: aws.String("/tmp")},
					Name: aws.String("tmp"),
				},
			},
		},
	}
	rName := acctest.RandomWithPrefix("tf-acc-test")
	resourceName := "aws_batch_job_definition.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t); testAccPreCheckAWSBatch(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBatchJobDefinitionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBatchJobDefinitionConfigContainerPropertiesAdvanced(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBatchJobDefinitionExists(resourceName, &jd),
					testAccCheckBatchJobDefinitionAttributes(&jd, &compare),
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

func TestAccAWSBatchJobDefinition_updateForcesNewResource(t *testing.T) {
	var before, after batch.JobDefinition
	rName := acctest.RandomWithPrefix("tf-acc-test")
	resourceName := "aws_batch_job_definition.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t); testAccPreCheckAWSBatch(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBatchJobDefinitionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBatchJobDefinitionConfigContainerPropertiesAdvanced(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBatchJobDefinitionExists(resourceName, &before),
					testAccCheckBatchJobDefinitionAttributes(&before, nil),
				),
			},
			{
				Config: testAccBatchJobDefinitionConfigContainerPropertiesAdvancedUpdate(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBatchJobDefinitionExists(resourceName, &after),
					testAccCheckJobDefinitionRecreated(t, &before, &after),
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

func testAccCheckBatchJobDefinitionExists(n string, jd *batch.JobDefinition) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Batch Job Queue ID is set")
		}

		conn := testAccProvider.Meta().(*AWSClient).batchconn
		arn := rs.Primary.Attributes["arn"]
		def, err := getJobDefinition(conn, arn)
		if err != nil {
			return err
		}
		if def == nil {
			return fmt.Errorf("Not found: %s", n)
		}
		*jd = *def

		return nil
	}
}

func testAccCheckBatchJobDefinitionAttributes(jd *batch.JobDefinition, compare *batch.JobDefinition) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != "aws_batch_job_definition" {
				continue
			}
			if *jd.JobDefinitionArn != rs.Primary.Attributes["arn"] {
				return fmt.Errorf("Bad Job Definition ARN\n\t expected: %s\n\tgot: %s\n", rs.Primary.Attributes["arn"], *jd.JobDefinitionArn)
			}
			if compare != nil {
				if compare.Parameters != nil && !reflect.DeepEqual(compare.Parameters, jd.Parameters) {
					return fmt.Errorf("Bad Job Definition Params\n\t expected: %v\n\tgot: %v\n", compare.Parameters, jd.Parameters)
				}
				if compare.RetryStrategy != nil && *compare.RetryStrategy.Attempts != *jd.RetryStrategy.Attempts {
					return fmt.Errorf("Bad Job Definition Retry Strategy\n\t expected: %d\n\tgot: %d\n", *compare.RetryStrategy.Attempts, *jd.RetryStrategy.Attempts)
				}
				if compare.ContainerProperties != nil && compare.ContainerProperties.Command != nil && !reflect.DeepEqual(compare.ContainerProperties, jd.ContainerProperties) {
					return fmt.Errorf("Bad Job Definition Container Properties\n\t expected: %s\n\tgot: %s\n", compare.ContainerProperties, jd.ContainerProperties)
				}
			}
		}
		return nil
	}
}

func testAccCheckJobDefinitionRecreated(t *testing.T,
	before, after *batch.JobDefinition) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if *before.Revision == *after.Revision {
			t.Fatalf("Expected change of JobDefinition Revisions, but both were %v", before.Revision)
		}
		return nil
	}
}

func testAccCheckBatchJobDefinitionDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aws_batch_job_definition" {
			continue
		}
		conn := testAccProvider.Meta().(*AWSClient).batchconn
		js, err := getJobDefinition(conn, rs.Primary.Attributes["arn"])
		if err == nil && js != nil {
			if *js.Status == "ACTIVE" {
				return fmt.Errorf("Error: Job Definition still active")
			}
		}
		return nil
	}
	return nil
}

func testAccBatchJobDefinitionConfigContainerPropertiesAdvanced(rName string) string {
	return fmt.Sprintf(`
resource "aws_batch_job_definition" "test" {
	name = %[1]q
	type = "container"
	parameters = {
		param1 = "val1"
		param2 = "val2"
	}
	retry_strategy {
		attempts = 1
	}
	timeout {
		attempt_duration_seconds = 60
	}
	container_properties = <<CONTAINER_PROPERTIES
{
	"command": ["ls", "-la"],
	"image": "busybox",
	"memory": 512,
	"vcpus": 1,
	"volumes": [
      {
        "host": {
          "sourcePath": "/tmp"
        },
        "name": "tmp"
      }
    ],
	"environment": [
		{"name": "VARNAME", "value": "VARVAL"}
	],
	"mountPoints": [
		{
          "sourceVolume": "tmp",
          "containerPath": "/tmp",
          "readOnly": false
        }
	],
    "ulimits": [
      {
        "hardLimit": 1024,
        "name": "nofile",
        "softLimit": 1024
      }
    ]
}
CONTAINER_PROPERTIES
}
`, rName)
}

func testAccBatchJobDefinitionConfigContainerPropertiesAdvancedUpdate(rName string) string {
	return fmt.Sprintf(`
resource "aws_batch_job_definition" "test" {
	name = %[1]q
	type = "container"
	container_properties = <<CONTAINER_PROPERTIES
{
	"command": ["ls", "-la"],
	"image": "busybox",
	"memory": 1024,
	"vcpus": 1,
	"volumes": [
      {
        "host": {
          "sourcePath": "/tmp"
        },
        "name": "tmp"
      }
    ],
	"environment": [
		{"name": "VARNAME", "value": "VARVAL"}
	],
	"mountPoints": [
		{
          "sourceVolume": "tmp",
          "containerPath": "/tmp",
          "readOnly": false
        }
	],
    "ulimits": [
      {
        "hardLimit": 1024,
        "name": "nofile",
        "softLimit": 1024
      }
    ]
}
CONTAINER_PROPERTIES
}
`, rName)
}

func testAccBatchJobDefinitionConfigName(rName string) string {
	return fmt.Sprintf(`
resource "aws_batch_job_definition" "test" {
  container_properties = jsonencode({
    command = ["echo", "test"]
    image   = "busybox"
    memory  = 128
    vcpus   = 1
  })
  name = %[1]q
  type = "container"
}
`, rName)
}
