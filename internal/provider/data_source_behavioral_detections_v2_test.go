package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccBehavioralDetectionsV2DataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `data "iru_behavioral_detections_v2" "test" {}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.iru_behavioral_detections_v2.test", "id"),
					resource.TestCheckResourceAttrSet("data.iru_behavioral_detections_v2.test", "results.#"),
				),
			},
		},
	})
}
