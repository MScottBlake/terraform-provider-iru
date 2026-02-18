# terraform-provider-iru

A Terraform provider for the Iru service.

## Requirements

- [Terraform](https://www.terraform.io/downloads.html) >= 1.0
- [Go](https://golang.org/doc/install) >= 1.24

## Building The Provider

1. Clone the repository
2. Enter the repository directory
3. Build the provider using the Go `install` command:

```sh
go install .
```

## Developing the Provider

If you wish to work on the provider, you'll first need
[Go](http://www.golang.org) installed on your machine (see
[Requirements](#requirements)).

To compile the provider, run `go install`. This will build the binary and put it
in your `$GOPATH/bin` directory.

To generate or update documentation, run `go generate`.

## Using the Provider

```hcl
provider "iru" {
  # configuration options
}

data "iru_example" "example" {}

resource "iru_example" "example" {
  configurable = "example-value"
}
```
