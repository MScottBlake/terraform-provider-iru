data "iru_ade_integration_devices" "example" {
  ade_token_id = "ade-token-uuid"
}

output "integration_device_count" {
  value = length(data.iru_ade_integration_devices.example.devices)
}
