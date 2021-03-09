# Specify the provider and access details
provider "ksyun" {
}

# Get  ScalingConfigurations
data "ksyun_scaling_configurations" "default" {
  output_file="output_result"
}
