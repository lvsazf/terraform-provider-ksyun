package ksyun

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"testing"
)

func TestAccKsyunRabbitmqSecurityRule_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRabbitmqSecurityRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRabbitmqSecurityRuleConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRabbitmqSecurityRuleExists("ksyun_rabbitmq_security_rule.default"),
				),
			},
		},
	})
}

func testAccCheckRabbitmqSecurityRuleExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("can't find resource or data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("rabbitmq instance is not exist")
		}

		client := testAccProvider.Meta().(*KsyunClient)
		securityRuleCheck := make(map[string]interface{})
		securityRuleCheck["InstanceId"] = rs.Primary.ID
		resp, err := client.rabbitmqconn.DescribeSecurityGroupRules(&securityRuleCheck)
		if err != nil {
			return fmt.Errorf("error on reading rabbitmq instance security rule %q, %s", rs.Primary.ID, err)
		}
		rules := (*resp)["Data"].([]interface{})
		if len(rules) == 0 {
			return fmt.Errorf("rabbitmq instance security rule is not exist")
		}

		return nil
	}
}

func testAccCheckRabbitmqSecurityRuleDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*KsyunClient)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "ksyun_rabbitmq_security_rule" {
			securityRuleCheck := make(map[string]interface{})
			securityRuleCheck["InstanceId"] = rs.Primary.ID
			resp, err := client.rabbitmqconn.DescribeSecurityGroupRules(&securityRuleCheck)

			if err != nil {
				return nil
			}

			rules := (*resp)["Data"].([]interface{})
			if len(rules) > 0 {
				return fmt.Errorf("delete rabbitmq security rule failure")
			}

			return nil
		}
	}

	return nil
}

const testAccRabbitmqSecurityRuleConfig = `
data "ksyun_availability_zones" "default" {
  output_file=""
  ids=[]
}
variable "available_zone" {
  default = "cn-beijing-6a"
}

resource "ksyun_vpc" "default" {
  vpc_name   = "lzs-ksyun-vpc-tf"
  cidr_block = "10.7.0.0/21"
}

variable "protocol" {
  default = "rabbitmq 3.7"
}

resource "ksyun_subnet" "default" {
  subnet_name      	= "lzs_ksyun_subnet_tf"
  cidr_block 		= "10.7.0.0/21"
  subnet_type 		= "Reserve"
  dhcp_ip_from 		= "10.7.0.2"
  dhcp_ip_to 		= "10.7.0.253"
  vpc_id 			= "${ksyun_vpc.default.id}"
  gateway_ip 		= "10.7.0.1"
  dns1 				= "198.18.254.41"
  dns2 				= "198.18.254.40"
  availability_zone = "${data.ksyun_availability_zones.default.availability_zones.0.availability_zone_name}"
}

resource "ksyun_rabbitmq_instance" "default" {
  availability_zone     = "${data.ksyun_availability_zones.default.availability_zones.0.availability_zone_name}"
  instance_name         = "lzs_rabbitmq_instance"
  instance_password     = "Shiwo1101"
  vpc_id                = "${ksyun_vpc.default.id}"
  mode                  = 1
  subnet_id             = "${ksyun_subnet.default.id}"
  engine_version        = "3.7"
  instance_type         = "2C4G"
  ssd_disk              = "5"
  node_num              = 3
  bill_type             = 87
  project_id            = 103800
  project_name          = "测试部"
}

resource "ksyun_rabbitmq_security_rule" "default" {
  instance_id = "${ksyun_rabbitmq_instance.default.id}"
  cidrs = "192.168.10.11/32"
}
`
