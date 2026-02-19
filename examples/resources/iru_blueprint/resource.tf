resource "iru_blueprint" "example" {
  name                   = "Standard Mac Blueprint"
  description            = "Managed by Terraform - Production"
  type                   = "classic"
  enrollment_code_active = true
}
