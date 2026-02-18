data "iru_prism_filevault" "status" {}

output "encrypted_devices" {
  value = [
    for d in data.iru_prism_filevault.status.results : d.device_name 
    if d.status == true
  ]
}
