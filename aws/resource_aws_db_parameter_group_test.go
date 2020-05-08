package aws

import (
	"fmt"
	"log"
	"regexp"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func init() {
	resource.AddTestSweepers("aws_db_parameter_group", &resource.Sweeper{
		Name: "aws_db_parameter_group",
		F:    testSweepRdsDbParameterGroups,
		Dependencies: []string{
			"aws_db_instance",
		},
	})
}

func testSweepRdsDbParameterGroups(region string) error {
	client, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting client: %s", err)
	}
	conn := client.(*AWSClient).rdsconn

	err = conn.DescribeDBParameterGroupsPages(&rds.DescribeDBParameterGroupsInput{}, func(out *rds.DescribeDBParameterGroupsOutput, lastPage bool) bool {
		for _, dbpg := range out.DBParameterGroups {
			if dbpg == nil {
				continue
			}

			input := &rds.DeleteDBParameterGroupInput{
				DBParameterGroupName: dbpg.DBParameterGroupName,
			}
			name := aws.StringValue(dbpg.DBParameterGroupName)

			if strings.HasPrefix(name, "default.") {
				log.Printf("[INFO] Skipping DB Parameter Group: %s", name)
				continue
			}

			log.Printf("[INFO] Deleting DB Parameter Group: %s", name)

			_, err := conn.DeleteDBParameterGroup(input)

			if err != nil {
				log.Printf("[ERROR] Failed to delete DB Parameter Group %s: %s", name, err)
				continue
			}
		}

		return !lastPage
	})

	if testSweepSkipSweepError(err) {
		log.Printf("[WARN] Skipping RDS DB Parameter Group sweep for %s: %s", region, err)
		return nil
	}

	if err != nil {
		return fmt.Errorf("error retrieving DB Parameter Groups: %s", err)
	}

	return nil
}

func TestAccAWSDBParameterGroup_basic(t *testing.T) {
	var v rds.DBParameterGroup
	resourceName := "aws_db_parameter_group.test"
	groupName := fmt.Sprintf("parameter-group-test-terraform-%d", acctest.RandInt())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:            func() { testAccPreCheck(t) },
		Providers:           testAccProviders,
		CheckDestroy:        testAccCheckAWSDBParameterGroupDestroy,
		DisableBinaryDriver: true,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSDBParameterGroupConfig(groupName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSDBParameterGroupExists(resourceName, &v),
					testAccCheckAWSDBParameterGroupAttributes(&v, groupName),
					resource.TestCheckResourceAttr(resourceName, "name", groupName),
					resource.TestCheckResourceAttr(resourceName, "family", "mysql5.6"),
					resource.TestCheckResourceAttr(resourceName, "parameter.1708034931.name", "character_set_results"),
					resource.TestCheckResourceAttr(resourceName, "parameter.1708034931.value", "utf8"),
					resource.TestCheckResourceAttr(resourceName, "parameter.2421266705.name", "character_set_server"),
					resource.TestCheckResourceAttr(resourceName, "parameter.2421266705.value", "utf8"),
					resource.TestCheckResourceAttr(resourceName, "parameter.2478663599.name", "character_set_client"),
					resource.TestCheckResourceAttr(resourceName, "parameter.2478663599.value", "utf8"),
					testAccMatchResourceAttrRegionalARN(resourceName, "arn", "rds", regexp.MustCompile(fmt.Sprintf("pg:%s$", groupName))),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccAWSDBParameterGroupAddParametersConfig(groupName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSDBParameterGroupExists(resourceName, &v),
					testAccCheckAWSDBParameterGroupAttributes(&v, groupName),
					resource.TestCheckResourceAttr(resourceName, "name", groupName),
					resource.TestCheckResourceAttr(resourceName, "family", "mysql5.6"),
					resource.TestCheckResourceAttr(resourceName, "parameter.1706463059.name", "collation_connection"),
					resource.TestCheckResourceAttr(resourceName, "parameter.1706463059.value", "utf8_unicode_ci"),
					resource.TestCheckResourceAttr(resourceName, "parameter.1708034931.name", "character_set_results"),
					resource.TestCheckResourceAttr(resourceName, "parameter.1708034931.value", "utf8"),
					resource.TestCheckResourceAttr(resourceName, "parameter.2421266705.name", "character_set_server"),
					resource.TestCheckResourceAttr(resourceName, "parameter.2421266705.value", "utf8"),
					resource.TestCheckResourceAttr(resourceName, "parameter.2475805061.name", "collation_server"),
					resource.TestCheckResourceAttr(resourceName, "parameter.2475805061.value", "utf8_unicode_ci"),
					resource.TestCheckResourceAttr(resourceName, "parameter.2478663599.name", "character_set_client"),
					resource.TestCheckResourceAttr(resourceName, "parameter.2478663599.value", "utf8"),
					testAccMatchResourceAttrRegionalARN(resourceName, "arn", "rds", regexp.MustCompile(fmt.Sprintf("pg:%s$", groupName))),
				),
			},
			{
				Config: testAccAWSDBParameterGroupConfig(groupName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSDBParameterGroupExists(resourceName, &v),
					testAccCheckAWSDBParameterGroupAttributes(&v, groupName),
					testAccCheckAWSDBParameterNotUserDefined(resourceName, "collation_connection"),
					testAccCheckAWSDBParameterNotUserDefined(resourceName, "collation_server"),
					resource.TestCheckNoResourceAttr(resourceName, "parameter.2475805061.value"),
					resource.TestCheckNoResourceAttr(resourceName, "parameter.1706463059.value"),
				),
			},
		},
	})
}

