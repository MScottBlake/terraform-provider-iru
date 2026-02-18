data "iru_prism_application_firewall" "example" {}

output "firewall_status" {
  value = data.iru_prism_application_firewall.example.results
}
