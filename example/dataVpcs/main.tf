# Specify the provider and access details
provider "ksyun" {
}

data "ksyun_vpcs" "default" {
  output_file="output_result"
  ids=["dc52ca0b-b3d5-4849-8a8b-ba567e3836b6"]
}

