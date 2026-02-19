package provider

import (
	"context"
	"fmt"

	"github.com/MScottBlake/terraform-provider-iru/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/action/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ action.Action = &deviceForceCheckInAction{}

func NewDeviceForceCheckInAction() action.Action {
	return &deviceForceCheckInAction{}
}

type deviceForceCheckInAction struct {
	client *client.Client
}

type deviceForceCheckInActionModel struct {
	DeviceID types.String `tfsdk:"device_id"`
}

func (a *deviceForceCheckInAction) Metadata(ctx context.Context, req action.MetadataRequest, resp *action.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_device_action_force_check_in"
}

func (a *deviceForceCheckInAction) Schema(ctx context.Context, req action.SchemaRequest, resp *action.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Forces a check-in for a specific device. This is an imperative action.",
		Attributes: map[string]schema.Attribute{
			"device_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The unique identifier for the Device.",
			},
		},
	}
}

func (a *deviceForceCheckInAction) Configure(ctx context.Context, req action.ConfigureRequest, resp *action.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	a.client = req.ProviderData.(*client.Client)
}

func (a *deviceForceCheckInAction) Invoke(ctx context.Context, req action.InvokeRequest, resp *action.InvokeResponse) {
	var data deviceForceCheckInActionModel
	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	deviceID := data.DeviceID.ValueString()
	err := a.client.DoRequest(ctx, "POST", fmt.Sprintf("/api/v1/devices/%s/action/force-check-in", deviceID), nil, nil)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to invoke force check-in, got error: %s", err))
		return
	}
}
