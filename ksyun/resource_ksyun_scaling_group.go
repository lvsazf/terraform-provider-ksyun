package ksyun

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-ksyun/logger"
	"strconv"
	"time"
)

func resourceKsyunScalingGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceKsyunScalingGroupCreate,
		Read:   resourceKsyunScalingGroupRead,
		Delete: resourceKsyunScalingGroupDelete,
		Update: resourceKsyunScalingGroupUpdate,
		Schema: map[string]*schema.Schema{

			"scaling_group_name": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "tf-scaling-group",
			},
			"scaling_configuration_id": {
				Type:     schema.TypeString,
				Required: true,
			},

			"min_size": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validateKecScalingGroupSize,
			},

			"max_size": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validateKecScalingGroupSize,
			},

			"desired_capacity": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validateKecScalingGroupDesiredCapacity,
			},

			"remove_policy": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "RemoveOldestInstance",
				ValidateFunc: validateKecScalingGroupRemovePolicy,
			},

			"subnet_strategy": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "balanced-distribution",
				ValidateFunc: validateKecScalingGroupSubnetStrategy,
			},

			"subnet_id_set": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set: schema.HashString,
			},

			"security_group_id": {
				Type:     schema.TypeString,
				Required: true,
			},

			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "Active",
				ValidateFunc: validateKecScalingGroupStatus,
			},

			"slb_config_set": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"slb_id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"listener_id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"weight": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  20,
						},
						"server_port_set": {
							Type:     schema.TypeSet,
							Optional: true,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeInt,
							},
							Set: schema.HashInt,
						},
					},
				},
			},

			"scaling_configuration_name": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"instance_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"vpc_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceKsyunScalingGroupExtra() map[string]SdkRequestMapping {
	var extra map[string]SdkRequestMapping
	extra = make(map[string]SdkRequestMapping)
	extra["slb_config_set"] = SdkRequestMapping{
		Field: "Slb.",
		FieldReqFunc: func(item interface{}, s string, source string, m *map[string]interface{}) error {
			if arr, ok := item.([]interface{}); ok {
				for i, value := range arr {
					if d, ok := value.(map[string]interface{}); ok {
						for k, v := range d {
							if k == "slb_id" {
								(*m)[s+strconv.Itoa(i+1)+".Id"] = v
							}
							if k == "listener_id" {
								(*m)[s+strconv.Itoa(i+1)+".ListenerId"] = v
							}
							if k == "weight" {
								(*m)[s+strconv.Itoa(i+1)+".Weight"] = v
							}
							if k == "server_port_set" {
								if x, ok := v.(*schema.Set); ok {
									for j, v1 := range (*x).List() {
										(*m)[s+strconv.Itoa(i+1)+".ServerPort."+strconv.Itoa(j+1)] = v1
									}
								}
							}
						}
					}
				}
			}
			return nil
		},
	}
	extra["subnet_id_set"] = SdkRequestMapping{
		Field: "SubnetId.",
		FieldReqFunc: func(item interface{}, s string, source string, m *map[string]interface{}) error {
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

func resourceKsyunScalingGroupCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*KsyunClient)
	conn := client.kecconn
	r := resourceKsyunScalingGroup()

	var resp *map[string]interface{}
	var err error

	req, err := SdkRequestAutoMapping(d, r, false, nil, resourceKsyunScalingGroupExtra())
	if err != nil {
		return fmt.Errorf("error on creating ScalingGroup, %s", err)
	}
	//zero process
	if _, ok := req["MinSize"]; !ok {
		req["MinSize"] = 0
	}
	if _, ok := req["MaxSize"]; !ok {
		req["MaxSize"] = 0
	}
	if _, ok := req["DesiredCapacity"]; !ok {
		req["DesiredCapacity"] = 0
	}

	action := "CreateScalingGroup"
	logger.Debug(logger.ReqFormat, action, req)
	resp, err = conn.CreateScalingGroup(&req)
	if err != nil {
		return fmt.Errorf("error on creating ScalingGroup, %s", err)
	}
	if resp != nil {
		d.SetId((*resp)["ReturnSet"].(map[string]interface{})["ScalingGroupId"].(string))
	}
	//set status
	if v, ok := d.GetOk("status"); ok {
		if v == "UnActive" {
			req = make(map[string]interface{})
			req["ScalingGroupId"] = d.Id()
			_, err = conn.DisableScalingGroup(&req)
			if err != nil {
				return fmt.Errorf("error on creating ScalingGroup, %s", err)
			}
		}
	}
	return resourceKsyunScalingGroupRead(d, meta)
}

func resourceKsyunScalingGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*KsyunClient)
	conn := client.kecconn
	r := resourceKsyunScalingGroup()

	var err error

	req, err := SdkRequestAutoMapping(d, r, true, nil, resourceKsyunScalingGroupExtra())
	if err != nil {
		return fmt.Errorf("error on modifying ScalingGroup, %s", err)
	}
	if len(req) > 0 {
		req1 := make(map[string]interface{})
		req1["ScalingGroupId"] = d.Id()
		if v, ok := req["Status"]; ok {
			if v == "Active" {
				_, err = conn.EnableScalingGroup(&req1)
				if err != nil {
					return fmt.Errorf("error on modifying ScalingGroup, %s", err)
				}
			} else {
				_, err = conn.DisableScalingGroup(&req1)
				if err != nil {
					return fmt.Errorf("error on modifying ScalingGroup, %s", err)
				}
			}
			delete(req, "Status")
		}
		if len(req) > 0 {
			req["ScalingGroupId"] = d.Id()
			action := "ModifyScalingGroup"
			logger.Debug(logger.ReqFormat, action, req)
			_, err = conn.ModifyScalingGroup(&req)
			if err != nil {
				return fmt.Errorf("error on modifying ScalingGroup, %s", err)
			}
		}

	}
	return resourceKsyunScalingGroupRead(d, meta)
}

func resourceKsyunScalingGroupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*KsyunClient)
	conn := client.kecconn

	req := make(map[string]interface{})
	req["ScalingGroupId.1"] = d.Id()
	action := "DescribeScalingGroup"
	logger.Debug(logger.ReqFormat, action, req)
	resp, err := conn.DescribeScalingGroup(&req)
	if err != nil {
		return fmt.Errorf("error on reading ScalingGroup %q, %s", d.Id(), err)
	}
	if resp != nil {
		items, ok := (*resp)["ScalingGroupSet"].([]interface{})
		if !ok || len(items) == 0 {
			d.SetId("")
			return nil
		}
		SdkResponseAutoResourceData(d, resourceKsyunScalingGroup(), items[0], nil)
	}
	return nil
}

func resourceKsyunScalingGroupDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*KsyunClient)
	conn := client.kecconn
	req := make(map[string]interface{})
	req["ScalingGroupId.1"] = d.Id()
	action := "DeleteScalingGroup"

	return resource.Retry(25*time.Minute, func() *resource.RetryError {
		logger.Debug(logger.ReqFormat, action, req)
		resp, err1 := conn.DeleteScalingGroup(&req)
		logger.Debug(logger.AllFormat, action, req, resp, err1)
		if err1 == nil {
			return nil
		} else if notFoundError(err1) {
			return nil
		} else {
			return resource.RetryableError(fmt.Errorf("error on  deleting ScalingGroup %q, %s", d.Id(), err1))
		}
	})

}
