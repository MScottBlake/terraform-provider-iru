data "iru_prism_desktop_screensaver" "example" {}

output "desktop_info" {
  value = data.iru_prism_desktop_screensaver.example.results
}
