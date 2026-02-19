data "iru_device_activity" "example" {
  device_id = "your-device-uuid"
  limit     = 20
}

output "device_activity" {
  value = data.iru_device_activity.example.activity
}
