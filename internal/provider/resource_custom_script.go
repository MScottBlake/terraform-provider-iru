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

var _ resource.Resource = &customScriptResource{}
var _ resource.ResourceWithImportState = &customScriptResource{}

func NewCustomScriptResource() resource.Resource {
	return &customScriptResource{}
}

type customScriptResource struct {
	client *client.Client
}

type customScriptResourceModel struct {
	ID                 types.String `tfsdk:"id"`
	Name               types.String `tfsdk:"name"`
	Active             types.Bool   `tfsdk:"active"`
	ExecutionFrequency types.String `tfsdk:"execution_frequency"`
	Restart            types.Bool   `tfsdk:"restart"`
	Script             types.String `tfsdk:"script"`
	RemediationScript  types.String `tfsdk:"remediation_script"`
	ShowInSelfService  types.Bool   `tfsdk:"show_in_self_service"`
}

func (r *customScriptResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_custom_script"
}

func (r *customScriptResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a Kandji Custom Script library item.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The unique identifier for the Custom Script.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The name of the Custom Script.",
			},
			"active": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Whether this library item is active.",
			},
			"execution_frequency": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The execution frequency. Valid values: once, every_15_min, every_day, no_enforcement.",
			},
			"restart": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Whether to restart the computer after script execution.",
			},
			"script": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The content of the script.",
			},
			"remediation_script": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "The content of the remediation script.",
			},
			"show_in_self_service": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Whether to show the script in Self Service.",
			},
		},
	}
}

func (r *customScriptResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *customScriptResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data customScriptResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	scriptRequest := client.CustomScript{
		Name:               data.Name.ValueString(),
		Active:             data.Active.ValueBool(),
		ExecutionFrequency: data.ExecutionFrequency.ValueString(),
		Restart:            data.Restart.ValueBool(),
		Script:             data.Script.ValueString(),
		RemediationScript:  data.RemediationScript.ValueString(),
		ShowInSelfService:  data.ShowInSelfService.ValueBool(),
	}

	var scriptResponse client.CustomScript
	err := r.client.DoRequest(ctx, "POST", "/library/custom-scripts", scriptRequest, &scriptResponse)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create custom script, got error: %s", err))
		return
	}

	data.ID = types.StringValue(scriptResponse.ID)
	data.Name = types.StringValue(scriptResponse.Name)
	data.Active = types.BoolValue(scriptResponse.Active)
	data.ExecutionFrequency = types.StringValue(scriptResponse.ExecutionFrequency)
	data.Restart = types.BoolValue(scriptResponse.Restart)
	data.Script = types.StringValue(scriptResponse.Script)
	data.RemediationScript = types.StringValue(scriptResponse.RemediationScript)
	data.ShowInSelfService = types.BoolValue(scriptResponse.ShowInSelfService)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *customScriptResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data customScriptResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var scriptResponse client.CustomScript
	err := r.client.DoRequest(ctx, "GET", "/library/custom-scripts/"+data.ID.ValueString(), nil, &scriptResponse)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read custom script, got error: %s", err))
		return
	}

	data.Name = types.StringValue(scriptResponse.Name)
	data.Active = types.BoolValue(scriptResponse.Active)
	data.ExecutionFrequency = types.StringValue(scriptResponse.ExecutionFrequency)
	data.Restart = types.BoolValue(scriptResponse.Restart)
	data.Script = types.StringValue(scriptResponse.Script)
	data.RemediationScript = types.StringValue(scriptResponse.RemediationScript)
	data.ShowInSelfService = types.BoolValue(scriptResponse.ShowInSelfService)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *customScriptResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data customScriptResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	scriptRequest := client.CustomScript{
		Name:               data.Name.ValueString(),
		Active:             data.Active.ValueBool(),
		ExecutionFrequency: data.ExecutionFrequency.ValueString(),
		Restart:            data.Restart.ValueBool(),
		Script:             data.Script.ValueString(),
		RemediationScript:  data.RemediationScript.ValueString(),
		ShowInSelfService:  data.ShowInSelfService.ValueBool(),
	}

	var scriptResponse client.CustomScript
	err := r.client.DoRequest(ctx, "PATCH", "/library/custom-scripts/"+data.ID.ValueString(), scriptRequest, &scriptResponse)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update custom script, got error: %s", err))
		return
	}

	data.Name = types.StringValue(scriptResponse.Name)
	data.Active = types.BoolValue(scriptResponse.Active)
	data.ExecutionFrequency = types.StringValue(scriptResponse.ExecutionFrequency)
	data.Restart = types.BoolValue(scriptResponse.Restart)
	data.Script = types.StringValue(scriptResponse.Script)
	data.RemediationScript = types.StringValue(scriptResponse.RemediationScript)
	data.ShowInSelfService = types.BoolValue(scriptResponse.ShowInSelfService)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *customScriptResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data customScriptResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DoRequest(ctx, "DELETE", "/library/custom-scripts/"+data.ID.ValueString(), nil, nil)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete custom script, got error: %s", err))
		return
	}
}

func (r *customScriptResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
