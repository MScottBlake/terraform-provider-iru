data "iru_blueprint_templates" "example" {}

output "blueprint_templates" {
  value = data.iru_blueprint_templates.example.templates
}
