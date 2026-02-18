resource "iru_custom_profile" "example" {
  name         = "Example Profile"
  active       = true
  profile_file = file("${path.module}/example.mobileconfig")
  runs_on_mac  = true
}
