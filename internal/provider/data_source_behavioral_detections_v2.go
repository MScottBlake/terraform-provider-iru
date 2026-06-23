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

var _ datasource.DataSource = &behavioralDetectionsV2DataSource{}

func NewBehavioralDetectionsV2DataSource() datasource.DataSource {
	return &behavioralDetectionsV2DataSource{}
}

type behavioralDetectionsV2DataSource struct {
	client *client.Client
}

type behavioralDetectionsV2DataSourceModel struct {
	ID              types.String               `tfsdk:"id"`
	Limit           types.Int64                `tfsdk:"limit"`
	Offset          types.Int64                `tfsdk:"offset"`
	Statuses        types.List                 `tfsdk:"statuses"`
	Severities      types.List                 `tfsdk:"severities"`
	ManagementState types.String               `tfsdk:"management_state"`
	DeviceID        types.String               `tfsdk:"device_id"`
	Results         []behavioralDetectionV2Model `tfsdk:"results"`
}

type behavioralDetectionV2Model struct {
	ID              types.String               `tfsdk:"id"`
	ThreatID        types.String               `tfsdk:"threat_id"`
	Description     types.String               `tfsdk:"description"`
	Classification  types.String               `tfsdk:"classification"`
	Status          types.String               `tfsdk:"status"`
	ManagementState types.String               `tfsdk:"management_state"`
	Severity        types.String               `tfsdk:"severity"`
	Tags            types.List                 `tfsdk:"tags"`
	DetectionDate   types.String               `tfsdk:"detection_date"`
	DeviceInfo      behavioralDetectionV2DeviceModel `tfsdk:"device_info"`
}

type behavioralDetectionV2DeviceModel struct {
	ID           types.String `tfsdk:"id"`
	Name         types.String `tfsdk:"name"`
	SerialNumber types.String `tfsdk:"serial_number"`
}

func (d *behavioralDetectionsV2DataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_behavioral_detections_v2"
}

func (d *behavioralDetectionsV2DataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "List behavioral detection events using the v2 API.",
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
				MarkdownDescription: "Filter by threat status (e.g., `informational`, `detected`).",
			},
			"severities": schema.ListAttribute{
				ElementType:         types.StringType,
				Optional:            true,
				MarkdownDescription: "Filter by severity.",
			},
			"management_state": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Filter by management state.",
			},
			"device_id": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Filter by device UUID.",
			},
			"results": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id":               schema.StringAttribute{Computed: true},
						"threat_id":        schema.StringAttribute{Computed: true},
						"description":      schema.StringAttribute{Computed: true},
						"classification":   schema.StringAttribute{Computed: true},
						"status":           schema.StringAttribute{Computed: true},
						"management_state": schema.StringAttribute{Computed: true},
						"severity":         schema.StringAttribute{Computed: true},
						"tags":             schema.ListAttribute{ElementType: types.StringType, Computed: true},
						"detection_date":   schema.StringAttribute{Computed: true},
						"device_info": schema.SingleNestedAttribute{
							Computed: true,
							Attributes: map[string]schema.Attribute{
								"id":            schema.StringAttribute{Computed: true},
								"name":          schema.StringAttribute{Computed: true},
								"serial_number": schema.StringAttribute{Computed: true},
							},
						},
					},
				},
			},
		},
	}
}

func (d *behavioralDetectionsV2DataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	d.client = req.ProviderData.(*client.Client)
}

func (d *behavioralDetectionsV2DataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data behavioralDetectionsV2DataSourceModel
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
	if !data.DeviceID.IsNull() {
		queryParams = append(queryParams, fmt.Sprintf("device_id=%s", data.DeviceID.ValueString()))
	}

	path := "/api/v2/threat/behavioral-detections/events"
	if len(queryParams) > 0 {
		path += "?" + strings.Join(queryParams, "&")
	}

	type behavioralResponse struct {
		Results []client.BehavioralDetectionV2 `json:"results"`
	}
	var listResp behavioralResponse

	err := d.client.DoRequest(ctx, "GET", path, nil, &listResp)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read behavioral detections v2, got error: %s", err))
		return
	}

	for _, item := range listResp.Results {
		tagsList, _ := types.ListValueFrom(ctx, types.StringType, item.Tags)
		
		data.Results = append(data.Results, behavioralDetectionV2Model{
			ID:              types.StringValue(item.ID),
			ThreatID:        types.StringValue(item.ThreatID),
			Description:     types.StringValue(item.Description),
			Classification:  types.StringValue(item.Classification),
			Status:          types.StringValue(item.Status),
			ManagementState: types.StringValue(item.ManagementState),
			Severity:        types.StringValue(item.Severity),
			Tags:            tagsList,
			DetectionDate:   types.StringValue(item.DetectionDate),
			DeviceInfo: behavioralDetectionV2DeviceModel{
				ID:           types.StringValue(item.DeviceInfo.ID),
				Name:         types.StringValue(item.DeviceInfo.Name),
				SerialNumber: types.StringValue(item.DeviceInfo.SerialNumber),
			},
		})
	}

	data.ID = types.StringValue("behavioral_detections_v2")

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
