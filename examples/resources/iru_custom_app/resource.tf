resource "iru_custom_app" "example" {
  name                = "Google Chrome"
  file_key            = "apps/chrome-v120.pkg"
  install_type        = "package"
  install_enforcement = "continuously_enforce"
  audit_script        = <<-EOT
    #!/bin/zsh
    if [ -d "/Applications/Google Chrome.app" ]; then
      exit 0
    else
      exit 1
    fi
  EOT
  active              = true
}
