package provider

import (
	"context"
	"fmt"
	"strings"

	"github.com/MScottBlake/terraform-provider-iru/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ datasource.DataSource = &threatsV2DataSource{}

func NewThreatsV2DataSource() datasource.DataSource {
	return &threatsV2DataSource{}
}

type threatsV2DataSource struct {
	client *client.Client
}

type threatsV2DataSourceModel struct {
	ID              types.String   `tfsdk:"id"`
	Limit           types.Int64    `tfsdk:"limit"`
	Offset          types.Int64    `tfsdk:"offset"`
	Statuses        types.List     `tfsdk:"statuses"`
	Severities      types.List     `tfsdk:"severities"`
	ManagementState types.String   `tfsdk:"management_state"`
	Results         []threatV2Model `tfsdk:"results"`
}

type threatV2Model struct {
	ID                 types.String `tfsdk:"id"`
	ThreatName         types.String `tfsdk:"threat_name"`
	Classification     types.String `tfsdk:"classification"`
	Status             types.String `tfsdk:"status"`
	ManagementState    types.String `tfsdk:"management_state"`
	Severity           types.String `tfsdk:"severity"`
	Tags               types.List   `tfsdk:"tags"`
	DeviceName         types.String `tfsdk:"device_name"`
	DeviceID           types.String `tfsdk:"device_id"`
	DetectionDate      types.String `tfsdk:"detection_date"`
	FilePath           types.String `tfsdk:"file_path"`
	FileHash           types.String `tfsdk:"file_hash"`
	DeviceSerialNumber types.String `tfsdk:"device_serial_number"`
}

func (d *threatsV2DataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_threats_v2"
}

func (d *threatsV2DataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "List detected threats using the v2 API, which provides enhanced performance and additional fields like severity and tags.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"limit": schema.Int64Attribute{
				Optional:            true,
				MarkdownDescription: "Maximum number of results to return.",
			},
			"offset": schema.Int64Attribute{
				Optional:            true,
				MarkdownDescription: "Number of results to skip.",
			},
			"statuses": schema.ListAttribute{
				ElementType:         types.StringType,
				Optional:            true,
				MarkdownDescription: "Filter by threat status (e.g., `quarantined`, `released`).",
			},
			"severities": schema.ListAttribute{
				ElementType:         types.StringType,
				Optional:            true,
				MarkdownDescription: "Filter by severity (e.g., `critical`, `high`, `medium`, `low`, `informational`).",
			},
			"management_state": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Filter by management state (e.g., `managed`, `unmanaged`).",
			},
			"results": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id":                   schema.StringAttribute{Computed: true},
						"threat_name":          schema.StringAttribute{Computed: true},
						"classification":       schema.StringAttribute{Computed: true},
						"status":               schema.StringAttribute{Computed: true},
						"management_state":     schema.StringAttribute{Computed: true},
						"severity":             schema.StringAttribute{Computed: true},
						"tags":                 schema.ListAttribute{ElementType: types.StringType, Computed: true},
						"device_name":          schema.StringAttribute{Computed: true},
						"device_id":            schema.StringAttribute{Computed: true},
						"detection_date":       schema.StringAttribute{Computed: true},
						"file_path":            schema.StringAttribute{Computed: true},
						"file_hash":            schema.StringAttribute{Computed: true},
						"device_serial_number": schema.StringAttribute{Computed: true},
					},
				},
			},
		},
	}
}

func (d *threatsV2DataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	d.client = req.ProviderData.(*client.Client)
}

func (d *threatsV2DataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data threatsV2DataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	queryParams := []string{}
	if !data.Limit.IsNull() {
		queryParams = append(queryParams, fmt.Sprintf("limit=%d", data.Limit.ValueInt64()))
	}
	if !data.Offset.IsNull() {
		queryParams = append(queryParams, fmt.Sprintf("offset=%d", data.Offset.ValueInt64()))
	}
	if !data.Statuses.IsNull() {
		var statuses []string
		data.Statuses.ElementsAs(ctx, &statuses, false)
		queryParams = append(queryParams, fmt.Sprintf("statuses=%s", strings.Join(statuses, ",")))
	}
	if !data.Severities.IsNull() {
		var severities []string
		data.Severities.ElementsAs(ctx, &severities, false)
		queryParams = append(queryParams, fmt.Sprintf("severities=%s", strings.Join(severities, ",")))
	}
	if !data.ManagementState.IsNull() {
		queryParams = append(queryParams, fmt.Sprintf("management_state=%s", data.ManagementState.ValueString()))
	}

	path := "/api/v2/threat/threat-details"
	if len(queryParams) > 0 {
		path += "?" + strings.Join(queryParams, "&")
	}

	type threatResponse struct {
		Results []client.ThreatV2 `json:"results"`
	}
	var listResp threatResponse

	err := d.client.DoRequest(ctx, "GET", path, nil, &listResp)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read threats v2, got error: %s", err))
		return
	}

	for _, item := range listResp.Results {
		tagsList, _ := types.ListValueFrom(ctx, types.StringType, item.Tags)
		
		data.Results = append(data.Results, threatV2Model{
			ID:                 types.StringValue(item.ID),
			ThreatName:         types.StringValue(item.ThreatName),
			Classification:     types.StringValue(item.Classification),
			Status:             types.StringValue(item.Status),
			ManagementState:    types.StringValue(item.ManagementState),
			Severity:           types.StringValue(item.Severity),
			Tags:               tagsList,
			DeviceName:         types.StringValue(item.DeviceName),
			DeviceID:           types.StringValue(item.DeviceID),
			DetectionDate:      types.StringValue(item.DetectionDate),
			FilePath:           types.StringValue(item.FilePath),
			FileHash:           types.StringValue(item.FileHash),
			DeviceSerialNumber: types.StringValue(item.DeviceSerialNumber),
		})
	}

	data.ID = types.StringValue("threats_v2")

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
