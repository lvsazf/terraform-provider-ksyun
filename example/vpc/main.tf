provider "ksyun" {
  region = "cn-beijing-6"
}

resource "ksyun_vpc" "test" {
  vpc_name   = "ksyun_vpc_tf_1"
  cidr_block = "10.1.0.0/16"
}
