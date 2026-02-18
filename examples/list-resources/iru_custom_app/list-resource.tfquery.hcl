# Query existing iru_custom_app resources
list "iru_custom_app" "all" {
  provider         = iru
  include_resource = true
  limit            = 100
}

# Use the queried results for bulk import into managed resources
/*
import {
  for_each = list.iru_custom_app.all
  to       = iru_custom_app.managed[each.key]
  id       = each.value.id
}

resource "iru_custom_app" "managed" {
  for_each = list.iru_custom_app.all
}
*/
