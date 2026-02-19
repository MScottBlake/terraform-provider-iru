# ADE Devices cannot be created via Terraform, only imported and managed.
# Use 'terraform import' or 'terraform query' to bring them into state.

resource "iru_ade_device" "example" {
  blueprint_id = "your-blueprint-uuid"
  asset_tag    = "LAB-1234"
  user_id      = "user-uuid"
}
