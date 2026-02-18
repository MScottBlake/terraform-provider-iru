data "iru_prism_device_information" "example" {}

output "device_details" {
  value = data.iru_prism_device_information.example.results
}
