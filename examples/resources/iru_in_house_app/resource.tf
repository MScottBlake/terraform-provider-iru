resource "iru_in_house_app" "example" {
  name           = "Internal Sales App"
  file_key       = "apps/sales-app-v2.ipa"
  runs_on_iphone = true
  runs_on_ipad   = true
  active         = true
}
