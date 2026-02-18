data "iru_prism_cellular" "example" {}

output "cellular_info" {
  value = data.iru_prism_cellular.example.results
}
