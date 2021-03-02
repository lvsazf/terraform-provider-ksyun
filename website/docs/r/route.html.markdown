---
layout: "ksyun"
page_title: "Ksyun: ksyun_route"
sidebar_current: "docs-ksyun-resource-route"
description: |-
  Provides a route resource under VPC resource.
---

# ksyun_route

Provides a route resource under VPC resource.

## Example Usage

```hcl
resource "ksyun_vpc" "example" {
  vpc_name   = "tf-example-vpc-01"
  cidr_block = "10.0.0.0/16"
}

resource "ksyun_route" "example" {
  destination_cidr_block = "10.0.0.0/16"
  route_type = "InternetGateway"
  vpc_id = "${ksyun_vpc.example.id}"
}
```

## Argument Reference

The following arguments are supported:

* `vpc_id` - (Required) The id of the vpc.
* `destination_cidr_block` - (Required) The CIDR block assigned to the route.
* `route_type ` - (Required) The type of route.Valid Values:'InternetGateway', 'Tunnel', 'Host', 'Peering', 'DirectConnect', 'Vpn'.
* `TunnelId` - (Optional) The id of the tunnel If route_type is Tunnel, This Field is Required.
* `InstanceId` - (Optional) The id of the VM , If route_type is Host, This Field is Required.
* `VpcPeeringConnectionId` - (Optional) The id of the Peering , If route_type is Peering, This Field is Required.
* `DirectConnectGatewayId` - (Optional) The id of the DirectConnectGateway , If route_type is DirectConnect, This Field is Required.
* `VpnTunnelId` - (Optional) The id of the Vpn , If route_type is Vpn, This Field is Required.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `create_time` - The time of creation of route, formatted in RFC3339 time string.

## Import

route can be imported using the `id`, e.g.

```
$ terraform import ksyun_route.example route-abc123456
```