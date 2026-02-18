data "iru_prism_installed_profiles" "example" {}

output "profiles" {
  value = data.iru_prism_installed_profiles.example.results
}
