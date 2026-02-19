resource "iru_ade_device" "routing_example" {
  use_blueprint_routing = true
  asset_tag            = "LAB-1234"
  user_id              = "user-uuid"
}
