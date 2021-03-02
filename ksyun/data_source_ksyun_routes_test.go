package ksyun

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"testing"
)

func TestAccKsyunRoutesDataSource_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataRouteConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIDExists("data.ksyun_routes.foo"),
				),
			},
		},
	})
}

const testAccDataRouteConfig = `
resource "ksyun_vpc" "test" {
  vpc_name   = "ksyun-vpc-tf"
  cidr_block = "10.7.0.0/21"
}
resource "ksyun_route" "foo" {
  destination_cidr_block = "10.0.0.0/16"
  route_type = "InternetGateway"
  vpc_id = "${ksyun_vpc.test.id}"
}
data "ksyun_routes" "foo" {
  output_file="output_result"
  ids = ["${ksyun_route.foo.id}"]
}
`
