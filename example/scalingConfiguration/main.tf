provider "ksyun" {
}

resource "ksyun_scaling_configuration" "foo" {
  scaling_configuration_name = "tf-xym-test-1"
  image_id = "IMG-5465174a-6d71-4770-b8e1-917a0dd92466"
  instance_type = "N3.1B"
  password = "Aa123456"
  key_id = [
    "70afa609-610c-4596-9a82-2bbd28967435"]
  data_disks {
    disk_type = "EHDD"
    disk_size = 50
    delete_with_instance = true
  }

  data_disks {
    disk_type = "EHDD"
    disk_size = 100
    delete_with_instance = false
  }
  address_band_width = 2

}