provider "ksyun" {
}

resource "ksyun_scaling_scheduled_task" "foo" {
  scaling_group_id = "541241314798505984"
  start_time = "2021-05-01 12:00:00"
}