func TestAccAWSDBParameterGroup_limit(t *testing.T) {
	var v rds.DBParameterGroup
	resourceName := "aws_db_parameter_group.test"
	groupName := fmt.Sprintf("parameter-group-test-terraform-%d", acctest.RandInt())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:            func() { testAccPreCheck(t) },
		Providers:           testAccProviders,
		CheckDestroy:        testAccCheckAWSDBParameterGroupDestroy,
		DisableBinaryDriver: true,
		Steps: []resource.TestStep{
			{
				Config: createAwsDbParameterGroupsExceedDefaultAwsLimit(groupName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSDBParameterGroupExists(resourceName, &v),
					testAccCheckAWSDBParameterGroupAttributes(&v, groupName),
					resource.TestCheckResourceAttr(resourceName, "name", groupName),
					resource.TestCheckResourceAttr(resourceName, "family", "mysql5.6"),
					resource.TestCheckResourceAttr(resourceName, "description", "RDS default parameter group: Exceed default AWS parameter group limit of twenty"),
					resource.TestCheckResourceAttr(resourceName, "parameter.2421266705.name", "character_set_server"),
					resource.TestCheckResourceAttr(resourceName, "parameter.2421266705.value", "utf8"),
					resource.TestCheckResourceAttr(resourceName, "parameter.2478663599.name", "character_set_client"),
					resource.TestCheckResourceAttr(resourceName, "parameter.2478663599.value", "utf8"),
					resource.TestCheckResourceAttr(resourceName, "parameter.1680942586.name", "collation_server"),
					resource.TestCheckResourceAttr(resourceName, "parameter.1680942586.value", "utf8_general_ci"),
					resource.TestCheckResourceAttr(resourceName, "parameter.2450940716.name", "collation_connection"),
					resource.TestCheckResourceAttr(resourceName, "parameter.2450940716.value", "utf8_general_ci"),
					resource.TestCheckResourceAttr(resourceName, "parameter.242489837.name", "join_buffer_size"),
					resource.TestCheckResourceAttr(resourceName, "parameter.242489837.value", "16777216"),
					resource.TestCheckResourceAttr(resourceName, "parameter.2026669454.name", "key_buffer_size"),
					resource.TestCheckResourceAttr(resourceName, "parameter.2026669454.value", "67108864"),
					resource.TestCheckResourceAttr(resourceName, "parameter.2705275319.name", "max_connections"),
					resource.TestCheckResourceAttr(resourceName, "parameter.2705275319.value", "3200"),
					resource.TestCheckResourceAttr(resourceName, "parameter.3512697936.name", "max_heap_table_size"),
					resource.TestCheckResourceAttr(resourceName, "parameter.3512697936.value", "67108864"),
					resource.TestCheckResourceAttr(resourceName, "parameter.780730667.name", "performance_schema"),
					resource.TestCheckResourceAttr(resourceName, "parameter.780730667.value", "1"),
					resource.TestCheckResourceAttr(resourceName, "parameter.2020346918.name", "performance_schema_users_size"),
					resource.TestCheckResourceAttr(resourceName, "parameter.2020346918.value", "1048576"),
					resource.TestCheckResourceAttr(resourceName, "parameter.1460834103.name", "query_cache_limit"),
					resource.TestCheckResourceAttr(resourceName, "parameter.1460834103.value", "2097152"),
					resource.TestCheckResourceAttr(resourceName, "parameter.484865451.name", "query_cache_size"),
					resource.TestCheckResourceAttr(resourceName, "parameter.484865451.value", "67108864"),
					resource.TestCheckResourceAttr(resourceName, "parameter.255276438.name", "sort_buffer_size"),
					resource.TestCheckResourceAttr(resourceName, "parameter.255276438.value", "16777216"),
					resource.TestCheckResourceAttr(resourceName, "parameter.2981725119.name", "table_open_cache"),
					resource.TestCheckResourceAttr(resourceName, "parameter.2981725119.value", "4096"),
					resource.TestCheckResourceAttr(resourceName, "parameter.2703661820.name", "tmp_table_size"),
					resource.TestCheckResourceAttr(resourceName, "parameter.2703661820.value", "67108864"),
					resource.TestCheckResourceAttr(resourceName, "parameter.2386583229.name", "binlog_cache_size"),
					resource.TestCheckResourceAttr(resourceName, "parameter.2386583229.value", "131072"),
					resource.TestCheckResourceAttr(resourceName, "parameter.4012389720.name", "innodb_flush_log_at_trx_commit"),
					resource.TestCheckResourceAttr(resourceName, "parameter.4012389720.value", "0"),
					resource.TestCheckResourceAttr(resourceName, "parameter.2688783017.name", "innodb_open_files"),
					resource.TestCheckResourceAttr(resourceName, "parameter.2688783017.value", "4000"),
					resource.TestCheckResourceAttr(resourceName, "parameter.782983977.name", "innodb_read_io_threads"),
					resource.TestCheckResourceAttr(resourceName, "parameter.782983977.value", "64"),
					resource.TestCheckResourceAttr(resourceName, "parameter.2809980413.name", "innodb_thread_concurrency"),
					resource.TestCheckResourceAttr(resourceName, "parameter.2809980413.value", "0"),
					resource.TestCheckResourceAttr(resourceName, "parameter.3599115250.name", "innodb_write_io_threads"),
					resource.TestCheckResourceAttr(resourceName, "parameter.3599115250.value", "64"),
					resource.TestCheckResourceAttr(resourceName, "parameter.2557156277.name", "character_set_connection"),
					resource.TestCheckResourceAttr(resourceName, "parameter.2557156277.value", "utf8"),
					resource.TestCheckResourceAttr(resourceName, "parameter.2475346812.name", "character_set_database"),
					resource.TestCheckResourceAttr(resourceName, "parameter.2475346812.value", "utf8"),
					resource.TestCheckResourceAttr(resourceName, "parameter.1986528518.name", "character_set_filesystem"),
					resource.TestCheckResourceAttr(resourceName, "parameter.1986528518.value", "utf8"),
					resource.TestCheckResourceAttr(resourceName, "parameter.1708034931.name", "character_set_results"),
					resource.TestCheckResourceAttr(resourceName, "parameter.1708034931.value", "utf8"),
					resource.TestCheckResourceAttr(resourceName, "parameter.1937131004.name", "event_scheduler"),
					resource.TestCheckResourceAttr(resourceName, "parameter.1937131004.value", "ON"),
					resource.TestCheckResourceAttr(resourceName, "parameter.3437079877.name", "innodb_buffer_pool_dump_at_shutdown"),
					resource.TestCheckResourceAttr(resourceName, "parameter.3437079877.value", "1"),
					resource.TestCheckResourceAttr(resourceName, "parameter.1092112861.name", "innodb_file_format"),
					resource.TestCheckResourceAttr(resourceName, "parameter.1092112861.value", "barracuda"),
					resource.TestCheckResourceAttr(resourceName, "parameter.615571931.name", "innodb_io_capacity"),
					resource.TestCheckResourceAttr(resourceName, "parameter.615571931.value", "2000"),
					resource.TestCheckResourceAttr(resourceName, "parameter.1065962799.name", "innodb_io_capacity_max"),
					resource.TestCheckResourceAttr(resourceName, "parameter.1065962799.value", "3000"),
					resource.TestCheckResourceAttr(resourceName, "parameter.1411161182.name", "innodb_lock_wait_timeout"),
					resource.TestCheckResourceAttr(resourceName, "parameter.1411161182.value", "120"),
					resource.TestCheckResourceAttr(resourceName, "parameter.3133315879.name", "innodb_max_dirty_pages_pct"),
					resource.TestCheckResourceAttr(resourceName, "parameter.3133315879.value", "90"),
					resource.TestCheckResourceAttr(resourceName, "parameter.950177639.name", "log_bin_trust_function_creators"),
					resource.TestCheckResourceAttr(resourceName, "parameter.950177639.value", "1"),
					resource.TestCheckResourceAttr(resourceName, "parameter.591700516.name", "log_warnings"),
					resource.TestCheckResourceAttr(resourceName, "parameter.591700516.value", "2"),
					resource.TestCheckResourceAttr(resourceName, "parameter.1918306725.name", "log_output"),
					resource.TestCheckResourceAttr(resourceName, "parameter.1918306725.value", "FILE"),
					resource.TestCheckResourceAttr(resourceName, "parameter.386204433.name", "max_allowed_packet"),
					resource.TestCheckResourceAttr(resourceName, "parameter.386204433.value", "1073741824"),
					resource.TestCheckResourceAttr(resourceName, "parameter.1700901269.name", "max_connect_errors"),
					resource.TestCheckResourceAttr(resourceName, "parameter.1700901269.value", "100"),
					resource.TestCheckResourceAttr(resourceName, "parameter.2839701698.name", "query_cache_min_res_unit"),
					resource.TestCheckResourceAttr(resourceName, "parameter.2839701698.value", "512"),
					resource.TestCheckResourceAttr(resourceName, "parameter.427634017.name", "slow_query_log"),
					resource.TestCheckResourceAttr(resourceName, "parameter.427634017.value", "1"),
					resource.TestCheckResourceAttr(resourceName, "parameter.881816039.name", "sync_binlog"),
					resource.TestCheckResourceAttr(resourceName, "parameter.881816039.value", "0"),
					resource.TestCheckResourceAttr(resourceName, "parameter.748684209.name", "tx_isolation"),
					resource.TestCheckResourceAttr(resourceName, "parameter.748684209.value", "REPEATABLE-READ"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: updateAwsDbParameterGroupsExceedDefaultAwsLimit(groupName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSDBParameterGroupExists(resourceName, &v),
					testAccCheckAWSDBParameterGroupAttributes(&v, groupName),
					resource.TestCheckResourceAttr(resourceName, "name", groupName),
					resource.TestCheckResourceAttr(resourceName, "family", "mysql5.6"),
					resource.TestCheckResourceAttr(resourceName, "description", "Updated RDS default parameter group: Exceed default AWS parameter group limit of twenty"),
					resource.TestCheckResourceAttr(resourceName, "parameter.2421266705.name", "character_set_server"),
					resource.TestCheckResourceAttr(resourceName, "parameter.2421266705.value", "utf8"),
					resource.TestCheckResourceAttr(resourceName, "parameter.2478663599.name", "character_set_client"),
					resource.TestCheckResourceAttr(resourceName, "parameter.2478663599.value", "utf8"),
					resource.TestCheckResourceAttr(resourceName, "parameter.1680942586.name", "collation_server"),
					resource.TestCheckResourceAttr(resourceName, "parameter.1680942586.value", "utf8_general_ci"),
					resource.TestCheckResourceAttr(resourceName, "parameter.2450940716.name", "collation_connection"),
					resource.TestCheckResourceAttr(resourceName, "parameter.2450940716.value", "utf8_general_ci"),
					resource.TestCheckResourceAttr(resourceName, "parameter.242489837.name", "join_buffer_size"),
					resource.TestCheckResourceAttr(resourceName, "parameter.242489837.value", "16777216"),
					resource.TestCheckResourceAttr(resourceName, "parameter.2026669454.name", "key_buffer_size"),
					resource.TestCheckResourceAttr(resourceName, "parameter.2026669454.value", "67108864"),
					resource.TestCheckResourceAttr(resourceName, "parameter.2705275319.name", "max_connections"),
					resource.TestCheckResourceAttr(resourceName, "parameter.2705275319.value", "3200"),
					resource.TestCheckResourceAttr(resourceName, "parameter.3512697936.name", "max_heap_table_size"),
					resource.TestCheckResourceAttr(resourceName, "parameter.3512697936.value", "67108864"),
					resource.TestCheckResourceAttr(resourceName, "parameter.780730667.name", "performance_schema"),
					resource.TestCheckResourceAttr(resourceName, "parameter.780730667.value", "1"),
					resource.TestCheckResourceAttr(resourceName, "parameter.2020346918.name", "performance_schema_users_size"),
					resource.TestCheckResourceAttr(resourceName, "parameter.2020346918.value", "1048576"),
					resource.TestCheckResourceAttr(resourceName, "parameter.1460834103.name", "query_cache_limit"),
					resource.TestCheckResourceAttr(resourceName, "parameter.1460834103.value", "2097152"),
					resource.TestCheckResourceAttr(resourceName, "parameter.484865451.name", "query_cache_size"),
					resource.TestCheckResourceAttr(resourceName, "parameter.484865451.value", "67108864"),
					resource.TestCheckResourceAttr(resourceName, "parameter.255276438.name", "sort_buffer_size"),
					resource.TestCheckResourceAttr(resourceName, "parameter.255276438.value", "16777216"),
					resource.TestCheckResourceAttr(resourceName, "parameter.2981725119.name", "table_open_cache"),
					resource.TestCheckResourceAttr(resourceName, "parameter.2981725119.value", "4096"),
					resource.TestCheckResourceAttr(resourceName, "parameter.2703661820.name", "tmp_table_size"),
					resource.TestCheckResourceAttr(resourceName, "parameter.2703661820.value", "67108864"),
					resource.TestCheckResourceAttr(resourceName, "parameter.2386583229.name", "binlog_cache_size"),
					resource.TestCheckResourceAttr(resourceName, "parameter.2386583229.value", "131072"),
					resource.TestCheckResourceAttr(resourceName, "parameter.4012389720.name", "innodb_flush_log_at_trx_commit"),
					resource.TestCheckResourceAttr(resourceName, "parameter.4012389720.value", "0"),
					resource.TestCheckResourceAttr(resourceName, "parameter.2688783017.name", "innodb_open_files"),
					resource.TestCheckResourceAttr(resourceName, "parameter.2688783017.value", "4000"),
					resource.TestCheckResourceAttr(resourceName, "parameter.782983977.name", "innodb_read_io_threads"),
					resource.TestCheckResourceAttr(resourceName, "parameter.782983977.value", "64"),
					resource.TestCheckResourceAttr(resourceName, "parameter.2809980413.name", "innodb_thread_concurrency"),
					resource.TestCheckResourceAttr(resourceName, "parameter.2809980413.value", "0"),
					resource.TestCheckResourceAttr(resourceName, "parameter.3599115250.name", "innodb_write_io_threads"),
					resource.TestCheckResourceAttr(resourceName, "parameter.3599115250.value", "64"),
					resource.TestCheckResourceAttr(resourceName, "parameter.2557156277.name", "character_set_connection"),
					resource.TestCheckResourceAttr(resourceName, "parameter.2557156277.value", "utf8"),
					resource.TestCheckResourceAttr(resourceName, "parameter.2475346812.name", "character_set_database"),
					resource.TestCheckResourceAttr(resourceName, "parameter.2475346812.value", "utf8"),
					resource.TestCheckResourceAttr(resourceName, "parameter.1986528518.name", "character_set_filesystem"),
					resource.TestCheckResourceAttr(resourceName, "parameter.1986528518.value", "utf8"),
					resource.TestCheckResourceAttr(resourceName, "parameter.1708034931.name", "character_set_results"),
					resource.TestCheckResourceAttr(resourceName, "parameter.1708034931.value", "utf8"),
					resource.TestCheckResourceAttr(resourceName, "parameter.1937131004.name", "event_scheduler"),
					resource.TestCheckResourceAttr(resourceName, "parameter.1937131004.value", "ON"),
					resource.TestCheckResourceAttr(resourceName, "parameter.3437079877.name", "innodb_buffer_pool_dump_at_shutdown"),
					resource.TestCheckResourceAttr(resourceName, "parameter.3437079877.value", "1"),
					resource.TestCheckResourceAttr(resourceName, "parameter.1092112861.name", "innodb_file_format"),
					resource.TestCheckResourceAttr(resourceName, "parameter.1092112861.value", "barracuda"),
					resource.TestCheckResourceAttr(resourceName, "parameter.615571931.name", "innodb_io_capacity"),
					resource.TestCheckResourceAttr(resourceName, "parameter.615571931.value", "2000"),
					resource.TestCheckResourceAttr(resourceName, "parameter.1065962799.name", "innodb_io_capacity_max"),
					resource.TestCheckResourceAttr(resourceName, "parameter.1065962799.value", "3000"),
					resource.TestCheckResourceAttr(resourceName, "parameter.1411161182.name", "innodb_lock_wait_timeout"),
					resource.TestCheckResourceAttr(resourceName, "parameter.1411161182.value", "120"),
					resource.TestCheckResourceAttr(resourceName, "parameter.3133315879.name", "innodb_max_dirty_pages_pct"),
					resource.TestCheckResourceAttr(resourceName, "parameter.3133315879.value", "90"),
					resource.TestCheckResourceAttr(resourceName, "parameter.950177639.name", "log_bin_trust_function_creators"),
					resource.TestCheckResourceAttr(resourceName, "parameter.950177639.value", "1"),
					resource.TestCheckResourceAttr(resourceName, "parameter.591700516.name", "log_warnings"),
					resource.TestCheckResourceAttr(resourceName, "parameter.591700516.value", "2"),
					resource.TestCheckResourceAttr(resourceName, "parameter.1918306725.name", "log_output"),
					resource.TestCheckResourceAttr(resourceName, "parameter.1918306725.value", "FILE"),
					resource.TestCheckResourceAttr(resourceName, "parameter.386204433.name", "max_allowed_packet"),
					resource.TestCheckResourceAttr(resourceName, "parameter.386204433.value", "1073741824"),
					resource.TestCheckResourceAttr(resourceName, "parameter.1700901269.name", "max_connect_errors"),
					resource.TestCheckResourceAttr(resourceName, "parameter.1700901269.value", "100"),
					resource.TestCheckResourceAttr(resourceName, "parameter.2839701698.name", "query_cache_min_res_unit"),
					resource.TestCheckResourceAttr(resourceName, "parameter.2839701698.value", "512"),
					resource.TestCheckResourceAttr(resourceName, "parameter.427634017.name", "slow_query_log"),
					resource.TestCheckResourceAttr(resourceName, "parameter.427634017.value", "1"),
					resource.TestCheckResourceAttr(resourceName, "parameter.881816039.name", "sync_binlog"),
					resource.TestCheckResourceAttr(resourceName, "parameter.881816039.value", "0"),
					resource.TestCheckResourceAttr(resourceName, "parameter.748684209.name", "tx_isolation"),
					resource.TestCheckResourceAttr(resourceName, "parameter.748684209.value", "REPEATABLE-READ"),
				),
			},
		},
	})
}

