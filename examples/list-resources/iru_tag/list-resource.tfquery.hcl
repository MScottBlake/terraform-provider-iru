# Query existing iru_tag resources
list "iru_tag" "all" {
  provider         = iru
  include_resource = true
  limit            = 100
}

# Use the queried results for bulk import into managed resources
/*
import {
  for_each = list.iru_tag.all
  to       = iru_tag.managed[each.key]
  id       = each.value.id
}

resource "iru_tag" "managed" {
  for_each = list.iru_tag.all
}
*/
