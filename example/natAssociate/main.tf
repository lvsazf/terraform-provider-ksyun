provider "ksyun" {
}


resource "ksyun_vpc" "test" {
  vpc_name = "ksyun-vpc-tf"
  cidr_block = "10.0.0.0/16"
}

resource "ksyun_nat" "test" {
  nat_name = "ksyun-nat-tf"
  nat_mode = "Subnet"
  nat_type = "public"
  band_width = 1
  charge_type = "DailyPaidByTransfer"
  vpc_id = "${ksyun_vpc.test.id}"
}

resource "ksyun_subnet" "test" {
  subnet_name      = "tf-acc-subnet1"
  cidr_block = "10.0.5.0/24"
  subnet_type = "Normal"
  dhcp_ip_from = "10.0.5.2"
  dhcp_ip_to = "10.0.5.253"
  vpc_id  = "${ksyun_vpc.test.id}"
  gateway_ip = "10.0.5.1"
  dns1 = "198.18.254.41"
  dns2 = "198.18.254.40"
  availability_zone = "cn-beijing-6a"
}

resource "ksyun_subnet" "test1" {
  subnet_name      = "tf-acc-subnet1"
  cidr_block = "10.0.6.0/24"
  subnet_type = "Normal"
  dhcp_ip_from = "10.0.6.2"
  dhcp_ip_to = "10.0.6.253"
  vpc_id  = "${ksyun_vpc.test.id}"
  gateway_ip = "10.0.6.1"
  dns1 = "198.18.254.41"
  dns2 = "198.18.254.40"
  availability_zone = "cn-beijing-6b"
}

resource "ksyun_nat_associate" "test" {
  nat_id = "${ksyun_nat.test.id}"
  subnet_id = "${ksyun_subnet.test.id}"
}

resource "ksyun_nat_associate" "test1" {
  nat_id = "${ksyun_nat.test.id}"
  subnet_id = "${ksyun_subnet.test1.id}"
}