# Query existing iru_ade_device resources
list "iru_ade_device" "all" {
  provider         = iru
  include_resource = true
  limit            = 100
}

# Use the queried results for bulk import into managed resources
/*
import {
  for_each = list.iru_ade_device.all
  to       = iru_ade_device.managed[each.key]
  id       = each.value.id
}

resource "iru_ade_device" "managed" {
  for_each = list.iru_ade_device.all
}
*/
