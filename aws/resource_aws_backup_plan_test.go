package aws

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/backup"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccAwsBackupPlan_basic(t *testing.T) {
	var plan backup.GetBackupPlanOutput
	ruleNameMap := map[string]string{}
	resourceName := "aws_backup_plan.test"
	rName := fmt.Sprintf("tf-testacc-backup-%s", acctest.RandStringFromCharSet(14, acctest.CharSetAlphaNum))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t); testAccPreCheckAWSBackup(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAwsBackupPlanDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAwsBackupPlanConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAwsBackupPlanExists(resourceName, &plan, &ruleNameMap),
					testAccMatchResourceAttrRegionalARN(resourceName, "arn", "backup", regexp.MustCompile(`backup-plan:.+`)),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "rule.#", "1"),
					testAccCheckAwsBackupPlanRuleAttr(resourceName, &ruleNameMap, rName, "rule_name", rName),
					testAccCheckAwsBackupPlanRuleAttr(resourceName, &ruleNameMap, rName, "target_vault_name", rName),
					testAccCheckAwsBackupPlanRuleAttr(resourceName, &ruleNameMap, rName, "schedule", "cron(0 12 * * ? *)"),
					testAccCheckAwsBackupPlanRuleAttr(resourceName, &ruleNameMap, rName, "lifecycle.#", "0"),
					testAccCheckAwsBackupPlanRuleAttr(resourceName, &ruleNameMap, rName, "recovery_point_tags.%", "0"),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "0"),
					resource.TestCheckResourceAttrSet(resourceName, "version"),
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

func TestAccAwsBackupPlan_withTags(t *testing.T) {
	var plan backup.GetBackupPlanOutput
	ruleNameMap := map[string]string{}
	resourceName := "aws_backup_plan.test"
	rName := fmt.Sprintf("tf-testacc-backup-%s", acctest.RandStringFromCharSet(14, acctest.CharSetAlphaNum))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t); testAccPreCheckAWSBackup(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAwsBackupPlanDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAwsBackupPlanConfig_tags(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAwsBackupPlanExists(resourceName, &plan, &ruleNameMap),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "rule.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "3"),
					resource.TestCheckResourceAttr(resourceName, "tags.Name", rName),
					resource.TestCheckResourceAttr(resourceName, "tags.Key1", "Value1"),
					resource.TestCheckResourceAttr(resourceName, "tags.Key2", "Value2a"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccAwsBackupPlanConfig_tagsUpdated(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAwsBackupPlanExists(resourceName, &plan, &ruleNameMap),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "rule.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "3"),
					resource.TestCheckResourceAttr(resourceName, "tags.Name", rName),
					resource.TestCheckResourceAttr(resourceName, "tags.Key2", "Value2b"),
					resource.TestCheckResourceAttr(resourceName, "tags.Key3", "Value3"),
				),
			},
			{
				Config: testAccAwsBackupPlanConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAwsBackupPlanExists(resourceName, &plan, &ruleNameMap),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "rule.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "0"),
				),
			},
		},
	})
}

