provider "ksyun" {
}

resource "ksyun_scaling_instance" "foo" {
  scaling_group_id = "541241314798505984"
  scaling_instance_id = "a4ef95c5-e8f1-43f8-912a-758f15064063"
  protected_from_detach = 1
}