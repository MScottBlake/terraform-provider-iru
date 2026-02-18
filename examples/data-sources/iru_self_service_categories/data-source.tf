data "iru_self_service_categories" "example" {}

output "categories" {
  value = data.iru_self_service_categories.example.results
}