func TestAccAwsBackupPlan_withRules(t *testing.T) {
	var plan backup.GetBackupPlanOutput
	ruleNameMap := map[string]string{}
	resourceName := "aws_backup_plan.test"
	rName := fmt.Sprintf("tf-testacc-backup-%s", acctest.RandStringFromCharSet(14, acctest.CharSetAlphaNum))
	rule1Name := fmt.Sprintf("%s_1", rName)
	rule2Name := fmt.Sprintf("%s_2", rName)
	rule3Name := fmt.Sprintf("%s_3", rName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t); testAccPreCheckAWSBackup(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAwsBackupPlanDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAwsBackupPlanConfig_twoRules(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAwsBackupPlanExists(resourceName, &plan, &ruleNameMap),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "rule.#", "2"),
					testAccCheckAwsBackupPlanRuleAttr(resourceName, &ruleNameMap, rule1Name, "rule_name", rule1Name),
					testAccCheckAwsBackupPlanRuleAttr(resourceName, &ruleNameMap, rule1Name, "target_vault_name", rName),
					testAccCheckAwsBackupPlanRuleAttr(resourceName, &ruleNameMap, rule1Name, "schedule", "cron(0 12 * * ? *)"),
					testAccCheckAwsBackupPlanRuleAttr(resourceName, &ruleNameMap, rule1Name, "lifecycle.#", "0"),
					testAccCheckAwsBackupPlanRuleAttr(resourceName, &ruleNameMap, rule1Name, "recovery_point_tags.%", "0"),
					testAccCheckAwsBackupPlanRuleAttr(resourceName, &ruleNameMap, rule2Name, "rule_name", rule2Name),
					testAccCheckAwsBackupPlanRuleAttr(resourceName, &ruleNameMap, rule2Name, "target_vault_name", rName),
					testAccCheckAwsBackupPlanRuleAttr(resourceName, &ruleNameMap, rule2Name, "schedule", "cron(0 6 * * ? *)"),
					testAccCheckAwsBackupPlanRuleAttr(resourceName, &ruleNameMap, rule2Name, "lifecycle.#", "0"),
					testAccCheckAwsBackupPlanRuleAttr(resourceName, &ruleNameMap, rule2Name, "recovery_point_tags.%", "0"),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "0"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccAwsBackupPlanConfig_threeRules(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAwsBackupPlanExists(resourceName, &plan, &ruleNameMap),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "rule.#", "3"),
					testAccCheckAwsBackupPlanRuleAttr(resourceName, &ruleNameMap, rule1Name, "rule_name", rule1Name),
					testAccCheckAwsBackupPlanRuleAttr(resourceName, &ruleNameMap, rule1Name, "target_vault_name", rName),
					testAccCheckAwsBackupPlanRuleAttr(resourceName, &ruleNameMap, rule1Name, "schedule", "cron(0 6 * * ? *)"),
					testAccCheckAwsBackupPlanRuleAttr(resourceName, &ruleNameMap, rule1Name, "lifecycle.#", "0"),
					testAccCheckAwsBackupPlanRuleAttr(resourceName, &ruleNameMap, rule1Name, "recovery_point_tags.%", "0"),
					testAccCheckAwsBackupPlanRuleAttr(resourceName, &ruleNameMap, rule2Name, "rule_name", rule2Name),
					testAccCheckAwsBackupPlanRuleAttr(resourceName, &ruleNameMap, rule2Name, "target_vault_name", rName),
					testAccCheckAwsBackupPlanRuleAttr(resourceName, &ruleNameMap, rule2Name, "schedule", "cron(0 12 * * ? *)"),
					testAccCheckAwsBackupPlanRuleAttr(resourceName, &ruleNameMap, rule2Name, "lifecycle.#", "0"),
					testAccCheckAwsBackupPlanRuleAttr(resourceName, &ruleNameMap, rule2Name, "recovery_point_tags.%", "0"),
					testAccCheckAwsBackupPlanRuleAttr(resourceName, &ruleNameMap, rule3Name, "rule_name", rule3Name),
					testAccCheckAwsBackupPlanRuleAttr(resourceName, &ruleNameMap, rule3Name, "target_vault_name", rName),
					testAccCheckAwsBackupPlanRuleAttr(resourceName, &ruleNameMap, rule3Name, "schedule", "cron(0 18 * * ? *)"),
					testAccCheckAwsBackupPlanRuleAttr(resourceName, &ruleNameMap, rule3Name, "lifecycle.#", "0"),
					testAccCheckAwsBackupPlanRuleAttr(resourceName, &ruleNameMap, rule3Name, "recovery_point_tags.%", "0"),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "0"),
				),
			},
			{
				Config: testAccAwsBackupPlanConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAwsBackupPlanExists(resourceName, &plan, &ruleNameMap),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "rule.#", "1"),
					testAccCheckAwsBackupPlanRuleAttr(resourceName, &ruleNameMap, rName, "rule_name", rName),
					testAccCheckAwsBackupPlanRuleAttr(resourceName, &ruleNameMap, rName, "target_vault_name", rName),
					testAccCheckAwsBackupPlanRuleAttr(resourceName, &ruleNameMap, rName, "schedule", "cron(0 12 * * ? *)"),
					testAccCheckAwsBackupPlanRuleAttr(resourceName, &ruleNameMap, rName, "lifecycle.#", "0"),
					testAccCheckAwsBackupPlanRuleAttr(resourceName, &ruleNameMap, rName, "recovery_point_tags.%", "0"),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "0"),
				),
			},
		},
	})
}

