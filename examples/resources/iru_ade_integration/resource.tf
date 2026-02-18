resource "iru_ade_integration" "example" {
  blueprint_id          = "your-blueprint-uuid"
  phone                 = "1234567890"
  email                 = "admin@example.com"
  mdm_server_token_file = file("${path.module}/token.p7m")
}
