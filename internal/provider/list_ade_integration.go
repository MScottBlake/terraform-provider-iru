package provider

import (
	"context"

	"github.com/MScottBlake/terraform-provider-iru/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/list"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

var _ list.ListResource = &adeIntegrationListResource{}

func NewADEIntegrationListResource() list.ListResource {
	return &adeIntegrationListResource{}
}

type adeIntegrationListResource struct {
	client *client.Client
}

func (r *adeIntegrationListResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_ade_integration"
}

func (r *adeIntegrationListResource) ListResourceConfigSchema(ctx context.Context, req list.ListResourceSchemaRequest, resp *list.ListResourceSchemaResponse) {
}

func (r *adeIntegrationListResource) Configure(ctx context.Context, req list.ConfigureRequest, resp *list.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.client = req.ProviderData.(*client.Client)
}

func (r *adeIntegrationListResource) List(ctx context.Context, req list.ListRequest, resp *list.ListResultsStream) {
	resp.Results = list.ListResultsStreamDiagnostics(nil)
}
