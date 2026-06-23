package provider

import (
	"context"
	"fmt"

	"github.com/MScottBlake/terraform-provider-iru/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/action/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ action.Action = &deviceResetWorkProfilePasscodeAction{}

func NewDeviceResetWorkProfilePasscodeAction() action.Action {
	return &deviceResetWorkProfilePasscodeAction{}
}

type deviceResetWorkProfilePasscodeAction struct {
	client *client.Client
}

type deviceResetWorkProfilePasscodeActionModel struct {
	DeviceID           types.String `tfsdk:"device_id"`
	NewPassword        types.String `tfsdk:"new_password"`
	ResetPasswordFlags types.List   `tfsdk:"reset_password_flags"`
}

func (a *deviceResetWorkProfilePasscodeAction) Metadata(ctx context.Context, req action.MetadataRequest, resp *action.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_device_action_reset_work_profile_passcode"
}

func (a *deviceResetWorkProfilePasscodeAction) Schema(ctx context.Context, req action.SchemaRequest, resp *action.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Resets the work profile passcode on an Android device. If `new_password` is omitted, the existing passcode will be removed (if allowed by policy).",
		Attributes: map[string]schema.Attribute{
			"device_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The unique identifier for the Device.",
			},
			"new_password": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The new work profile passcode. Must be at least 4 characters (6 for Android 14+).",
			},
			"reset_password_flags": schema.ListAttribute{
				ElementType:         types.StringType,
				Optional:            true,
				MarkdownDescription: "Flags for the reset command, e.g., `LOCK_NOW`.",
			},
		},
	}
}

func (a *deviceResetWorkProfilePasscodeAction) Configure(ctx context.Context, req action.ConfigureRequest, resp *action.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	a.client = req.ProviderData.(*client.Client)
}

func (a *deviceResetWorkProfilePasscodeAction) Invoke(ctx context.Context, req action.InvokeRequest, resp *action.InvokeResponse) {
	var data deviceResetWorkProfilePasscodeActionModel
	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	deviceID := data.DeviceID.ValueString()
	payload := map[string]interface{}{}
	
	if !data.NewPassword.IsNull() {
		payload["newPassword"] = data.NewPassword.ValueString()
	}
	
	if !data.ResetPasswordFlags.IsNull() {
		var flags []string
		diags := data.ResetPasswordFlags.ElementsAs(ctx, &flags, false)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		payload["resetPasswordFlags"] = flags
	} else {
		payload["resetPasswordFlags"] = []string{}
	}

	err := a.client.DoRequest(ctx, "POST", fmt.Sprintf("/api/v1/devices/%s/action/resetworkprofilepasscode", deviceID), payload, nil)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to reset work profile passcode, got error: %s", err))
		return
	}
}
