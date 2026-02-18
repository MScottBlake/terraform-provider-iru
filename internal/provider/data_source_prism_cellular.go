package provider

import (
	"context"
	"fmt"

	"github.com/MScottBlake/terraform-provider-iru/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ datasource.DataSource = &prismCellularDataSource{}

func NewPrismCellularDataSource() datasource.DataSource {
	return &prismCellularDataSource{}
}

type prismCellularDataSource struct {
	client *client.Client
}

type prismCellularDataSourceModel struct {
	Results []prismCellularModel `tfsdk:"results"`
}

type prismCellularModel struct {
	DeviceID     types.String `tfsdk:"device_id"`
	DeviceName   types.String `tfsdk:"device_name"`
	SerialNumber types.String `tfsdk:"serial_number"`
	Carrier      types.String `tfsdk:"carrier"`
	PhoneNumber  types.String `tfsdk:"phone_number"`
	IMEI         types.String `tfsdk:"imei"`
	ICCID        types.String `tfsdk:"iccid"`
}

func (d *prismCellularDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_prism_cellular"
}

func (d *prismCellularDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "List cellular information from Prism.",
		Attributes: map[string]schema.Attribute{
			"results": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"device_id":     schema.StringAttribute{Computed: true},
						"device_name":   schema.StringAttribute{Computed: true},
						"serial_number": schema.StringAttribute{Computed: true},
						"carrier":       schema.StringAttribute{Computed: true},
						"phone_number":  schema.StringAttribute{Computed: true},
						"imei":          schema.StringAttribute{Computed: true},
						"iccid":         schema.StringAttribute{Computed: true},
					},
				},
			},
		},
	}
}

func (d *prismCellularDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	d.client = req.ProviderData.(*client.Client)
}

func (d *prismCellularDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data prismCellularDataSourceModel

	var all []client.PrismEntry
	offset := 0
	limit := 300
	
	for {
		type prismResponse struct {
			Data []client.PrismEntry `json:"data"`
		}
		var listResp prismResponse
		
		path := fmt.Sprintf("/prism/cellular/?limit=%d&offset=%d", limit, offset)
		err := d.client.DoRequest(ctx, "GET", path, nil, &listResp)
		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read prism cellular, got error: %s", err))
			return
		}

		all = append(all, listResp.Data...)
		
		if len(listResp.Data) < limit {
			break
		}
		offset += limit
	}

	for _, item := range all {
		carrier, _ := item["carrier"].(string)
		phone, _ := item["phone_number"].(string)
		imei, _ := item["imei"].(string)
		iccid, _ := item["iccid"].(string)

		data.Results = append(data.Results, prismCellularModel{
			DeviceID:     types.StringValue(fmt.Sprintf("%v", item["device_id"])),
			DeviceName:   types.StringValue(fmt.Sprintf("%v", item["device__name"])),
			SerialNumber: types.StringValue(fmt.Sprintf("%v", item["serial_number"])),
			Carrier:      types.StringValue(carrier),
			PhoneNumber:  types.StringValue(phone),
			IMEI:         types.StringValue(imei),
			ICCID:        types.StringValue(iccid),
		})
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