func TestAccAWSDBParameterGroup_Disappears(t *testing.T) {
	var v rds.DBParameterGroup
	resourceName := "aws_db_parameter_group.test"
	groupName := fmt.Sprintf("parameter-group-test-terraform-%d", acctest.RandInt())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAWSDBParameterGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSDBParameterGroupConfig(groupName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSDBParameterGroupExists(resourceName, &v),
					testAccCheckAWSDbParamaterGroupDisappears(&v),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccAWSDBParameterGroup_namePrefix(t *testing.T) {
	var v rds.DBParameterGroup

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAWSDBParameterGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDBParameterGroupConfig_namePrefix,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSDBParameterGroupExists("aws_db_parameter_group.test", &v),
					resource.TestMatchResourceAttr("aws_db_parameter_group.test", "name", regexp.MustCompile("^tf-test-")),
				),
			},
		},
	})
}

func TestAccAWSDBParameterGroup_generatedName(t *testing.T) {
	var v rds.DBParameterGroup

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAWSDBParameterGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDBParameterGroupConfig_generatedName,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSDBParameterGroupExists("aws_db_parameter_group.test", &v),
				),
			},
		},
	})
}

func TestAccAWSDBParameterGroup_withApplyMethod(t *testing.T) {
	var v rds.DBParameterGroup
	resourceName := "aws_db_parameter_group.test"
	groupName := fmt.Sprintf("parameter-group-test-terraform-%d", acctest.RandInt())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:            func() { testAccPreCheck(t) },
		Providers:           testAccProviders,
		CheckDestroy:        testAccCheckAWSDBParameterGroupDestroy,
		DisableBinaryDriver: true,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSDBParameterGroupConfigWithApplyMethod(groupName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSDBParameterGroupExists(resourceName, &v),
					testAccCheckAWSDBParameterGroupAttributes(&v, groupName),
					resource.TestCheckResourceAttr(resourceName, "name", groupName),
					resource.TestCheckResourceAttr(resourceName, "family", "mysql5.6"),
					resource.TestCheckResourceAttr(resourceName, "description", "Managed by Terraform"),
					resource.TestCheckResourceAttr(resourceName, "parameter.2421266705.name", "character_set_server"),
					resource.TestCheckResourceAttr(resourceName, "parameter.2421266705.value", "utf8"),
					resource.TestCheckResourceAttr(resourceName, "parameter.2421266705.apply_method", "immediate"),
					resource.TestCheckResourceAttr(resourceName, "parameter.2478663599.name", "character_set_client"),
					resource.TestCheckResourceAttr(resourceName, "parameter.2478663599.value", "utf8"),
					resource.TestCheckResourceAttr(resourceName, "parameter.2478663599.apply_method", "pending-reboot"),
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

func TestAccAWSDBParameterGroup_Only(t *testing.T) {
	var v rds.DBParameterGroup
	resourceName := "aws_db_parameter_group.test"
	groupName := fmt.Sprintf("parameter-group-test-terraform-%d", acctest.RandInt())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAWSDBParameterGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSDBParameterGroupOnlyConfig(groupName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSDBParameterGroupExists(resourceName, &v),
					testAccCheckAWSDBParameterGroupAttributes(&v, groupName),
					resource.TestCheckResourceAttr(resourceName, "name", groupName),
					resource.TestCheckResourceAttr(resourceName, "family", "mysql5.6"),
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

func TestAccAWSDBParameterGroup_MatchDefault(t *testing.T) {
	var v rds.DBParameterGroup
	resourceName := "aws_db_parameter_group.test"
	groupName := fmt.Sprintf("parameter-group-test-terraform-%d", acctest.RandInt())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAWSDBParameterGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSDBParameterGroupIncludeDefaultConfig(groupName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSDBParameterGroupExists(resourceName, &v),
					resource.TestCheckResourceAttr(resourceName, "name", groupName),
					resource.TestCheckResourceAttr(resourceName, "family", "postgres9.4"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"parameter"},
			},
		},
	})
}

func testAccCheckAWSDbParamaterGroupDisappears(v *rds.DBParameterGroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := testAccProvider.Meta().(*AWSClient).rdsconn

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "aws_db_parameter_group" {
				continue
			}
			_, err := conn.DeleteDBParameterGroup(&rds.DeleteDBParameterGroupInput{
				DBParameterGroupName: v.DBParameterGroupName,
			})
			return err
		}
		return nil
	}
}

func testAccCheckAWSDBParameterGroupDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*AWSClient).rdsconn

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aws_db_parameter_group" {
			continue
		}

		// Try to find the Group
		resp, err := conn.DescribeDBParameterGroups(
			&rds.DescribeDBParameterGroupsInput{
				DBParameterGroupName: aws.String(rs.Primary.ID),
			})

		if err == nil {
			if len(resp.DBParameterGroups) != 0 &&
				*resp.DBParameterGroups[0].DBParameterGroupName == rs.Primary.ID {
				return fmt.Errorf("DB Parameter Group still exists")
			}
		}

		// Verify the error
		newerr, ok := err.(awserr.Error)
		if !ok {
			return err
		}
		if newerr.Code() != "DBParameterGroupNotFound" {
			return err
		}
	}

	return nil
}

func testAccCheckAWSDBParameterGroupAttributes(v *rds.DBParameterGroup, name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if *v.DBParameterGroupName != name {
			return fmt.Errorf("Bad Parameter Group name, expected (%s), got (%s)", name, *v.DBParameterGroupName)
		}

		if *v.DBParameterGroupFamily != "mysql5.6" {
			return fmt.Errorf("bad family: %#v", v.DBParameterGroupFamily)
		}

		return nil
	}
}

func testAccCheckAWSDBParameterGroupExists(n string, v *rds.DBParameterGroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No DB Parameter Group ID is set")
		}

		conn := testAccProvider.Meta().(*AWSClient).rdsconn

		opts := rds.DescribeDBParameterGroupsInput{
			DBParameterGroupName: aws.String(rs.Primary.ID),
		}

		resp, err := conn.DescribeDBParameterGroups(&opts)

		if err != nil {
			return err
		}

		if len(resp.DBParameterGroups) != 1 ||
			*resp.DBParameterGroups[0].DBParameterGroupName != rs.Primary.ID {
			return fmt.Errorf("DB Parameter Group not found")
		}

		*v = *resp.DBParameterGroups[0]

		return nil
	}
}

