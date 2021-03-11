---
layout: "ksyun"
page_title: "Ksyun: ksyun_scaling_scheduled_task"
sidebar_current: "docs-ksyun-resource-scaling-scheduled-task"
description: |-
  Provides a ScalingScheduledTask resource.
---

# ksyun_scaling_scheduled_task

Provides a ScalingScheduledTask resource.

## Example Usage

```hcl
resource "ksyun_scaling_scheduled_task" "foo" {
  scaling_group_id = "541241314798505984"
  start_time = "2021-05-01 12:00:00"
}

```

## Argument Reference

The following arguments are supported:

* `scaling_group_id` - (Required, ForceNew)The ScalingGroup ID of the desired ScalingScheduledTask belong to.
* `scaling_scheduled_task_name` - (Optional) The Name of the desired ScalingScheduledTask.
* `readjust_max_size` - (Optional) The Readjust Max Size of the desired ScalingScheduledTask.
* `readjust_min_size` - (Optional) The Readjust Min Size of the desired ScalingScheduledTask.
* `readjust_expect_size` - (Optional) The Readjust Expect Size of the desired ScalingScheduledTask.
* `start_time` - (Required) The Start Time of the desired ScalingScheduledTask.
* `end_time` -  (Optional) The End Time Operator of the desired ScalingScheduledTask.
* `recurrence` - (Optional) The Recurrence of the desired ScalingScheduledTask.
* `repeat_unit` - (Optional) The Repeat Unit of the desired ScalingScheduledTask.
* `repeat_cycle` - (Optional) The Repeat Cycle the desired ScalingScheduledTask.


## Attributes Reference

In addition to all arguments above, the following attributes are exported:


```
$ terraform import ksyun_scaling_scheduled_task.example scaling-scheduled-task-abc123456
```