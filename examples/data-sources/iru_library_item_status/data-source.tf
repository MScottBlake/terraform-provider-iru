data "iru_library_item_status" "example" {
  library_item_id = "your-library-item-uuid"
}

output "item_statuses" {
  value = data.iru_library_item_status.example.statuses
}
