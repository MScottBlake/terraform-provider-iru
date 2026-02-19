resource "iru_custom_profile" "example" {
  name           = "Corporate WiFi Profile"
  active         = true
  profile_file   = file("${path.module}/corporate_wifi.mobileconfig")
  runs_on_mac    = true
  runs_on_iphone = true
  runs_on_ipad   = true
}
