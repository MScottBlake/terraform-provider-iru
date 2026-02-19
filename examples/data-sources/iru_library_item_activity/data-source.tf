data "iru_library_item_activity" "example" {
  library_item_id = "your-library-item-uuid"
}

output "item_activity" {
  value = data.iru_library_item_activity.example.activity
}
