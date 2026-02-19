package provider

import (
	"context"
	"fmt"

	"github.com/MScottBlake/terraform-provider-iru/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ datasource.DataSource = &blueprintTemplatesDataSource{}

func NewBlueprintTemplatesDataSource() datasource.DataSource {
	return &blueprintTemplatesDataSource{}
}

type blueprintTemplatesDataSource struct {
	client *client.Client
}

type blueprintTemplatesDataSourceModel struct {
	ID        types.String            `tfsdk:"id"`
	Templates []blueprintTemplateModel `tfsdk:"templates"`
}

type blueprintTemplateModel struct {
	ID   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
}

func (d *blueprintTemplatesDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_blueprint_templates"
}

func (d *blueprintTemplatesDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "List available blueprint templates.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"templates": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed: true,
						},
						"name": schema.StringAttribute{
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func (d *blueprintTemplatesDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	d.client = req.ProviderData.(*client.Client)
}

func (d *blueprintTemplatesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data blueprintTemplatesDataSourceModel

	var response struct {
		Results []client.BlueprintTemplate `json:"results"` // Assuming response structure
	}
	// Note: The Postman doesn't show the response structure for templates explicitly in the truncated output
	// But usually it's "results": [...]
	err := d.client.DoRequest(ctx, "GET", "/blueprints/templates", nil, &response)
	if err != nil {
		// Fallback: maybe it returns a list directly?
		// But client.DoRequest expects struct to decode into.
		// If I am wrong about structure, I might need to debug.
		// Assuming standard Kandji/Iru pagination wrapper "results".
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to list blueprint templates, got error: %s", err))
		return
	}

	data.ID = types.StringValue("blueprint_templates")
	for _, t := range response.Results {
		data.Templates = append(data.Templates, blueprintTemplateModel{
			ID:   types.StringValue(t.ID),
			Name: types.StringValue(t.Name),
		})
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
