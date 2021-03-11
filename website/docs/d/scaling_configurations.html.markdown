---
layout: "ksyun"
page_title: "Ksyun: ksyun_scaling_configurations"
sidebar_current: "docs-ksyun-datasource-scaling-configurations"
description: |-
  Provides a list of ScalingConfiguration resources in the current region.
---

# ksyun_scaling_configurations

This data source provides a list of ScalingConfiguration resources.

## Example Usage

```hcl
data "ksyun_nats" "default" {
  output_file="output_result"
  ids=[]
  project_ids=[]
  scaling_configuration_name= "test"
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional) A list of ScalingConfiguration IDs, all the ScalingConfiguration resources belong to this region will be retrieved if the ID is `""`.
* `project_ids` - (Optional) A list of Project id that the desired ScalingConfiguration belongs to .
* `scaling_configuration_name` - (Optional) The Name of ScalingConfiguration .  
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `scaling_configurations` - It is a nested type which documented below.
* `total_count` - Total number of ScalingConfiguration resources that satisfy the condition.

The attribute (`scaling_configurations`) support the following:

* `id` - The ID of ScalingConfiguration.
* `scaling_configuration_name` - The Name of the desired ScalingConfiguration.
* `cpu` - The CPU core size of the desired ScalingConfiguration. 
* `mem` - The Memory GB size of the desired ScalingConfiguration.
* `data_disk_gb` - The Local Volume GB size of the desired ScalingConfiguration.
* `gpu` - The GPU core size the desired ScalingConfiguration.
* `image_id` - The System Image Id of the desired ScalingConfiguration.  
* `need_monitor_agent` - The Monitor agent flag desired ScalingConfiguration.
* `need_security_agent` - The Security agent flag desired ScalingConfiguration.
* `instance_type` - The KEC instance type of the desired ScalingConfiguration.
* `instance_name` - The KEC instance name of the desired ScalingConfiguration.
* `instance_name_suffix` - The kec instance name suffix of the desired ScalingConfiguration.
* `project_id` - The Project Id of the desired ScalingConfiguration belong to.
* `keep_image_login` - The Flag with image login set of the desired ScalingConfiguration.  
* `system_disk_type` - The subnet associate list of the desired ScalingConfiguration.
* `system_disk_size` - The subnet associate list of the desired ScalingConfiguration.
* `key_id` - The SSH key set of the desired ScalingConfiguration.  
* `data_disks` - It is a nested type which documented below.
* `instance_name_time_suffix` -  The kec instance name time suffix of the desired ScalingConfiguration.
* `user_data` - The user data of the desired ScalingConfiguration.  
* `create_time` - The time of creation of ScalingGroup, formatted in RFC3339 time string.

The attribute (`data_disks`) support the following:

* `disk_type` - The EBS Data Disk Type of the desired data_disk.
* `disk_size` - The EBS Data Disk Size of the desired data_disk.
* `delete_with_instance` - The Flag with delete EBS Data Disk when KEC Instance destroy.