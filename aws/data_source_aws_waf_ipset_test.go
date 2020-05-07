package aws

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccDataSourceAwsWafIPSet_Basic(t *testing.T) {
	name := acctest.RandomWithPrefix("tf-acc-test")
	resourceName := "aws_waf_ipset.ipset"
	datasourceName := "data.aws_waf_ipset.ipset"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testAccDataSourceAwsWafIPSet_NonExistent,
				ExpectError: regexp.MustCompile(`WAF IP Set not found`),
			},
			{
				Config: testAccDataSourceAwsWafIPSet_Name(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(datasourceName, "id", resourceName, "id"),
					resource.TestCheckResourceAttrPair(datasourceName, "name", resourceName, "name"),
				),
			},
		},
	})
}

func testAccDataSourceAwsWafIPSet_Name(name string) string {
	return fmt.Sprintf(`
resource "aws_waf_ipset" "ipset" {
  name = %[1]q
}
data "aws_waf_ipset" "ipset" {
  name = "${aws_waf_ipset.ipset.name}"
}
`, name)
}

const testAccDataSourceAwsWafIPSet_NonExistent = `
data "aws_waf_ipset" "ipset" {
  name = "tf-acc-test-does-not-exist"
}
`
