# Specify the provider and access details
provider "ksyun" {
  access_key = "your ak"
  secret_key = "your sk"
  region = "cn-beijing-6"
}

data "ksyun_rabbitmqs" "default" {
  output_file = "output_result"
  project_id = ""
  instance_id = ""
  instance_name = ""
  subnet_id = ""
  vpc_id = ""
  vip = ""
}
