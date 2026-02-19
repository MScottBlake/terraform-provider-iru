package provider

import (
	"context"
	"fmt"

	"github.com/MScottBlake/terraform-provider-iru/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/action/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ action.Action = &deviceRenameAction{}

func NewDeviceRenameAction() action.Action {
	return &deviceRenameAction{}
}

type deviceRenameAction struct {
	client *client.Client
}

type deviceRenameActionModel struct {
	DeviceID types.String `tfsdk:"device_id"`
	NewName  types.String `tfsdk:"new_name"`
}

func (a *deviceRenameAction) Metadata(ctx context.Context, req action.MetadataRequest, resp *action.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_device_action_rename"
}

func (a *deviceRenameAction) Schema(ctx context.Context, req action.SchemaRequest, resp *action.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Renames a specific device. This is an imperative action.",
		Attributes: map[string]schema.Attribute{
			"device_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The unique identifier for the Device.",
			},
			"new_name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The new name for the device.",
			},
		},
	}
}

func (a *deviceRenameAction) Configure(ctx context.Context, req action.ConfigureRequest, resp *action.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	a.client = req.ProviderData.(*client.Client)
}

func (a *deviceRenameAction) Invoke(ctx context.Context, req action.InvokeRequest, resp *action.InvokeResponse) {
	var data deviceRenameActionModel
	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	deviceID := data.DeviceID.ValueString()
	payload := map[string]string{
		"device_name": data.NewName.ValueString(),
	}
	
	err := a.client.DoRequest(ctx, "POST", fmt.Sprintf("/api/v1/devices/%s/action/rename", deviceID), payload, nil)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to invoke rename, got error: %s", err))
		return
	}
}
