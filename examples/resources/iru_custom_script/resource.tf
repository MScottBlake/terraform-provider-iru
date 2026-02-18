resource "iru_custom_script" "example" {
  name                = "Example Script"
  active              = true
  execution_frequency = "once"
  script              = "#!/bin/zsh
echo 'Hello from Terraform'"
}
