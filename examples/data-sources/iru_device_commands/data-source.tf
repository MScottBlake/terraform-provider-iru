data "iru_device_commands" "example" {
  device_id = "your-device-uuid"
  limit     = 10
}

output "recent_commands" {
  value = data.iru_device_commands.example.commands
}
