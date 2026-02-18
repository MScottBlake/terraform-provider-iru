package provider

import (
	"context"

	"github.com/MScottBlake/terraform-provider-iru/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/list"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

var _ list.ListResource = &customAppListResource{}

func NewCustomAppListResource() list.ListResource {
	return &customAppListResource{}
}

type customAppListResource struct {
	client *client.Client
}

func (r *customAppListResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_custom_app"
}

func (r *customAppListResource) ListResourceConfigSchema(ctx context.Context, req list.ListResourceSchemaRequest, resp *list.ListResourceSchemaResponse) {
}

func (r *customAppListResource) Configure(ctx context.Context, req list.ConfigureRequest, resp *list.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.client = req.ProviderData.(*client.Client)
}

func (r *customAppListResource) List(ctx context.Context, req list.ListRequest, resp *list.ListResultsStream) {
	resp.Results = list.ListResultsStreamDiagnostics(nil)
}
