package ksyun

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"strings"
	"testing"
)

func TestAccKsyunScalingInstance_basic(t *testing.T) {
	var val map[string]interface{}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		IDRefreshName: "ksyun_scaling_instance.foo",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckScalingInstanceDestroy,

		Steps: []resource.TestStep{
			{
				Config: testAccScalingInstanceConfig,

				Check: resource.ComposeTestCheckFunc(
					testAccCheckScalingInstanceExists("ksyun_scaling_instance.foo", &val),
					testAccCheckScalingInstanceAttributes(&val),
				),
			},
		},
	})
}

func TestAccKsyunScalingInstance_update(t *testing.T) {
	var val map[string]interface{}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		IDRefreshName: "ksyun_scaling_instance.foo",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckScalingInstanceDestroy,

		Steps: []resource.TestStep{
			{
				Config: testAccScalingInstanceConfig,

				Check: resource.ComposeTestCheckFunc(
					testAccCheckScalingInstanceExists("ksyun_scaling_instance.foo", &val),
					testAccCheckScalingInstanceAttributes(&val),
				),
			},
			{
				Config: testAccScalingInstanceConfigUpdate,

				Check: resource.ComposeTestCheckFunc(
					testAccCheckScalingInstanceExists("ksyun_scaling_instance.foo", &val),
					testAccCheckScalingInstanceAttributes(&val),
				),
			},
		},
	})
}

func testAccCheckScalingInstanceExists(n string, val *map[string]interface{}) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf(" ScalingInstance id is empty ")
		}

		client := testAccProvider.Meta().(*KsyunClient)
		req := make(map[string]interface{})
		req["ScalingGroupId"] = strings.Split(rs.Primary.ID, ":")[1]
		req["ScalingInstanceId.1"] = strings.Split(rs.Primary.ID, ":")[0]
		ptr, err := client.kecconn.DescribeScalingInstance(&req)

		if err != nil {
			return err
		}
		if ptr != nil {
			l := (*ptr)["ScalingInstanceSet"].([]interface{})
			if len(l) == 0 {
				return err
			}
		}

		*val = *ptr
		return nil
	}
}

func testAccCheckScalingInstanceAttributes(val *map[string]interface{}) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if val != nil {
			l := (*val)["ScalingInstanceSet"].([]interface{})
			if len(l) == 0 {
				return fmt.Errorf(" ScalingInstance id is empty ")
			}
		}
		return nil
	}
}

func testAccCheckScalingInstanceDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ksyun_scaling_instance" {
			continue
		}

		client := testAccProvider.Meta().(*KsyunClient)
		req := make(map[string]interface{})
		req["ScalingGroupId"] = strings.Split(rs.Primary.ID, ":")[1]
		req["ScalingInstanceId.1"] = strings.Split(rs.Primary.ID, ":")[0]
		ptr, err := client.kecconn.DescribeScalingInstance(&req)

		// Verify the error is what we want
		if err != nil {
			return err
		}
		if ptr != nil && (*ptr)["ScalingInstanceSet"] != nil {
			l := (*ptr)["ScalingInstanceSet"].([]interface{})
			if len(l) == 0 {
				continue
			} else {
				return fmt.Errorf(" ScalingInstance still exist ")
			}
		}
	}

	return nil
}

const testAccScalingInstanceConfig = `

provider "ksyun" {
  region = "cn-beijing-6"
}

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

resource "ksyun_instance" "foo" {
  image_id="IMG-5465174a-6d71-4770-b8e1-917a0dd92466"
  instance_type="N3.1B"
  subnet_id= ksyun_subnet.foo.id
  instance_password="Aa123456"
  charge_type="Daily"
  security_group_id=[ksyun_security_group.foo.id]
  instance_name="test-tf"
}

resource "ksyun_scaling_instance" "foo" {
  scaling_group_id = ksyun_scaling_group.foo.id
  scaling_instance_id = ksyun_instance.foo.id
  protected_from_detach = 1
}
`

const testAccScalingInstanceConfigUpdate = `
provider "ksyun" {
  region = "cn-beijing-6"
}

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

resource "ksyun_instance" "foo" {
  image_id="IMG-5465174a-6d71-4770-b8e1-917a0dd92466"
  instance_type="N3.1B"
  subnet_id= ksyun_subnet.foo.id
  instance_password="Aa123456"
  charge_type="Daily"
  security_group_id=[ksyun_security_group.foo.id]
  instance_name="test-tf"
}

resource "ksyun_scaling_instance" "foo" {
  scaling_group_id = ksyun_scaling_group.foo.id
  scaling_instance_id = ksyun_instance.foo.id
  protected_from_detach = 0
}

`
