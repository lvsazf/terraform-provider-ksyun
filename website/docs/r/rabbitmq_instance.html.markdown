---
layout: "ksyun"
page_title: "Ksyun: ksyun_rabbitmq_instance"
sidebar_current: "docs-ksyun-resource-rabbitmq-instance"
description: |-
  Provides an replica set Rabbitmq resource.
---

# ksyun_rabbitmq_instance

Provides an replica set Rabbitmq resource.

## Example Usage

```hcl
resource "ksyun_rabbitmq_instance" "default" {
  availability_zone     = "cn-beijing-6a"
  instance_name         = "my_rabbitmq_instance"
  instance_password     = "Shiwo1101"
  instance_type         = "2C4G"
  vpc_id                = "VpcId"
  subnet_id             = "VnetId"
  mode                  = 1
  engine_version        = "3.7"
  ssd_disk              = "5"
  node_num              = 3
  bill_type             = 87
  project_id            = 103800
  project_name          = "测试部"
  duration              = ""
}
```

## Argument Reference

The following arguments are supported:

* `instance_name` - (Required) The name of instance, which contains 6-64 characters and only support Chinese, English, numbers, '-', '_'.
* `instance_password` - (Required) The administrator password of instance.
* `instance_type` - (Required) The class of instance cpu and memory.
* `ssd_disk` - (Required) The size of instance disk, measured in GB (GigaByte).
* `mode` - (Required) The mode of instance.
* `vpc_id` - (Required) The id of VPC linked to the instance.
* `subnet_id` - (Required) The id of subnet linked to the instance.
* `engine_version` - (Required) The version of instance engine.
* `bill_type` - (Required) Instance charge type,Valid values are 1 (Monthly), 87(UsageInstantSettlement).
* `duration` - (Optional) The duration of instance use, if `bill_type` is `1`, the duration is required.
* `node_num` - (Optional) the number of instance node, if not defined `node_num`, the instance will use `3`
* `project_id` - (Optional) The project id of instance belong, if not defined `project_id`, the instance will use `0`.
* `project_name` - (Optional) The project name of instance belong, if not defined `project_name`, the instance will use ``.
* `availability_zone` - (Required) Availability zone where instance is located.


