package provider

import (
	"context"
	"fmt"

	"github.com/MScottBlake/terraform-provider-iru/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ datasource.DataSource = &prismDesktopScreensaverDataSource{}

func NewPrismDesktopScreensaverDataSource() datasource.DataSource {
	return &prismDesktopScreensaverDataSource{}
}

type prismDesktopScreensaverDataSource struct {
	client *client.Client
}

type prismDesktopScreensaverDataSourceModel struct {
	Results []prismDesktopScreensaverModel `tfsdk:"results"`
}

type prismDesktopScreensaverModel struct {
	DeviceID     types.String `tfsdk:"device_id"`
	DeviceName   types.String `tfsdk:"device_name"`
	SerialNumber types.String `tfsdk:"serial_number"`
	HostName     types.String `tfsdk:"host_name"`
	OSVersion    types.String `tfsdk:"os_version"`
	MarketingName types.String `tfsdk:"marketing_name"`
}

func (d *prismDesktopScreensaverDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_prism_desktop_screensaver"
}

func (d *prismDesktopScreensaverDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "List desktop and screensaver information from Prism.",
		Attributes: map[string]schema.Attribute{
			"results": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"device_id":     schema.StringAttribute{Computed: true},
						"device_name":   schema.StringAttribute{Computed: true},
						"serial_number": schema.StringAttribute{Computed: true},
						"host_name":     schema.StringAttribute{Computed: true},
						"os_version":    schema.StringAttribute{Computed: true},
						"marketing_name": schema.StringAttribute{Computed: true},
					},
				},
			},
		},
	}
}

func (d *prismDesktopScreensaverDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	d.client = req.ProviderData.(*client.Client)
}

func (d *prismDesktopScreensaverDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data prismDesktopScreensaverDataSourceModel

	var all []client.PrismEntry
	offset := 0
	limit := 300
	
	for {
		type prismResponse struct {
			Data []client.PrismEntry `json:"data"`
		}
		var listResp prismResponse
		
		path := fmt.Sprintf("/prism/desktop_screensaver/?limit=%d&offset=%d", limit, offset)
		err := d.client.DoRequest(ctx, "GET", path, nil, &listResp)
		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read prism desktop_screensaver, got error: %s", err))
			return
		}

		all = append(all, listResp.Data...)
		if len(listResp.Data) < limit {
			break
		}
		offset += limit
	}

	for _, item := range all {
		data.Results = append(data.Results, prismDesktopScreensaverModel{
			DeviceID:     types.StringValue(fmt.Sprintf("%v", item["device_id"])),
			DeviceName:   types.StringValue(fmt.Sprintf("%v", item["device__name"])),
			SerialNumber: types.StringValue(fmt.Sprintf("%v", item["serial_number"])),
			HostName:     types.StringValue(fmt.Sprintf("%v", item["host_name"])),
			OSVersion:    types.StringValue(fmt.Sprintf("%v", item["os_version"])),
			MarketingName: types.StringValue(fmt.Sprintf("%v", item["marketing_name"])),
		})
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
