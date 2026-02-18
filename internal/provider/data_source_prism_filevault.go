package provider

import (
	"context"
	"fmt"

	"github.com/MScottBlake/terraform-provider-iru/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ datasource.DataSource = &prismFileVaultDataSource{}

func NewPrismFileVaultDataSource() datasource.DataSource {
	return &prismFileVaultDataSource{}
}

type prismFileVaultDataSource struct {
	client *client.Client
}

type prismFileVaultDataSourceModel struct {
	Results []prismFileVaultModel `tfsdk:"results"`
}

type prismFileVaultModel struct {
	DeviceID     types.String `tfsdk:"device_id"`
	DeviceName   types.String `tfsdk:"device_name"`
	SerialNumber types.String `tfsdk:"serial_number"`
	Status       types.Bool   `tfsdk:"status"`
	KeyType      types.String `tfsdk:"key_type"`
	KeyEscrowed  types.Bool   `tfsdk:"key_escrowed"`
}

func (d *prismFileVaultDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_prism_filevault"
}

func (d *prismFileVaultDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "List FileVault status for macOS devices from Prism.",
		Attributes: map[string]schema.Attribute{
			"results": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"device_id": schema.StringAttribute{
							Computed: true,
						},
						"device_name": schema.StringAttribute{
							Computed: true,
						},
						"serial_number": schema.StringAttribute{
							Computed: true,
						},
						"status": schema.BoolAttribute{
							Computed: true,
						},
						"key_type": schema.StringAttribute{
							Computed: true,
						},
						"key_escrowed": schema.BoolAttribute{
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func (d *prismFileVaultDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	d.client = req.ProviderData.(*client.Client)
}

func (d *prismFileVaultDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data prismFileVaultDataSourceModel

	var all []client.PrismFileVault
	offset := 0
	limit := 300
	
	for {
		type prismResponse struct {
			Data []client.PrismFileVault `json:"data"`
		}
		var listResp prismResponse
		
		path := fmt.Sprintf("/prism/filevault/?limit=%d&offset=%d", limit, offset)
		err := d.client.DoRequest(ctx, "GET", path, nil, &listResp)
		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read prism filevault, got error: %s", err))
			return
		}

		all = append(all, listResp.Data...)
		
		if len(listResp.Data) < limit {
			break
		}
		offset += limit
	}

	for _, item := range all {
		data.Results = append(data.Results, prismFileVaultModel{
			DeviceID:     types.StringValue(item.DeviceID),
			DeviceName:   types.StringValue(item.DeviceName),
			SerialNumber: types.StringValue(item.SerialNumber),
			Status:       types.BoolValue(item.Status),
			KeyType:      types.StringValue(item.KeyType),
			KeyEscrowed:  types.BoolValue(item.KeyEscrowed),
		})
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
