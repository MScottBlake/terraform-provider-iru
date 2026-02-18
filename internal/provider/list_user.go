package provider

import (
	"context"

	"github.com/MScottBlake/terraform-provider-iru/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/list"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

var _ list.ListResource = &userListResource{}

func NewUserListResource() list.ListResource {
	return &userListResource{}
}

type userListResource struct {
	client *client.Client
}

func (r *userListResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_user"
}

func (r *userListResource) ListResourceConfigSchema(ctx context.Context, req list.ListResourceSchemaRequest, resp *list.ListResourceSchemaResponse) {
}

func (r *userListResource) Configure(ctx context.Context, req list.ConfigureRequest, resp *list.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.client = req.ProviderData.(*client.Client)
}

func (r *userListResource) List(ctx context.Context, req list.ListRequest, resp *list.ListResultsStream) {
	resp.Results = list.ListResultsStreamDiagnostics(nil)
}
