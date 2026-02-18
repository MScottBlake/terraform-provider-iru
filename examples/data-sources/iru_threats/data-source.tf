data "iru_threats" "example" {}

output "active_threats" {
  value = data.iru_threats.example.results
}
