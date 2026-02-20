# Example of a Custom Script available in Self Service
resource "iru_custom_script" "self_service_example" {
  name                = "Clear Local Cache"
  active              = true
  execution_frequency = "no_enforcement" # Only runs when user clicks in Self Service
  
  script = <<-EOT
    #!/bin/zsh
    echo "Clearing caches..."
    rm -rf ~/Library/Caches/*
  EOT

  show_in_self_service = true
}
