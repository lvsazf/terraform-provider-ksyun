package ksyun

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-ksyun/logger"
	"time"
)

func resourceKsyunRoute() *schema.Resource {
	return &schema.Resource{
		Create: resourceKsyunRouteCreate,
		Read:   resourceKsyunRouteRead,
		Delete: resourceKsyunRouteDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"vpc_id": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},

			"destination_cidr_block": {
				Type:         schema.TypeString,
				ForceNew:     true,
				Required:     true,
				ValidateFunc: validateCIDRNetworkAddress,
			},

			"route_type": {
				Type:         schema.TypeString,
				ForceNew:     true,
				Required:     true,
				ValidateFunc: validateRouteType,
			},
			"tunnel_id": {
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
			},
			"vpc_peering_connection_id": {
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
			},
			"direct_connect_gateway_id": {
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
			},
			"vpn_tunnel_id": {
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
			},
			"next_hop_set": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"gateway_id": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"gateway_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},

			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceKsyunRouteCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*KsyunClient)
	conn := client.vpcconn

	var resp *map[string]interface{}
	var err error
	creates := []string{
		"vpc_id",
		"destination_cidr_block",
		"route_type",
		"tunnel_id",
		"instance_id",
		"vpc_peering_connection_id",
		"direct_connect_gateway_id",
		"vpn_tunnel_id",
	}
	createRoute := GetSdkParam(d, creates)

	if d.Get("route_type").(string) == "Tunnel" && (d.Get("tunnel_id") == nil ||
		d.Get("tunnel_id").(string) == "") {
		return fmt.Errorf("route_type is Tunnel ,Must set tunnel_id")
	}
	if d.Get("route_type").(string) == "Host" && (d.Get("instance_id") == nil ||
		d.Get("instance_id").(string) == "") {
		return fmt.Errorf("route_type is Host ,Must set instance_id")
	}
	if d.Get("route_type").(string) == "Peering" && (d.Get("vpc_peering_connection_id") == nil ||
		d.Get("instance_id").(string) == "") {
		return fmt.Errorf("route_type is Peering ,Must set vpc_peering_connection_id")
	}
	if d.Get("route_type").(string) == "DirectConnect" && (d.Get("direct_connect_gateway_id") == nil ||
		d.Get("instance_id").(string) == "") {
		return fmt.Errorf("route_type is DirectConnect ,Must set direct_connect_gateway_id")
	}
	if d.Get("route_type").(string) == "Vpn" && (d.Get("vpn_tunnel_id") == nil ||
		d.Get("instance_id").(string) == "") {
		return fmt.Errorf("route_type is Vpn ,Must set vpn_tunnel_id")
	}

	action := "CreateRoute"
	logger.Debug(logger.ReqFormat, action, createRoute)
	resp, err = conn.CreateRoute(&createRoute)
	logger.Debug(logger.AllFormat, action, createRoute, resp, err)
	if err != nil {
		return fmt.Errorf("error on creating Route, %s", err)
	}
	if resp != nil {
		d.SetId((*resp)["RouteId"].(string))
	}
	return resourceKsyunRouteRead(d, meta)
}

func resourceKsyunRouteRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*KsyunClient)
	conn := client.vpcconn

	readRoute := make(map[string]interface{})
	readRoute["RouteId.1"] = d.Id()
	action := "DescribeRoutes"
	logger.Debug(logger.ReqFormat, action, readRoute)
	resp, err := conn.DescribeRoutes(&readRoute)
	logger.Debug(logger.AllFormat, action, readRoute, resp, err)
	if err != nil {
		return fmt.Errorf("error on reading Route %q, %s", d.Id(), err)
	}
	if resp != nil {
		items, ok := (*resp)["RouteSet"].([]interface{})
		if !ok || len(items) == 0 {
			d.SetId("")
			return nil
		}
		SetResourceDataByResp(d, items[0], routeKeys)
	}
	return nil
}

func resourceKsyunRouteUpdate(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceKsyunRouteDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*KsyunClient)
	conn := client.vpcconn
	deleteRoute := make(map[string]interface{})
	deleteRoute["RouteId"] = d.Id()
	action := "DeleteRoute"

	return resource.Retry(25*time.Minute, func() *resource.RetryError {
		logger.Debug(logger.ReqFormat, action, deleteRoute)
		resp, err1 := conn.DeleteRoute(&deleteRoute)
		logger.Debug(logger.AllFormat, action, deleteRoute, resp, err1)
		if err1 == nil {
			return nil
		} else if notFoundError(err1) {
			return nil
		} else {
			return resource.RetryableError(fmt.Errorf("error on  deleting Route %q, %s", d.Id(), err1))
		}
	})
}
