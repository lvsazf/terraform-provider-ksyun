package ksyun

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"testing"
)

func TestAccKsyunRoute_basic(t *testing.T) {
	var val map[string]interface{}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		IDRefreshName: "ksyun_route.foo",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckRouteDestroy,

		Steps: []resource.TestStep{
			{
				Config: testAccRouteConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRouteExists("ksyun_route.foo", &val),
					testAccCheckRouteAttributes(&val),
				),
			},
		},
	})
}

func testAccCheckRouteExists(n string, val *map[string]interface{}) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("route id is empty")
		}

		client := testAccProvider.Meta().(*KsyunClient)
		Route := make(map[string]interface{})
		Route["RouteId.1"] = rs.Primary.ID
		ptr, err := client.vpcconn.DescribeRoutes(&Route)

		if err != nil {
			return err
		}
		if ptr != nil {
			l := (*ptr)["RouteSet"].([]interface{})
			if len(l) == 0 {
				return err
			}
		}

		*val = *ptr
		return nil
	}
}

func testAccCheckRouteAttributes(val *map[string]interface{}) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if val != nil {
			l := (*val)["RouteSet"].([]interface{})
			if len(l) == 0 {
				return fmt.Errorf("route id is empty")
			}
		}
		return nil
	}
}

func testAccCheckRouteDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ksyun_route" {
			continue
		}

		client := testAccProvider.Meta().(*KsyunClient)
		Route := make(map[string]interface{})
		Route["RouteId.1"] = rs.Primary.ID
		ptr, err := client.vpcconn.DescribeRoutes(&Route)

		// Verify the error is what we want
		if err != nil {
			return err
		}
		if ptr != nil {
			l := (*ptr)["RouteSet"].([]interface{})
			if len(l) == 0 {
				continue
			} else {
				return fmt.Errorf("route still exist")
			}
		}
	}

	return nil
}

const testAccRouteConfig = `
resource "ksyun_vpc" "test" {
  vpc_name   = "ksyun-vpc-tf"
  cidr_block = "10.7.0.0/21"
}
resource "ksyun_route" "foo" {
  destination_cidr_block = "10.0.0.0/16"
  route_type = "InternetGateway"
  vpc_id = "${ksyun_vpc.test.id}"
}
`