func testAccCheckAWSDBParameterNotUserDefined(n, paramName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No DB Parameter Group ID is set")
		}

		conn := testAccProvider.Meta().(*AWSClient).rdsconn

		opts := rds.DescribeDBParametersInput{
			DBParameterGroupName: aws.String(rs.Primary.ID),
			Source:               aws.String("user"),
		}

		userDefined := false
		err := conn.DescribeDBParametersPages(&opts, func(page *rds.DescribeDBParametersOutput, lastPage bool) bool {
			for _, param := range page.Parameters {
				if *param.ParameterName == paramName {
					userDefined = true
					return false
				}
			}
			return true
		})

		if userDefined {
			return fmt.Errorf("DB Parameter is user defined")
		}

		return err
	}
}

func testAccAWSDBParameterGroupConfig(n string) string {
	return fmt.Sprintf(`
resource "aws_db_parameter_group" "test" {
  name   = "%s"
  family = "mysql5.6"

  parameter {
    name  = "character_set_server"
    value = "utf8"
  }

  parameter {
    name  = "character_set_client"
    value = "utf8"
  }

  parameter {
    name  = "character_set_results"
    value = "utf8"
  }
}
`, n)
}

func testAccAWSDBParameterGroupConfigWithApplyMethod(n string) string {
	return fmt.Sprintf(`
resource "aws_db_parameter_group" "test" {
  name   = "%s"
  family = "mysql5.6"

  parameter {
    name  = "character_set_server"
    value = "utf8"
  }

  parameter {
    name         = "character_set_client"
    value        = "utf8"
    apply_method = "pending-reboot"
  }

  tags = {
    foo = "test"
  }
}
`, n)
}

