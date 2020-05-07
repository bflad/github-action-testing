package aws

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/waf"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccAWSWafXssMatchSet_basic(t *testing.T) {
	var v waf.XssMatchSet
	rName := acctest.RandomWithPrefix("tf-acc-test")
	resourceName := "aws_waf_xss_match_set.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t); testAccPreCheckAWSWaf(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAWSWafXssMatchSetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSWafXssMatchSetConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSWafXssMatchSetExists(resourceName, &v),
					testAccMatchResourceAttrGlobalARN(resourceName, "arn", "waf", regexp.MustCompile(`xssmatchset/.+`)),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "xss_match_tuples.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "xss_match_tuples.41660541.field_to_match.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "xss_match_tuples.41660541.field_to_match.0.data", ""),
					resource.TestCheckResourceAttr(resourceName, "xss_match_tuples.41660541.field_to_match.0.type", "URI"),
					resource.TestCheckResourceAttr(resourceName, "xss_match_tuples.41660541.text_transformation", "NONE"),
					resource.TestCheckResourceAttr(resourceName, "xss_match_tuples.599421078.field_to_match.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "xss_match_tuples.599421078.field_to_match.0.data", ""),
					resource.TestCheckResourceAttr(resourceName, "xss_match_tuples.599421078.field_to_match.0.type", "QUERY_STRING"),
					resource.TestCheckResourceAttr(resourceName, "xss_match_tuples.599421078.text_transformation", "NONE"),
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

func TestAccAWSWafXssMatchSet_changeNameForceNew(t *testing.T) {
	var before, after waf.XssMatchSet
	rName := acctest.RandomWithPrefix("tf-acc-test")
	xssMatchSetNewName := fmt.Sprintf("xssMatchSetNewName-%s", acctest.RandString(5))
	resourceName := "aws_waf_xss_match_set.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t); testAccPreCheckAWSWaf(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAWSWafXssMatchSetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSWafXssMatchSetConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSWafXssMatchSetExists(resourceName, &before),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "xss_match_tuples.#", "2"),
				),
			},
			{
				Config: testAccAWSWafXssMatchSetConfigChangeName(xssMatchSetNewName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSWafXssMatchSetExists(resourceName, &after),
					resource.TestCheckResourceAttr(resourceName, "name", xssMatchSetNewName),
					resource.TestCheckResourceAttr(resourceName, "xss_match_tuples.#", "2"),
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

func TestAccAWSWafXssMatchSet_disappears(t *testing.T) {
	var v waf.XssMatchSet
	rName := acctest.RandomWithPrefix("tf-acc-test")
	resourceName := "aws_waf_xss_match_set.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t); testAccPreCheckAWSWaf(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAWSWafXssMatchSetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSWafXssMatchSetConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSWafXssMatchSetExists(resourceName, &v),
					testAccCheckAWSWafXssMatchSetDisappears(&v),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccAWSWafXssMatchSet_changeTuples(t *testing.T) {
	var before, after waf.XssMatchSet
	setName := acctest.RandomWithPrefix("tf-acc-test")
	resourceName := "aws_waf_xss_match_set.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t); testAccPreCheckAWSWaf(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAWSWafXssMatchSetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSWafXssMatchSetConfig(setName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckAWSWafXssMatchSetExists(resourceName, &before),
					resource.TestCheckResourceAttr(resourceName, "name", setName),
					resource.TestCheckResourceAttr(resourceName, "xss_match_tuples.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "xss_match_tuples.599421078.field_to_match.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "xss_match_tuples.599421078.field_to_match.0.data", ""),
					resource.TestCheckResourceAttr(resourceName, "xss_match_tuples.599421078.field_to_match.0.type", "QUERY_STRING"),
					resource.TestCheckResourceAttr(resourceName, "xss_match_tuples.599421078.text_transformation", "NONE"),
					resource.TestCheckResourceAttr(resourceName, "xss_match_tuples.41660541.field_to_match.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "xss_match_tuples.41660541.field_to_match.0.data", ""),
					resource.TestCheckResourceAttr(resourceName, "xss_match_tuples.41660541.field_to_match.0.type", "URI"),
					resource.TestCheckResourceAttr(resourceName, "xss_match_tuples.41660541.text_transformation", "NONE"),
				),
			},
			{
				Config: testAccAWSWafXssMatchSetConfig_changeTuples(setName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckAWSWafXssMatchSetExists(resourceName, &after),
					resource.TestCheckResourceAttr(resourceName, "name", setName),
					resource.TestCheckResourceAttr(resourceName, "xss_match_tuples.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "xss_match_tuples.42378128.field_to_match.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "xss_match_tuples.42378128.field_to_match.0.data", "GET"),
					resource.TestCheckResourceAttr(resourceName, "xss_match_tuples.42378128.field_to_match.0.type", "METHOD"),
					resource.TestCheckResourceAttr(resourceName, "xss_match_tuples.42378128.text_transformation", "HTML_ENTITY_DECODE"),
					resource.TestCheckResourceAttr(resourceName, "xss_match_tuples.3815294338.field_to_match.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "xss_match_tuples.3815294338.field_to_match.0.data", ""),
					resource.TestCheckResourceAttr(resourceName, "xss_match_tuples.3815294338.field_to_match.0.type", "BODY"),
					resource.TestCheckResourceAttr(resourceName, "xss_match_tuples.3815294338.text_transformation", "CMD_LINE"),
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

func TestAccAWSWafXssMatchSet_noTuples(t *testing.T) {
	var ipset waf.XssMatchSet
	setName := acctest.RandomWithPrefix("tf-acc-test")
	resourceName := "aws_waf_xss_match_set.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t); testAccPreCheckAWSWaf(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAWSWafXssMatchSetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSWafXssMatchSetConfig_noTuples(setName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckAWSWafXssMatchSetExists(resourceName, &ipset),
					resource.TestCheckResourceAttr(resourceName, "name", setName),
					resource.TestCheckResourceAttr(resourceName, "xss_match_tuples.#", "0"),
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

func testAccCheckAWSWafXssMatchSetDisappears(v *waf.XssMatchSet) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := testAccProvider.Meta().(*AWSClient).wafconn

		wr := newWafRetryer(conn)
		_, err := wr.RetryWithToken(func(token *string) (interface{}, error) {
			req := &waf.UpdateXssMatchSetInput{
				ChangeToken:   token,
				XssMatchSetId: v.XssMatchSetId,
			}

			for _, xssMatchTuple := range v.XssMatchTuples {
				xssMatchTupleUpdate := &waf.XssMatchSetUpdate{
					Action: aws.String(waf.ChangeActionDelete),
					XssMatchTuple: &waf.XssMatchTuple{
						FieldToMatch:       xssMatchTuple.FieldToMatch,
						TextTransformation: xssMatchTuple.TextTransformation,
					},
				}
				req.Updates = append(req.Updates, xssMatchTupleUpdate)
			}
			return conn.UpdateXssMatchSet(req)
		})
		if err != nil {
			return fmt.Errorf("Error updating XssMatchSet: %s", err)
		}

		_, err = wr.RetryWithToken(func(token *string) (interface{}, error) {
			opts := &waf.DeleteXssMatchSetInput{
				ChangeToken:   token,
				XssMatchSetId: v.XssMatchSetId,
			}
			return conn.DeleteXssMatchSet(opts)
		})
		if err != nil {
			return fmt.Errorf("Error deleting WAF XSS Match Set: %s", err)
		}
		return nil
	}
}

func testAccCheckAWSWafXssMatchSetExists(n string, v *waf.XssMatchSet) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No WAF XSS Match Set ID is set")
		}

		conn := testAccProvider.Meta().(*AWSClient).wafconn
		resp, err := conn.GetXssMatchSet(&waf.GetXssMatchSetInput{
			XssMatchSetId: aws.String(rs.Primary.ID),
		})

		if err != nil {
			return err
		}

		if *resp.XssMatchSet.XssMatchSetId == rs.Primary.ID {
			*v = *resp.XssMatchSet
			return nil
		}

		return fmt.Errorf("WAF XssMatchSet (%s) not found", rs.Primary.ID)
	}
}

func testAccCheckAWSWafXssMatchSetDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aws_waf_xss_match_set" {
			continue
		}

		conn := testAccProvider.Meta().(*AWSClient).wafconn
		resp, err := conn.GetXssMatchSet(
			&waf.GetXssMatchSetInput{
				XssMatchSetId: aws.String(rs.Primary.ID),
			})

		if err == nil {
			if *resp.XssMatchSet.XssMatchSetId == rs.Primary.ID {
				return fmt.Errorf("WAF XssMatchSet %s still exists", rs.Primary.ID)
			}
		}

		// Return nil if the XssMatchSet is already destroyed
		if isAWSErr(err, waf.ErrCodeNonexistentItemException, "") {
			return nil
		}

		return err
	}

	return nil
}

