# iru_blueprint_routing (Resource)

Manages Blueprint Routing settings. This is a singleton resource.

## Example Usage

```terraform
resource "iru_blueprint_routing" "example" {
  enrollment_code_active = true
}
```

<!-- schema -->
### Required

- `enrollment_code_active` (Boolean) Whether the enrollment code for Blueprint Routing is active.

### Computed

- `enrollment_code` (String) The enrollment code for Blueprint Routing.
- `id` (String) A fixed identifier for the singleton resource.
