provider "ksyun" {
}

resource "ksyun_scaling_notification" "foo" {
  scaling_group_id = "541241314798505984"
  scaling_notification_types = ["1","3"]
}
