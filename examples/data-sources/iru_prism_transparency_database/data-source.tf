data "iru_prism_transparency_database" "example" {}

output "tcc" {
  value = data.iru_prism_transparency_database.example.results
}
