# Imperative action to restart a device
action "iru_device_action_restart" "now" {
  device_id = "8a9f88d9-e7f4-47e6-9326-fd4b39534c4e"
}

# Imperative action to erase a Windows device with specific mode
action "iru_device_action_erase" "wipe_windows" {
  device_id  = "c0148e35-c734-4402-b2fb-1c61aab72550"
  erase_mode = "WIPE"
}
