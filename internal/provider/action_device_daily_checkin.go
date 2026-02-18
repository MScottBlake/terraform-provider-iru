package provider

import (
	"context"
	"fmt"

	"github.com/MScottBlake/terraform-provider-iru/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/action/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ action.Action = &deviceDailyCheckinAction{}

func NewDeviceDailyCheckinAction() action.Action {
	return &deviceDailyCheckinAction{}
}

type deviceDailyCheckinAction struct {
	client *client.Client
}

type deviceDailyCheckinActionModel struct {
	DeviceID types.String `tfsdk:"device_id"`
}

func (a *deviceDailyCheckinAction) Metadata(ctx context.Context, req action.MetadataRequest, resp *action.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_device_action_daily_checkin"
}

func (a *deviceDailyCheckinAction) Schema(ctx context.Context, req action.SchemaRequest, resp *action.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Initiates a daily check-in for a device.",
		Attributes: map[string]schema.Attribute{
			"device_id": schema.StringAttribute{
				Required: true,
			},
		},
	}
}

func (a *deviceDailyCheckinAction) Configure(ctx context.Context, req action.ConfigureRequest, resp *action.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	a.client = req.ProviderData.(*client.Client)
}

func (a *deviceDailyCheckinAction) Invoke(ctx context.Context, req action.InvokeRequest, resp *action.InvokeResponse) {
	var data deviceDailyCheckinActionModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := a.client.DoRequest(ctx, "POST", fmt.Sprintf("/devices/%s/action/dailycheckin", data.DeviceID.ValueString()), nil, nil)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to invoke daily checkin, got error: %s", err))
		return
	}
}
