package ksyun

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-ksyun/logger"
	"strings"
	"time"
)

func resourceKsyunScalingInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceKsyunScalingInstanceCreate,
		Read:   resourceKsyunScalingInstanceRead,
		Delete: resourceKsyunScalingInstanceDelete,
		Update: resourceKsyunScalingInstanceUpdate,
		Schema: map[string]*schema.Schema{
			"scaling_group_id": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},

			"scaling_instance_id": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},

			"protected_from_detach": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validateKecInstanceAgent,
			},

			"scaling_instance_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"health_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"add_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"creation_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceKsyunScalingInstanceExtra() map[string]SdkRequestMapping {
	var extra map[string]SdkRequestMapping
	extra = make(map[string]SdkRequestMapping)
	extra["scaling_instance_id"] = SdkRequestMapping{
		Field: "ScalingInstanceId.1",
	}
	return extra
}

func resourceKsyunScalingInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*KsyunClient)
	conn := client.kecconn
	r := resourceKsyunScalingInstance()

	var err error

	var only map[string]SdkReqTransform
	only = map[string]SdkReqTransform{
		"scaling_group_id":    {},
		"scaling_instance_id": {Type: TransformSingleN},
	}

	req, err := SdkRequestAutoMapping(d, r, false, only, nil)
	if err != nil {
		return fmt.Errorf("error on creating ScalingInstance, %s", err)
	}

	action := "AttachInstance"
	logger.Debug(logger.ReqFormat, action, req)
	_, err = conn.AttachInstance(&req)
	if err != nil {
		return fmt.Errorf("error on creating ScalingInstance, %s", err)
	}
	d.SetId(d.Get("scaling_instance_id").(string) + ":" + d.Get("scaling_group_id").(string))

	if _, ok := d.GetOk("protected_from_detach"); ok {
		req["ProtectedFromDetach"] = d.Get("protected_from_detach").(int)
		action = "SetKvmProtectedDetach"
		logger.Debug(logger.ReqFormat, action, req)
		_, err = conn.SetKvmProtectedDetach(&req)
		if err != nil {
			return fmt.Errorf("error on creating ScalingInstance, %s", err)
		}
	}

	return resourceKsyunScalingInstanceRead(d, meta)
}
func resourceKsyunScalingInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*KsyunClient)
	conn := client.kecconn
	r := resourceKsyunScalingInstance()

	var err error

	var only map[string]SdkReqTransform
	only = map[string]SdkReqTransform{
		"protected_from_detach": {},
		"scaling_group_id":      {forceUpdateParam: true},
		"scaling_instance_id":   {forceUpdateParam: true, Type: TransformSingleN},
	}

	req, err := SdkRequestAutoMapping(d, r, true, only, nil)
	if err != nil {
		return fmt.Errorf("error on updating ScalingInstance, %s", err)
	}

	//zero process
	if _, ok := req["ProtectedFromDetach"]; !ok {
		req["ProtectedFromDetach"] = 0
	}

	action := "SetKvmProtectedDetach"
	logger.Debug(logger.ReqFormat, action, req)
	_, err = conn.SetKvmProtectedDetach(&req)
	if err != nil {
		return fmt.Errorf("error on updating ScalingInstance, %s", err)
	}

	return resourceKsyunScalingInstanceRead(d, meta)
}

func resourceKsyunScalingInstanceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*KsyunClient)
	conn := client.kecconn

	req := make(map[string]interface{})
	req["ScalingGroupId"] = strings.Split(d.Id(), ":")[1]
	req["ScalingInstanceId.1"] = strings.Split(d.Id(), ":")[0]
	action := "DescribeScalingInstance"
	logger.Debug(logger.ReqFormat, action, req)
	resp, err := conn.DescribeScalingInstance(&req)
	if err != nil {
		return fmt.Errorf("error on reading ScalingInstance %q, %s", d.Id(), err)
	}
	if resp != nil {
		items, ok := (*resp)["ScalingInstanceSet"].([]interface{})
		if !ok || len(items) == 0 {
			d.SetId("")
			return nil
		}
		SdkResponseAutoResourceData(d, resourceKsyunScalingInstance(), items[0], scalingInstanceSpecialMapping())
	}
	return nil
}

func resourceKsyunScalingInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*KsyunClient)
	conn := client.kecconn
	req := make(map[string]interface{})
	req["ScalingGroupId"] = strings.Split(d.Id(), ":")[1]
	req["ScalingInstanceId.1"] = strings.Split(d.Id(), ":")[0]
	action := "DetachInstance"
	otherErrorRetry := 10

	return resource.Retry(25*time.Minute, func() *resource.RetryError {
		logger.Debug(logger.ReqFormat, action, req)
		resp, err1 := conn.DetachInstance(&req)
		logger.Debug(logger.AllFormat, action, req, resp, err1)
		if err1 == nil {
			return nil
		} else if notFoundError(err1) {
			return nil
		} else {
			return OtherErrorProcess(&otherErrorRetry, fmt.Errorf("error on  deleting ScalingInstance %q, %s", d.Id(), err1))
		}
	})

}
