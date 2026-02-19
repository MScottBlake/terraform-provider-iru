data "iru_custom_scripts" "all" {}

output "custom_scripts" {
  value = data.iru_custom_scripts.all.scripts
}
