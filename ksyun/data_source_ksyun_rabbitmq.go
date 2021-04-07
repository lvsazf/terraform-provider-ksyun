package ksyun

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-ksyun/logger"
	"strconv"
)

// instance List
func dataSourceKsyunRabbitmqInstances() *schema.Resource {
	return &schema.Resource{
		// Instance List Query Function
		Read: dataSourceRabbitmqInstancesRead,
		// Define input and output parameters
		Schema: map[string]*schema.Schema{
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"total_count": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"project_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"instance_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"subnet_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"vip": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"instances": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"instance_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"instance_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"project_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"instance_password": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vpc_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"subnet_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"engine_version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"project_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"bill_type": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"duration": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"mode": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"instance_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ssd_disk": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"node_num": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"availability_zone": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"engine": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"user_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"region": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"web_vip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"protocol": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"security_group_id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"network_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"product_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"create_date": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"expiration_date": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"product_what": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"mode_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"eip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"web_eip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"eip_egress": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"port": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceRabbitmqInstancesRead(d *schema.ResourceData, meta interface{}) error {

	var (
		allInstances []interface{}
		limit        = 100
		nextToken    string
	)

	readReq := make(map[string]interface{})
	filters := []string{"instance_id", "vpc_id", "instance_name", "vip", "project_id", "subnet_id"}
	for _, v := range filters {
		if value, ok := d.GetOk(v); ok {
			readReq[Downline2Hump(v)] = fmt.Sprintf("%v", value)
		}
	}
	readReq["limit"] = fmt.Sprintf("%v", limit)

	conn := meta.(*KsyunClient).rabbitmqconn

	for {
		if nextToken != "" {
			readReq["offset"] = nextToken
		}
		logger.Debug(logger.ReqFormat, "DescribeRabbitmqInstances", readReq)

		resp, err := conn.DescribeInstances(&readReq)
		if err != nil {
			return fmt.Errorf("error on reading instance list req(%v):%s", readReq, err)
		}
		logger.Debug(logger.RespFormat, "DescribeRabbitmqInstances", readReq, *resp)

		result, ok := (*resp)["Data"]
		if !ok {
			break
		}
		item, ok := result.(map[string]interface{})
		if !ok {
			break
		}
		items, ok := item["Instances"].([]interface{})
		if !ok {
			break
		}
		if len(items) < 1 {
			break
		}
		allInstances = append(allInstances, items...)
		if len(items) < limit {
			break
		}
		nextToken = strconv.Itoa(int(item["limit"].(float64)) + int(item["Offset"].(float64)))
	}

	values := GetSubSliceDByRep(allInstances, rabbitmqInstanceKeys)

	if err := dataSourceKscSave(d, "instances", []string{}, values); err != nil {
		return fmt.Errorf("error on save instance list, %s", err)
	}

	return nil
}
