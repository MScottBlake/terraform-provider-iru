data "iru_blueprint_templates" "all" {}

output "blueprint_templates" {
  value = data.iru_blueprint_templates.all.templates
}
