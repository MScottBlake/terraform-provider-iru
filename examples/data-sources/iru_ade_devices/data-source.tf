data "iru_ade_devices" "all" {}

data "iru_ade_devices" "by_blueprint" {
  blueprint_id = "your-blueprint-uuid"
}

output "all_ade_device_serials" {
  value = [for d in data.iru_ade_devices.all.devices : d.serial_number]
}
