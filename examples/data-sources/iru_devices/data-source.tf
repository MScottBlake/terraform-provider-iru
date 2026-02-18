data "iru_devices" "all" {}

output "total_devices" {
  value = length(data.iru_devices.all.devices)
}
