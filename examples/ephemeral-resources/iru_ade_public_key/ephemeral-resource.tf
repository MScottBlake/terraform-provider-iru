ephemeral "iru_ade_public_key" "example" {}

output "public_key" {
  value = ephemeral.iru_ade_public_key.example.public_key
}
