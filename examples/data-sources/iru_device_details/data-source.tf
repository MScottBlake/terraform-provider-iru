data "iru_device_details" "example" {
  device_id = "your-device-uuid"
}

output "device_full_name" {
  value = data.iru_device_details.example.device_name
}
