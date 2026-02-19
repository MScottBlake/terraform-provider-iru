package provider

import (
	"context"
	"fmt"

	"github.com/MScottBlake/terraform-provider-iru/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/action/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ action.Action = &deviceDisableLostModeAction{}

func NewDeviceDisableLostModeAction() action.Action {
	return &deviceDisableLostModeAction{}
}

type deviceDisableLostModeAction struct {
	client *client.Client
}

type deviceDisableLostModeActionModel struct {
	DeviceID types.String `tfsdk:"device_id"`
}

func (a *deviceDisableLostModeAction) Metadata(ctx context.Context, req action.MetadataRequest, resp *action.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_device_action_disable_lost_mode"
}

func (a *deviceDisableLostModeAction) Schema(ctx context.Context, req action.SchemaRequest, resp *action.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Disables Lost Mode on a specific device.",
		Attributes: map[string]schema.Attribute{
			"device_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The unique identifier for the Device.",
			},
		},
	}
}

func (a *deviceDisableLostModeAction) Configure(ctx context.Context, req action.ConfigureRequest, resp *action.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	a.client = req.ProviderData.(*client.Client)
}

func (a *deviceDisableLostModeAction) Invoke(ctx context.Context, req action.InvokeRequest, resp *action.InvokeResponse) {
	var data deviceDisableLostModeActionModel
	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	deviceID := data.DeviceID.ValueString()
	err := a.client.DoRequest(ctx, "POST", fmt.Sprintf("/api/v1/devices/%s/action/disablelostmode", deviceID), nil, nil)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to disable lost mode, got error: %s", err))
		return
	}
}
