---
layout: "ksyun"
page_title: "Ksyun: ksyun_scaling_groups"
sidebar_current: "docs-ksyun-datasource-scaling-groups"
description: |-
  Provides a list of ScalingGroup resources in the current region.
---

# ksyun_scaling_groups

This data source provides a list of ScalingGroup resources .

## Example Usage

```hcl
data "ksyun_scaling_groups" "default" {
  output_file="output_result"
  vpc_id = "246b37be-5213-49da-a971-8748d73029c2"
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional) A list of ScalingGroup IDs, all the ScalingGroup resources belong to this region will be retrieved if the ID is `""`.
* `vpc_id` - (Optional) A list of vpc id that the desired ScalingGroup set to .
* `scaling_configuration_id` -  (Optional) A list of scaling configuration id that the desired ScalingGroup set to .
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `scaling_groups` - It is a nested type which documented below.
* `total_count` - Total number of ScalingGroup resources that satisfy the condition.

The attribute (`scaling_groups`) support the following:

* `id` - The ID of ScalingGroup.
* `scaling_group_name` - The Name of the desired ScalingGroup.
* `scaling_configuration_id` - The Scaling Configuration ID of the desired ScalingGroup set to. 
* `scaling_configuration_name` - The Scaling Configuration Name of the desired ScalingGroup set to.
* `min_size` - The Min KEC instance size of the desired ScalingGroup set to.
* `max_size` - The Min KEC instance size of the desired ScalingGroup set to.
* `desired_capacity` - The Desire Capacity KEC instance count of the desired ScalingGroup set to.
* `instance_num` - The KEC instance Number of the desired ScalingGroup set to. 
* `remove_policy` -The KEC instance remove policy of the desired ScalingGroup set to.
* `vpc_id` - The VPC ID of the desired ScalingGroup set to.
* `security_group_id` - The Security Group ID of the desired ScalingGroup set to.
* `status` - The Status of the desired ScalingGroup.
* `subnet_strategy` -The Subnet Strategy of the desired ScalingGroup set to.
* `subnet_id_set` - The Subnet ID Set of the desired ScalingGroup set to.
* `slb_config_set` - The SLB Config Set of the desired ScalingGroup set to.
* `create_time` - The time of creation of ScalingGroup, formatted in RFC3339 time string.

The attribute (`slb_config_set`) support the following:

* `slb_id` - The SLB ID of the desired ScalingGroup set to.
* `listener_id` - The Listener ID of the desired ScalingGroup set to.
* `weight` - The weight of the desired ScalingGroup set to.
* `server_port_set` - The Server Port Set of the desired ScalingGroup set to.