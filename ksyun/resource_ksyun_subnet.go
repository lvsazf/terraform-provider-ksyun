package ksyun

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-ksyun/logger"
	"time"
)

func resourceKsyunSubnet() *schema.Resource {
	return &schema.Resource{
		Create: resourceKsyunSubnetCreate,
		Update: resourceKsyunSubnetUpdate,
		Read:   resourceKsyunSubnetRead,
		Delete: resourceKsyunSubnetDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"availability_zone": {
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
				Computed: true,
			},
			"subnet_name": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     false,
				Computed:     true,
				ValidateFunc: validateName,
			},

			"cidr_block": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateCIDRNetworkAddress,
			},

			"subnet_type": {
				Type:         schema.TypeString,
				ForceNew:     true,
				Required:     true,
				ValidateFunc: validateSubnetType,
			},

			"dhcp_ip_to": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Computed:     true,
				ValidateFunc: validateIpAddress,
			},

			"dhcp_ip_from": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Computed:     true,
				ValidateFunc: validateIpAddress,
			},

			"vpc_id": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},

			"gateway_ip": {
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
				Computed: true,
			},

			"dns1": {
				Type:         schema.TypeString,
				ForceNew:     false,
				Optional:     true,
				ValidateFunc: validateIpAddress,
				Computed:     true,
			},

			"dns2": {
				Type:         schema.TypeString,
				ForceNew:     false,
				Optional:     true,
				ValidateFunc: validateIpAddress,
				Computed:     true,
			},
			"network_acl_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"nat_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"availability_zone_name": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"subnet_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"availble_i_p_number": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceKsyunSubnetReq(req *map[string]interface{}) {
	for k, v := range *req {
		if k == "SubnetType" && v.(string) != "Reserve" {
			gw, start, end := getCidrIpRange((*req)["CidrBlock"].(string))
			if _, ok := (*req)["GatewayIp"]; !ok {
				(*req)["GatewayIp"] = gw
			}
			if _, ok := (*req)["DhcpIpFrom"]; !ok {
				(*req)["DhcpIpFrom"] = start
			}
			if _, ok := (*req)["DhcpIpTo"]; !ok {
				(*req)["DhcpIpTo"] = end
			}
		}
	}
}

func resourceKsyunSubnetCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*KsyunClient)
	conn := client.vpcconn
	r := resourceKsyunSubnet()

	var resp *map[string]interface{}
	var err error

	req, err := SdkRequestAutoMapping(d, r, false, nil, nil)
	if err != nil {
		return fmt.Errorf("error on creating ScalingPolicy, %s", err)
	}
	resourceKsyunSubnetReq(&req)

	action := "CreateSubnet"
	logger.Debug(logger.ReqFormat, action, req)
	resp, err = conn.CreateSubnet(&req)
	logger.Debug(logger.AllFormat, action, req, resp, err)
	if err != nil {
		return fmt.Errorf("error on creating Subnet, %s", err)
	}
	if resp != nil {
		Subnet := (*resp)["Subnet"].(map[string]interface{})
		d.SetId(Subnet["SubnetId"].(string))
	}
	return resourceKsyunSubnetRead(d, meta)
}

func resourceKsyunSubnetRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*KsyunClient)
	conn := client.vpcconn

	req := make(map[string]interface{})
	req["SubnetId.1"] = d.Id()
	action := "DescribeSubnets"
	logger.Debug(logger.ReqFormat, action, req)
	resp, err := conn.DescribeSubnets(&req)
	if err != nil {
		return fmt.Errorf("error on reading DescribeSubnets %q, %s", d.Id(), err)
	}
	if resp != nil {
		items, ok := (*resp)["SubnetSet"].([]interface{})
		if !ok || len(items) == 0 {
			d.SetId("")
			return nil
		}
		SdkResponseAutoResourceData(d, resourceKsyunSubnet(), items[0], nil)
	}
	return nil
}

func resourceKsyunSubnetUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*KsyunClient)
	conn := client.vpcconn
	r := resourceKsyunSubnet()

	var resp *map[string]interface{}
	var err error

	req, err := SdkRequestAutoMapping(d, r, true, nil, nil)
	if err != nil {
		return fmt.Errorf("error on creating ScalingPolicy, %s", err)
	}
	if len(req) > 0 {
		req["SubnetId"] = d.Id()
		action := "ModifySubnet"
		logger.Debug(logger.ReqFormat, action, req)
		resp, err = conn.ModifySubnet(&req)
		logger.Debug(logger.AllFormat, action, req, resp, err)
		if err != nil {
			return fmt.Errorf("error on updating Subnet, %s", err)
		}
	}
	return resourceKsyunSubnetRead(d, meta)
}

func resourceKsyunSubnetDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*KsyunClient)
	conn := client.vpcconn
	deleteSubnet := make(map[string]interface{})
	deleteSubnet["SubnetId"] = d.Id()
	action := "DeleteSubnet"

	return resource.Retry(25*time.Minute, func() *resource.RetryError {
		logger.Debug(logger.ReqFormat, action, deleteSubnet)
		resp, err1 := conn.DeleteSubnet(&deleteSubnet)
		logger.Debug(logger.AllFormat, action, deleteSubnet, *resp, err1)
		if err1 == nil || (err1 != nil && notFoundError(err1)) {
			return nil
		}
		if err1 != nil && inUseError(err1) {
			return resource.RetryableError(err1)
		}
		readSubnet := make(map[string]interface{})
		readSubnet["SubnetId.1"] = d.Id()
		action = "DescribeSubnets"
		logger.Debug(logger.ReqFormat, action, readSubnet)
		resp, err := conn.DescribeSubnets(&readSubnet)
		logger.Debug(logger.AllFormat, action, readSubnet, *resp, err)
		if err != nil && notFoundError(err1) {
			return nil
		}
		if err != nil {
			return resource.NonRetryableError(fmt.Errorf("error on  reading SubnetS when delete %q, %s", d.Id(), err))
		}
		itemset, ok := (*resp)["SubnetSet"]
		if !ok {
			return nil
		}
		item, ok := itemset.([]interface{})
		if !ok || len(item) == 0 {
			return nil
		}
		return resource.RetryableError(fmt.Errorf("error on  deleting SubnetS %q, %s", d.Id(), err1))
	})

}
