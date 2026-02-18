package provider

import (
	"context"

	"github.com/MScottBlake/terraform-provider-iru/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/list"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

var _ list.ListResource = &inHouseAppListResource{}

func NewInHouseAppListResource() list.ListResource {
	return &inHouseAppListResource{}
}

type inHouseAppListResource struct {
	client *client.Client
}

func (r *inHouseAppListResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_in_house_app"
}

func (r *inHouseAppListResource) ListResourceConfigSchema(ctx context.Context, req list.ListResourceSchemaRequest, resp *list.ListResourceSchemaResponse) {
}

func (r *inHouseAppListResource) Configure(ctx context.Context, req list.ConfigureRequest, resp *list.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.client = req.ProviderData.(*client.Client)
}

func (r *inHouseAppListResource) List(ctx context.Context, req list.ListRequest, resp *list.ListResultsStream) {
	resp.Results = list.ListResultsStreamDiagnostics(nil)
}
