package provider

import (
	"context"
	"fmt"

	"github.com/MScottBlake/terraform-provider-iru/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ datasource.DataSource = &customProfilesDataSource{}

func NewCustomProfilesDataSource() datasource.DataSource {
	return &customProfilesDataSource{}
}

type customProfilesDataSource struct {
	client *client.Client
}

type customProfilesDataSourceModel struct {
	Profiles []customProfileDataSourceModel `tfsdk:"profiles"`
}

type customProfileDataSourceModel struct {
	ID            types.String `tfsdk:"id"`
	Name          types.String `tfsdk:"name"`
	Active        types.Bool   `tfsdk:"active"`
	MDMIdentifier types.String `tfsdk:"mdm_identifier"`
	RunsOnMac     types.Bool   `tfsdk:"runs_on_mac"`
	RunsOnIPhone  types.Bool   `tfsdk:"runs_on_iphone"`
	RunsOnIPad    types.Bool   `tfsdk:"runs_on_ipad"`
	RunsOnTV      types.Bool   `tfsdk:"runs_on_tv"`
	RunsOnVision  types.Bool   `tfsdk:"runs_on_vision"`
}

func (d *customProfilesDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_custom_profiles"
}

func (d *customProfilesDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "List all custom profiles in the Kandji instance.",
		Attributes: map[string]schema.Attribute{
			"profiles": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The unique identifier for the Custom Profile.",
						},
						"name": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The name of the Custom Profile.",
						},
						"active": schema.BoolAttribute{
							Computed:            true,
							MarkdownDescription: "Whether this library item is active.",
						},
						"mdm_identifier": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The MDM identifier of the profile.",
						},
						"runs_on_mac": schema.BoolAttribute{
							Computed:            true,
							MarkdownDescription: "Whether the profile runs on macOS.",
						},
						"runs_on_iphone": schema.BoolAttribute{
							Computed:            true,
							MarkdownDescription: "Whether the profile runs on iOS.",
						},
						"runs_on_ipad": schema.BoolAttribute{
							Computed:            true,
							MarkdownDescription: "Whether the profile runs on iPadOS.",
						},
						"runs_on_tv": schema.BoolAttribute{
							Computed:            true,
							MarkdownDescription: "Whether the profile runs on tvOS.",
						},
						"runs_on_vision": schema.BoolAttribute{
							Computed:            true,
							MarkdownDescription: "Whether the profile runs on visionOS.",
						},
					},
				},
			},
		},
	}
}

func (d *customProfilesDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *customProfilesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data customProfilesDataSourceModel

	var allProfiles []client.CustomProfile
	offset := 0
	limit := 300
	
	for {
		type listProfilesResponse struct {
			Results []client.CustomProfile `json:"results"`
		}
		var listResp listProfilesResponse
		
		path := fmt.Sprintf("/library/custom-profiles?limit=%d&offset=%d", limit, offset)
		err := d.client.DoRequest(ctx, "GET", path, nil, &listResp)
		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read custom profiles, got error: %s", err))
			return
		}

		allProfiles = append(allProfiles, listResp.Results...)
		
		if len(listResp.Results) < limit {
			break
		}
		offset += limit
	}

	for _, resp := range allProfiles {
		data.Profiles = append(data.Profiles, customProfileDataSourceModel{
			ID:            types.StringValue(resp.ID),
			Name:          types.StringValue(resp.Name),
			Active:        types.BoolValue(resp.Active),
			MDMIdentifier: types.StringValue(resp.MDMIdentifier),
			RunsOnMac:     types.BoolValue(resp.RunsOnMac),
			RunsOnIPhone:  types.BoolValue(resp.RunsOnIPhone),
			RunsOnIPad:    types.BoolValue(resp.RunsOnIPad),
			RunsOnTV:      types.BoolValue(resp.RunsOnTV),
			RunsOnVision:  types.BoolValue(resp.RunsOnVision),
		})
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
