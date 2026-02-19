data "iru_device_apps" "example" {
  device_id = "your-device-uuid"
}

output "installed_apps" {
  value = data.iru_device_apps.example.apps
}
