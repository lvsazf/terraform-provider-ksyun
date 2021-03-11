---
layout: "ksyun"
page_title: "Ksyun: ksyun_scaling_group"
sidebar_current: "docs-ksyun-resource-scaling-group"
description: |-
  Provides a ScalingGroup resource.
---

# ksyun_scaling_group

Provides a ScalingGroup resource.

## Example Usage

```hcl
resource "ksyun_scaling_group" "foo" {
  subnet_id_set = [ksyun_subnet.foo.id]
  security_group_id = ksyun_security_group.foo.id
  scaling_configuration_id = ksyun_scaling_configuration.foo.id
  min_size = 0
  max_size = 2
  desired_capacity = 0
  status = "Active"
  slb_config_set  {
    slb_id = ksyun_lb.foo.id}
    listener_id = ksyun_lb_listener.foo.id
    server_port_set = [80]
  }
}
```

## Argument Reference

The following arguments are supported:

* `scaling_group_name` - (Optional) The Name of the desired ScalingGroup.
* `scaling_configuration_id` - (Required) The Scaling Configuration ID of the desired ScalingGroup set to.
* `scaling_configuration_name` - (Optional) The Scaling Configuration Name of the desired ScalingGroup set to.
* `min_size` - (Optional) The Min KEC instance size of the desired ScalingGroup set to.Valid Value 0-10.
* `max_size` - (Optional) The Min KEC instance size of the desired ScalingGroup set to.Valid Value 0-10.
* `desired_capacity` (Optional) - The Desire Capacity KEC instance count of the desired ScalingGroup set to.Valid Value 0-10.
* `instance_num` - (Optional) The KEC instance Number of the desired ScalingGroup set to.Valid Value 0-10.
* `remove_policy` - (Optional) The KEC instance remove policy of the desired ScalingGroup set to.Valid Values:'RemoveOldestInstance', 'RemoveNewestInstance'.
* `security_group_id` - (Required) The Security Group ID of the desired ScalingGroup set to.
* `status` -  (Optional) The Status of the desired ScalingGroup.Valid Values:'Active', 'UnActive'.
* `subnet_strategy` - (Optional) The Subnet Strategy of the desired ScalingGroup set to.Valid Values:'balanced-distribution', 'choice-first'.
* `subnet_id_set` - (Optional) The Subnet ID Set of the desired ScalingGroup set to.
* `slb_config_set` - (Optional) The SLB Config Set of the desired ScalingGroup set to.

The attribute (`slb_config_set`) support the following:

* `slb_id` - (Required) The SLB ID of the desired ScalingGroup set to.
* `listener_id` - (Required) The Listener ID of the desired ScalingGroup set to.
* `weight` - (Optional) The weight of the desired ScalingGroup set to.Valid Values 1-100.
* `server_port_set` - (Optional) The Server Port Set of the desired ScalingGroup set to.Valid Values 1-65535.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `create_time` - The time of creation of ScalingGroup, formatted in RFC3339 time string.

## Import

scalingGroup can be imported using the `id`, e.g.

```
$ terraform import ksyun_scaling_group.example scaling-group-abc123456
```