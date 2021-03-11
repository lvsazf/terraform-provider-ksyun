package ksyun

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-ksyun/logger"
	"strconv"
	"time"
)

func resourceKsyunScalingConfiguration() *schema.Resource {
	return &schema.Resource{
		Create: resourceKsyunScalingConfigurationCreate,
		Read:   resourceKsyunScalingConfigurationRead,
		Delete: resourceKsyunScalingConfigurationDelete,
		Update: resourceKsyunScalingConfigurationUpdate,
		Schema: map[string]*schema.Schema{

			"scaling_configuration_name": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "tf-scaling-config",
				ForceNew: false,
			},

			"image_id": {
				Type:     schema.TypeString,
				ForceNew: false,
				Required: true,
			},

			"instance_type": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "I1.1A",
				ForceNew: false,
			},

			"password": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},

			"system_disk_type": {
				Type:         schema.TypeString,
				Computed:     true,
				Optional:     true,
				ValidateFunc: validateKecSystemDiskType,
			},

			"system_disk_size": {
				Type:         schema.TypeInt,
				Computed:     true,
				Optional:     true,
				ValidateFunc: validateKecSystemDiskSize,
			},

			"data_disk_gb": {
				Type:     schema.TypeInt,
				Computed: true,
				Optional: true,
			},

			"data_disks": {
				Type:     schema.TypeList,
				Computed: true,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"disk_type": {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validateKecDataDiskType,
						},
						"disk_size": {
							Type:         schema.TypeInt,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validateKecDataDiskSize,
						},
						"delete_with_instance": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
					},
				},
			},

			"key_id": {
				Type:     schema.TypeSet,
				Computed: true,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set: schema.HashString,
			},

			"project_id": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},

			"keep_image_login": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},

			"instance_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"instance_name_suffix": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"user_data": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"instance_name_time_suffix": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},

			"need_monitor_agent": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validateKecInstanceAgent,
			},

			"need_security_agent": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validateKecInstanceAgent,
			},

			"charge_type": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"cpu": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"gpu": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"mem": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceKsyunScalingConfigurationExtra() map[string]SdkRequestMapping {
	var extra map[string]SdkRequestMapping
	extra = make(map[string]SdkRequestMapping)
	extra["data_disks"] = SdkRequestMapping{
		Field: "DataDisk.",
		FieldReqFunc: func(item interface{}, s string, m *map[string]interface{}) error {
			if arr, ok := item.([]interface{}); ok {
				for i, value := range arr {
					if d, ok := value.(map[string]interface{}); ok {
						if d["disk_type"] == "" || d["disk_size"] == 0 {
							return fmt.Errorf(" if set data_disks, disk_type and disk_size must set value ")
						}
						for k, v := range d {
							if k == "disk_type" {
								(*m)[s+strconv.Itoa(i+1)+".Type"] = v
							}
							if k == "disk_size" {
								(*m)[s+strconv.Itoa(i+1)+".Size"] = v
							}
							if k == "delete_with_instance" {
								(*m)[s+strconv.Itoa(i+1)+".DeleteWithInstance"] = v
							}
						}
					}
				}
			}
			return nil
		},
	}
	extra["key_id"] = SdkRequestMapping{
		Field: "KeyId.",
		FieldReqFunc: func(item interface{}, s string, m *map[string]interface{}) error {
			if x, ok := item.(*schema.Set); ok {
				for i, value := range (*x).List() {
					if d, ok := value.(string); ok {
						(*m)[s+strconv.Itoa(i+1)] = d
					}
				}
			}
			return nil
		},
	}
	return extra
}

func resourceKsyunScalingConfigurationCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*KsyunClient)
	conn := client.kecconn
	scalingConfiguration := resourceKsyunScalingConfiguration()

	var resp *map[string]interface{}
	var err error

	createScalingConfiguration, err := SdkRequestAutoMapping(d, scalingConfiguration, false, nil, resourceKsyunScalingConfigurationExtra())
	if err != nil {
		return fmt.Errorf("error on creating ScalingConfiguration, %s", err)
	}

	action := "CreateScalingConfiguration"
	logger.Debug(logger.ReqFormat, action, createScalingConfiguration)
	resp, err = conn.CreateScalingConfiguration(&createScalingConfiguration)
	if err != nil {
		return fmt.Errorf("error on creating ScalingConfiguration, %s", err)
	}
	if resp != nil {
		d.SetId((*resp)["ScalingConfigurationId"].(string))
	}
	return resourceKsyunScalingConfigurationRead(d, meta)
}

func resourceKsyunScalingConfigurationUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*KsyunClient)
	conn := client.kecconn
	scalingConfiguration := resourceKsyunScalingConfiguration()

	var err error

	modifyScalingConfiguration, err := SdkRequestAutoMapping(d, scalingConfiguration, true, nil, resourceKsyunScalingConfigurationExtra())
	if err != nil {
		return fmt.Errorf("error on modifying ScalingConfiguration, %s", err)
	}
	if len(modifyScalingConfiguration) > 0 {
		modifyScalingConfiguration["ScalingConfigurationId"] = d.Id()
		action := "ModifyScalingConfiguration"
		logger.Debug(logger.ReqFormat, action, modifyScalingConfiguration)
		_, err = conn.ModifyScalingConfiguration(&modifyScalingConfiguration)
		if err != nil {
			return fmt.Errorf("error on modifying ScalingConfiguration, %s", err)
		}
	}
	return resourceKsyunScalingConfigurationRead(d, meta)
}

func resourceKsyunScalingConfigurationRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*KsyunClient)
	conn := client.kecconn

	readScalingConfiguration := make(map[string]interface{})
	readScalingConfiguration["ScalingConfigurationId.1"] = d.Id()
	if pj, ok := d.GetOk("project_id"); ok {
		readScalingConfiguration["ProjectId.1"] = fmt.Sprintf("%v", pj)
	} else {
		projectErr := GetProjectInfo(&readScalingConfiguration, client)
		if projectErr != nil {
			return projectErr
		}
	}
	action := "DescribeScalingConfiguration"
	logger.Debug(logger.ReqFormat, action, readScalingConfiguration)
	resp, err := conn.DescribeScalingConfiguration(&readScalingConfiguration)
	if err != nil {
		return fmt.Errorf("error on reading ScalingConfiguration %q, %s", d.Id(), err)
	}
	if resp != nil {
		items, ok := (*resp)["ScalingConfigurationSet"].([]interface{})
		if !ok || len(items) == 0 {
			d.SetId("")
			return nil
		}
		SdkResponseAutoResourceData(d, resourceKsyunScalingConfiguration(), items[0], scalingConfigurationSpecialMapping())
	}
	return nil
}

func resourceKsyunScalingConfigurationDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*KsyunClient)
	conn := client.kecconn
	deleteScalingConfiguration := make(map[string]interface{})
	deleteScalingConfiguration["ScalingConfigurationId.1"] = d.Id()
	action := "DeleteScalingConfiguration"

	return resource.Retry(25*time.Minute, func() *resource.RetryError {
		logger.Debug(logger.ReqFormat, action, deleteScalingConfiguration)
		resp, err1 := conn.DeleteScalingConfiguration(&deleteScalingConfiguration)
		logger.Debug(logger.AllFormat, action, deleteScalingConfiguration, resp, err1)
		if err1 == nil {
			return nil
		} else if notFoundError(err1) {
			return nil
		} else {
			return resource.RetryableError(fmt.Errorf("error on  deleting ScalingConfiguration %q, %s", d.Id(), err1))
		}
	})

}
