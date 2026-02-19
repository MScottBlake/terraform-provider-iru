data "iru_device_library_items" "example" {
  device_id = "your-device-uuid"
}

output "device_library_items" {
  value = data.iru_device_library_items.example.library_items
}
