terraform {
  required_version = ">= 1.14.0"

  required_providers {
    iru = {
      source  = "MScottBlake/iru"
      version = "~> 0.0.9"
    }
  }
}

provider "iru" {
  api_url   = "your-subdomain.api.kandji.io" # or IRU_API_URL env var
  api_token = "your-api-token"               # or IRU_API_TOKEN env var
}
