data "iru_licensing" "example" {}

output "tenant_status" {
  value = data.iru_licensing.example.tenant_over_license_limit
}
