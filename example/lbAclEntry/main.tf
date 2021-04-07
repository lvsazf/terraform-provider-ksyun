# Specify the provider and access details
//provider "ksyun" {
//  access_key = "你的ak"
//  secret_key = "你的sk"
//  region = "eu-east-1"
//}

# Create Load Balancer Listener Acl Entry
resource "ksyun_lb_acl_entry" "default" {
  load_balancer_acl_id = "db946207-31ae-4b9c-abd6-bd3cfc276721"
  cidr_block = "192.168.11.2/32"
  rule_number = 20
  rule_action = "allow"
  protocol = "ip"
}