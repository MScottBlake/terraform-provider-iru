# Using the custom function to parse metadata from a profile file
output "profile_metadata" {
  value = provider::iru::parse_profile(file("my_profile.mobileconfig"))
}

# Result will be a map: { identifier = "...", uuid = "..." }
