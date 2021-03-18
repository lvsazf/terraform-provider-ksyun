# Specify the provider and access details
provider "ksyun" {
}

# Get  ScalingConfigurations
data "ksyun_scaling_groups" "default" {
  output_file="output_result"
  ids = ["569972821092278272"]
}
