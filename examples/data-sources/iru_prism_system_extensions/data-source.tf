data "iru_prism_system_extensions" "example" {}

output "extensions" {
  value = data.iru_prism_system_extensions.example.results
}
