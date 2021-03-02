provider "ksyun" {
}


resource "ksyun_vpc" "test" {
  vpc_name = "ksyun-vpc-tf"
  cidr_block = "10.7.0.0/21"
}
resource "ksyun_nat" "foo" {
  nat_name = "ksyun-nat-tf"
  nat_mode = "Vpc"
  project_id = "0"
  nat_type = "public"
  band_width = 1
  charge_type = "DailyPaidByTransfer"
  vpc_id = "${ksyun_vpc.test.id}"
}