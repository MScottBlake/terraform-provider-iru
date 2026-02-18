data "iru_prism_kernel_extensions" "example" {}

output "kexts" {
  value = data.iru_prism_kernel_extensions.example.results
}
