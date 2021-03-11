# Specify the provider and access details
provider "ksyun" {
}
resource "ksyun_vpc" "test" {
  vpc_name   = "tf-example-vpc-02"
  cidr_block = "10.0.0.0/16"
}

resource "ksyun_subnet" "test" {
  subnet_name      = "tf-acc-subnet1"
  cidr_block = "10.0.1.0/24"
  subnet_type = "Reserve"
  availability_zone = "cn-beijing-6a"
  vpc_id  = "${ksyun_vpc.test.id}"
}
# Create Load Balancer
resource "ksyun_lb" "default" {
  vpc_id  = "${ksyun_vpc.test.id}"
  load_balancer_name = "tf-xun1"
  type = "internal"
  subnet_id = "${ksyun_subnet.test.id}"
  load_balancer_state = "start"
  private_ip_address = "10.0.1.2"
}
