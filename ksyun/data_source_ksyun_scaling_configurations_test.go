package ksyun

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"testing"
)

func TestAccKsyunScalingConfigurationsDataSource_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataScalingConfigurationsConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIDExists("data.ksyun_scaling_configurations.foo"),
				),
			},
		},
	})
}

const testAccDataScalingConfigurationsConfig = `
data "ksyun_scaling_configurations" "foo" {
  output_file="output_result"
}
`
