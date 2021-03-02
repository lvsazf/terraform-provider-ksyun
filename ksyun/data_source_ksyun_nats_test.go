package ksyun

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"testing"
)

func TestAccKsyunNatsDataSource_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataNatConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIDExists("data.ksyun_nats.foo"),
				),
			},
		},
	})
}

const testAccDataNatConfig = `
data "ksyun_nats" "foo" {
  output_file="output_result"
  ids=["24e8534e-bee7-4737-bb2a-9d62d3ac8842"]
}
`
