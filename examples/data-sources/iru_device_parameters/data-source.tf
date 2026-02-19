data "iru_device_parameters" "example" {
  device_id = "your-device-uuid"
}

output "device_parameters" {
  value = data.iru_device_parameters.example.parameters
}