func TestAccAwsBackupPlan_withLifecycle(t *testing.T) {
	var plan backup.GetBackupPlanOutput
	ruleNameMap := map[string]string{}
	resourceName := "aws_backup_plan.test"
	rName := fmt.Sprintf("tf-testacc-backup-%s", acctest.RandStringFromCharSet(14, acctest.CharSetAlphaNum))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t); testAccPreCheckAWSBackup(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAwsBackupPlanDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAwsBackupPlanConfig_lifecycleColdStorageAfterOnly(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAwsBackupPlanExists(resourceName, &plan, &ruleNameMap),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "rule.#", "1"),
					testAccCheckAwsBackupPlanRuleAttr(resourceName, &ruleNameMap, rName, "rule_name", rName),
					testAccCheckAwsBackupPlanRuleAttr(resourceName, &ruleNameMap, rName, "lifecycle.#", "1"),
					testAccCheckAwsBackupPlanRuleAttr(resourceName, &ruleNameMap, rName, "lifecycle.0.cold_storage_after", "7"),
					testAccCheckAwsBackupPlanRuleAttr(resourceName, &ruleNameMap, rName, "lifecycle.0.delete_after", "0"),
					testAccCheckAwsBackupPlanRuleAttr(resourceName, &ruleNameMap, rName, "recovery_point_tags.%", "0"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccAwsBackupPlanConfig_lifecycleDeleteAfterOnly(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAwsBackupPlanExists(resourceName, &plan, &ruleNameMap),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "rule.#", "1"),
					testAccCheckAwsBackupPlanRuleAttr(resourceName, &ruleNameMap, rName, "rule_name", rName),
					testAccCheckAwsBackupPlanRuleAttr(resourceName, &ruleNameMap, rName, "lifecycle.#", "1"),
					testAccCheckAwsBackupPlanRuleAttr(resourceName, &ruleNameMap, rName, "lifecycle.0.cold_storage_after", "0"),
					testAccCheckAwsBackupPlanRuleAttr(resourceName, &ruleNameMap, rName, "lifecycle.0.delete_after", "120"),
					testAccCheckAwsBackupPlanRuleAttr(resourceName, &ruleNameMap, rName, "recovery_point_tags.%", "0"),
				),
			},
			{
				Config: testAccAwsBackupPlanConfig_lifecycleColdStorageAfterAndDeleteAfter(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAwsBackupPlanExists(resourceName, &plan, &ruleNameMap),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "rule.#", "1"),
					testAccCheckAwsBackupPlanRuleAttr(resourceName, &ruleNameMap, rName, "rule_name", rName),
					testAccCheckAwsBackupPlanRuleAttr(resourceName, &ruleNameMap, rName, "lifecycle.#", "1"),
					testAccCheckAwsBackupPlanRuleAttr(resourceName, &ruleNameMap, rName, "lifecycle.0.cold_storage_after", "30"),
					testAccCheckAwsBackupPlanRuleAttr(resourceName, &ruleNameMap, rName, "lifecycle.0.delete_after", "180"),
					testAccCheckAwsBackupPlanRuleAttr(resourceName, &ruleNameMap, rName, "recovery_point_tags.%", "0"),
				),
			},
			{
				Config: testAccAwsBackupPlanConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAwsBackupPlanExists(resourceName, &plan, &ruleNameMap),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "rule.#", "1"),
					testAccCheckAwsBackupPlanRuleAttr(resourceName, &ruleNameMap, rName, "rule_name", rName),
					testAccCheckAwsBackupPlanRuleAttr(resourceName, &ruleNameMap, rName, "lifecycle.#", "0"),
					testAccCheckAwsBackupPlanRuleAttr(resourceName, &ruleNameMap, rName, "recovery_point_tags.%", "0"),
				),
			},
		},
	})
}

