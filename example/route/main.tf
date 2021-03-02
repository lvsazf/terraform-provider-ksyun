provider "ksyun" {
}
resource "ksyun_vpc" "test" {
  vpc_name   = "ksyun-vpc-tf"
  cidr_block = "10.7.0.0/21"
}
resource "ksyun_route" "foo" {
  destination_cidr_block = "10.0.0.0/16"
  route_type = "InternetGateway"
  vpc_id = "${ksyun_vpc.test.id}"
}
