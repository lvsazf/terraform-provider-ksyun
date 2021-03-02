package ksyun

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-ksyun/logger"
	"time"
)

func resourceKsyunNat() *schema.Resource {
	return &schema.Resource{
		Create: resourceKsyunNatCreate,
		Update: resourceKsyunNatUpdate,
		Read:   resourceKsyunNatRead,
		Delete: resourceKsyunNatDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:     schema.TypeString,
				ForceNew: false,
				Optional: true,
				Computed: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"nat_name": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     false,
				ValidateFunc: validateName,
				Computed:     true,
			},

			"nat_mode": {
				Type:         schema.TypeString,
				ForceNew:     true,
				Required:     true,
				ValidateFunc: validateNatMode,
			},

			"nat_type": {
				Type:         schema.TypeString,
				ForceNew:     true,
				Required:     true,
				ValidateFunc: validateNatType,
			},

			"nat_ip_number": {
				Type:         schema.TypeInt,
				ForceNew:     false,
				Optional:     true,
				Default:      1,
				ValidateFunc: validateNatIpNumber,
			},

			"band_width": {
				Type:         schema.TypeInt,
				ForceNew:     false,
				Required:     true,
				ValidateFunc: validateNatBandWidth,
			},

			"charge_type": {
				Type:         schema.TypeString,
				ForceNew:     false,
				Optional:     true,
				Default:      "DailyPaidByTransfer",
				ValidateFunc: validateNatChargeType,
			},

			"purchase_time": {
				Type:     schema.TypeInt,
				ForceNew: false,
				Optional: true,
			},

			"nat_ip_set": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"nat_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"nat_ip_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},

			"associate_nat_set": {
				Type:     schema.TypeList,
				Computed: true,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"subnet_id": {
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

func resourceKsyunNatCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*KsyunClient)
	conn := client.vpcconn

	var resp *map[string]interface{}
	var err error
	creates := []string{
		"project_id",
		"vpc_id",
		"nat_name",
		"nat_mode",
		"nat_type",
		"nat_ip_number",
		"band_width",
		"charge_type",
		"purchase_time",
	}
	createNat := make(map[string]interface{})
	for _, v := range creates {
		if v1, ok := d.GetOk(v); ok {
			vv := Downline2Hump(v)
			createNat[vv] = fmt.Sprintf("%v", v1)
		}
	}

	if d.Get("charge_type") == "Monthly" && (d.Get("purchase_time") == nil ||
		d.Get("purchase_time").(int) < 1 || d.Get("purchase_time").(int) > 15000) {
		return fmt.Errorf("purchase_time must set on charge_type is Monthly and in 1-15000 ")
	}
	action := "CreateNat"
	logger.Debug(logger.ReqFormat, action, createNat)
	resp, err = conn.CreateNat(&createNat)
	if err != nil {
		return fmt.Errorf("error on creating nat, %s", err)
	}
	if resp != nil {
		d.SetId((*resp)["NatId"].(string))
	}
	return resourceKsyunNatRead(d, meta)
}

func resourceKsyunNatRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*KsyunClient)
	conn := client.vpcconn

	readNat := make(map[string]interface{})
	readNat["NatId.1"] = d.Id()
	if pj, ok := d.GetOk("project_id"); ok {
		readNat["ProjectId.1"] = fmt.Sprintf("%v", pj)
	} else {
		projectErr := GetProjectInfo(&readNat, client)
		if projectErr != nil {
			return projectErr
		}
	}
	action := "DescribeNats"
	logger.Debug(logger.ReqFormat, action, readNat)
	resp, err := conn.DescribeNats(&readNat)
	if err != nil {
		return fmt.Errorf("error on reading nat %q, %s", d.Id(), err)
	}
	if resp != nil {
		items, ok := (*resp)["NatSet"].([]interface{})
		if !ok || len(items) == 0 {
			d.SetId("")
			return nil
		}
		SetResourceDataByResp(d, items[0], natKeys)
	}
	return nil
}

func resourceKsyunNatUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*KsyunClient)
	conn := client.vpcconn
	//d.Partial(true)
	attributeUpdate := false
	modifyNat := make(map[string]interface{})
	modifyNat["NatId"] = d.Id()

	if d.HasChange("nat_name") && !d.IsNewResource() {
		modifyNat["NatName"] = fmt.Sprintf("%v", d.Get("nat_name"))
		attributeUpdate = true
	}
	if d.HasChange("band_width") && !d.IsNewResource() {
		modifyNat["BandWidth"] = d.Get("band_width")
		attributeUpdate = true
	}
	if attributeUpdate {
		action := "ModifyNat"
		logger.Debug(logger.ReqFormat, action, modifyNat)
		resp, err := conn.ModifyNat(&modifyNat)
		logger.Debug(logger.AllFormat, action, modifyNat, resp, err)
		if err != nil {
			return fmt.Errorf("error on updating Subnet, %s", err)
		}
		//d.SetPartial("nat_name")
		//d.SetPartial("band_width")
	}
	//d.Partial(false)
	return resourceKsyunNatRead(d, meta)
}

func resourceKsyunNatDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*KsyunClient)
	conn := client.vpcconn
	deleteNat := make(map[string]interface{})
	deleteNat["NatId"] = d.Id()
	action := "DeleteNat"

	return resource.Retry(25*time.Minute, func() *resource.RetryError {
		logger.Debug(logger.ReqFormat, action, deleteNat)
		resp, err1 := conn.DeleteNat(&deleteNat)
		logger.Debug(logger.AllFormat, action, deleteNat, resp, err1)
		if err1 == nil {
			return nil
		} else if notFoundError(err1) {
			return nil
		} else {
			return resource.RetryableError(fmt.Errorf("error on  deleting nat %q, %s", d.Id(), err1))
		}
	})

}
