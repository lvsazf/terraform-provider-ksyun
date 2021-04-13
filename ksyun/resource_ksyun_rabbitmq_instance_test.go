package ksyun

import (
	"errors"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-ksyun/logger"
	"strings"
	"testing"
)

func TestAccKsyunRabbitmqInstance_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRabbitmqInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRabbitmqInstanceConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRabbitmqInstanceExists("ksyun_rabbitmq_instance.default"),
				),
				ImportStateVerify: false,
			},
			{
				Config: testRabbitmqUpdateAccKcsConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKcsInstanceExists("ksyun_rabbitmq_instance.default"),
				),
			},
		},
	})
}

func testAccCheckRabbitmqInstanceDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*KsyunClient)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "ksyun_rabbitmq_instance" {
			instanceCheck := make(map[string]interface{})
			instanceCheck["instanceId"] = rs.Primary.ID
			resp, err := client.rabbitmqconn.DescribeInstance(&instanceCheck)

			if err != nil {
				if strings.Contains(err.Error(), "InstanceNotFound") {
					return nil
				}
				return err
			}
			if resp != nil {
				if (*resp)["Data"] != nil {
					return errors.New("delete instance failure")
				}
			}
		}
	}

	return nil
}

func testAccCheckRabbitmqInstanceExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("can't find resource or data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("rabbitmq instance create failure")
		}

		client := testAccProvider.Meta().(*KsyunClient)
		readReq := make(map[string]interface{})
		readReq["instanceId"] = rs.Primary.ID

		logger.Debug(logger.ReqFormat, "DescribeRabbitmqInstance", readReq)
		_, err := client.rabbitmqconn.DescribeInstance(&readReq)
		if err != nil {
			return fmt.Errorf("error on reading instance %q, %s", rs.Primary.ID, err)
		}

		return nil
	}
}

const testAccRabbitmqInstanceConfig = `
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
  availability_zone = "${var.available_zone}"
}

resource "ksyun_rabbitmq_instance" "default" {
  availability_zone     = "${var.available_zone}"
  instance_name         = "my_rabbitmq_instance"
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
`

const testRabbitmqUpdateAccKcsConfig = `
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
  availability_zone = "${var.available_zone}"
}

resource "ksyun_rabbitmq_instance" "default" {
  availability_zone     = "${var.available_zone}"
  instance_name         = "my_rabbitmq_instance_haha"
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
`
