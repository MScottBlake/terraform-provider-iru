data "iru_ade_device" "example" {
  id = "ade-device-uuid"
}

output "ade_device_model" {
  value = data.iru_ade_device.example.model
}
