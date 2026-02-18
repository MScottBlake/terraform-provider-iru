data "iru_prism_local_users" "example" {}

output "users" {
  value = data.iru_prism_local_users.example.results
}
