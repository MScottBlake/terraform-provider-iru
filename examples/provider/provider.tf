terraform {
  required_version = ">= 1.14.0"
}

provider "iru" {
  api_url   = "your-subdomain.api.kandji.io"
  api_token = "your-api-token"
}
