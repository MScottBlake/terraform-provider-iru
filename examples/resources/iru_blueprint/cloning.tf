# Example of cloning an existing blueprint
resource "iru_blueprint" "cloned" {
  name                   = "Cloned Blueprint"
  description            = "Cloned from an existing production blueprint"
  type                   = "classic"
  enrollment_code_active = false
  
  # Note: source_id and source_type are only used during resource creation.
  # Terraform will not detect changes to these after the resource is created.
  source_id   = "existing-blueprint-uuid"
  source_type = "classic"
}
