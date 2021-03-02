# Specify the provider and access details
provider "ksyun" {
}

# Get  routes
data "ksyun_routes" "default" {
  output_file="output_result"
}
