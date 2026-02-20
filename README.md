# terraform-provider-iru

A Terraform provider for managing resources in the Iru ecosystem.

## Requirements

- [Terraform](https://www.terraform.io/downloads.html) >= 1.14.0
- [Go](https://golang.org/doc/install) >= 1.24

## Building The Provider

1. Clone the repository
2. Enter the repository directory
3. Build the provider:

```sh
go build -o terraform-provider-iru
```

### Local Development Installation

To use the provider locally for testing without a registry, you can build it into the Terraform local plugin directory.

```sh
# Example for darwin_arm64 (Apple Silicon Mac)
go build -o ~/.terraform.d/plugins/github.com/MScottBlake/iru/0.0.1/darwin_arm64/terraform-provider-iru -ldflags="-X main.version=0.0.1"
```

## Developing the Provider

To compile and install the provider to your `$GOPATH/bin` directory:

```sh
go install .
```

### Documentation

Documentation is generated using `terraform-plugin-docs`. To generate or update documentation, run:

```sh
go generate ./...
```

### Testing

Unit tests:

```sh
go test ./internal/client -v
```

Acceptance tests (requires `TF_ACC=1` and credentials):

```sh
export IRU_API_URL="..."
export IRU_API_TOKEN="..."
export TF_ACC=1
go test ./internal/provider -v
```

## Using the Provider

```hcl
terraform {
  required_providers {
    iru = {
      source  = "MScottBlake/iru"
      version = "~> 0"
    }
  }
}

provider "iru" {
  api_url   = "your-subdomain.api.kandji.io" # or IRU_API_URL env var
  api_token = "your-api-token"               # or IRU_API_TOKEN env var
}

# Data Source
data "iru_device" "example" {
  device_id = "device-uuid-here"
}

# Resource (Note: Devices are imported, not created)
resource "iru_device" "managed" {
  id = "device-uuid-here"
  asset_tag = "MAC-001"
}

# List Resource
list "iru_device" "all" {
  provider = iru
}

# Action Block (Imperative commands)
action "iru_device_action_restart" "restart_mac" {
  device_id = "device-uuid-here"
}
```