func TestAccAwsBackupPlan_withRecoveryPointTags(t *testing.T) {
	var plan backup.GetBackupPlanOutput
	ruleNameMap := map[string]string{}
	resourceName := "aws_backup_plan.test"
	rName := fmt.Sprintf("tf-testacc-backup-%s", acctest.RandStringFromCharSet(14, acctest.CharSetAlphaNum))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t); testAccPreCheckAWSBackup(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAwsBackupPlanDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAwsBackupPlanConfig_recoveryPointTags(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAwsBackupPlanExists(resourceName, &plan, &ruleNameMap),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "rule.#", "1"),
					testAccCheckAwsBackupPlanRuleAttr(resourceName, &ruleNameMap, rName, "rule_name", rName),
					testAccCheckAwsBackupPlanRuleAttr(resourceName, &ruleNameMap, rName, "target_vault_name", rName),
					testAccCheckAwsBackupPlanRuleAttr(resourceName, &ruleNameMap, rName, "schedule", "cron(0 12 * * ? *)"),
					testAccCheckAwsBackupPlanRuleAttr(resourceName, &ruleNameMap, rName, "lifecycle.#", "0"),
					testAccCheckAwsBackupPlanRuleAttr(resourceName, &ruleNameMap, rName, "recovery_point_tags.%", "3"),
					testAccCheckAwsBackupPlanRuleAttr(resourceName, &ruleNameMap, rName, "recovery_point_tags.Name", rName),
					testAccCheckAwsBackupPlanRuleAttr(resourceName, &ruleNameMap, rName, "recovery_point_tags.Key1", "Value1"),
					testAccCheckAwsBackupPlanRuleAttr(resourceName, &ruleNameMap, rName, "recovery_point_tags.Key2", "Value2a"),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "0"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccAwsBackupPlanConfig_recoveryPointTagsUpdated(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAwsBackupPlanExists(resourceName, &plan, &ruleNameMap),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "rule.#", "1"),
					testAccCheckAwsBackupPlanRuleAttr(resourceName, &ruleNameMap, rName, "rule_name", rName),
					testAccCheckAwsBackupPlanRuleAttr(resourceName, &ruleNameMap, rName, "target_vault_name", rName),
					testAccCheckAwsBackupPlanRuleAttr(resourceName, &ruleNameMap, rName, "schedule", "cron(0 12 * * ? *)"),
					testAccCheckAwsBackupPlanRuleAttr(resourceName, &ruleNameMap, rName, "lifecycle.#", "0"),
					testAccCheckAwsBackupPlanRuleAttr(resourceName, &ruleNameMap, rName, "recovery_point_tags.%", "3"),
					testAccCheckAwsBackupPlanRuleAttr(resourceName, &ruleNameMap, rName, "recovery_point_tags.Name", rName),
					testAccCheckAwsBackupPlanRuleAttr(resourceName, &ruleNameMap, rName, "recovery_point_tags.Key2", "Value2b"),
					testAccCheckAwsBackupPlanRuleAttr(resourceName, &ruleNameMap, rName, "recovery_point_tags.Key3", "Value3"),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "0"),
				),
			},
			{
				Config: testAccAwsBackupPlanConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAwsBackupPlanExists(resourceName, &plan, &ruleNameMap),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "rule.#", "1"),
					testAccCheckAwsBackupPlanRuleAttr(resourceName, &ruleNameMap, rName, "rule_name", rName),
					testAccCheckAwsBackupPlanRuleAttr(resourceName, &ruleNameMap, rName, "target_vault_name", rName),
					testAccCheckAwsBackupPlanRuleAttr(resourceName, &ruleNameMap, rName, "schedule", "cron(0 12 * * ? *)"),
					testAccCheckAwsBackupPlanRuleAttr(resourceName, &ruleNameMap, rName, "lifecycle.#", "0"),
					testAccCheckAwsBackupPlanRuleAttr(resourceName, &ruleNameMap, rName, "recovery_point_tags.%", "0"),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "0"),
				),
			},
		},
	})
}

