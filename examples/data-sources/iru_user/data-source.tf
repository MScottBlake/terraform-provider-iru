data "iru_user" "example" {
  id = "your-user-uuid"
}

output "user_email" {
  value = data.iru_user.example.email
}
