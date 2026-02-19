# iru_blueprint_routing (Data Source)

Provides Blueprint Routing settings.

## Example Usage

```terraform
data "iru_blueprint_routing" "current" {}
```

<!-- schema -->
### Computed

- `enrollment_code` (String) The enrollment code for Blueprint Routing.
- `enrollment_code_active` (Boolean) Whether the enrollment code for Blueprint Routing is active.
- `id` (String) A fixed identifier for the singleton data source.
