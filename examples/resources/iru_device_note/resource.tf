resource "iru_device_note" "example" {
  device_id = "your-device-uuid"
  content   = "This is a note managed by Terraform."
}
