data "iru_prism_certificates" "example" {}

output "certs" {
  value = data.iru_prism_certificates.example.results
}
