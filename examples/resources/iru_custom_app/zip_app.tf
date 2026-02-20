# Example of a Custom App using a ZIP file
resource "iru_custom_app" "zip_example" {
  name                = "Internal Tool"
  file_key            = "internal/tools/mytool.zip"
  install_type        = "zip"
  unzip_location      = "/usr/local/bin"
  install_enforcement = "install_once"
  active              = true
  
  preinstall_script = <<-EOT
    #!/bin/zsh
    echo "Preparing for installation..."
  EOT

  postinstall_script = <<-EOT
    #!/bin/zsh
    chmod +x /usr/local/bin/mytool
  EOT
}
