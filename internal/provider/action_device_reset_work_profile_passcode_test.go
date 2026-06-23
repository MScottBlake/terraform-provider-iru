package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDeviceResetWorkProfilePasscodeAction(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
action "iru_device_action_reset_work_profile_passcode" "test" {
  device_id    = "PLACEHOLDER"
  new_password = "TestPassword123!"
  reset_password_flags = ["LOCK_NOW"]
}
`,
			},
		},
	})
}
