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

var _ datasource.DataSource = &prismActivationLockDataSource{}

func NewPrismActivationLockDataSource() datasource.DataSource {
	return &prismActivationLockDataSource{}
}

type prismActivationLockDataSource struct {
	client *client.Client
}

type prismActivationLockDataSourceModel struct {
	ID      types.String               `tfsdk:"id"`
	Limit   types.Int64                `tfsdk:"limit"`
	Offset  types.Int64                `tfsdk:"offset"`
	Results []prismActivationLockModel `tfsdk:"results"`
}

type prismActivationLockModel struct {
	DeviceID              types.String `tfsdk:"device_id"`
	DeviceName            types.String `tfsdk:"device_name"`
	SerialNumber          types.String `tfsdk:"serial_number"`
	ActivationLockEnabled types.Bool   `tfsdk:"activation_lock_enabled"`
}

func (d *prismActivationLockDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_prism_activation_lock"
}

func (d *prismActivationLockDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "List Activation Lock status from Prism.",
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
						"device_id":               schema.StringAttribute{Computed: true},
						"device_name":             schema.StringAttribute{Computed: true},
						"serial_number":           schema.StringAttribute{Computed: true},
						"activation_lock_enabled": schema.BoolAttribute{Computed: true},
					},
				},
			},
		},
	}
}

func (d *prismActivationLockDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	d.client = req.ProviderData.(*client.Client)
}

func (d *prismActivationLockDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data prismActivationLockDataSourceModel
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

		path := "/prism/activation_lock?" + params.Encode()
		type prismResponse struct {
			Data []client.PrismEntry `json:"data"`
		}
		var listResp prismResponse

		err := d.client.DoRequest(ctx, "GET", path, nil, &listResp)
		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read prism activation_lock, got error: %s", err))
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

	data.ID = types.StringValue("prism_activation_lock")
	data.Results = make([]prismActivationLockModel, 0, len(all))
	for _, item := range all {
		enabled, _ := item["activation_lock_enabled"].(bool)

		data.Results = append(data.Results, prismActivationLockModel{
			DeviceID:              types.StringValue(fmt.Sprintf("%v", item["device_id"])),
			DeviceName:            types.StringValue(fmt.Sprintf("%v", item["device__name"])),
			SerialNumber:          types.StringValue(fmt.Sprintf("%v", item["serial_number"])),
			ActivationLockEnabled: types.BoolValue(enabled),
		})
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
