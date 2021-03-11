---
layout: "ksyun"
page_title: "Ksyun: ksyun_scaling_configuration"
sidebar_current: "docs-ksyun-resource-scaling-configuration"
description: |-
  Provides a ScalingConfiguration resource.
---

# ksyun_scaling_configuration

Provides a ScalingConfiguration resource.

## Example Usage

```hcl
resource "ksyun_scaling_configuration" "foo" {
  scaling_configuration_name = "tf-xym-test-1"
  image_id = "IMG-5465174a-6d71-4770-b8e1-917a0dd92466"
  instance_type = "N3.1B"
  password = "Aa123456"
}
```

## Argument Reference

The following arguments are supported:

* `scaling_configuration_name` - (Optional) The Name of the desired ScalingConfiguration.
* `data_disk_gb` - (Optional) The Local Volume GB size of the desired ScalingConfiguration.
* `image_id` - (Required) The System Image Id of the desired ScalingConfiguration.
* `need_monitor_agent` - (Optional) The Monitor agent flag desired ScalingConfiguration.
* `need_security_agent` - (Optional) The Security agent flag desired ScalingConfiguration.
* `instance_type` - (Required) The KEC instance type of the desired ScalingConfiguration.
* `instance_name` - (Optional) The KEC instance name of the desired ScalingConfiguration.
* `instance_name_suffix` - (Optional) The kec instance name suffix of the desired ScalingConfiguration.
* `project_id` - (Optional) The Project Id of the desired ScalingConfiguration belong to.
* `keep_image_login` - (Optional) The Flag with image login set of the desired ScalingConfiguration.
* `system_disk_type` - (Optional) The subnet associate list of the desired ScalingConfiguration.Valid Values:'Local_SSD', 'SSD3.0', 'EHDD'.
* `system_disk_size` - (Optional) The subnet associate list of the desired ScalingConfiguration.
* `key_id` - (Optional) The SSH key set of the desired ScalingConfiguration.
* `data_disks` - (Optional) It is a nested type which documented below.
* `instance_name_time_suffix` -  (Optional) The kec instance name time suffix of the desired ScalingConfiguration.
* `user_data` - (Optional) The user data of the desired ScalingConfiguration.

The attribute (`data_disks`) support the following:

* `disk_type` - (Required) The EBS Data Disk Type of the desired data_disk.Valid Values: 'SSD3.0', 'EHDD'.
* `disk_size` - (Required) The EBS Data Disk Size of the desired data_disk.
* `delete_with_instance` - (Optional) The Flag with delete EBS Data Disk when KEC Instance destroy.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `create_time` - The time of creation of ScalingConfiguration, formatted in RFC3339 time string.

## Import

scalingConfiguration can be imported using the `id`, e.g.

```
$ terraform import ksyun_scaling_configuration.example scaling-configuration-abc123456
```