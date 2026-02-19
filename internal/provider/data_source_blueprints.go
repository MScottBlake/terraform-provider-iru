package provider

import (
	"context"
	"fmt"

	"github.com/MScottBlake/terraform-provider-iru/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ datasource.DataSource = &blueprintsDataSource{}

func NewBlueprintsDataSource() datasource.DataSource {
	return &blueprintsDataSource{}
}

type blueprintsDataSource struct {
	client *client.Client
}

type blueprintsDataSourceModel struct {
	ID         types.String                 `tfsdk:"id"`
	Blueprints []blueprintDataSourceModel `tfsdk:"blueprints"`
}

type blueprintDataSourceModel struct {
	ID             types.String `tfsdk:"id"`
	Name           types.String `tfsdk:"name"`
	Description    types.String `tfsdk:"description"`
	Icon           types.String `tfsdk:"icon"`
	Color          types.String `tfsdk:"color"`
	EnrollmentCode types.String `tfsdk:"enrollment_code"`
}

func (d *blueprintsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_blueprints"
}

func (d *blueprintsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "List all blueprints in the Kandji instance.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"blueprints": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The unique identifier for the Blueprint.",
						},
						"name": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The name of the Blueprint.",
						},
						"description": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The description of the Blueprint.",
						},
						"icon": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The icon of the Blueprint.",
						},
						"color": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The color of the Blueprint.",
						},
						"enrollment_code": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The enrollment code for the Blueprint.",
						},
					},
				},
			},
		},
	}
}

func (d *blueprintsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected DataSource Configure Type",
			fmt.Sprintf("Expected *client.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	d.client = client
}

func (d *blueprintsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data blueprintsDataSourceModel

	var allBlueprints []client.Blueprint
	offset := 0
	limit := 300
	
	for {
		// Note: List Blueprints API might return {results: []} or just [].
		// Postman shows {results: []}.
		type listBlueprintsResponse struct {
			Results []client.Blueprint `json:"results"`
		}
		var listResp listBlueprintsResponse
		
		path := fmt.Sprintf("/blueprints?limit=%d&offset=%d", limit, offset)
		err := d.client.DoRequest(ctx, "GET", path, nil, &listResp)
		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read blueprints, got error: %s", err))
			return
		}

		allBlueprints = append(allBlueprints, listResp.Results...)
		
		if len(listResp.Results) < limit {
			break
		}
		offset += limit
	}

	for _, blueprint := range allBlueprints {
		data.Blueprints = append(data.Blueprints, blueprintDataSourceModel{
			ID:             types.StringValue(blueprint.ID),
			Name:           types.StringValue(blueprint.Name),
			Description:    types.StringValue(blueprint.Description),
			Icon:           types.StringValue(blueprint.Icon),
			Color:          types.StringValue(blueprint.Color),
			EnrollmentCode: types.StringValue(blueprint.EnrollmentCode.Code),
		})
	}

	data.ID = types.StringValue("blueprints")

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
