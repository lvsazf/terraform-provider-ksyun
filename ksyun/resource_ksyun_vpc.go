package ksyun

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-ksyun/logger"
	"time"
)

func resourceKsyunVpc() *schema.Resource {
	return &schema.Resource{
		Create: resourceKsyunVpcCreate,
		Update: resourceKsyunVpcUpdate,
		Read:   resourceKsyunVpcRead,
		Delete: resourceKsyunVpcDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"vpc_name": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     false,
				ValidateFunc: validateName,
			},

			"cidr_block": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validateCIDRNetworkAddress,
			},

			"is_default": {
				Type:     schema.TypeBool,
				ForceNew: true,
				Default:  false,
				Optional: true,
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceKsyunVpcCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*KsyunClient)
	conn := client.vpcconn
	r := resourceKsyunVpc()

	var resp *map[string]interface{}
	var err error

	req, err := SdkRequestAutoMapping(d, r, false, nil, nil)
	if err != nil {
		return fmt.Errorf("error on creating vpc, %s", err)
	}

	action := "CreateVpc"
	logger.Debug(logger.ReqFormat, action, req)
	resp, err = conn.CreateVpc(&req)
	if err != nil {
		return fmt.Errorf("error on creating vpc, %s", err)
	}
	if resp != nil {
		vpcId, err := getSdkValue("Vpc.VpcId", *resp)
		if err != nil {
			return fmt.Errorf("error on creating vpc, %s", err)
		}
		d.SetId(vpcId.(string))
	}
	return resourceKsyunVpcRead(d, meta)
}

func resourceKsyunVpcRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*KsyunClient)
	conn := client.vpcconn

	req := make(map[string]interface{})
	req["VpcId.1"] = d.Id()
	action := "DescribeVpcs"
	logger.Debug(logger.ReqFormat, action, req)
	resp, err := conn.DescribeVpcs(&req)
	if err != nil {
		return fmt.Errorf("error on reading vpc %q, %s", d.Id(), err)
	}
	if resp != nil {
		items, ok := (*resp)["VpcSet"].([]interface{})
		if !ok || len(items) == 0 {
			d.SetId("")
			return nil
		}
		SdkResponseAutoResourceData(d, resourceKsyunEip(), items[0], nil)
	}
	return nil
}

func resourceKsyunVpcUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*KsyunClient)
	conn := client.vpcconn
	r := resourceKsyunVpc()

	var err error

	req, err := SdkRequestAutoMapping(d, r, true, nil, nil)
	if err != nil {
		return fmt.Errorf("error on updating Vpc, %s", err)
	}
	if len(req) > 0 {
		req["VpcId"] = d.Id()
		action := "ModifyVpc"
		logger.Debug(logger.ReqFormat, action, req)
		_, err = conn.ModifyVpc(&req)
		if err != nil {
			return fmt.Errorf("error on modifying Vpc, %s", err)
		}
	}
	return resourceKsyunVpcRead(d, meta)
}

func resourceKsyunVpcDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*KsyunClient)
	conn := client.vpcconn
	req := make(map[string]interface{})
	req["VpcId"] = d.Id()
	action := "DeleteVpc"

	return resource.Retry(25*time.Minute, func() *resource.RetryError {
		logger.Debug(logger.ReqFormat, action, req)
		resp, err1 := conn.DeleteVpc(&req)
		logger.Debug(logger.AllFormat, action, req, resp, err1)
		if err1 == nil {
			return nil
		} else if notFoundError(err1) {
			return nil
		} else {
			return resource.RetryableError(fmt.Errorf("error on  deleting vpc %q, %s", d.Id(), err1))
		}
	})
}
