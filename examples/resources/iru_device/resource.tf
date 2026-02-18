# Devices cannot be created via Terraform. 
# Use 'terraform import' to manage existing devices.

# terraform import iru_device.example <device_uuid>

resource "iru_device" "example" {
  asset_tag    = "IT-MAC-001"
  blueprint_id = "c0148e35-c734-4402-b2fb-1c61aab72550"
  user_id      = "8a9f88d9-e7f4-47e6-9326-fd4b39534c4e"
}
