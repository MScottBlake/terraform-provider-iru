data "iru_prism_startup_settings" "example" {}

output "startup" {
  value = data.iru_prism_startup_settings.example.results
}
