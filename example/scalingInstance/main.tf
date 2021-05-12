provider "ksyun" {
}

resource "ksyun_scaling_instance" "foo" {
  scaling_group_id = "572862736620621824"
  scaling_instance_id = "dcf8793d-5ce3-4565-820a-9cfa88b65d86"
  protected_from_detach = 0
}