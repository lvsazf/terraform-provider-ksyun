# Specify the provider and access details
provider "ksyun" {
}

# Get  ScalingConfigurations
data "ksyun_scaling_instances" "default" {
  output_file="output_result"
  scaling_group_id = "541241314798505984"
  protected_from_detach = 0
}
