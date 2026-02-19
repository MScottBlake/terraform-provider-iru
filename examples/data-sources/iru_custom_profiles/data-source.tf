data "iru_custom_profiles" "all" {}

output "custom_profiles" {
  value = data.iru_custom_profiles.all.profiles
}
