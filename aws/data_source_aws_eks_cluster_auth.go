package aws

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-aws/aws/internal/service/eks/token"
)

func dataSourceAwsEksClusterAuth() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAwsEksClusterAuthRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.NoZeroValues,
			},

			"token": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
		},
	}
}

func dataSourceAwsEksClusterAuthRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AWSClient).stsconn
	name := d.Get("name").(string)
	generator, err := token.NewGenerator(false, false)
	if err != nil {
		return fmt.Errorf("error getting token generator: %v", err)
	}
	token, err := generator.GetWithSTS(name, conn)
	if err != nil {
		return fmt.Errorf("error getting token: %v", err)
	}

	d.SetId(time.Now().UTC().String())
	d.Set("token", token.Token)

	return nil
}
