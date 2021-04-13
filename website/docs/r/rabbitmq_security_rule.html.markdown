---
layout: "ksyun"
page_title: "Ksyun: ksyun_rabbitmq_security_rule"
sidebar_current: "docs-ksyun-resource-rabbitmq-security-rule"
description: |-
  Provides an Rabbitmq Security Rule resource.
---

# ksyun_rabbitmq_security_rule

Provides a Rabbitmq Security Rule resource.

## Example Usage

```hcl
resource "ksyun_rabbitmq_security_rule" "default" {
  instance_id = "InstanceId"
  cidrs = "192.168.10.1/32"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required) The id of instance.
* `cidrs` - (Required) The cidr block of source for the instance, multiple cidr separated by comma.

