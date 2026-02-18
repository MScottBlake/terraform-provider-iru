# Query existing iru_blueprint resources
list "iru_blueprint" "all" {
  provider         = iru
  include_resource = true
  limit            = 100
}

# Use the queried results for bulk import into managed resources
/*
import {
  for_each = list.iru_blueprint.all
  to       = iru_blueprint.managed[each.key]
  id       = each.value.id
}

resource "iru_blueprint" "managed" {
  for_each = list.iru_blueprint.all
}
*/
