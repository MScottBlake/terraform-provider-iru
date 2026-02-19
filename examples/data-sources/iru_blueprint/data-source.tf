data "iru_blueprint" "example" {
  id = "your-blueprint-uuid"
}

output "blueprint_name" {
  value = data.iru_blueprint.example.name
}
