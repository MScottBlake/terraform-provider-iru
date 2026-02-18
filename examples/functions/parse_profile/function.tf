output "profile_metadata" {
  value = provider::iru::parse_profile(file("my_profile.mobileconfig"))
}
