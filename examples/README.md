# Iru (Kandji) Provider Examples

This directory contains examples for using the Iru Terraform Provider.

## Structure

- `provider/`: Basic provider configuration.
- `resources/`: Standard CRUD resources (Blueprints, Tags, Scripts, etc.).
- `data-sources/`: Inventory and reporting data sources.
- `actions/`: Imperative actions (Restart, Erase, etc.).
- `ephemeral/`: Secure data retrieval without state persistence (Device Secrets).
- `functions/`: Provider-specific HCL functions.

## Modern Features (Terraform 1.14+)

### Imperative Actions

Use the `action` block to perform one-time commands:

```hcl
action "iru_device_action_restart" "now" {
  device_id = "..."
}
```

### Ephemeral Resources

Use `ephemeral` to handle sensitive data that should never touch your state file:

```hcl
ephemeral "iru_device_secrets" "example" {
  device_id = "..."
}
```

### List Resources & `terraform query`

You can discover existing objects using the `terraform query` command:

```sh
# Query all devices
terraform query iru_device

# Filter devices by platform
terraform query 'iru_device.filter(d => d.platform == "macOS")'

# Query specific tags
terraform query iru_tag
```

### Provider Functions

Use custom functions to process data locally:

```hcl
output "meta" {
  value = provider::iru::parse_profile(local.xml)
}
```
