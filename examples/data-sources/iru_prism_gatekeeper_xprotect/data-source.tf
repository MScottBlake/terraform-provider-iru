data "iru_prism_gatekeeper_xprotect" "example" {}

output "security_status" {
  value = data.iru_prism_gatekeeper_xprotect.example.results
}
