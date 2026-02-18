data "iru_users" "example" {}

output "all_users" {
  value = data.iru_users.example.users
}
