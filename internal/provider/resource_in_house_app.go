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

var _ resource.Resource = &inHouseAppResource{}
var _ resource.ResourceWithImportState = &inHouseAppResource{}

func NewInHouseAppResource() resource.Resource {
	return &inHouseAppResource{}
}

type inHouseAppResource struct {
	client *client.Client
}

type inHouseAppResourceModel struct {
	ID           types.String `tfsdk:"id"`
	Name         types.String `tfsdk:"name"`
	FileKey      types.String `tfsdk:"file_key"`
	RunsOnIPhone types.Bool   `tfsdk:"runs_on_iphone"`
	RunsOnIPad   types.Bool   `tfsdk:"runs_on_ipad"`
	RunsOnTV     types.Bool   `tfsdk:"runs_on_tv"`
	Active       types.Bool   `tfsdk:"active"`
}

func (r *inHouseAppResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_in_house_app"
}

func (r *inHouseAppResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a Kandji In-House App (.ipa) library item. Note: You must handle the file upload to S3 out-of-band and provide the file_key.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The unique identifier for the In-House App.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The name for this In-House App.",
			},
			"file_key": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The S3 key from the upload endpoint.",
			},
			"runs_on_iphone": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
			},
			"runs_on_ipad": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
			},
			"runs_on_tv": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
			},
			"active": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
			},
		},
	}
}

func (r *inHouseAppResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.client = req.ProviderData.(*client.Client)
}

func (r *inHouseAppResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data inHouseAppResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	appRequest := r.mapToClient(&data)
	var appResponse client.InHouseApp
	err := r.client.DoRequest(ctx, "POST", "/library/ipa-apps", appRequest, &appResponse)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create in-house app, got error: %s", err))
		return
	}

	r.mapFromClient(&data, &appResponse)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *inHouseAppResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data inHouseAppResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var appResponse client.InHouseApp
	err := r.client.DoRequest(ctx, "GET", "/library/ipa-apps/"+data.ID.ValueString(), nil, &appResponse)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read in-house app, got error: %s", err))
		return
	}

	r.mapFromClient(&data, &appResponse)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *inHouseAppResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data inHouseAppResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	appRequest := r.mapToClient(&data)
	var appResponse client.InHouseApp
	err := r.client.DoRequest(ctx, "PATCH", "/library/ipa-apps/"+data.ID.ValueString(), appRequest, &appResponse)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update in-house app, got error: %s", err))
		return
	}

	r.mapFromClient(&data, &appResponse)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *inHouseAppResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data inHouseAppResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DoRequest(ctx, "DELETE", "/library/ipa-apps/"+data.ID.ValueString(), nil, nil)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete in-house app, got error: %s", err))
		return
	}
}

func (r *inHouseAppResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *inHouseAppResource) mapToClient(data *inHouseAppResourceModel) client.InHouseApp {
	return client.InHouseApp{
		Name:         data.Name.ValueString(),
		FileKey:      data.FileKey.ValueString(),
		RunsOnIPhone: data.RunsOnIPhone.ValueBool(),
		RunsOnIPad:   data.RunsOnIPad.ValueBool(),
		RunsOnTV:     data.RunsOnTV.ValueBool(),
		Active:       data.Active.ValueBool(),
	}
}

func (r *inHouseAppResource) mapFromClient(data *inHouseAppResourceModel, resp *client.InHouseApp) {
	data.ID = types.StringValue(resp.ID)
	data.Name = types.StringValue(resp.Name)
	data.FileKey = types.StringValue(resp.FileKey)
	data.RunsOnIPhone = types.BoolValue(resp.RunsOnIPhone)
	data.RunsOnIPad = types.BoolValue(resp.RunsOnIPad)
	data.RunsOnTV = types.BoolValue(resp.RunsOnTV)
	data.Active = types.BoolValue(resp.Active)
}
