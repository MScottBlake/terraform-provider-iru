package provider

import (
	"context"
	"fmt"

	"github.com/MScottBlake/terraform-provider-iru/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ datasource.DataSource = &prismLocalUsersDataSource{}

func NewPrismLocalUsersDataSource() datasource.DataSource {
	return &prismLocalUsersDataSource{}
}

type prismLocalUsersDataSource struct {
	client *client.Client
}

type prismLocalUsersDataSourceModel struct {
	Results []prismLocalUserModel `tfsdk:"results"`
}

type prismLocalUserModel struct {
	DeviceID      types.String `tfsdk:"device_id"`
	DeviceName    types.String `tfsdk:"device_name"`
	SerialNumber  types.String `tfsdk:"serial_number"`
	Username      types.String `tfsdk:"username"`
	FullName      types.String `tfsdk:"full_name"`
	UserType      types.String `tfsdk:"user_type"`
	UID           types.Int64  `tfsdk:"uid"`
	LoggedIn      types.Bool   `tfsdk:"logged_in"`
	HiddenUser    types.Bool   `tfsdk:"hidden_user"`
	FileVaultUser types.Bool   `tfsdk:"filevault_user"`
}

func (d *prismLocalUsersDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_prism_local_users"
}

func (d *prismLocalUsersDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "List local users from Prism.",
		Attributes: map[string]schema.Attribute{
			"results": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"device_id":      schema.StringAttribute{Computed: true},
						"device_name":    schema.StringAttribute{Computed: true},
						"serial_number":  schema.StringAttribute{Computed: true},
						"username":       schema.StringAttribute{Computed: true},
						"full_name":      schema.StringAttribute{Computed: true},
						"user_type":      schema.StringAttribute{Computed: true},
						"uid":            schema.Int64Attribute{Computed: true},
						"logged_in":      schema.BoolAttribute{Computed: true},
						"hidden_user":    schema.BoolAttribute{Computed: true},
						"filevault_user": schema.BoolAttribute{Computed: true},
					},
				},
			},
		},
	}
}

func (d *prismLocalUsersDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	d.client = req.ProviderData.(*client.Client)
}

func (d *prismLocalUsersDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data prismLocalUsersDataSourceModel

	var all []client.PrismEntry
	offset := 0
	limit := 300
	
	for {
		type prismResponse struct {
			Data []client.PrismEntry `json:"data"`
		}
		var listResp prismResponse
		
		path := fmt.Sprintf("/prism/local_users?limit=%d&offset=%d", limit, offset)
		err := d.client.DoRequest(ctx, "GET", path, nil, &listResp)
		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read prism local_users, got error: %s", err))
			return
		}

		all = append(all, listResp.Data...)
		
		if len(listResp.Data) < limit {
			break
		}
		offset += limit
	}

	for _, item := range all {
		username, _ := item["username"].(string)
		fullName, _ := item["full_name"].(string)
		userType, _ := item["type"].(string)
		uid, _ := item["uid"].(float64)
		loggedIn, _ := item["logged_in"].(bool)
		hidden, _ := item["hidden_user"].(bool)
		fvUser, _ := item["filevault_user"].(bool)

		data.Results = append(data.Results, prismLocalUserModel{
			DeviceID:      types.StringValue(fmt.Sprintf("%v", item["device_id"])),
			DeviceName:    types.StringValue(fmt.Sprintf("%v", item["device__name"])),
			SerialNumber:  types.StringValue(fmt.Sprintf("%v", item["serial_number"])),
			Username:      types.StringValue(username),
			FullName:      types.StringValue(fullName),
			UserType:      types.StringValue(userType),
			UID:           types.Int64Value(int64(uid)),
			LoggedIn:      types.BoolValue(loggedIn),
			HiddenUser:    types.BoolValue(hidden),
			FileVaultUser: types.BoolValue(fvUser),
		})
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
