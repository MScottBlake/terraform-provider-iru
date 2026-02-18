package provider

import (
	"context"
	"fmt"

	"github.com/MScottBlake/terraform-provider-iru/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ datasource.DataSource = &prismDeviceInformationDataSource{}

func NewPrismDeviceInformationDataSource() datasource.DataSource {
	return &prismDeviceInformationDataSource{}
}

type prismDeviceInformationDataSource struct {
	client *client.Client
}

type prismDeviceInformationDataSourceModel struct {
	Results []prismDeviceInformationModel `tfsdk:"results"`
}

type prismDeviceInformationModel struct {
	DeviceID       types.String `tfsdk:"device_id"`
	DeviceName     types.String `tfsdk:"device_name"`
	SerialNumber   types.String `tfsdk:"serial_number"`
	DeviceCapacity types.Float64 `tfsdk:"device_capacity"`
	ModelName      types.String `tfsdk:"model_name"`
	OSVersion      types.String `tfsdk:"os_version"`
	MDMEnabled     types.Bool   `tfsdk:"mdm_enabled"`
	AgentInstalled types.Bool   `tfsdk:"agent_installed"`
}

func (d *prismDeviceInformationDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_prism_device_information"
}

func (d *prismDeviceInformationDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "List device information from Prism.",
		Attributes: map[string]schema.Attribute{
			"results": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"device_id":       schema.StringAttribute{Computed: true},
						"device_name":     schema.StringAttribute{Computed: true},
						"serial_number":   schema.StringAttribute{Computed: true},
						"device_capacity": schema.Float64Attribute{Computed: true},
						"model_name":      schema.StringAttribute{Computed: true},
						"os_version":      schema.StringAttribute{Computed: true},
						"mdm_enabled":     schema.BoolAttribute{Computed: true},
						"agent_installed": schema.BoolAttribute{Computed: true},
					},
				},
			},
		},
	}
}

func (d *prismDeviceInformationDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	d.client = req.ProviderData.(*client.Client)
}

func (d *prismDeviceInformationDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data prismDeviceInformationDataSourceModel

	var all []client.PrismEntry
	offset := 0
	limit := 300
	
	for {
		type prismResponse struct {
			Data []client.PrismEntry `json:"data"`
		}
		var listResp prismResponse
		
		path := fmt.Sprintf("/prism/device_information/?limit=%d&offset=%d", limit, offset)
		err := d.client.DoRequest(ctx, "GET", path, nil, &listResp)
		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read prism device_information, got error: %s", err))
			return
		}

		all = append(all, listResp.Data...)
		if len(listResp.Data) < limit {
			break
		}
		offset += limit
	}

	for _, item := range all {
		capacity, _ := item["device_capacity"].(float64)
		mdmEnabled, _ := item["mdm_enabled"].(bool)
		agentInstalled, _ := item["agent_installed"].(bool)

		data.Results = append(data.Results, prismDeviceInformationModel{
			DeviceID:       types.StringValue(fmt.Sprintf("%v", item["device_id"])),
			DeviceName:     types.StringValue(fmt.Sprintf("%v", item["device__name"])),
			SerialNumber:   types.StringValue(fmt.Sprintf("%v", item["serial_number"])),
			DeviceCapacity: types.Float64Value(capacity),
			ModelName:      types.StringValue(fmt.Sprintf("%v", item["model_name"])),
			OSVersion:      types.StringValue(fmt.Sprintf("%v", item["os_version"])),
			MDMEnabled:     types.BoolValue(mdmEnabled),
			AgentInstalled: types.BoolValue(agentInstalled),
		})
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
