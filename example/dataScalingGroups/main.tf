# Specify the provider and access details
provider "ksyun" {
}

# Get  ScalingConfigurations
data "ksyun_scaling_groups" "default" {
  output_file="output_result"
  vpc_id = "246b37be-5213-49da-a971-8748d73029c2"
}
