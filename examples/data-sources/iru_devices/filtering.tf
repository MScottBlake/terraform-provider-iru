# Example of filtering devices by platform and blueprint
data "iru_devices" "macs" {
  platform     = "Mac"
  blueprint_id = "your-blueprint-uuid"
}

output "mac_device_names" {
  value = [for d in data.iru_devices.macs.devices : d.device_name]
}

# Example of searching for a specific device by serial number
data "iru_devices" "by_serial" {
  serial_number = "C02XXXXXXX"
}
