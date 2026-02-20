package provider

import (
	"context"
	"fmt"

	"github.com/MScottBlake/terraform-provider-iru/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/action/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ action.Action = &deviceReinstallAgentAction{}

func NewDeviceReinstallAgentAction() action.Action {
	return &deviceReinstallAgentAction{}
}

type deviceReinstallAgentAction struct {
	client *client.Client
}

type deviceReinstallAgentActionModel struct {
	DeviceID types.String `tfsdk:"device_id"`
}

func (a *deviceReinstallAgentAction) Metadata(ctx context.Context, req action.MetadataRequest, resp *action.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_device_action_reinstall_agent"
}

func (a *deviceReinstallAgentAction) Schema(ctx context.Context, req action.SchemaRequest, resp *action.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Reinstalls the Iru Agent on macOS devices.",
		Attributes: map[string]schema.Attribute{
			"device_id": schema.StringAttribute{
				Required: true,
			},
		},
	}
}

func (a *deviceReinstallAgentAction) Configure(ctx context.Context, req action.ConfigureRequest, resp *action.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	a.client = req.ProviderData.(*client.Client)
}

func (a *deviceReinstallAgentAction) Invoke(ctx context.Context, req action.InvokeRequest, resp *action.InvokeResponse) {
	var data deviceReinstallAgentActionModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := a.client.DoRequest(ctx, "POST", fmt.Sprintf("/api/v1/devices/%s/action/reinstallagent", data.DeviceID.ValueString()), nil, nil)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to invoke reinstall agent, got error: %s", err))
		return
	}
}
