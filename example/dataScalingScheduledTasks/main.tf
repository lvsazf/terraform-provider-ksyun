# Specify the provider and access details
provider "ksyun" {
}

# Get  ScalingConfigurations
data "ksyun_scaling_scheduled_tasks" "default" {
  output_file="output_result"
  scaling_group_id = "541241314798505984"
}
