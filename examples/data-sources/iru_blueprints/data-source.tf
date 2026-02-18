data "iru_blueprints" "all" {}

output "blueprint_names" {
  value = [for b in data.iru_blueprints.all.blueprints : b.name]
}
