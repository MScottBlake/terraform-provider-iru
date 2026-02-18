package provider

import (
	"context"
	"fmt"

	"github.com/MScottBlake/terraform-provider-iru/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ datasource.DataSource = &prismLaunchAgentsDaemonsDataSource{}

func NewPrismLaunchAgentsDaemonsDataSource() datasource.DataSource {
	return &prismLaunchAgentsDaemonsDataSource{}
}

type prismLaunchAgentsDaemonsDataSource struct {
	client *client.Client
}

type prismLaunchAgentsDaemonsDataSourceModel struct {
	Results []prismLaunchAgentDaemonModel `tfsdk:"results"`
}

type prismLaunchAgentDaemonModel struct {
	DeviceID     types.String `tfsdk:"device_id"`
	DeviceName   types.String `tfsdk:"device_name"`
	SerialNumber types.String `tfsdk:"serial_number"`
	Label        types.String `tfsdk:"label"`
	Path         types.String `tfsdk:"path"`
	IsLoaded     types.Bool   `tfsdk:"is_loaded"`
}

func (d *prismLaunchAgentsDaemonsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_prism_launch_agents_daemons"
}

func (d *prismLaunchAgentsDaemonsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "List launch agents and daemons from Prism.",
		Attributes: map[string]schema.Attribute{
			"results": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"device_id":     schema.StringAttribute{Computed: true},
						"device_name":   schema.StringAttribute{Computed: true},
						"serial_number": schema.StringAttribute{Computed: true},
						"label":         schema.StringAttribute{Computed: true},
						"path":          schema.StringAttribute{Computed: true},
						"is_loaded":     schema.BoolAttribute{Computed: true},
					},
				},
			},
		},
	}
}

func (d *prismLaunchAgentsDaemonsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	d.client = req.ProviderData.(*client.Client)
}

func (d *prismLaunchAgentsDaemonsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data prismLaunchAgentsDaemonsDataSourceModel

	var all []client.PrismEntry
	offset := 0
	limit := 300
	
	for {
		type prismResponse struct {
			Data []client.PrismEntry `json:"data"`
		}
		var listResp prismResponse
		
		path := fmt.Sprintf("/prism/launch_agents_and_daemons/?limit=%d&offset=%d", limit, offset)
		err := d.client.DoRequest(ctx, "GET", path, nil, &listResp)
		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read prism launch_agents_daemons, got error: %s", err))
			return
		}

		all = append(all, listResp.Data...)
		if len(listResp.Data) < limit {
			break
		}
		offset += limit
	}

	for _, item := range all {
		isLoaded, _ := item["is_loaded"].(bool)

		data.Results = append(data.Results, prismLaunchAgentDaemonModel{
			DeviceID:     types.StringValue(fmt.Sprintf("%v", item["device_id"])),
			DeviceName:   types.StringValue(fmt.Sprintf("%v", item["device__name"])),
			SerialNumber: types.StringValue(fmt.Sprintf("%v", item["serial_number"])),
			Label:        types.StringValue(fmt.Sprintf("%v", item["label"])),
			Path:         types.StringValue(fmt.Sprintf("%v", item["path"])),
			IsLoaded:     types.BoolValue(isLoaded),
		})
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
