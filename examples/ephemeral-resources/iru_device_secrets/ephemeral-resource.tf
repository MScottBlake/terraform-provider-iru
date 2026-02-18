# Ephemeral resource to fetch secrets without storing them in state
ephemeral "iru_device_secrets" "example" {
  device_id = "8a9f88d9-e7f4-47e6-9326-fd4b39534c4e"
}

# Accessing values during the plan/apply phase
output "recovery_key" {
  value     = ephemeral.iru_device_secrets.example.filevault_recovery_key
  sensitive = true
}

# Fetching ADE public key for ABM configuration
ephemeral "iru_ade_public_key" "abm" {}

output "public_key" {
  value = ephemeral.iru_ade_public_key.abm.public_key
}
