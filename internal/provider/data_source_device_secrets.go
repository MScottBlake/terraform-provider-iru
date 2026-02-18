package provider

import (
	"context"
	"fmt"

	"github.com/MScottBlake/terraform-provider-iru/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ datasource.DataSource = &deviceSecretsDataSource{}

func NewDeviceSecretsDataSource() datasource.DataSource {
	return &deviceSecretsDataSource{}
}

type deviceSecretsDataSource struct {
	client *client.Client
}

type deviceSecretsDataSourceModel struct {
	DeviceID               types.String `tfsdk:"device_id"`
	UserBasedALBC          types.String `tfsdk:"user_based_albc"`
	DeviceBasedALBC        types.String `tfsdk:"device_based_albc"`
	FileVaultRecoveryKey   types.String `tfsdk:"filevault_recovery_key"`
	UnlockPin              types.String `tfsdk:"unlock_pin"`
	RecoveryLockPassword   types.String `tfsdk:"recovery_lock_password"`
}

func (d *deviceSecretsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_device_secrets"
}

func (d *deviceSecretsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Fetch secrets for a specific device. Warning: These values are sensitive and will be stored in the Terraform state.",
		Attributes: map[string]schema.Attribute{
			"device_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The unique identifier for the Device.",
			},
			"user_based_albc": schema.StringAttribute{
				Computed:            true,
				Sensitive:           true,
				MarkdownDescription: "User-based Activation Lock Bypass Code.",
			},
			"device_based_albc": schema.StringAttribute{
				Computed:            true,
				Sensitive:           true,
				MarkdownDescription: "Device-based Activation Lock Bypass Code.",
			},
			"filevault_recovery_key": schema.StringAttribute{
				Computed:            true,
				Sensitive:           true,
				MarkdownDescription: "FileVault Personal Recovery Key.",
			},
			"unlock_pin": schema.StringAttribute{
				Computed:            true,
				Sensitive:           true,
				MarkdownDescription: "Device Unlock PIN.",
			},
			"recovery_lock_password": schema.StringAttribute{
				Computed:            true,
				Sensitive:           true,
				MarkdownDescription: "Recovery Lock Password.",
			},
		},
	}
}

func (d *deviceSecretsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	d.client = req.ProviderData.(*client.Client)
}

func (d *deviceSecretsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data deviceSecretsDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	deviceID := data.DeviceID.ValueString()

	// ALBC
	var albc client.DeviceSecretsALBC
	err := d.client.DoRequest(ctx, "GET", fmt.Sprintf("/devices/%s/secrets/bypasscode", deviceID), nil, &albc)
	if err == nil {
		data.UserBasedALBC = types.StringValue(albc.UserBasedALBC)
		data.DeviceBasedALBC = types.StringValue(albc.DeviceBasedALBC)
	}

	// FileVault
	var fv client.DeviceSecretsFileVault
	err = d.client.DoRequest(ctx, "GET", fmt.Sprintf("/devices/%s/secrets/filevaultkey", deviceID), nil, &fv)
	if err == nil {
		data.FileVaultRecoveryKey = types.StringValue(fv.Key)
	}

	// Unlock Pin
	var pin client.DeviceSecretsUnlockPin
	err = d.client.DoRequest(ctx, "GET", fmt.Sprintf("/devices/%s/secrets/unlockpin", deviceID), nil, &pin)
	if err == nil {
		data.UnlockPin = types.StringValue(pin.Pin)
	}

	// Recovery Lock
	var rl client.DeviceSecretsRecoveryLock
	err = d.client.DoRequest(ctx, "GET", fmt.Sprintf("/devices/%s/secrets/recoverypassword", deviceID), nil, &rl)
	if err == nil {
		data.RecoveryLockPassword = types.StringValue(rl.RecoveryPassword)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
