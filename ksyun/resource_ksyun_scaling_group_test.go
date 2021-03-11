package ksyun

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"testing"
)

func TestAccKsyunScalingGroup_basic(t *testing.T) {
	var val map[string]interface{}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		IDRefreshName: "ksyun_scaling_group.foo",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckScalingGroupDestroy,

		Steps: []resource.TestStep{
			{
				Config: testAccScalingGroupConfig,

				Check: resource.ComposeTestCheckFunc(
					testAccCheckScalingGroupExists("ksyun_scaling_group.foo", &val),
					testAccCheckScalingGroupAttributes(&val),
				),
			},
		},
	})
}

func TestAccKsyunScalingGroup_update(t *testing.T) {
	var val map[string]interface{}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		IDRefreshName: "ksyun_scaling_group.foo",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckScalingGroupDestroy,

		Steps: []resource.TestStep{
			{
				Config: testAccScalingGroupConfig,

				Check: resource.ComposeTestCheckFunc(
					testAccCheckScalingGroupExists("ksyun_scaling_group.foo", &val),
					testAccCheckScalingGroupAttributes(&val),
				),
			},
			{
				Config: testAccScalingGroupConfigUpdate,

				Check: resource.ComposeTestCheckFunc(
					testAccCheckScalingGroupExists("ksyun_scaling_group.foo", &val),
					testAccCheckScalingGroupAttributes(&val),
				),
			},
		},
	})
}

func testAccCheckScalingGroupExists(n string, val *map[string]interface{}) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf(" ScalingGroup id is empty ")
		}

		client := testAccProvider.Meta().(*KsyunClient)
		req := make(map[string]interface{})
		req["ScalingGroupId.1"] = rs.Primary.ID
		ptr, err := client.kecconn.DescribeScalingGroup(&req)

		if err != nil {
			return err
		}
		if ptr != nil {
			l := (*ptr)["ScalingGroupSet"].([]interface{})
			if len(l) == 0 {
				return err
			}
		}

		*val = *ptr
		return nil
	}
}

func testAccCheckScalingGroupAttributes(val *map[string]interface{}) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if val != nil {
			l := (*val)["ScalingGroupSet"].([]interface{})
			if len(l) == 0 {
				return fmt.Errorf(" ScalingGroup id is empty ")
			}
		}
		return nil
	}
}

func testAccCheckScalingGroupDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ksyun_scaling_group" {
			continue
		}

		client := testAccProvider.Meta().(*KsyunClient)
		req := make(map[string]interface{})
		req["ScalingGroupId.1"] = rs.Primary.ID
		ptr, err := client.kecconn.DescribeScalingGroup(&req)

		// Verify the error is what we want
		if err != nil {
			return err
		}
		if ptr != nil && (*ptr)["ScalingGroupSet"] != nil {
			l := (*ptr)["ScalingGroupSet"].([]interface{})
			if len(l) == 0 {
				continue
			} else {
				return fmt.Errorf(" ScalingGroup still exist ")
			}
		}
	}

	return nil
}

const testAccScalingGroupConfig = `
resource "ksyun_vpc" "foo" {
  vpc_name = "tf-example-vpc-01"
  cidr_block = "10.0.0.0/16"
}

resource "ksyun_subnet" "foo" {
  subnet_name = "tf-acc-subnet1"
  cidr_block = "10.0.5.0/24"
  subnet_type = "Normal"
  dhcp_ip_from = "10.0.5.2"
  dhcp_ip_to = "10.0.5.253"
  vpc_id = ksyun_vpc.foo.id
  gateway_ip = "10.0.5.1"
  dns1 = "198.18.254.41"
  dns2 = "198.18.254.40"
  availability_zone = "cn-beijing-6b"
}

resource "ksyun_security_group" "foo" {
  vpc_id = ksyun_vpc.foo.id
  security_group_name = "tf-acc-sg"
}

resource "ksyun_lb" "foo" {
  vpc_id = ksyun_vpc.foo.id
  load_balancer_name = "tf-acc-lb"
  type = "public"
  load_balancer_state = "start"
}

resource "ksyun_lb_listener" "foo" {
  listener_name = "tf-acc-listener"
  listener_port = "80"
  listener_protocol = "HTTP"
  listener_state = "stop"
  load_balancer_id = ksyun_lb.foo.id
  method = "RoundRobin"
  session {
    session_state = "stop"
    session_persistence_period = 3600
  }
}

resource "ksyun_scaling_configuration" "foo" {
  scaling_configuration_name = "tf-xym-sc"
  image_id = "IMG-5465174a-6d71-4770-b8e1-917a0dd92466"
  instance_type = "N3.1B"
  password = "Aa123456"
  data_disks  {
      disk_type = "EHDD"
      disk_size = 50
      delete_with_instance = true
  }
}

resource "ksyun_scaling_group" "foo" {
  subnet_id_set = [ksyun_subnet.foo.id]
  security_group_id = ksyun_security_group.foo.id
  scaling_configuration_id = ksyun_scaling_configuration.foo.id
  min_size = 0
  max_size = 2
  desired_capacity = 0
  status = "UnActive"
  slb_config_set  {
    slb_id = ksyun_lb.foo.id
    listener_id = ksyun_lb_listener.foo.id
    server_port_set = [80]
  }
}
`

const testAccScalingGroupConfigUpdate = `
resource "ksyun_vpc" "foo" {
  vpc_name = "tf-example-vpc-01"
  cidr_block = "10.0.0.0/16"
}

resource "ksyun_subnet" "foo" {
  subnet_name = "tf-acc-subnet1"
  cidr_block = "10.0.5.0/24"
  subnet_type = "Normal"
  dhcp_ip_from = "10.0.5.2"
  dhcp_ip_to = "10.0.5.253"
  vpc_id = ksyun_vpc.foo.id
  gateway_ip = "10.0.5.1"
  dns1 = "198.18.254.41"
  dns2 = "198.18.254.40"
  availability_zone = "cn-beijing-6b"
}

resource "ksyun_security_group" "foo" {
  vpc_id = ksyun_vpc.foo.id
  security_group_name = "tf-acc-sg"
}

resource "ksyun_lb" "foo" {
  vpc_id = ksyun_vpc.foo.id
  load_balancer_name = "tf-acc-lb"
  type = "public"
  load_balancer_state = "start"
}

resource "ksyun_lb_listener" "foo" {
  listener_name = "tf-acc-listener"
  listener_port = "80"
  listener_protocol = "HTTP"
  listener_state = "stop"
  load_balancer_id = ksyun_lb.foo.id
  method = "RoundRobin"
  session {
    session_state = "stop"
    session_persistence_period = 3600
  }
}

resource "ksyun_scaling_configuration" "foo" {
  scaling_configuration_name = "tf-xym-sc"
  image_id = "IMG-5465174a-6d71-4770-b8e1-917a0dd92466"
  instance_type = "N3.1B"
  password = "Aa123456"
  data_disks  {
      disk_type = "EHDD"
      disk_size = 50
      delete_with_instance = true
  }
}

resource "ksyun_scaling_group" "foo" {
  subnet_id_set = [ksyun_subnet.foo.id]
  security_group_id = ksyun_security_group.foo.id
  scaling_configuration_id = ksyun_scaling_configuration.foo.id
  min_size = 0
  max_size = 2
  desired_capacity = 0
  status = "Active"
  slb_config_set  {
    slb_id = ksyun_lb.foo.id
    listener_id = ksyun_lb_listener.foo.id
    server_port_set = [80]
  }
}
`
