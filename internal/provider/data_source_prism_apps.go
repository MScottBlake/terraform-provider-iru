package provider

import (
	"context"
	"fmt"

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

	var all []client.PrismApp
	offset := 0
	limit := 300
	
	for {
		type prismAppResponse struct {
			Data []client.PrismApp `json:"data"`
		}
		var listResp prismAppResponse
		
		path := fmt.Sprintf("/prism/apps?limit=%d&offset=%d", limit, offset)
		err := d.client.DoRequest(ctx, "GET", path, nil, &listResp)
		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read prism apps, got error: %s", err))
			return
		}

		all = append(all, listResp.Data...)
		
		if len(listResp.Data) < limit {
			break
		}
		offset += limit
	}

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

	data.ID = types.StringValue("prism_apps")

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
