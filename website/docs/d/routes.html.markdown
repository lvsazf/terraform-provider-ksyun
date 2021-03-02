---
layout: "ksyun"
page_title: "Ksyun: ksyun_routes"
sidebar_current: "docs-ksyun-datasource-routes"
description: |-
  Provides a list of Route resources in the current region.
---

# ksyun_routes

This data source provides a list of Route resources according to their Route ID, cidr and the VPC they belong to.

## Example Usage

```hcl
data "ksyun_routes" "default" {
  output_file="output_result"
  ids=[]
  vpc_ids=[]
  instance_ids=[]
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional) A list of Route IDs, all the Route resources belong to this region will be retrieved if the ID is `""`.
* `vpc_ids` - (Optional) A list of VPC id that the desired Route belongs to .
* `instance_ids` - (Optional) A list of the Route target id .  
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `routes` - It is a nested type which documented below.
* `total_count` - Total number of Route resources that satisfy the condition.

The attribute (`routes`) support the following:

* `id` - The ID of Route.
* `destination_cidr_block` - The cidr block of the desired Route.
* `route_type` - The route type of the desired Route.  
* `create_time` - The time of creation of Route, formatted in RFC3339 time string.