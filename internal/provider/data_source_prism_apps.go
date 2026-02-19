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

var _ datasource.DataSource = &prismAppsDataSource{}

func NewPrismAppsDataSource() datasource.DataSource {
	return &prismAppsDataSource{}
}

type prismAppsDataSource struct {
	client *client.Client
}

type prismAppsDataSourceModel struct {
	ID      types.String    `tfsdk:"id"`
	Limit   types.Int64     `tfsdk:"limit"`
	Offset  types.Int64     `tfsdk:"offset"`
	Results []prismAppModel `tfsdk:"results"`
}

type prismAppModel struct {
	DeviceID     types.String `tfsdk:"device_id"`
	DeviceName   types.String `tfsdk:"device_name"`
	SerialNumber types.String `tfsdk:"serial_number"`
	Name         types.String `tfsdk:"name"`
	Version      types.String `tfsdk:"version"`
	BundleID     types.String `tfsdk:"bundle_id"`
	Path         types.String `tfsdk:"path"`
}

func (d *prismAppsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_prism_apps"
}

func (d *prismAppsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "List applications installed on devices from Prism.",
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
						"device_id":     schema.StringAttribute{Computed: true},
						"device_name":   schema.StringAttribute{Computed: true},
						"serial_number": schema.StringAttribute{Computed: true},
						"name":          schema.StringAttribute{Computed: true},
						"version":       schema.StringAttribute{Computed: true},
						"bundle_id":     schema.StringAttribute{Computed: true},
						"path":          schema.StringAttribute{Computed: true},
					},
				},
			},
		},
	}
}

func (d *prismAppsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	d.client = req.ProviderData.(*client.Client)
}

func (d *prismAppsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data prismAppsDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var all []client.PrismApp
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

		path := "/prism/apps?" + params.Encode()
		type prismAppResponse struct {
			Data []client.PrismApp `json:"data"`
		}
		var listResp prismAppResponse

		err := d.client.DoRequest(ctx, "GET", path, nil, &listResp)
		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read prism apps, got error: %s", err))
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

	data.ID = types.StringValue("prism_apps")
	data.Results = make([]prismAppModel, 0, len(all))
	for _, item := range all {
		data.Results = append(data.Results, prismAppModel{
			DeviceID:     types.StringValue(item.DeviceID),
			DeviceName:   types.StringValue(item.DeviceName),
			SerialNumber: types.StringValue(item.SerialNumber),
			Name:         types.StringValue(item.Name),
			Version:      types.StringValue(item.Version),
			BundleID:     types.StringValue(item.BundleID),
			Path:         types.StringValue(item.Path),
		})
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
