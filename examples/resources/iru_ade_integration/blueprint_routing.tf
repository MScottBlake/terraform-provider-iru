resource "iru_ade_integration" "routing_example" {
  use_blueprint_routing = true
  phone                 = "1234567890"
  email                 = "admin@example.com"
  mdm_server_token_file = file("${path.module}/token.p7m")
}
