package provider

import (
	"bytes"
	"context"
	"fmt"

	"github.com/MScottBlake/terraform-provider-iru/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ resource.Resource = &adeIntegrationResource{}
var _ resource.ResourceWithImportState = &adeIntegrationResource{}

func NewADEIntegrationResource() resource.Resource {
	return &adeIntegrationResource{}
}

type adeIntegrationResource struct {
	client *client.Client
}

type adeIntegrationResourceModel struct {
	ID                 types.String `tfsdk:"id"`
	BlueprintID        types.String `tfsdk:"blueprint_id"`
	Phone              types.String `tfsdk:"phone"`
	Email              types.String `tfsdk:"email"`
	MDMServerTokenFile types.String `tfsdk:"mdm_server_token_file"`
}

func (r *adeIntegrationResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_ade_integration"
}

func (r *adeIntegrationResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a Kandji ADE Integration.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The unique identifier for the ADE Integration.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"blueprint_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The UUID of the default blueprint to associate with the integration.",
			},
			"phone": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "A phone number for the integration.",
			},
			"email": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "An email address for the integration.",
			},
			"mdm_server_token_file": schema.StringAttribute{
				Required:            true,
				Sensitive:           true,
				MarkdownDescription: "The content of the MDM server token file (.p7m) downloaded from Apple Business Manager.",
			},
		},
	}
}

func (r *adeIntegrationResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *client.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	r.client = client
}

func (r *adeIntegrationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data adeIntegrationResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	fields := map[string]string{
		"blueprint_id": data.BlueprintID.ValueString(),
		"phone":        data.Phone.ValueString(),
		"email":        data.Email.ValueString(),
	}

	fileContent := []byte(data.MDMServerTokenFile.ValueString())
	
	var adeResponse client.ADEIntegration
	err := r.client.DoMultipartRequest(ctx, "POST", "/integrations/apple/ade/", fields, "file", "token.p7m", bytes.NewReader(fileContent), &adeResponse)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create ADE integration, got error: %s", err))
		return
	}

	data.ID = types.StringValue(adeResponse.ID)
	// Response contains blueprint object, but we store blueprint_id
	if adeResponse.Blueprint != nil {
		data.BlueprintID = types.StringValue(adeResponse.Blueprint.ID)
	}
	data.Phone = types.StringValue(adeResponse.Phone)
	data.Email = types.StringValue(adeResponse.Email)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *adeIntegrationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data adeIntegrationResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var adeResponse client.ADEIntegration
	err := r.client.DoRequest(ctx, "GET", "/integrations/apple/ade/"+data.ID.ValueString(), nil, &adeResponse)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read ADE integration, got error: %s", err))
		return
	}

	if adeResponse.Blueprint != nil {
		data.BlueprintID = types.StringValue(adeResponse.Blueprint.ID)
	}
	data.Phone = types.StringValue(adeResponse.Phone)
	data.Email = types.StringValue(adeResponse.Email)
	// MDMServerTokenFile is not returned by API, keep from state (handled by Terraform automatically if not set here?)
	// Actually, if we don't set it, it might show as null if we overwrite `data` completely.
	// But we are updating `data` fields. `data.MDMServerTokenFile` preserves previous value if we don't touch it.

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *adeIntegrationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state adeIntegrationResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if !plan.MDMServerTokenFile.Equal(state.MDMServerTokenFile) {
		// Token changed, use Renew endpoint
		fields := map[string]string{
			"blueprint_id": plan.BlueprintID.ValueString(),
			"phone":        plan.Phone.ValueString(),
			"email":        plan.Email.ValueString(),
		}
		fileContent := []byte(plan.MDMServerTokenFile.ValueString())
		
		var adeResponse client.ADEIntegration
		err := r.client.DoMultipartRequest(ctx, "POST", "/integrations/apple/ade/"+plan.ID.ValueString()+"/renew", fields, "file", "token.p7m", bytes.NewReader(fileContent), &adeResponse)
		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to renew ADE integration, got error: %s", err))
			return
		}
		
		if adeResponse.Blueprint != nil {
			plan.BlueprintID = types.StringValue(adeResponse.Blueprint.ID)
		}
		plan.Phone = types.StringValue(adeResponse.Phone)
		plan.Email = types.StringValue(adeResponse.Email)
	} else {
		// Normal update
		updateRequest := client.ADEIntegration{
			BlueprintID: plan.BlueprintID.ValueString(),
			Phone:       plan.Phone.ValueString(),
			Email:       plan.Email.ValueString(),
		}
		
		var adeResponse client.ADEIntegration
		err := r.client.DoRequest(ctx, "PATCH", "/integrations/apple/ade/"+plan.ID.ValueString(), updateRequest, &adeResponse)
		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update ADE integration, got error: %s", err))
			return
		}

		if adeResponse.Blueprint != nil {
			plan.BlueprintID = types.StringValue(adeResponse.Blueprint.ID)
		}
		plan.Phone = types.StringValue(adeResponse.Phone)
		plan.Email = types.StringValue(adeResponse.Email)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *adeIntegrationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data adeIntegrationResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DoRequest(ctx, "DELETE", "/integrations/apple/ade/"+data.ID.ValueString(), nil, nil)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete ADE integration, got error: %s", err))
		return
	}
}

func (r *adeIntegrationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
