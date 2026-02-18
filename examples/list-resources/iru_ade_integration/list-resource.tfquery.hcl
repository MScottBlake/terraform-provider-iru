# Query existing iru_ade_integration resources
list "iru_ade_integration" "all" {
  provider         = iru
  include_resource = true
  limit            = 100
}

# Use the queried results for bulk import into managed resources
/*
import {
  for_each = list.iru_ade_integration.all
  to       = iru_ade_integration.managed[each.key]
  id       = each.value.id
}

resource "iru_ade_integration" "managed" {
  for_each = list.iru_ade_integration.all
}
*/
