# Specify the provider and access details
provider "ksyun" {
  access_key = ""
  secret_key = ""
  region = "cn-beijing-6"
}

data "ksyun_availability_zones" "default" {
  output_file = ""
  ids = []
}
resource "ksyun_vpc" "default" {
  vpc_name = "lzs-ksyun-vpc-tf"
  cidr_block = "10.7.0.0/21"
}
resource "ksyun_subnet" "default" {
  subnet_name = "lzs_ksyun_subnet_tf"
  cidr_block = "10.7.0.0/21"
  subnet_type = "Reserve"
  dhcp_ip_from = "10.7.0.2"
  dhcp_ip_to = "10.7.0.253"
  vpc_id = "${ksyun_vpc.default.id}"
  gateway_ip = "10.7.0.1"
  dns1 = "198.18.254.41"
  dns2 = "198.18.254.40"
  availability_zone = "${data.ksyun_availability_zones.default.availability_zones.0.availability_zone_name}"
}

resource "ksyun_rabbitmq_instance" "default" {
  availability_zone = "${var.available_zone}"
  instance_name = "my_rabbitmq_instance"
  instance_password = "Shiwo1101"
  vpc_id = "${ksyun_vpc.default.id}"
  mode = 1
  subnet_id = "${ksyun_subnet.default.id}"
  engine_version = "3.7"
  instance_type = "2C4G"
  ssd_disk = "5"
  node_num = 3
  bill_type = 87
  project_id = 103800
  project_name = "测试部"
}




