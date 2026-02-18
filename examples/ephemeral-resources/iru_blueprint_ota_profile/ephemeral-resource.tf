ephemeral "iru_blueprint_ota_profile" "example" {
  blueprint_id = "c0148e35-c734-4402-b2fb-1c61aab72550"
}

output "manual_profile" {
  value = ephemeral.iru_blueprint_ota_profile.example.profile_xml
}
