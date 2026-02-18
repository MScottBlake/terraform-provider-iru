# Query existing iru_custom_script resources
list "iru_custom_script" "all" {
  provider         = iru
  include_resource = true
  limit            = 100
}

# Use the queried results for bulk import into managed resources
/*
import {
  for_each = list.iru_custom_script.all
  to       = iru_custom_script.managed[each.key]
  id       = each.value.id
}

resource "iru_custom_script" "managed" {
  for_each = list.iru_custom_script.all
}
*/
