---
layout: "ksyun"
page_title: "Ksyun: ksyun_scaling_instances"
sidebar_current: "docs-ksyun-datasource-scaling-instances"
description: |-
  Provides a list of ScalingInstance resources in the current region belong a ScalingGroup.
---

# ksyun_scaling_instances

This data source provides a list of ScalingInstance resources in a ScalingGroup.

## Example Usage

```hcl
data "ksyun_scaling_instances" "default" {
  output_file="output_result"
  scaling_group_id = "246b37be-5213-49da-a971-8748d73029c2"
}
```

## Argument Reference

The following arguments are supported:

* `scaling_group_id` -  (Required) A scaling group id that the desired ScalingInstance belong to .
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `scaling_instances` - It is a nested type which documented below.
* `total_count` - Total number of ScalingInstance resources that satisfy the condition.

The attribute (`scaling_instances`) support the following:

* `scaling_instance_id` - The KEC Instance ID of the desired ScalingInstance.
* `scaling_instance_name` - The KEC Instance Name of the desired ScalingInstance. 
* `health_status` - The Health Status of the desired ScalingInstance.
* `creation_type` - The Creation Type of the desired ScalingInstance.
* `protected_from_detach` - The KEC Instance Protected Model of the desired ScalingInstance.
* `add_time` - The time of creation of ScalingInstance, formatted in RFC3339 time string.