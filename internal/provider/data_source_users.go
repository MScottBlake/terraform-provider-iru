package provider

import (
	"context"
	"fmt"

	"github.com/MScottBlake/terraform-provider-iru/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ datasource.DataSource = &usersDataSource{}

func NewUsersDataSource() datasource.DataSource {
	return &usersDataSource{}
}

type usersDataSource struct {
	client *client.Client
}

type usersDataSourceModel struct {
	Users []userDataSourceModel `tfsdk:"users"`
}

type userDataSourceModel struct {
	ID         types.String `tfsdk:"id"`
	Name       types.String `tfsdk:"name"`
	Email      types.String `tfsdk:"email"`
	IsArchived types.Bool   `tfsdk:"is_archived"`
}

func (d *usersDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_users"
}

func (d *usersDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "List all users in the Kandji instance.",
		Attributes: map[string]schema.Attribute{
			"users": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The unique identifier for the User.",
						},
						"name": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The name of the User.",
						},
						"email": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The email of the User.",
						},
						"is_archived": schema.BoolAttribute{
							Computed:            true,
							MarkdownDescription: "Whether the user is archived.",
						},
					},
				},
			},
		},
	}
}

func (d *usersDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *usersDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data usersDataSourceModel

	var allUsers []client.User
	offset := 0
	limit := 300
	
	for {
		type listUsersResponse struct {
			Results []client.User `json:"results"`
		}
		var listResp listUsersResponse
		
		path := fmt.Sprintf("/users/?limit=%d&offset=%d", limit, offset)
		err := d.client.DoRequest(ctx, "GET", path, nil, &listResp)
		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read users, got error: %s", err))
			return
		}

		allUsers = append(allUsers, listResp.Results...)
		
		if len(listResp.Results) < limit {
			break
		}
		offset += limit
	}

	for _, user := range allUsers {
		data.Users = append(data.Users, userDataSourceModel{
			ID:         types.StringValue(user.ID),
			Name:       types.StringValue(user.Name),
			Email:      types.StringValue(user.Email),
			IsArchived: types.BoolValue(user.IsArchived),
		})
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
