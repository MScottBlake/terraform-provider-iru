package provider

import (
	"context"
	"os"

	"github.com/MScottBlake/terraform-provider-iru/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/list"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure IruProvider satisfies various provider interfaces.
var _ provider.Provider = &IruProvider{}
var _ provider.ProviderWithActions = &IruProvider{}
var _ provider.ProviderWithFunctions = &IruProvider{}
var _ provider.ProviderWithEphemeralResources = &IruProvider{}
var _ provider.ProviderWithListResources = &IruProvider{}

// IruProvider defines the provider implementation.
type IruProvider struct {
	// version is set to the provider version on creation, and optionally used in HTTP headers.
	version string
}

// IruProviderModel describes the provider data model.
type IruProviderModel struct {
	APIURL   types.String `tfsdk:"api_url"`
	APIToken types.String `tfsdk:"api_token"`
}

func (p *IruProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "iru"
	resp.Version = p.version
}

func (p *IruProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"api_url": schema.StringAttribute{
				MarkdownDescription: "The API URL for Iru.",
				Optional:            true,
			},
			"api_token": schema.StringAttribute{
				MarkdownDescription: "The API Token for authentication.",
				Optional:            true,
				Sensitive:           true,
			},
		},
	}
}

func (p *IruProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data IruProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	apiURL := os.Getenv("IRU_API_URL")
	apiToken := os.Getenv("IRU_API_TOKEN")

	if !data.APIURL.IsNull() {
		apiURL = data.APIURL.ValueString()
	}

	if !data.APIToken.IsNull() {
		apiToken = data.APIToken.ValueString()
	}

	if apiURL == "" {
		resp.Diagnostics.AddError("Missing API URL", "The 'api_url' provider attribute or IRU_API_URL environment variable must be set.")
		return
	}

	if apiToken == "" {
		resp.Diagnostics.AddError("Missing API Token", "The 'api_token' provider attribute or IRU_API_TOKEN environment variable must be set.")
		return
	}

	c := client.NewClient(apiURL, apiToken)

	resp.DataSourceData = c
	resp.ResourceData = c
}

func (p *IruProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewBlueprintResource,
		NewADEIntegrationResource,
		NewDeviceResource,
		NewTagResource,
		NewCustomScriptResource,
		NewCustomProfileResource,
		NewCustomAppResource,
		NewInHouseAppResource,
	}
}

func (p *IruProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewDevicesDataSource,
		NewBlueprintsDataSource,
		NewTagsDataSource,
		NewCustomScriptsDataSource,
		NewCustomProfilesDataSource,
		NewUsersDataSource,
		NewPrismFileVaultDataSource,
		NewPrismAppFirewallDataSource,
		NewVulnerabilitiesDataSource,
		NewAuditEventsDataSource,
		NewLicensingDataSource,
		NewPrismAppsDataSource,
		NewPrismCertificatesDataSource,
		NewPrismLocalUsersDataSource,
		NewPrismSystemExtensionsDataSource,
		NewPrismActivationLockDataSource,
		NewPrismCellularDataSource,
		NewPrismDesktopScreensaverDataSource,
		NewPrismDeviceInformationDataSource,
		NewPrismGatekeeperXProtectDataSource,
		NewPrismInstalledProfilesDataSource,
		NewPrismKernelExtensionsDataSource,
		NewPrismLaunchAgentsDaemonsDataSource,
		NewPrismStartupSettingsDataSource,
		NewPrismTransparencyDatabaseDataSource,
		NewThreatsDataSource,
		NewBehavioralDetectionsDataSource,
		NewSelfServiceCategoriesDataSource,
	}
}

func (p *IruProvider) Actions(ctx context.Context) []func() action.Action {
	return []func() action.Action{
		NewDeviceRestartAction,
		NewDeviceShutdownAction,
		NewDeviceLockAction,
		NewDeviceEraseAction,
		NewDeviceBlankPushAction,
		NewDeviceRenameAction,
		NewDeviceEnableRemoteDesktopAction,
		NewDeviceForceCheckInAction,
		NewDeviceClearPasscodeAction,
		NewDeviceBypassActivationLockAction,
		NewDeviceUnlockAccountAction,
		NewDeviceReinstallAgentAction,
		NewDeviceDailyCheckinAction,
	}
}

func (p *IruProvider) Functions(ctx context.Context) []func() function.Function {
	return []func() function.Function{
		NewParseProfileFunction,
	}
}

func (p *IruProvider) EphemeralResources(ctx context.Context) []func() ephemeral.EphemeralResource {
	return []func() ephemeral.EphemeralResource{
		NewDeviceSecretsEphemeralResource,
		NewADEPublicKeyEphemeralResource,
		NewBlueprintOTAProfileEphemeralResource,
	}
}

func (p *IruProvider) ListResources(ctx context.Context) []func() list.ListResource {
	return []func() list.ListResource{
		NewDeviceListResource,
		NewTagListResource,
		NewBlueprintListResource,
		NewCustomAppListResource,
		NewCustomProfileListResource,
		NewCustomScriptListResource,
		NewInHouseAppListResource,
		NewADEIntegrationListResource,
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &IruProvider{
			version: version,
		}
	}
}
