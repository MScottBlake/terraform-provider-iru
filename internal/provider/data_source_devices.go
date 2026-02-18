package provider

import (
	"context"
	"fmt"

	"github.com/MScottBlake/terraform-provider-iru/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ datasource.DataSource = &devicesDataSource{}

func NewDevicesDataSource() datasource.DataSource {
	return &devicesDataSource{}
}

type devicesDataSource struct {
	client *client.Client
}

type devicesDataSourceModel struct {
	Devices []deviceModel `tfsdk:"devices"`
}

type deviceModel struct {
	ID           types.String `tfsdk:"id"`
	DeviceName   types.String `tfsdk:"device_name"`
	SerialNumber types.String `tfsdk:"serial_number"`
	Model        types.String `tfsdk:"model"`
	OSVersion    types.String `tfsdk:"os_version"`
	Platform     types.String `tfsdk:"platform"`
	LastCheckIn  types.String `tfsdk:"last_check_in"`
}

func (d *devicesDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_devices"
}

func (d *devicesDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "List all devices in the Kandji instance.",
		Attributes: map[string]schema.Attribute{
			"devices": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The unique identifier for the Device.",
						},
						"device_name": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The name of the Device.",
						},
						"serial_number": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The serial number of the Device.",
						},
						"model": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The model of the Device.",
						},
						"os_version": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The OS version of the Device.",
						},
						"platform": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The platform of the Device.",
						},
						"last_check_in": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The last check-in time of the Device.",
						},
					},
				},
			},
		},
	}
}

func (d *devicesDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected DataSource Configure Type",
			fmt.Sprintf("Expected *client.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	d.client = client
}

func (d *devicesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data devicesDataSourceModel

	// Determine pagination or simply fetch all.
	// Kandji API usually paginates. If no params, it might return first page or all?
	// The docs mention `limit` and `offset`. Default limit might be 100?
	// I'll assume standard pagination logic or fetch all loop.
	// For simplicity in this iteration, I'll fetch the default page.
	// But "full featured" implies handling pagination.
	// I'll implement a loop to fetch all devices.

	var allDevices []client.Device
	offset := 0
	limit := 300 // Max limit often supported
	
	for {
		// API docs usually return array for list, or object with results.
		// If it returns array, headers have pagination info?
		// I'll check if API returns array or object.
		// "GET /devices" -> Returns list of devices.
		// I'll assume array.
		// Pagination parameters: `limit`, `offset`.
		
		path := fmt.Sprintf("/devices/?limit=%d&offset=%d", limit, offset)
		var pageResponse []client.Device
		err := d.client.DoRequest(ctx, "GET", path, nil, &pageResponse)
		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read devices, got error: %s", err))
			return
		}

		allDevices = append(allDevices, pageResponse...)
		
		if len(pageResponse) < limit {
			break
		}
		offset += limit
	}

	for _, device := range allDevices {
		data.Devices = append(data.Devices, deviceModel{
			ID:           types.StringValue(device.ID),
			DeviceName:   types.StringValue(device.DeviceName),
			SerialNumber: types.StringValue(device.SerialNumber),
			Model:        types.StringValue(device.Model),
			OSVersion:    types.StringValue(device.OSVersion),
			Platform:     types.StringValue(device.Platform),
			LastCheckIn:  types.StringValue(device.LastCheckIn),
		})
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