func testAccAWSDBParameterGroupAddParametersConfig(n string) string {
	return fmt.Sprintf(`
resource "aws_db_parameter_group" "test" {
  name        = "%s"
  family      = "mysql5.6"

  parameter {
    name  = "character_set_server"
    value = "utf8"
  }

  parameter {
    name  = "character_set_client"
    value = "utf8"
  }

  parameter {
    name  = "character_set_results"
    value = "utf8"
  }

  parameter {
    name  = "collation_server"
    value = "utf8_unicode_ci"
  }

  parameter {
    name  = "collation_connection"
    value = "utf8_unicode_ci"
  }
}
`, n)
}

func testAccAWSDBParameterGroupOnlyConfig(n string) string {
	return fmt.Sprintf(`
resource "aws_db_parameter_group" "test" {
  name        = "%s"
  family      = "mysql5.6"
  description = "Test parameter group for terraform"
}
`, n)
}

func createAwsDbParameterGroupsExceedDefaultAwsLimit(n string) string {
	return fmt.Sprintf(`
resource "aws_db_parameter_group" "test" {
  name        = "%s"
  family      = "mysql5.6"
  description = "RDS default parameter group: Exceed default AWS parameter group limit of twenty"

  parameter {
    name  = "binlog_cache_size"
    value = 131072
  }

  parameter {
    name  = "character_set_client"
    value = "utf8"
  }

  parameter {
    name  = "character_set_connection"
    value = "utf8"
  }

  parameter {
    name  = "character_set_database"
    value = "utf8"
  }

  parameter {
    name  = "character_set_filesystem"
    value = "utf8"
  }

  parameter {
    name  = "character_set_results"
    value = "utf8"
  }

  parameter {
    name  = "character_set_server"
    value = "utf8"
  }

  parameter {
    name  = "collation_connection"
    value = "utf8_general_ci"
  }

  parameter {
    name  = "collation_server"
    value = "utf8_general_ci"
  }

  parameter {
    name  = "event_scheduler"
    value = "ON"
  }

  parameter {
    name  = "innodb_buffer_pool_dump_at_shutdown"
    value = 1
  }

  parameter {
    name  = "innodb_file_format"
    value = "barracuda"
  }

  parameter {
    name  = "innodb_flush_log_at_trx_commit"
    value = 0
  }

  parameter {
    name  = "innodb_io_capacity"
    value = 2000
  }

  parameter {
    name  = "innodb_io_capacity_max"
    value = 3000
  }

  parameter {
    name  = "innodb_lock_wait_timeout"
    value = 120
  }

  parameter {
    name  = "innodb_max_dirty_pages_pct"
    value = 90
  }

  parameter {
    name         = "innodb_open_files"
    value        = 4000
    apply_method = "pending-reboot"
  }

  parameter {
    name         = "innodb_read_io_threads"
    value        = 64
    apply_method = "pending-reboot"
  }

  parameter {
    name  = "innodb_thread_concurrency"
    value = 0
  }

  parameter {
    name         = "innodb_write_io_threads"
    value        = 64
    apply_method = "pending-reboot"
  }

  parameter {
    name  = "join_buffer_size"
    value = 16777216
  }

  parameter {
    name  = "key_buffer_size"
    value = 67108864
  }

  parameter {
    name  = "log_bin_trust_function_creators"
    value = 1
  }

  parameter {
    name  = "log_warnings"
    value = 2
  }

  parameter {
    name  = "log_output"
    value = "FILE"
  }

  parameter {
    name  = "max_allowed_packet"
    value = 1073741824
  }

  parameter {
    name  = "max_connect_errors"
    value = 100
  }

  parameter {
    name  = "max_connections"
    value = 3200
  }

  parameter {
    name  = "max_heap_table_size"
    value = 67108864
  }

  parameter {
    name         = "performance_schema"
    value        = 1
    apply_method = "pending-reboot"
  }

  parameter {
    name         = "performance_schema_users_size"
    value        = 1048576
    apply_method = "pending-reboot"
  }

  parameter {
    name  = "query_cache_limit"
    value = 2097152
  }

  parameter {
    name  = "query_cache_min_res_unit"
    value = 512
  }

  parameter {
    name  = "query_cache_size"
    value = 67108864
  }

  parameter {
    name  = "slow_query_log"
    value = 1
  }

  parameter {
    name  = "sort_buffer_size"
    value = 16777216
  }

  parameter {
    name  = "sync_binlog"
    value = 0
  }

  parameter {
    name  = "table_open_cache"
    value = 4096
  }

  parameter {
    name  = "tmp_table_size"
    value = 67108864
  }

  parameter {
    name  = "tx_isolation"
    value = "REPEATABLE-READ"
  }
}
`, n)
}

