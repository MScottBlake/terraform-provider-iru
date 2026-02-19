# iru_blueprint_routing_activity (Data Source)

Provides Blueprint Routing activity information.

## Example Usage

```terraform
data "iru_blueprint_routing_activity" "example" {
  limit = 100
}
```

<!-- schema -->
### Optional

- `limit` (Number) Maximum number of results to return. Default is 300.

### Computed

- `activities` (List of Object) (see [below for nested schema](#nestedatt--activities))

<a id="nestedatt--activities"></a>
### Nested Schema for `activities`

Computed:

- `activity_time` (String)
- `activity_type` (String)
- `device_id` (String)
- `id` (Number)