func TestAccAwsBackupPlan_Rule_CopyAction_SameRegion(t *testing.T) {
	var plan backup.GetBackupPlanOutput
	ruleNameMap := map[string]string{}
	resourceName := "aws_backup_plan.test"
	rName := acctest.RandomWithPrefix("tf-acc-test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t); testAccPreCheckAWSBackup(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAwsBackupPlanDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAwsBackupPlanConfigRuleCopyAction(rName, 30, 180),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAwsBackupPlanExists(resourceName, &plan, &ruleNameMap),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "rule.#", "1"),
					testAccCheckAwsBackupPlanRuleAttr(resourceName, &ruleNameMap, rName, "rule_name", rName),
					testAccCheckAwsBackupPlanRuleAttr(resourceName, &ruleNameMap, rName, "lifecycle.#", "1"),
					testAccCheckAwsBackupPlanRuleAttr(resourceName, &ruleNameMap, rName, "lifecycle.0.cold_storage_after", "30"),
					testAccCheckAwsBackupPlanRuleAttr(resourceName, &ruleNameMap, rName, "lifecycle.0.delete_after", "180"),
					testAccCheckAwsBackupPlanRuleAttr(resourceName, &ruleNameMap, rName, "copy_action.#", "1"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccAwsBackupPlanConfigRuleCopyAction(rName, 60, 365),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAwsBackupPlanExists(resourceName, &plan, &ruleNameMap),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "rule.#", "1"),
					testAccCheckAwsBackupPlanRuleAttr(resourceName, &ruleNameMap, rName, "rule_name", rName),
					testAccCheckAwsBackupPlanRuleAttr(resourceName, &ruleNameMap, rName, "lifecycle.#", "1"),
					testAccCheckAwsBackupPlanRuleAttr(resourceName, &ruleNameMap, rName, "lifecycle.0.cold_storage_after", "30"),
					testAccCheckAwsBackupPlanRuleAttr(resourceName, &ruleNameMap, rName, "lifecycle.0.delete_after", "180"),
					testAccCheckAwsBackupPlanRuleAttr(resourceName, &ruleNameMap, rName, "copy_action.#", "1"),
				),
			},
			{
				Config: testAccAwsBackupPlanConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAwsBackupPlanExists(resourceName, &plan, &ruleNameMap),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "rule.#", "1"),
					testAccCheckAwsBackupPlanRuleAttr(resourceName, &ruleNameMap, rName, "rule_name", rName),
					testAccCheckAwsBackupPlanRuleAttr(resourceName, &ruleNameMap, rName, "lifecycle.#", "0"),
					testAccCheckAwsBackupPlanRuleAttr(resourceName, &ruleNameMap, rName, "copy_action.#", "0"),
				),
			},
		},
	})
}

