package provider

import (
	"context"
	"fmt"

	"github.com/MScottBlake/terraform-provider-iru/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/action/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ action.Action = &deviceShutdownAction{}

func NewDeviceShutdownAction() action.Action {
	return &deviceShutdownAction{}
}

type deviceShutdownAction struct {
	client *client.Client
}

type deviceShutdownActionModel struct {
	DeviceID types.String `tfsdk:"device_id"`
}

func (a *deviceShutdownAction) Metadata(ctx context.Context, req action.MetadataRequest, resp *action.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_device_action_shutdown"
}

func (a *deviceShutdownAction) Schema(ctx context.Context, req action.SchemaRequest, resp *action.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Shuts down a specific device. This is an imperative action.",
		Attributes: map[string]schema.Attribute{
			"device_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The unique identifier for the Device.",
			},
		},
	}
}

func (a *deviceShutdownAction) Configure(ctx context.Context, req action.ConfigureRequest, resp *action.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	a.client = req.ProviderData.(*client.Client)
}

func (a *deviceShutdownAction) Invoke(ctx context.Context, req action.InvokeRequest, resp *action.InvokeResponse) {
	var data deviceShutdownActionModel
	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	deviceID := data.DeviceID.ValueString()
	err := a.client.DoRequest(ctx, "POST", fmt.Sprintf("/api/v1/devices/%s/action/shutdown", deviceID), nil, nil)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to invoke shutdown, got error: %s", err))
		return
	}
}
