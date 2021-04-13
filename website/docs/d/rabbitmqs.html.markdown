---
layout: "ksyun"
page_title: "Ksyun: ksyun_rabbitmqs"
sidebar_current: "docs-ksyun-datasource-mongodbs"
description: |-
  Provides a list of Rabbitmq resources in the current region.
---

# ksyun_rabbitmqs

This data source provides a list of Rabbitmq resources according to their name, Instance ID, Subnet ID, VPC ID and the Project ID they belong to .

## Example Usage

```hcl
# Get  rabbitmqs
data "ksyun_rabbitmqs" "default" {
  output_file = "output_result"
  project_id = ""
  instance_id = ""
  instance_name = ""
  subnet_id = ""
  vpc_id = ""
  vip = ""
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Optional)  The name of Rabbitmq, all the Rabbitmqs belong to this region will be retrieved if the name is `""`.
* `instance_id` - (Optional)  The id of Rabbitmq, all the Rabbitmqs belong to this region will be retrieved if the instance_id is `""`.
* `project_id` - (Optional)  The project instance belongs to.
* `vpc_id` - (Optional)   Used to retrieve instances belong to specified VPC .
* `subnet_id` - (Optional) The ID of subnet. the instance will use the subnet in the current region.
* `vip` - (Optional) The vip of instances. 
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `instances` - It is a nested type which documented below.
* `total_count` - Total number of Rabbitmqs that satisfy the condition.

