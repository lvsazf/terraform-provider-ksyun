package ksyun

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-ksyun/logger"
	"strings"
	"time"
)

func resourceKsyunLoadBalancerAclEntry() *schema.Resource {
	return &schema.Resource{
		Create: resourceKsyunLoadBalancerAclEntryCreate,
		Delete: resourceKsyunLoadBalancerAclEntryDelete,
		Update: resourceKsyunLoadBalancerAclEntryUpdate,
		Read:   resourceKsyunLoadBalancerAclEntryRead,
		Schema: map[string]*schema.Schema{
			"load_balancer_acl_entry_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"load_balancer_acl_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"cidr_block": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"rule_number": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      1,
				ValidateFunc: validation.IntBetween(1, 32766),
			},
			"rule_action": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "allow",
				ValidateFunc: validation.StringInSlice([]string{
					"allow",
					"deny",
				}, false),
			},
			"protocol": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  "ip",
				ValidateFunc: validation.StringInSlice([]string{
					"ip",
				}, false),
			},
		},
	}
}

func resourceKsyunLoadBalancerAclEntryGet(d *schema.ResourceData, meta interface{}) (interface{}, error) {
	client := meta.(*KsyunClient)
	conn := client.slbconn

	req := make(map[string]interface{})
	ids := strings.Split(d.Id(), ":")
	if len(ids) != 2 {
		return nil, fmt.Errorf("error id:%v", d.Id())
	}
	req["LoadBalancerAclId.1"] = ids[0]
	action := "DescribeLoadBalancerAcls"
	logger.Debug(logger.ReqFormat, action, req)
	resp, err := conn.DescribeLoadBalancerAcls(&req)
	if err != nil {
		return nil, fmt.Errorf("error on reading LoadBalancerAcls %q, %s", d.Id(), err)
	}
	var result interface{}
	if resp != nil {
		items, ok := (*resp)["LoadBalancerAclSet"].([]interface{})
		if !ok || len(items) == 0 {
			return nil, nil
		}
		entries, ok := items[0].(map[string]interface{})["LoadBalancerAclEntrySet"].([]interface{})
		if !ok || len(entries) == 0 {
			return nil, nil
		}
		for _, entry := range entries {
			loadBalancerAclEntryId := entry.(map[string]interface{})["LoadBalancerAclEntryId"]
			if loadBalancerAclEntryId.(string) == ids[1] {
				result = entry
				break
			}
		}

	}
	return result, nil
}

func resourceKsyunLoadBalancerAclEntryRead(d *schema.ResourceData, meta interface{}) error {
	result, _ := resourceKsyunLoadBalancerAclEntryGet(d, meta)
	if result == nil {
		d.SetId("")
	} else {
		SdkResponseAutoResourceData(d, resourceKsyunLoadBalancerAclEntry(), result, nil)
	}
	return nil
}

func resourceKsyunLoadBalancerAclEntryCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*KsyunClient)
	conn := client.slbconn
	r := resourceKsyunLoadBalancerAclEntry()

	var resp *map[string]interface{}
	var err error

	req, err := SdkRequestAutoMapping(d, r, false, nil, nil)
	if err != nil {
		return fmt.Errorf("error on creating LoadBalancerAclEntry, %s", err)
	}

	action := "CreateLoadBalancerAclEntry"
	logger.Debug(logger.ReqFormat, action, req)
	resp, err = conn.CreateLoadBalancerAclEntry(&req)
	if err != nil {
		return fmt.Errorf("error on creating LoadBalancerAclEntry, %s", err)
	}
	if resp != nil {
		loadBalancerAclEntryId, err := getSdkValue("LoadBalancerAclEntry.LoadBalancerAclEntryId", *resp)
		if err != nil {
			return fmt.Errorf("error on creating LoadBalancerAclEntry, %s", err)
		}
		if loadBalancerAclEntryId == nil {
			return fmt.Errorf("error on creating LoadBalancerAclEntry,loadBalancerAclEntryId not get ")
		}
		d.SetId(fmt.Sprintf("%s:%s", d.Get("load_balancer_acl_id").(string), loadBalancerAclEntryId.(string)))
	}
	return resourceKsyunLoadBalancerAclEntryRead(d, meta)
}

func resourceKsyunLoadBalancerAclEntryUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*KsyunClient)
	conn := client.slbconn
	ids := strings.Split(d.Id(), ":")
	if len(ids) != 2 {
		return fmt.Errorf("error id:%v", d.Id())
	}
	r := resourceKsyunLoadBalancerAclEntry()

	var err error

	req, err := SdkRequestAutoMapping(d, r, true, nil, nil)
	if err != nil {
		return fmt.Errorf("error on updating LoadBalancerAclEntry, %s", err)
	}

	if len(req) > 0 {
		req["LoadBalancerAclEntryId"] = ids[1]
		action := "ModifyLoadBalancerAclEntry"
		logger.Debug(logger.ReqFormat, action, req)
		_, err = conn.ModifyLoadBalancerAclEntry(&req)
		if err != nil {
			return fmt.Errorf("error on updating LoadBalancerAclEntry, %s", err)
		}
	}
	return resourceKsyunLoadBalancerAclEntryRead(d, meta)
}

func resourceKsyunLoadBalancerAclEntryDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*KsyunClient)
	conn := client.slbconn
	ids := strings.Split(d.Id(), ":")
	if len(ids) != 2 {
		return fmt.Errorf("error id:%v", d.Id())
	}
	req := make(map[string]interface{})
	req["LoadBalancerAclEntryId"] = ids[1]
	req["LoadBalancerAclId"] = ids[0]
	return resource.Retry(25*time.Minute, func() *resource.RetryError {
		action := "DeleteLoadBalancerAclEntry"
		logger.Debug(logger.ReqFormat, action, req)
		_, err1 := conn.DeleteLoadBalancerAclEntry(&req)
		if err1 == nil {
			return nil
		} else {
			//if delete error try to read and retry
			result, err := resourceKsyunLoadBalancerAclEntryGet(d, meta)
			if err != nil {
				return resource.NonRetryableError(fmt.Errorf("error on  reading LoadBalancerAclEntry when delete %q, %s", d.Id(), err))
			}
			if result == nil {
				return nil
			}
			return resource.RetryableError(fmt.Errorf("error on  deleting LoadBalancerAclEntry %q, %s", d.Id(), err1))
		}
	})
}
