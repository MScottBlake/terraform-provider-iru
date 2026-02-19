package provider

import (
	"context"
	"fmt"

	"github.com/MScottBlake/terraform-provider-iru/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/action/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ action.Action = &deviceUpdateInventoryAction{}

func NewDeviceUpdateInventoryAction() action.Action {
	return &deviceUpdateInventoryAction{}
}

type deviceUpdateInventoryAction struct {
	client *client.Client
}

type deviceUpdateInventoryActionModel struct {
	DeviceID types.String `tfsdk:"device_id"`
}

func (a *deviceUpdateInventoryAction) Metadata(ctx context.Context, req action.MetadataRequest, resp *action.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_device_action_update_inventory"
}

func (a *deviceUpdateInventoryAction) Schema(ctx context.Context, req action.SchemaRequest, resp *action.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Updates inventory for a specific device.",
		Attributes: map[string]schema.Attribute{
			"device_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The unique identifier for the Device.",
			},
		},
	}
}

func (a *deviceUpdateInventoryAction) Configure(ctx context.Context, req action.ConfigureRequest, resp *action.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	a.client = req.ProviderData.(*client.Client)
}

func (a *deviceUpdateInventoryAction) Invoke(ctx context.Context, req action.InvokeRequest, resp *action.InvokeResponse) {
	var data deviceUpdateInventoryActionModel
	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	deviceID := data.DeviceID.ValueString()
	err := a.client.DoRequest(ctx, "POST", fmt.Sprintf("/api/v1/devices/%s/action/updateinventory", deviceID), nil, nil)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update inventory, got error: %s", err))
		return
	}
}