func TestAccAwsBackupPlan_Rule_CopyAction_Multiple(t *testing.T) {
	var plan backup.GetBackupPlanOutput
	ruleNameMap := map[string]string{}
	resourceName := "aws_backup_plan.test"
	rName := acctest.RandomWithPrefix("tf-acc-test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t); testAccPreCheckAWSBackup(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAwsBackupPlanDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAwsBackupPlanConfigRuleCopyActionMultiple(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAwsBackupPlanExists(resourceName, &plan, &ruleNameMap),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "rule.#", "1"),
					testAccCheckAwsBackupPlanRuleAttr(resourceName, &ruleNameMap, rName, "rule_name", rName),
					testAccCheckAwsBackupPlanRuleAttr(resourceName, &ruleNameMap, rName, "lifecycle.#", "1"),
					testAccCheckAwsBackupPlanRuleAttr(resourceName, &ruleNameMap, rName, "lifecycle.0.cold_storage_after", "30"),
					testAccCheckAwsBackupPlanRuleAttr(resourceName, &ruleNameMap, rName, "lifecycle.0.delete_after", "180"),
					testAccCheckAwsBackupPlanRuleAttr(resourceName, &ruleNameMap, rName, "copy_action.#", "2"),
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

func TestAccAwsBackupPlan_Rule_CopyAction_CrossRegion(t *testing.T) {
	var providers []*schema.Provider
	var plan backup.GetBackupPlanOutput
	ruleNameMap := map[string]string{}
	resourceName := "aws_backup_plan.test"
	rName := acctest.RandomWithPrefix("tf-acc-test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckAWSBackup(t)
			testAccMultipleRegionsPreCheck(t)
			testAccAlternateRegionPreCheck(t)
		},
		ProviderFactories: testAccProviderFactories(&providers),
		CheckDestroy:      testAccCheckAwsBackupPlanDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAwsBackupPlanConfigRuleCopyActionCrossRegion(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAwsBackupPlanExists(resourceName, &plan, &ruleNameMap),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "rule.#", "1"),
					testAccCheckAwsBackupPlanRuleAttr(resourceName, &ruleNameMap, rName, "rule_name", rName),
					testAccCheckAwsBackupPlanRuleAttr(resourceName, &ruleNameMap, rName, "lifecycle.#", "1"),
					testAccCheckAwsBackupPlanRuleAttr(resourceName, &ruleNameMap, rName, "lifecycle.0.cold_storage_after", "30"),
					testAccCheckAwsBackupPlanRuleAttr(resourceName, &ruleNameMap, rName, "lifecycle.0.delete_after", "180"),
					testAccCheckAwsBackupPlanRuleAttr(resourceName, &ruleNameMap, rName, "copy_action.#", "1"),
				),
			},
			{
				Config:            testAccAwsBackupPlanConfigRuleCopyActionCrossRegion(rName),
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAwsBackupPlan_disappears(t *testing.T) {
	var plan backup.GetBackupPlanOutput
	ruleNameMap := map[string]string{}
	resourceName := "aws_backup_plan.test"
	rName := fmt.Sprintf("tf-testacc-backup-%s", acctest.RandStringFromCharSet(14, acctest.CharSetAlphaNum))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t); testAccPreCheckAWSBackup(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAwsBackupPlanDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAwsBackupPlanConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAwsBackupPlanExists(resourceName, &plan, &ruleNameMap),
					testAccCheckResourceDisappears(testAccProvider, resourceAwsBackupPlan(), resourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccCheckAwsBackupPlanDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*AWSClient).backupconn
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aws_backup_plan" {
			continue
		}

		input := &backup.GetBackupPlanInput{
			BackupPlanId: aws.String(rs.Primary.ID),
		}

		resp, err := conn.GetBackupPlan(input)

		if err == nil {
			if *resp.BackupPlanId == rs.Primary.ID {
				return fmt.Errorf("Plane '%s' was not deleted properly", rs.Primary.ID)
			}
		}
	}

	return nil
}

func testAccCheckAwsBackupPlanExists(name string, plan *backup.GetBackupPlanOutput, ruleNameMap *map[string]string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := testAccProvider.Meta().(*AWSClient).backupconn

		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		output, err := conn.GetBackupPlan(&backup.GetBackupPlanInput{
			BackupPlanId: aws.String(rs.Primary.ID),
		})
		if err != nil {
			return err
		}

		*plan = *output

		// Build map of rule name to hash value.
		re := regexp.MustCompile(`^rule\.(\d+)\.rule_name$`)
		for k, v := range rs.Primary.Attributes {
			matches := re.FindStringSubmatch(k)
			if matches != nil {
				(*ruleNameMap)[v] = matches[1]
			}
		}

		return nil
	}
}

func testAccCheckAwsBackupPlanRuleAttr(name string, ruleNameMap *map[string]string, ruleName, key, value string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		return resource.TestCheckResourceAttr(name, fmt.Sprintf("rule.%s.%s", (*ruleNameMap)[ruleName], key), value)(s)
	}
}

func testAccAwsBackupPlanConfig_basic(rName string) string {
	return fmt.Sprintf(`
resource "aws_backup_vault" "test" {
  name = %[1]q
}

resource "aws_backup_plan" "test" {
  name = %[1]q

  rule {
    rule_name         = %[1]q
    target_vault_name = "${aws_backup_vault.test.name}"
    schedule          = "cron(0 12 * * ? *)"
  }
}
`, rName)
}

