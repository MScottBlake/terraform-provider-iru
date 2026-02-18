# Query existing iru_in_house_app resources
list "iru_in_house_app" "all" {
  provider         = iru
  include_resource = true
  limit            = 100
}

# Use the queried results for bulk import into managed resources
/*
import {
  for_each = list.iru_in_house_app.all
  to       = iru_in_house_app.managed[each.key]
  id       = each.value.id
}

resource "iru_in_house_app" "managed" {
  for_each = list.iru_in_house_app.all
}
*/
