resource "iru_custom_script" "example" {
  name                = "Audit and Remediate SSH"
  active              = true
  execution_frequency = "every_day"
  script              = <<-EOT
    #!/bin/zsh
    # Audit: check if SSH is enabled
    if systemsetup -getremotelogin | grep -q "On"; then
      exit 1
    fi
    exit 0
  EOT
  remediation_script = <<-EOT
    #!/bin/zsh
    # Remediate: disable SSH
    systemsetup -setremotelogin off
  EOT
  show_in_self_service = false
}
