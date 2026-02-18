data "iru_prism_apps" "example" {}

output "installed_apps" {
  value = data.iru_prism_apps.example.results
}
