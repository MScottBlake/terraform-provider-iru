package provider

import (
	"context"
	"fmt"

	"github.com/MScottBlake/terraform-provider-iru/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ datasource.DataSource = &prismTransparencyDatabaseDataSource{}

func NewPrismTransparencyDatabaseDataSource() datasource.DataSource {
	return &prismTransparencyDatabaseDataSource{}
}

type prismTransparencyDatabaseDataSource struct {
	client *client.Client
}

type prismTransparencyDatabaseDataSourceModel struct {
	Results []prismTransparencyDatabaseModel `tfsdk:"results"`
}

type prismTransparencyDatabaseModel struct {
	DeviceID     types.String `tfsdk:"device_id"`
	DeviceName   types.String `tfsdk:"device_name"`
	SerialNumber types.String `tfsdk:"serial_number"`
	Service      types.String `tfsdk:"service"`
	Application  types.String `tfsdk:"application"`
	Status       types.String `tfsdk:"status"`
	LocalUser    types.String `tfsdk:"local_user"`
}

func (d *prismTransparencyDatabaseDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_prism_transparency_database"
}

func (d *prismTransparencyDatabaseDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "List TCC (Transparency, Consent, and Control) database entries from Prism.",
		Attributes: map[string]schema.Attribute{
			"results": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"device_id":     schema.StringAttribute{Computed: true},
						"device_name":   schema.StringAttribute{Computed: true},
						"serial_number": schema.StringAttribute{Computed: true},
						"service":       schema.StringAttribute{Computed: true},
						"application":   schema.StringAttribute{Computed: true},
						"status":        schema.StringAttribute{Computed: true},
						"local_user":    schema.StringAttribute{Computed: true},
					},
				},
			},
		},
	}
}

func (d *prismTransparencyDatabaseDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	d.client = req.ProviderData.(*client.Client)
}

func (d *prismTransparencyDatabaseDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data prismTransparencyDatabaseDataSourceModel

	var all []client.PrismEntry
	offset := 0
	limit := 300
	
	for {
		type prismResponse struct {
			Data []client.PrismEntry `json:"data"`
		}
		var listResp prismResponse
		
		path := fmt.Sprintf("/prism/transparency_database/?limit=%d&offset=%d", limit, offset)
		err := d.client.DoRequest(ctx, "GET", path, nil, &listResp)
		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read prism transparency_database, got error: %s", err))
			return
		}

		all = append(all, listResp.Data...)
		if len(listResp.Data) < limit {
			break
		}
		offset += limit
	}

	for _, item := range all {
		service, _ := item["service"].(string)
		app, _ := item["application"].(string)
		status, _ := item["status"].(string)
		user, _ := item["local_user"].(string)

		data.Results = append(data.Results, prismTransparencyDatabaseModel{
			DeviceID:     types.StringValue(fmt.Sprintf("%v", item["device_id"])),
			DeviceName:   types.StringValue(fmt.Sprintf("%v", item["device__name"])),
			SerialNumber: types.StringValue(fmt.Sprintf("%v", item["serial_number"])),
			Service:      types.StringValue(service),
			Application:  types.StringValue(app),
			Status:       types.StringValue(status),
			LocalUser:    types.StringValue(user),
		})
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
