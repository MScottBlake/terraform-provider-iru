data "iru_device_lost_mode" "example" {
  device_id = "your-device-uuid"
}

output "lost_mode_status" {
  value = data.iru_device_lost_mode.example.status
}
