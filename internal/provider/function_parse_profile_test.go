package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccParseProfileFunction(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
output "test" {
  value = provider::iru::parse_profile("mock-xml")
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckOutput("test", "{identifier = "extracted-id", uuid = "extracted-uuid"}"),
				),
			},
		},
	})
}
