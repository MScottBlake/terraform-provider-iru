data "iru_blueprint_library_items" "example" {
  blueprint_id = "your-blueprint-uuid"
}

output "assigned_library_items" {
  value = data.iru_blueprint_library_items.example.library_items
}