func updateAwsDbParameterGroupsExceedDefaultAwsLimit(n string) string {
	return fmt.Sprintf(`
resource "aws_db_parameter_group" "test" {
  name        = "%s"
  family      = "mysql5.6"
  description = "Updated RDS default parameter group: Exceed default AWS parameter group limit of twenty"

  parameter {
    name  = "binlog_cache_size"
    value = 131072
  }

  parameter {
    name  = "character_set_client"
    value = "utf8"
  }

  parameter {
    name  = "character_set_connection"
    value = "utf8"
  }

  parameter {
    name  = "character_set_database"
    value = "utf8"
  }

  parameter {
    name  = "character_set_filesystem"
    value = "utf8"
  }

  parameter {
    name  = "character_set_results"
    value = "utf8"
  }

  parameter {
    name  = "character_set_server"
    value = "utf8"
  }

  parameter {
    name  = "collation_connection"
    value = "utf8_general_ci"
  }

  parameter {
    name  = "collation_server"
    value = "utf8_general_ci"
  }

  parameter {
    name  = "event_scheduler"
    value = "ON"
  }

  parameter {
    name  = "innodb_buffer_pool_dump_at_shutdown"
    value = 1
  }

  parameter {
    name  = "innodb_file_format"
    value = "barracuda"
  }

  parameter {
    name  = "innodb_flush_log_at_trx_commit"
    value = 0
  }

  parameter {
    name  = "innodb_io_capacity"
    value = 2000
  }

  parameter {
    name  = "innodb_io_capacity_max"
    value = 3000
  }

  parameter {
    name  = "innodb_lock_wait_timeout"
    value = 120
  }

  parameter {
    name  = "innodb_max_dirty_pages_pct"
    value = 90
  }

  parameter {
    name         = "innodb_open_files"
    value        = 4000
    apply_method = "pending-reboot"
  }

  parameter {
    name         = "innodb_read_io_threads"
    value        = 64
    apply_method = "pending-reboot"
  }

  parameter {
    name  = "innodb_thread_concurrency"
    value = 0
  }

  parameter {
    name         = "innodb_write_io_threads"
    value        = 64
    apply_method = "pending-reboot"
  }

  parameter {
    name  = "join_buffer_size"
    value = 16777216
  }

  parameter {
    name  = "key_buffer_size"
    value = 67108864
  }

  parameter {
    name  = "log_bin_trust_function_creators"
    value = 1
  }

  parameter {
    name  = "log_warnings"
    value = 2
  }

  parameter {
    name  = "log_output"
    value = "FILE"
  }

  parameter {
    name  = "max_allowed_packet"
    value = 1073741824
  }

  parameter {
    name  = "max_connect_errors"
    value = 100
  }

  parameter {
    name  = "max_connections"
    value = 3200
  }

  parameter {
    name  = "max_heap_table_size"
    value = 67108864
  }

  parameter {
    name         = "performance_schema"
    value        = 1
    apply_method = "pending-reboot"
  }

  parameter {
    name         = "performance_schema_users_size"
    value        = 1048576
    apply_method = "pending-reboot"
  }

  parameter {
    name  = "query_cache_limit"
    value = 2097152
  }

  parameter {
    name  = "query_cache_min_res_unit"
    value = 512
  }

  parameter {
    name  = "query_cache_size"
    value = 67108864
  }

  parameter {
    name  = "slow_query_log"
    value = 1
  }

  parameter {
    name  = "sort_buffer_size"
    value = 16777216
  }

  parameter {
    name  = "sync_binlog"
    value = 0
  }

  parameter {
    name  = "table_open_cache"
    value = 4096
  }

  parameter {
    name  = "tmp_table_size"
    value = 67108864
  }

  parameter {
    name  = "tx_isolation"
    value = "REPEATABLE-READ"
  }
}
`, n)
}

func testAccAWSDBParameterGroupIncludeDefaultConfig(n string) string {
	return fmt.Sprintf(`
resource "aws_db_parameter_group" "test" {
  name   = "%s"
  family = "postgres9.4"

  parameter {
    name         = "client_encoding"
    value        = "UTF8"
    apply_method = "pending-reboot"
  }
}
`, n)
}

const testAccDBParameterGroupConfig_namePrefix = `
resource "aws_db_parameter_group" "test" {
	name_prefix = "tf-test-"
	family = "mysql5.6"

	parameter {
		name = "sync_binlog"
		value = 0
	}
}
`

const testAccDBParameterGroupConfig_generatedName = `
resource "aws_db_parameter_group" "test" {
	family = "mysql5.6"

	parameter {
		name = "sync_binlog"
		value = 0
	}
}
`
