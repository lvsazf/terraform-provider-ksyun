---
layout: "ksyun"
page_title: "Ksyun: ksyun_nats"
sidebar_current: "docs-ksyun-datasource-nats"
description: |-
  Provides a list of Nat resources in the current region.
---

# ksyun_nats

This data source provides a list of Nat resources according to their Nat ID and the VPC they belong to.

## Example Usage

```hcl
data "ksyun_nats" "default" {
  output_file="output_result"
  ids=[]
  vpc_ids=[]
  project_ids=[]
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional) A list of Nat IDs, all the Nat resources belong to this region will be retrieved if the ID is `""`.
* `vpc_ids` - (Optional) A list of VPC id that the desired Nat belongs to .
* `project_ids` - (Optional) A list of Project id that the desired Nat belongs to .  
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `nats` - It is a nested type which documented below.
* `total_count` - Total number of Route resources that satisfy the condition.

The attribute (`nats`) support the following:

* `id` - The ID of Nat.
* `nat_type` - The nat type of the desired Nat.
* `vpc_id` - The VPC ID of the desired Nat belongs to. 
* `nat_name` - The nat name of the desired Nat.
* `nat_mode` - The nat mode of the desired Nat.
* `nat_ip_set` - The nat ip list of the desired Nat.
* `nat_ip_number` - The nat ip count of the desired Nat.  
* `band_width` - The nat ip band width of the desired Nat.  
* `associate_nat_set` - The subnet associate list of the desired Nat.
* `create_time` - The time of creation of Nat, formatted in RFC3339 time string.