func testAccAwsBackupPlanConfig_tags(rName string) string {
	return fmt.Sprintf(`
resource "aws_backup_vault" "test" {
  name = %[1]q
}

resource "aws_backup_plan" "test" {
  name = %[1]q

  rule {
    rule_name         = %[1]q
    target_vault_name = "${aws_backup_vault.test.name}"
    schedule          = "cron(0 12 * * ? *)"
  }

  tags = {
    Name = %[1]q
    Key1 = "Value1"
    Key2 = "Value2a"
  }
}
`, rName)
}

func testAccAwsBackupPlanConfig_tagsUpdated(rName string) string {
	return fmt.Sprintf(`
resource "aws_backup_vault" "test" {
  name = %[1]q
}

resource "aws_backup_plan" "test" {
  name = %[1]q

  rule {
    rule_name         = %[1]q
    target_vault_name = "${aws_backup_vault.test.name}"
    schedule          = "cron(0 12 * * ? *)"
  }

  tags = {
    Name = %[1]q
    Key2 = "Value2b"
    Key3 = "Value3"
  }
}
`, rName)
}

func testAccAwsBackupPlanConfig_twoRules(rName string) string {
	return fmt.Sprintf(`
resource "aws_backup_vault" "test" {
  name = %[1]q
}

resource "aws_backup_plan" "test" {
  name = %[1]q

  rule {
    rule_name         = "%[1]s_1"
    target_vault_name = "${aws_backup_vault.test.name}"
    schedule          = "cron(0 12 * * ? *)"
  }
  rule {
    rule_name         = "%[1]s_2"
    target_vault_name = "${aws_backup_vault.test.name}"
    schedule          = "cron(0 6 * * ? *)"
  }
}
`, rName)
}

func testAccAwsBackupPlanConfig_threeRules(rName string) string {
	return fmt.Sprintf(`
resource "aws_backup_vault" "test" {
  name = %[1]q
}

resource "aws_backup_plan" "test" {
  name = %[1]q

  rule {
    rule_name         = "%[1]s_1"
    target_vault_name = "${aws_backup_vault.test.name}"
    schedule          = "cron(0 6 * * ? *)"
  }
  rule {
    rule_name         = "%[1]s_2"
    target_vault_name = "${aws_backup_vault.test.name}"
    schedule          = "cron(0 12 * * ? *)"
  }
  rule {
    rule_name         = "%[1]s_3"
    target_vault_name = "${aws_backup_vault.test.name}"
    schedule          = "cron(0 18 * * ? *)"
  }
}
`, rName)
}

func testAccAwsBackupPlanConfig_lifecycleColdStorageAfterOnly(rName string) string {
	return fmt.Sprintf(`
resource "aws_backup_vault" "test" {
  name = %[1]q
}

resource "aws_backup_plan" "test" {
  name = %[1]q

  rule {
    rule_name         = %[1]q
    target_vault_name = "${aws_backup_vault.test.name}"
    schedule          = "cron(0 12 * * ? *)"

    lifecycle {
      cold_storage_after = 7
    }
  }
}
`, rName)
}

func testAccAwsBackupPlanConfig_lifecycleDeleteAfterOnly(rName string) string {
	return fmt.Sprintf(`
resource "aws_backup_vault" "test" {
  name = %[1]q
}

resource "aws_backup_plan" "test" {
  name = %[1]q

  rule {
    rule_name         = %[1]q
    target_vault_name = "${aws_backup_vault.test.name}"
    schedule          = "cron(0 12 * * ? *)"

    lifecycle {
      delete_after = 120
    }
  }
}
`, rName)
}

func testAccAwsBackupPlanConfig_lifecycleColdStorageAfterAndDeleteAfter(rName string) string {
	return fmt.Sprintf(`
resource "aws_backup_vault" "test" {
  name = %[1]q
}

resource "aws_backup_plan" "test" {
  name = %[1]q

  rule {
    rule_name         = %[1]q
    target_vault_name = "${aws_backup_vault.test.name}"
    schedule          = "cron(0 12 * * ? *)"

    lifecycle {
      cold_storage_after = 30
      delete_after       = 180
    }
  }
}
`, rName)
}

