data "iru_device_notes" "example" {
  device_id = "your-device-uuid"
}

output "device_notes" {
  value = data.iru_device_notes.example.notes
}
