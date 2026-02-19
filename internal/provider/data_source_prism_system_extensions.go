package provider

import (
	"context"
	"fmt"
	"net/url"

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
	ID      types.String                `tfsdk:"id"`
	Limit   types.Int64                 `tfsdk:"limit"`
	Offset  types.Int64                 `tfsdk:"offset"`
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
			"id": schema.StringAttribute{
				Computed: true,
			},
			"limit": schema.Int64Attribute{
				Optional:            true,
				MarkdownDescription: "Maximum number of results to return.",
			},
			"offset": schema.Int64Attribute{
				Optional:            true,
				MarkdownDescription: "Number of results to skip.",
			},
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
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var all []client.PrismEntry
	offset := 0
	if !data.Offset.IsNull() {
		offset = int(data.Offset.ValueInt64())
	}
	limit := 300
	if !data.Limit.IsNull() {
		limit = int(data.Limit.ValueInt64())
	}

	for {
		params := url.Values{}
		params.Add("limit", fmt.Sprintf("%d", limit))
		params.Add("offset", fmt.Sprintf("%d", offset))

		path := "/prism/system_extensions?" + params.Encode()
		type prismResponse struct {
			Data []client.PrismEntry `json:"data"`
		}
		var listResp prismResponse

		err := d.client.DoRequest(ctx, "GET", path, nil, &listResp)
		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read prism system_extensions, got error: %s", err))
			return
		}

		all = append(all, listResp.Data...)

		if !data.Limit.IsNull() && len(all) >= limit {
			all = all[:limit]
			break
		}

		if len(listResp.Data) < limit {
			break
		}
		offset += len(listResp.Data)
	}

	data.ID = types.StringValue("prism_system_extensions")
	data.Results = make([]prismSystemExtensionModel, 0, len(all))
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
