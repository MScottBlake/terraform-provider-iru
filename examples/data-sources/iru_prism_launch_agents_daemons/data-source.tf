data "iru_prism_launch_agents_daemons" "example" {}

output "persistence_items" {
  value = data.iru_prism_launch_agents_daemons.example.results
}
