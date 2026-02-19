data "iru_device" "example" {
  id = "your-device-uuid"
}

output "device_name" {
  value = data.iru_device.example.device_name
}