func testAccAWSWafXssMatchSetConfig(name string) string {
	return fmt.Sprintf(`
resource "aws_waf_xss_match_set" "test" {
  name = %[1]q

  xss_match_tuples {
    text_transformation = "NONE"

    field_to_match {
      type = "URI"
    }
  }

  xss_match_tuples {
    text_transformation = "NONE"

    field_to_match {
      type = "QUERY_STRING"
    }
  }
}
`, name)
}

func testAccAWSWafXssMatchSetConfigChangeName(name string) string {
	return fmt.Sprintf(`
resource "aws_waf_xss_match_set" "test" {
  name = %[1]q

  xss_match_tuples {
    text_transformation = "NONE"

    field_to_match {
      type = "URI"
    }
  }

  xss_match_tuples {
    text_transformation = "NONE"

    field_to_match {
      type = "QUERY_STRING"
    }
  }
}
`, name)
}

func testAccAWSWafXssMatchSetConfig_changeTuples(name string) string {
	return fmt.Sprintf(`
resource "aws_waf_xss_match_set" "test" {
  name = %[1]q

  xss_match_tuples {
    text_transformation = "CMD_LINE"

    field_to_match {
      type = "BODY"
    }
  }

  xss_match_tuples {
    text_transformation = "HTML_ENTITY_DECODE"

    field_to_match {
      type = "METHOD"
      data = "GET"
    }
  }
}
`, name)
}

func testAccAWSWafXssMatchSetConfig_noTuples(name string) string {
	return fmt.Sprintf(`
resource "aws_waf_xss_match_set" "test" {
  name = %[1]q
}
`, name)
}
