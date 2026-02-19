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

var _ datasource.DataSource = &prismAppFirewallDataSource{}

func NewPrismAppFirewallDataSource() datasource.DataSource {
	return &prismAppFirewallDataSource{}
}

type prismAppFirewallDataSource struct {
	client *client.Client
}

type prismAppFirewallDataSourceModel struct {
	ID      types.String            `tfsdk:"id"`
	Limit   types.Int64             `tfsdk:"limit"`
	Offset  types.Int64             `tfsdk:"offset"`
	Results []prismAppFirewallModel `tfsdk:"results"`
}

type prismAppFirewallModel struct {
	DeviceID                types.String `tfsdk:"device_id"`
	DeviceName              types.String `tfsdk:"device_name"`
	SerialNumber            types.String `tfsdk:"serial_number"`
	Status                  types.Bool   `tfsdk:"status"`
	BlockAllIncoming        types.Bool   `tfsdk:"block_all_incoming"`
	StealthMode             types.Bool   `tfsdk:"stealth_mode"`
	AllowSignedApplications types.Bool   `tfsdk:"allow_signed_applications"`
}

func (d *prismAppFirewallDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_prism_application_firewall"
}

func (d *prismAppFirewallDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "List Application Firewall status for macOS devices from Prism.",
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
						"device_id": schema.StringAttribute{
							Computed: true,
						},
						"device_name": schema.StringAttribute{
							Computed: true,
						},
						"serial_number": schema.StringAttribute{
							Computed: true,
						},
						"status": schema.BoolAttribute{
							Computed: true,
						},
						"block_all_incoming": schema.BoolAttribute{
							Computed: true,
						},
						"stealth_mode": schema.BoolAttribute{
							Computed: true,
						},
						"allow_signed_applications": schema.BoolAttribute{
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func (d *prismAppFirewallDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	d.client = req.ProviderData.(*client.Client)
}

func (d *prismAppFirewallDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data prismAppFirewallDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var all []client.PrismAppFirewall
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

		path := "/api/v1/prism/application_firewall?" + params.Encode()
		type prismResponse struct {
			Data []client.PrismAppFirewall `json:"data"`
		}
		var listResp prismResponse

		err := d.client.DoRequest(ctx, "GET", path, nil, &listResp)
		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read prism application_firewall, got error: %s", err))
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

	data.ID = types.StringValue("prism_application_firewall")
	data.Results = make([]prismAppFirewallModel, 0, len(all))
	for _, item := range all {
		data.Results = append(data.Results, prismAppFirewallModel{
			DeviceID:                types.StringValue(item.DeviceID),
			DeviceName:              types.StringValue(item.DeviceName),
			SerialNumber:            types.StringValue(item.SerialNumber),
			Status:                  types.BoolValue(item.Status),
			BlockAllIncoming:        types.BoolValue(item.BlockAllIncoming),
			StealthMode:             types.BoolValue(item.StealthMode),
			AllowSignedApplications: types.BoolValue(item.AllowSignedApplications),
		})
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
