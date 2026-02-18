package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccADEIntegrationListResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				// In Terraform 1.14, list resources are queried. 
				// Acceptance tests for them can verify they exist.
				Config: `
# List resources are primarily for 'terraform query'
`,
			},
		},
	})
}
