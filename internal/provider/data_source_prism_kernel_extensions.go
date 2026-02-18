package provider

import (
	"context"
	"fmt"

	"github.com/MScottBlake/terraform-provider-iru/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ datasource.DataSource = &prismKernelExtensionsDataSource{}

func NewPrismKernelExtensionsDataSource() datasource.DataSource {
	return &prismKernelExtensionsDataSource{}
}

type prismKernelExtensionsDataSource struct {
	client *client.Client
}

type prismKernelExtensionsDataSourceModel struct {
	Results []prismKernelExtensionModel `tfsdk:"results"`
}

type prismKernelExtensionModel struct {
	DeviceID     types.String `tfsdk:"device_id"`
	DeviceName   types.String `tfsdk:"device_name"`
	SerialNumber types.String `tfsdk:"serial_number"`
	BundleID     types.String `tfsdk:"bundle_id"`
	Version      types.String `tfsdk:"version"`
	Path         types.String `tfsdk:"path"`
}

func (d *prismKernelExtensionsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_prism_kernel_extensions"
}

func (d *prismKernelExtensionsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "List kernel extensions from Prism.",
		Attributes: map[string]schema.Attribute{
			"results": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"device_id":     schema.StringAttribute{Computed: true},
						"device_name":   schema.StringAttribute{Computed: true},
						"serial_number": schema.StringAttribute{Computed: true},
						"bundle_id":     schema.StringAttribute{Computed: true},
						"version":       schema.StringAttribute{Computed: true},
						"path":          schema.StringAttribute{Computed: true},
					},
				},
			},
		},
	}
}

func (d *prismKernelExtensionsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	d.client = req.ProviderData.(*client.Client)
}

func (d *prismKernelExtensionsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data prismKernelExtensionsDataSourceModel

	var all []client.PrismEntry
	offset := 0
	limit := 300
	
	for {
		type prismResponse struct {
			Data []client.PrismEntry `json:"data"`
		}
		var listResp prismResponse
		
		path := fmt.Sprintf("/prism/kernel_extensions/?limit=%d&offset=%d", limit, offset)
		err := d.client.DoRequest(ctx, "GET", path, nil, &listResp)
		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read prism kernel_extensions, got error: %s", err))
			return
		}

		all = append(all, listResp.Data...)
		if len(listResp.Data) < limit {
			break
		}
		offset += limit
	}

	for _, item := range all {
		data.Results = append(data.Results, prismKernelExtensionModel{
			DeviceID:     types.StringValue(fmt.Sprintf("%v", item["device_id"])),
			DeviceName:   types.StringValue(fmt.Sprintf("%v", item["device__name"])),
			SerialNumber: types.StringValue(fmt.Sprintf("%v", item["serial_number"])),
			BundleID:     types.StringValue(fmt.Sprintf("%v", item["bundle_id"])),
			Version:      types.StringValue(fmt.Sprintf("%v", item["version"])),
			Path:         types.StringValue(fmt.Sprintf("%v", item["path"])),
		})
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
