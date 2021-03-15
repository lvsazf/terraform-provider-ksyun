# Specify the provider and access details
provider "ksyun" {
  region = "cn-beijing-6"
}

resource "ksyun_instance" "default" {
  image_id="e22d048a-e0b8-465e-a692-d3358981eeff"
  instance_type="S4.1A"
//  key_id=["6e3dee9c-291c-4647-bfc2-4c1eaa93fb80"]
//  system_disk{
//    disk_type="SSD3.0"
//    disk_size=30
//  }
  data_disks {
    disk_type = "SSD3.0"
    disk_size = 100
  }
  data_disk_gb=0
  #only support part type
  subnet_id="55dcbce0-d052-4556-aaf1-b17972d3f5e2"
  instance_password="Aa123456"
  keep_image_login=false
  charge_type="Daily"
  purchase_time=1
  security_group_id=["90877b57-cb42-4635-89fc-633d0355f46b","855b74e3-cc4c-476c-b08f-551af4009c35"]
  private_ip_address=""
  instance_name="xym-tf"
  instance_name_suffix=""
  sriov_net_support="false"
  project_id=0
  data_guard_id=""
  d_n_s1 =""
  d_n_s2 =""
  force_delete =true
  user_data=""
}
