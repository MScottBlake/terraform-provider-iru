package provider

import (
	"context"
	"fmt"

	"github.com/MScottBlake/terraform-provider-iru/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ datasource.DataSource = &prismStartupSettingsDataSource{}

func NewPrismStartupSettingsDataSource() datasource.DataSource {
	return &prismStartupSettingsDataSource{}
}

type prismStartupSettingsDataSource struct {
	client *client.Client
}

type prismStartupSettingsDataSourceModel struct {
	Results []prismStartupSettingModel `tfsdk:"results"`
}

type prismStartupSettingModel struct {
	DeviceID     types.String `tfsdk:"device_id"`
	DeviceName   types.String `tfsdk:"device_name"`
	SerialNumber types.String `tfsdk:"serial_number"`
	SIP          types.Bool   `tfsdk:"sip"`
	SSV          types.Bool   `tfsdk:"ssv"`
	SecureBoot   types.String `tfsdk:"secure_boot"`
}

func (d *prismStartupSettingsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_prism_startup_settings"
}

func (d *prismStartupSettingsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "List startup settings from Prism.",
		Attributes: map[string]schema.Attribute{
			"results": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"device_id":     schema.StringAttribute{Computed: true},
						"device_name":   schema.StringAttribute{Computed: true},
						"serial_number": schema.StringAttribute{Computed: true},
						"sip":           schema.BoolAttribute{Computed: true},
						"ssv":           schema.BoolAttribute{Computed: true},
						"secure_boot":   schema.StringAttribute{Computed: true},
					},
				},
			},
		},
	}
}

func (d *prismStartupSettingsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	d.client = req.ProviderData.(*client.Client)
}

func (d *prismStartupSettingsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data prismStartupSettingsDataSourceModel

	var all []client.PrismEntry
	offset := 0
	limit := 300
	
	for {
		type prismResponse struct {
			Data []client.PrismEntry `json:"data"`
		}
		var listResp prismResponse
		
		path := fmt.Sprintf("/prism/startup_settings?limit=%d&offset=%d", limit, offset)
		err := d.client.DoRequest(ctx, "GET", path, nil, &listResp)
		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read prism startup_settings, got error: %s", err))
			return
		}

		all = append(all, listResp.Data...)
		if len(listResp.Data) < limit {
			break
		}
		offset += limit
	}

	for _, item := range all {
		sip, _ := item["sip"].(bool)
		ssv, _ := item["ssv"].(bool)
		secureBoot, _ := item["secure_boot_level"].(string)

		data.Results = append(data.Results, prismStartupSettingModel{
			DeviceID:     types.StringValue(fmt.Sprintf("%v", item["device_id"])),
			DeviceName:   types.StringValue(fmt.Sprintf("%v", item["device__name"])),
			SerialNumber: types.StringValue(fmt.Sprintf("%v", item["serial_number"])),
			SIP:          types.BoolValue(sip),
			SSV:          types.BoolValue(ssv),
			SecureBoot:   types.StringValue(secureBoot),
		})
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
