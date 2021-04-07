---
layout: "ksyun"
page_title: "Ksyun: ksyun_nat"
sidebar_current: "docs-ksyun-resource-nat"
description: |-
  Provides a Nat resource under VPC resource.
---

# ksyun_nat

Provides a Nat resource under VPC resource.

## Example Usage

```hcl
resource "ksyun_vpc" "test" {
  vpc_name = "ksyun-vpc-tf"
  cidr_block = "10.0.0.0/16"
}
resource "ksyun_nat" "foo" {
  nat_name = "ksyun-nat-tf"
  nat_mode = "Vpc"
  nat_type = "public"
  band_width = 1
  charge_type = "DailyPaidByTransfer"
  vpc_id = "${ksyun_vpc.test.id}"
}
```

## Argument Reference

The following arguments are supported:

* `vpc_id` - (Required) The id of the vpc.
* `nat_name` - (Optional) The Name of the Nat.  
* `nat_mode` - (Required) The Mode of the Nat. Valid Values: 'Vpc', 'Subnet'.
* `nat_type ` - (Required) The Type of Nat.Valid Values:'public'.
* `nat_ip_number` - (Optional) The Counts of Nat Ip, Default is 1.
* `band_width` - (Optional) The BandWidth of Nat Ip, Default is 1.
* `charge_type` - (Optional) The ChargeType of the Nat, Valid Values: 'DailyPaidByTransfer','Daily', 'Peak', 'PostPaidByAdvanced95Peak' .
* `purchase_time` - (Optional) The PurchaseTime of the Nat, in 1-36 ,If charge_type is Monthly this Field is Required.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `create_time` - The time of creation of nat, formatted in RFC3339 time string.

## Import

nat can be imported using the `id`, e.g.

```
$ terraform import ksyun_nat.example nat-abc123456
```