data "iru_device_status" "example" {
  device_id = "your-device-uuid"
}

output "full_device_status" {
  value = data.iru_device_status.example
}