func testAccAwsBackupPlanConfig_recoveryPointTags(rName string) string {
	return fmt.Sprintf(`
resource "aws_backup_vault" "test" {
  name = %[1]q
}

resource "aws_backup_plan" "test" {
  name = %[1]q

  rule {
    rule_name         = %[1]q
    target_vault_name = "${aws_backup_vault.test.name}"
    schedule          = "cron(0 12 * * ? *)"

    recovery_point_tags = {
      Name = %[1]q
      Key1 = "Value1"
      Key2 = "Value2a"
    }
  }
}
`, rName)
}

func testAccAwsBackupPlanConfig_recoveryPointTagsUpdated(rName string) string {
	return fmt.Sprintf(`
resource "aws_backup_vault" "test" {
  name = %[1]q
}

resource "aws_backup_plan" "test" {
  name = %[1]q

  rule {
    rule_name         = %[1]q
    target_vault_name = "${aws_backup_vault.test.name}"
    schedule          = "cron(0 12 * * ? *)"

    recovery_point_tags = {
      Name = %[1]q
      Key2 = "Value2b"
      Key3 = "Value3"
    }
  }
}
`, rName)
}

func testAccAwsBackupPlanConfigRuleCopyAction(rName string, coldStorageAfter int, deleteAfter int) string {
	return fmt.Sprintf(`
resource "aws_backup_vault" "test" {
  name = "%[1]s-1"
}

resource "aws_backup_vault" "test2" {
  name = "%[1]s-2"
}

resource "aws_backup_plan" "test" {
  name = %[1]q

  rule {
    rule_name         = %[1]q
    target_vault_name = aws_backup_vault.test.name
    schedule          = "cron(0 12 * * ? *)"

    lifecycle {
      cold_storage_after = 30
      delete_after       = 180
    }

    copy_action {
      lifecycle {
        cold_storage_after = %[2]d
        delete_after       = %[3]d
      }

      destination_vault_arn = aws_backup_vault.test2.arn
    }
  }
}
`, rName, coldStorageAfter, deleteAfter)
}

func testAccAwsBackupPlanConfigRuleCopyActionMultiple(rName string) string {
	return fmt.Sprintf(`
resource "aws_backup_vault" "test" {
  name = "%[1]s-1"
}

resource "aws_backup_vault" "test2" {
  name = "%[1]s-2"
}

resource "aws_backup_vault" "test3" {
  name = "%[1]s-3"
}

resource "aws_backup_plan" "test" {
  name = %[1]q

  rule {
    rule_name         = %[1]q
    target_vault_name = aws_backup_vault.test.name
    schedule          = "cron(0 12 * * ? *)"

    lifecycle {
      cold_storage_after = 30
      delete_after       = 180
    }

    copy_action {
      lifecycle {
        cold_storage_after = 30
        delete_after       = 180
      }

      destination_vault_arn = aws_backup_vault.test2.arn
    }

    copy_action {
      lifecycle {
        cold_storage_after = 60
        delete_after       = 365
      }

      destination_vault_arn = aws_backup_vault.test3.arn
    }
  }
}
`, rName)
}

func testAccAwsBackupPlanConfigRuleCopyActionCrossRegion(rName string) string {
	return testAccAlternateRegionProviderConfig() + fmt.Sprintf(`
resource "aws_backup_vault" "test" {
  name = "%[1]s-1"
}

resource "aws_backup_vault" "test2" {
  provider = "aws.alternate"
  name     = "%[1]s-2"
}

resource "aws_backup_plan" "test" {
  name = %[1]q

  rule {
    rule_name         = %[1]q
    target_vault_name = "${aws_backup_vault.test.name}"
    schedule          = "cron(0 12 * * ? *)"

    lifecycle {
      cold_storage_after = 30
      delete_after       = 180
    }

    copy_action {
      lifecycle {
        cold_storage_after = 30
        delete_after       = 180
      }

      destination_vault_arn = "${aws_backup_vault.test2.arn}"
    }
  }
}
`, rName)
}
