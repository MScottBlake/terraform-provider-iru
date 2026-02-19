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

var _ datasource.DataSource = &prismGatekeeperXProtectDataSource{}

func NewPrismGatekeeperXProtectDataSource() datasource.DataSource {
	return &prismGatekeeperXProtectDataSource{}
}

type prismGatekeeperXProtectDataSource struct {
	client *client.Client
}

type prismGatekeeperXProtectDataSourceModel struct {
	ID      types.String                   `tfsdk:"id"`
	Limit   types.Int64                    `tfsdk:"limit"`
	Offset  types.Int64                    `tfsdk:"offset"`
	Results []prismGatekeeperXProtectModel `tfsdk:"results"`
}

type prismGatekeeperXProtectModel struct {
	DeviceID                  types.String `tfsdk:"device_id"`
	DeviceName                types.String `tfsdk:"device_name"`
	SerialNumber              types.String `tfsdk:"serial_number"`
	GatekeeperStatus          types.Bool   `tfsdk:"gatekeeper_status"`
	TrustedDevelopers         types.Bool   `tfsdk:"trusted_developers"`
	GatekeeperVersion         types.String `tfsdk:"gatekeeper_version"`
	XProtectVersion           types.String `tfsdk:"xprotect_version"`
	MalwareRemovalToolVersion types.String `tfsdk:"malware_removal_tool_version"`
}

func (d *prismGatekeeperXProtectDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_prism_gatekeeper_xprotect"
}

func (d *prismGatekeeperXProtectDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "List Gatekeeper and XProtect information from Prism.",
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
						"device_id":                    schema.StringAttribute{Computed: true},
						"device_name":                  schema.StringAttribute{Computed: true},
						"serial_number":                schema.StringAttribute{Computed: true},
						"gatekeeper_status":            schema.BoolAttribute{Computed: true},
						"trusted_developers":           schema.BoolAttribute{Computed: true},
						"gatekeeper_version":           schema.StringAttribute{Computed: true},
						"xprotect_version":             schema.StringAttribute{Computed: true},
						"malware_removal_tool_version": schema.StringAttribute{Computed: true},
					},
				},
			},
		},
	}
}

func (d *prismGatekeeperXProtectDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	d.client = req.ProviderData.(*client.Client)
}

func (d *prismGatekeeperXProtectDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data prismGatekeeperXProtectDataSourceModel
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

		path := "/api/v1/prism/gatekeeper_and_xprotect?" + params.Encode()
		type prismResponse struct {
			Data []client.PrismEntry `json:"data"`
		}
		var listResp prismResponse

		err := d.client.DoRequest(ctx, "GET", path, nil, &listResp)
		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read prism gatekeeper_xprotect, got error: %s", err))
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

	data.ID = types.StringValue("prism_gatekeeper_xprotect")
	data.Results = make([]prismGatekeeperXProtectModel, 0, len(all))
	for _, item := range all {
		gkStatus, _ := item["gatekeeper_status"].(bool)
		trustedDevs, _ := item["trusted_developers"].(bool)

		data.Results = append(data.Results, prismGatekeeperXProtectModel{
			DeviceID:                  types.StringValue(fmt.Sprintf("%v", item["device_id"])),
			DeviceName:                types.StringValue(fmt.Sprintf("%v", item["device__name"])),
			SerialNumber:              types.StringValue(fmt.Sprintf("%v", item["serial_number"])),
			GatekeeperStatus:          types.BoolValue(gkStatus),
			TrustedDevelopers:         types.BoolValue(trustedDevs),
			GatekeeperVersion:         types.StringValue(fmt.Sprintf("%v", item["gatekeeper_version"])),
			XProtectVersion:           types.StringValue(fmt.Sprintf("%v", item["xprotect_version"])),
			MalwareRemovalToolVersion: types.StringValue(fmt.Sprintf("%v", item["malware_removal_tool_version"])),
		})
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
