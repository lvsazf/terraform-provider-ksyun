provider "ksyun" {
}

resource "ksyun_scaling_policy" "foo" {
  scaling_group_id = "541241314798505984"
  threshold = 20
}