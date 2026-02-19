data "iru_prism_count" "apps" {
  category = "apps"
}

output "total_apps" {
  value = data.iru_prism_count.apps.count
}
