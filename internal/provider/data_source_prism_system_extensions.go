package provider

import (
	"context"
	"fmt"

	"github.com/MScottBlake/terraform-provider-iru/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ datasource.DataSource = &prismSystemExtensionsDataSource{}

func NewPrismSystemExtensionsDataSource() datasource.DataSource {
	return &prismSystemExtensionsDataSource{}
}

type prismSystemExtensionsDataSource struct {
	client *client.Client
}

type prismSystemExtensionsDataSourceModel struct {
	Results []prismSystemExtensionModel `tfsdk:"results"`
}

type prismSystemExtensionModel struct {
	DeviceID     types.String `tfsdk:"device_id"`
	DeviceName   types.String `tfsdk:"device_name"`
	SerialNumber types.String `tfsdk:"serial_number"`
	Identifier   types.String `tfsdk:"identifier"`
	Name         types.String `tfsdk:"name"`
	State        types.String `tfsdk:"state"`
	TeamID       types.String `tfsdk:"team_id"`
	BundlePath   types.String `tfsdk:"bundle_path"`
	IsMDMManaged types.Bool   `tfsdk:"is_mdm_managed"`
}

func (d *prismSystemExtensionsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_prism_system_extensions"
}

func (d *prismSystemExtensionsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "List system extensions from Prism.",
		Attributes: map[string]schema.Attribute{
			"results": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"device_id":      schema.StringAttribute{Computed: true},
						"device_name":    schema.StringAttribute{Computed: true},
						"serial_number":  schema.StringAttribute{Computed: true},
						"identifier":     schema.StringAttribute{Computed: true},
						"name":           schema.StringAttribute{Computed: true},
						"state":          schema.StringAttribute{Computed: true},
						"team_id":        schema.StringAttribute{Computed: true},
						"bundle_path":    schema.StringAttribute{Computed: true},
						"is_mdm_managed": schema.BoolAttribute{Computed: true},
					},
				},
			},
		},
	}
}

func (d *prismSystemExtensionsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	d.client = req.ProviderData.(*client.Client)
}

func (d *prismSystemExtensionsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data prismSystemExtensionsDataSourceModel

	var all []client.PrismEntry
	offset := 0
	limit := 300
	
	for {
		type prismResponse struct {
			Data []client.PrismEntry `json:"data"`
		}
		var listResp prismResponse
		
		path := fmt.Sprintf("/prism/system_extensions?limit=%d&offset=%d", limit, offset)
		err := d.client.DoRequest(ctx, "GET", path, nil, &listResp)
		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read prism system_extensions, got error: %s", err))
			return
		}

		all = append(all, listResp.Data...)
		
		if len(listResp.Data) < limit {
			break
		}
		offset += limit
	}

	for _, item := range all {
		identifier, _ := item["identifier"].(string)
		name, _ := item["name"].(string)
		state, _ := item["state"].(string)
		teamID, _ := item["team_id"].(string)
		bundlePath, _ := item["bundle_path"].(string)
		isMDM, _ := item["is_mdm_managed"].(bool)

		data.Results = append(data.Results, prismSystemExtensionModel{
			DeviceID:     types.StringValue(fmt.Sprintf("%v", item["device_id"])),
			DeviceName:   types.StringValue(fmt.Sprintf("%v", item["device__name"])),
			SerialNumber: types.StringValue(fmt.Sprintf("%v", item["serial_number"])),
			Identifier:   types.StringValue(identifier),
			Name:         types.StringValue(name),
			State:        types.StringValue(state),
			TeamID:       types.StringValue(teamID),
			BundlePath:   types.StringValue(bundlePath),
			IsMDMManaged: types.BoolValue(isMDM),
		})
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
