# ADE Devices cannot be created via Terraform. 
# Use 'terraform import' to manage existing ADE devices.

# terraform import iru_ade_device.example <device_uuid>

resource "iru_ade_device" "example" {
  asset_tag               = "IT-MAC-001"
  blueprint_id            = "c0148e35-c734-4402-b2fb-1c61aab72550"
  user_id                 = "8a9f88d9-e7f4-47e6-9326-fd4b39534c4e"
  override_dep_profile_id = "8429280a-13ae-40cf-8d1e-53f8b80df62f"
}
