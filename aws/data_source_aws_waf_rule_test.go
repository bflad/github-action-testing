package aws

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccDataSourceAwsWafRule_Basic(t *testing.T) {
	name := acctest.RandomWithPrefix("tf-acc-test")
	resourceName := "aws_waf_rule.wafrule"
	datasourceName := "data.aws_waf_rule.wafrule"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testAccDataSourceAwsWafRuleConfig_NonExistent,
				ExpectError: regexp.MustCompile(`WAF Rules not found`),
			},
			{
				Config: testAccDataSourceAwsWafRuleConfig_Name(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(datasourceName, "id", resourceName, "id"),
					resource.TestCheckResourceAttrPair(datasourceName, "name", resourceName, "name"),
				),
			},
		},
	})
}

func testAccDataSourceAwsWafRuleConfig_Name(name string) string {
	return fmt.Sprintf(`
resource "aws_waf_rule" "wafrule" {
  name        = %[1]q
  metric_name = "WafruleTest"
}

data "aws_waf_rule" "wafrule" {
  name = "${aws_waf_rule.wafrule.name}"
}
`, name)
}

const testAccDataSourceAwsWafRuleConfig_NonExistent = `
data "aws_waf_rule" "wafrule" {
  name = "tf-acc-test-does-not-exist"
}
`
