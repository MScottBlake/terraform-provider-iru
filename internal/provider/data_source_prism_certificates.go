package provider

import (
	"context"
	"fmt"

	"github.com/MScottBlake/terraform-provider-iru/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ datasource.DataSource = &prismCertificatesDataSource{}

func NewPrismCertificatesDataSource() datasource.DataSource {
	return &prismCertificatesDataSource{}
}

type prismCertificatesDataSource struct {
	client *client.Client
}

type prismCertificatesDataSourceModel struct {
	Results []prismCertificateModel `tfsdk:"results"`
}

type prismCertificateModel struct {
	DeviceID            types.String `tfsdk:"device_id"`
	DeviceName          types.String `tfsdk:"device_name"`
	SerialNumber        types.String `tfsdk:"serial_number"`
	CommonName          types.String `tfsdk:"common_name"`
	IdentityCertificate types.Bool   `tfsdk:"identity_certificate"`
}

func (d *prismCertificatesDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_prism_certificates"
}

func (d *prismCertificatesDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "List certificates from Prism.",
		Attributes: map[string]schema.Attribute{
			"results": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"device_id":            schema.StringAttribute{Computed: true},
						"device_name":          schema.StringAttribute{Computed: true},
						"serial_number":        schema.StringAttribute{Computed: true},
						"common_name":          schema.StringAttribute{Computed: true},
						"identity_certificate": schema.BoolAttribute{Computed: true},
					},
				},
			},
		},
	}
}

func (d *prismCertificatesDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	d.client = req.ProviderData.(*client.Client)
}

func (d *prismCertificatesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data prismCertificatesDataSourceModel

	var all []client.PrismEntry
	offset := 0
	limit := 300
	
	for {
		type prismResponse struct {
			Data []client.PrismEntry `json:"data"`
		}
		var listResp prismResponse
		
		path := fmt.Sprintf("/prism/certificates/?limit=%d&offset=%d", limit, offset)
		err := d.client.DoRequest(ctx, "GET", path, nil, &listResp)
		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read prism certificates, got error: %s", err))
			return
		}

		all = append(all, listResp.Data...)
		
		if len(listResp.Data) < limit {
			break
		}
		offset += limit
	}

	for _, item := range all {
		commonName, _ := item["common_name"].(string)
		identityCert, _ := item["identity_certificate"].(bool)

		data.Results = append(data.Results, prismCertificateModel{
			DeviceID:            types.StringValue(fmt.Sprintf("%v", item["device_id"])),
			DeviceName:          types.StringValue(fmt.Sprintf("%v", item["device__name"])),
			SerialNumber:        types.StringValue(fmt.Sprintf("%v", item["serial_number"])),
			CommonName:          types.StringValue(commonName),
			IdentityCertificate: types.BoolValue(identityCert),
		})
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
