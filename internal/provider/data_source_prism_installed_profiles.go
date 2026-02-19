package provider

import (
	"context"
	"fmt"

	"github.com/MScottBlake/terraform-provider-iru/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ datasource.DataSource = &prismInstalledProfilesDataSource{}

func NewPrismInstalledProfilesDataSource() datasource.DataSource {
	return &prismInstalledProfilesDataSource{}
}

type prismInstalledProfilesDataSource struct {
	client *client.Client
}

type prismInstalledProfilesDataSourceModel struct {
	Results []prismInstalledProfileModel `tfsdk:"results"`
}

type prismInstalledProfileModel struct {
	DeviceID           types.String `tfsdk:"device_id"`
	DeviceName         types.String `tfsdk:"device_name"`
	SerialNumber       types.String `tfsdk:"serial_number"`
	ProfileDisplayName types.String `tfsdk:"profile_display_name"`
	PayloadIdentifier  types.String `tfsdk:"payload_identifier"`
	PayloadUUID        types.String `tfsdk:"payload_uuid"`
	Managed            types.Bool   `tfsdk:"managed"`
}

func (d *prismInstalledProfilesDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_prism_installed_profiles"
}

func (d *prismInstalledProfilesDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "List installed profiles from Prism.",
		Attributes: map[string]schema.Attribute{
			"results": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"device_id":            schema.StringAttribute{Computed: true},
						"device_name":          schema.StringAttribute{Computed: true},
						"serial_number":        schema.StringAttribute{Computed: true},
						"profile_display_name": schema.StringAttribute{Computed: true},
						"payload_identifier":   schema.StringAttribute{Computed: true},
						"payload_uuid":         schema.StringAttribute{Computed: true},
						"managed":              schema.BoolAttribute{Computed: true},
					},
				},
			},
		},
	}
}

func (d *prismInstalledProfilesDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	d.client = req.ProviderData.(*client.Client)
}

func (d *prismInstalledProfilesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data prismInstalledProfilesDataSourceModel

	var all []client.PrismEntry
	offset := 0
	limit := 300
	
	for {
		type prismResponse struct {
			Data []client.PrismEntry `json:"data"`
		}
		var listResp prismResponse
		
		path := fmt.Sprintf("/prism/installed_profiles?limit=%d&offset=%d", limit, offset)
		err := d.client.DoRequest(ctx, "GET", path, nil, &listResp)
		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read prism installed_profiles, got error: %s", err))
			return
		}

		all = append(all, listResp.Data...)
		if len(listResp.Data) < limit {
			break
		}
		offset += limit
	}

	for _, item := range all {
		managed, _ := item["managed"].(bool)

		data.Results = append(data.Results, prismInstalledProfileModel{
			DeviceID:           types.StringValue(fmt.Sprintf("%v", item["device_id"])),
			DeviceName:         types.StringValue(fmt.Sprintf("%v", item["device__name"])),
			SerialNumber:       types.StringValue(fmt.Sprintf("%v", item["serial_number"])),
			ProfileDisplayName: types.StringValue(fmt.Sprintf("%v", item["profile_display_name"])),
			PayloadIdentifier:  types.StringValue(fmt.Sprintf("%v", item["payload_identifier"])),
			PayloadUUID:        types.StringValue(fmt.Sprintf("%v", item["payload_uuid"])),
			Managed:            types.BoolValue(managed),
		})
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
