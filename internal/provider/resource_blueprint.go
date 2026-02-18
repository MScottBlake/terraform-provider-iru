package provider

import (
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

var _ resource.Resource = &blueprintResource{}
var _ resource.ResourceWithImportState = &blueprintResource{}

func NewBlueprintResource() resource.Resource {
	return &blueprintResource{}
}

type blueprintResource struct {
	client *client.Client
}

type blueprintResourceModel struct {
	ID             types.String `tfsdk:"id"`
	Name           types.String `tfsdk:"name"`
	Description    types.String `tfsdk:"description"`
	Icon           types.String `tfsdk:"icon"`
	Color          types.String `tfsdk:"color"`
	EnrollmentCode types.String `tfsdk:"enrollment_code"`
}

func (r *blueprintResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_blueprint"
}

func (r *blueprintResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a Kandji Blueprint.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The unique identifier for the Blueprint.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The name of the Blueprint.",
			},
			"description": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The description of the Blueprint.",
			},
			"icon": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "The icon of the Blueprint.",
			},
			"color": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "The color of the Blueprint.",
			},
			"enrollment_code": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The enrollment code for the Blueprint.",
			},
		},
	}
}

func (r *blueprintResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *blueprintResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data blueprintResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	blueprintRequest := client.Blueprint{
		Name:        data.Name.ValueString(),
		Description: data.Description.ValueString(),
		Icon:        data.Icon.ValueString(),
		Color:       data.Color.ValueString(),
	}

	var blueprintResponse client.Blueprint
	err := r.client.DoRequest(ctx, "POST", "/blueprints", blueprintRequest, &blueprintResponse)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create blueprint, got error: %s", err))
		return
	}

	data.ID = types.StringValue(blueprintResponse.ID)
	data.Name = types.StringValue(blueprintResponse.Name)
	data.Description = types.StringValue(blueprintResponse.Description)
	data.Icon = types.StringValue(blueprintResponse.Icon)
	data.Color = types.StringValue(blueprintResponse.Color)
	data.EnrollmentCode = types.StringValue(blueprintResponse.EnrollmentCode)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *blueprintResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data blueprintResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var blueprintResponse client.Blueprint
	err := r.client.DoRequest(ctx, "GET", "/blueprints/"+data.ID.ValueString(), nil, &blueprintResponse)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read blueprint, got error: %s", err))
		return
	}

	data.Name = types.StringValue(blueprintResponse.Name)
	data.Description = types.StringValue(blueprintResponse.Description)
	data.Icon = types.StringValue(blueprintResponse.Icon)
	data.Color = types.StringValue(blueprintResponse.Color)
	data.EnrollmentCode = types.StringValue(blueprintResponse.EnrollmentCode)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *blueprintResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data blueprintResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	blueprintRequest := client.Blueprint{
		Name:        data.Name.ValueString(),
		Description: data.Description.ValueString(),
		Icon:        data.Icon.ValueString(),
		Color:       data.Color.ValueString(),
	}

	var blueprintResponse client.Blueprint
	err := r.client.DoRequest(ctx, "PATCH", "/blueprints/"+data.ID.ValueString(), blueprintRequest, &blueprintResponse)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update blueprint, got error: %s", err))
		return
	}

	data.Name = types.StringValue(blueprintResponse.Name)
	data.Description = types.StringValue(blueprintResponse.Description)
	data.Icon = types.StringValue(blueprintResponse.Icon)
	data.Color = types.StringValue(blueprintResponse.Color)
	data.EnrollmentCode = types.StringValue(blueprintResponse.EnrollmentCode)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *blueprintResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data blueprintResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DoRequest(ctx, "DELETE", "/blueprints/"+data.ID.ValueString(), nil, nil)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete blueprint, got error: %s", err))
		return
	}
}

func (r *blueprintResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
