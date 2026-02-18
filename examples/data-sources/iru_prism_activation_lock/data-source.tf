data "iru_prism_activation_lock" "example" {}

output "lock_status" {
  value = data.iru_prism_activation_lock.example.results
}
