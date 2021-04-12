package ksyun

import (
	"errors"
	"fmt"
	"github.com/KscSDK/ksc-sdk-go/service/rabbitmq"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-ksyun/logger"
	"strings"
	"time"
)

func resourceKsyunRabbitmq() *schema.Resource {
	return &schema.Resource{
		Create: resourceRabbitmqInstanceCreate,
		Read:   resourceRabbitmqInstanceRead,
		Update: resourceRabbitmqInstanceUpdate,
		Delete: resourceRabbitmqInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(3 * time.Hour),
			Delete: schema.DefaultTimeout(3 * time.Hour),
			Update: schema.DefaultTimeout(3 * time.Hour),
		},
		Schema: map[string]*schema.Schema{
			"instance_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"instance_password": {
				Type:     schema.TypeString,
				Required: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"subnet_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"engine_version": {
				Type:     schema.TypeString,
				Required: true,
			},
			"project_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"project_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"bill_type": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"duration": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"mode": {
				Type:     schema.TypeString,
				Required: true,
			},
			"instance_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"ssd_disk": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"node_num": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"availability_zone": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enable_plugins": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"engine": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"user_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"region": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"status_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"vip": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"web_vip": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"protocol": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"security_group_id": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"network_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"product_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"create_date": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"expiration_date": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"product_what": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"mode_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"eip": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"web_eip": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"eip_egress": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"port": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}

}

func resourceRabbitmqInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	var (
		resp *map[string]interface{}
		err  error
		az   string
	)

	conn := meta.(*KsyunClient).rabbitmqconn
	r := resourceKsyunRabbitmq()
	req, err := SdkRequestAutoMapping(d, r, false, nil, nil)
	action := "CreateInstance"
	logger.Debug(logger.ReqFormat, action, req)
	if resp, err = conn.CreateInstance(&req); err != nil {
		return fmt.Errorf("error on creating instance: %s", err)
	}
	logger.Debug(logger.RespFormat, action, req, *resp)
	if resp != nil {
		d.SetId((*resp)["Data"].(map[string]interface{})["InstanceId"].(string))
	}
	stateConf := &resource.StateChangeConf{
		Pending:    []string{"creating"},
		Target:     []string{"running"},
		Refresh:    rabbitmqStateRefreshForCreateFunc(conn, az, d.Id(), []string{"running"}),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      20 * time.Second,
		MinTimeout: 1 * time.Minute,
	}
	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("error on create Instance: %s", err)
	}

	return resourceRabbitmqInstanceRead(d, meta)
}

func rabbitmqStateRefreshForCreateFunc(conn *rabbitmq.Rabbitmq, az string, instanceId string, target []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {

		queryReq := map[string]interface{}{"InstanceId": instanceId}
		logger.Debug(logger.ReqFormat, "DescribeRabbitmqInstance", queryReq)

		resp, err := conn.DescribeInstance(&queryReq)
		if err != nil {
			return nil, "", err
		}
		logger.Debug(logger.RespFormat, "DescribeRabbitmqInstance", queryReq, *resp)

		item, ok := (*resp)["Data"].(map[string]interface{})

		if !ok {
			return nil, "", fmt.Errorf("no instance information was queried. InstanceId:%s", instanceId)
		}
		status := item["Status"].(string)
		if status == "error" {
			return nil, "", fmt.Errorf("instance create error, status:%v", status)
		}

		for k, v := range target {
			if v == status {
				return resp, status, nil
			}
			if k == len(target)-1 {
				status = "creating"
			}
		}
		return resp, status, nil
	}
}

func resourceRabbitmqInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*KsyunClient).rabbitmqconn

	deleteReq := make(map[string]interface{})
	deleteReq["InstanceId"] = d.Id()

	logger.Debug(logger.ReqFormat, "DeleteRabbitmqInstance", deleteReq)

	err := resource.Retry(5*time.Minute, func() *resource.RetryError {
		_, err := conn.DeleteInstance(&deleteReq)
		if err != nil {
			return resource.RetryableError(errors.New(""))
		} else {
			return nil
		}
	})

	if err != nil {
		return fmt.Errorf("error on deleting instance %q, %s", d.Id(), err)
	}

	return resource.Retry(20*time.Minute, func() *resource.RetryError {

		queryReq := make(map[string]interface{})
		queryReq["InstanceId"] = d.Id()

		logger.Debug(logger.ReqFormat, "DescribeRabbitmqInstance", queryReq)
		resp, err := conn.DescribeInstance(&queryReq)
		logger.Debug(logger.RespFormat, "DescribeRabbitmqInstance", queryReq, resp)

		if err != nil {
			if strings.Contains(err.Error(), "InstanceNotFound") {
				return nil
			} else {
				return resource.NonRetryableError(err)
			}
		}

		_, ok := (*resp)["Data"].(map[string]interface{})

		if !ok {
			return nil
		}

		return resource.RetryableError(errors.New("deleting"))
	})
}

func resourceRabbitmqInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	d.Partial(true)
	defer d.Partial(false)
	conn := meta.(*KsyunClient).rabbitmqconn
	// rename
	if d.HasChange("instance_name") {
		d.SetPartial("instance_name")
		v, ok := d.GetOk("instance_name")
		if !ok {
			return fmt.Errorf("cann't change instance_name to empty string")
		}
		rename := make(map[string]interface{})
		rename["instanceId"] = d.Id()
		rename["instanceName"] = v.(string)
		logger.Debug(logger.ReqFormat, "RenameRabbitmqName", rename)
		resp, err := conn.Rename(&rename)
		if err != nil {
			return fmt.Errorf("error on rename instance %q, %s", d.Id(), err)
		}
		logger.Debug(logger.RespFormat, "RenameRabbitmqName", rename, *resp)
	}

	return resourceRabbitmqInstanceRead(d, meta)
}

func resourceRabbitmqInstanceRead(d *schema.ResourceData, meta interface{}) error {
	var (
		item map[string]interface{}
		resp *map[string]interface{}
		ok   bool
		err  error
	)

	conn := meta.(*KsyunClient).rabbitmqconn
	queryReq := make(map[string]interface{})
	queryReq["instanceId"] = d.Id()
	action := "DescribeInstance"
	logger.Debug(logger.ReqFormat, action, queryReq)
	if resp, err = conn.DescribeInstance(&queryReq); err != nil {
		return fmt.Errorf("error on reading instance %q, %s", d.Id(), err)
	}
	logger.Debug(logger.RespFormat, action, queryReq, *resp)
	if item, ok = (*resp)["Data"].(map[string]interface{}); !ok {
		return nil
	}

	SdkResponseAutoResourceData(d, resourceKsyunEip(), item, nil)

	return nil
}
