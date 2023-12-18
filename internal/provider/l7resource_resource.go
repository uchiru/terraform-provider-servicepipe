package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/types"

	v1 "terraform-provider-servicepipe/internal/pkg/sdkv1"
	l7resource "terraform-provider-servicepipe/internal/pkg/sdkv1/l7resource"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource              = &l7resourceResource{}
	_ resource.ResourceWithConfigure = &l7resourceResource{}
)

// Newl7resourceResource is a helper function to simplify the provider implementation.
func NewL7resourceResource() resource.Resource {
	return &l7resourceResource{}
}

// l7resourceResource is the resource implementation.
type l7resourceResource struct {
	client *v1.Client
}

// l7resourceResourceModel maps the resource schema data.
type l7resourceResourceModel struct {
	L7ResourceID          types.Int64  `tfsdk:"l7_resource_id"`
	L7ResourceName        types.String `tfsdk:"l7_resource_name"`
	L7ResourceIsActive    types.Int64  `tfsdk:"l7_resource_is_active"`
	L7ProtectionDisable   types.Int64  `tfsdk:"l7_protection_disable"`
	UseCustomSsl          types.Int64  `tfsdk:"use_custom_ssl"`
	UseLetsencryptSsl     types.Int64  `tfsdk:"use_letsencrypt_ssl"`
	CustomSslKey          types.String `tfsdk:"custom_ssl_key"`
	CustomSslCrt          types.String `tfsdk:"custom_ssl_crt"`
	Forcessl              types.Int64  `tfsdk:"force_ssl"`
	ServiceHttp2          types.Int64  `tfsdk:"service_http2"`
	GeoipMode             types.Int64  `tfsdk:"geoip_mode"`
	GeoipList             types.String `tfsdk:"geoip_list"`
	GlobalWhitelistActive types.Int64  `tfsdk:"global_whitelist_active"`
	Http2https            types.Int64  `tfsdk:"http_2_https"`
	Https2http            types.Int64  `tfsdk:"https_2_http"`
	ProtectedIp           types.String `tfsdk:"protected_ip"`
	Wwwredir              types.Int64  `tfsdk:"www_redir"`
	Cdn                   types.Int64  `tfsdk:"cdn"`
	CdnHost               types.String `tfsdk:"cdn_host"`
	CdnProxyHost          types.String `tfsdk:"cdn_proxy_host"`

	// TODO надо убрать этот кусок и реализовать через ресура l7origin
	OriginData types.String `tfsdk:"origin_data"`

	LastUpdated types.String `tfsdk:"last_updated"`
}

// Metadata returns the resource type name.
func (r *l7resourceResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_l7resource"
}

// Schema defines the schema for the resource.
func (r *l7resourceResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"l7_resource_id": schema.Int64Attribute{
				Computed: true,
			},
			"l7_resource_name": schema.StringAttribute{
				Required: true,
			},
			"l7_resource_is_active": schema.Int64Attribute{
				Optional: true,
				Computed: true,
				Default:  int64default.StaticInt64(1),
			},
			"l7_protection_disable": schema.Int64Attribute{
				Optional: true,
				Computed: true,
				Default:  int64default.StaticInt64(0),
			},
			"use_custom_ssl": schema.Int64Attribute{
				Optional: true,
				Computed: true,
				Default:  int64default.StaticInt64(0),
			},
			"use_letsencrypt_ssl": schema.Int64Attribute{
				Optional: true,
				Computed: true,
				Default:  int64default.StaticInt64(0),
			},
			"custom_ssl_key": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString(""),
			},
			"custom_ssl_crt": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString(""),
			},
			"force_ssl": schema.Int64Attribute{
				Optional: true,
				Computed: true,
				Default:  int64default.StaticInt64(0),
			},
			"service_http2": schema.Int64Attribute{
				Optional: true,
				Computed: true,
				Default:  int64default.StaticInt64(0),
			},
			"geoip_mode": schema.Int64Attribute{
				Optional: true,
				Computed: true,
				Default:  int64default.StaticInt64(0),
			},
			"geoip_list": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString(""),
			},
			"global_whitelist_active": schema.Int64Attribute{
				Optional: true,
				Computed: true,
				Default:  int64default.StaticInt64(1),
			},
			"http_2_https": schema.Int64Attribute{
				Optional: true,
				Computed: true,
				Default:  int64default.StaticInt64(0),
			},
			"https_2_http": schema.Int64Attribute{
				Optional: true,
				Computed: true,
				Default:  int64default.StaticInt64(0),
			},
			"protected_ip": schema.StringAttribute{
				Computed: true,
			},
			"www_redir": schema.Int64Attribute{
				Optional: true,
				Computed: true,
				Default:  int64default.StaticInt64(0),
			},
			"cdn": schema.Int64Attribute{
				Optional: true,
				Computed: true,
				Default:  int64default.StaticInt64(0),
			},
			"cdn_host": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString(""),
			},
			"cdn_proxy_host": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString(""),
			},
			"origin_data": schema.StringAttribute{
				Required: true,
			},
			"last_updated": schema.StringAttribute{
				Computed: true,
			},
		},
	}
}

// Configure adds the provider configured client to the resource.
func (r *l7resourceResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*v1.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *hashicups.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

// Create creates the resource and sets the initial Terraform state.
func (r *l7resourceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan *l7resourceResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	createOpts := &l7resource.CreateOpts{
		L7ResourceName: plan.L7ResourceName.ValueString(),
		OriginData:     plan.OriginData.ValueString(),
	}

	orig := plan.OriginData
	response, _, err := l7resource.Create(ctx, r.client, createOpts)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating l7resource",
			"Could not create l7resource, unexpected error: "+err.Error(),
		)
		return
	}

	update := false

	if !plan.L7ResourceName.IsNull() || !plan.L7ResourceName.IsUnknown() {
		response.Data.Result.L7ResourceName = plan.L7ResourceName.ValueString()
		update = true
	}

	if !plan.L7ResourceIsActive.IsNull() || !plan.L7ResourceIsActive.IsUnknown() {
		response.Data.Result.L7ResourceIsActive = int(plan.L7ResourceIsActive.ValueInt64())
		update = true
	}

	if !plan.L7ProtectionDisable.IsNull() || !plan.L7ProtectionDisable.IsUnknown() {
		response.Data.Result.L7ProtectionDisable = int(plan.L7ProtectionDisable.ValueInt64())
		update = true
	}

	if !plan.UseCustomSsl.IsNull() || !plan.UseCustomSsl.IsUnknown() {
		response.Data.Result.UseCustomSsl = int(plan.UseCustomSsl.ValueInt64())
		update = true
	}

	if !plan.UseLetsencryptSsl.IsNull() || !plan.UseLetsencryptSsl.IsUnknown() {
		response.Data.Result.UseLetsencryptSsl = int(plan.UseLetsencryptSsl.ValueInt64())
		update = true
	}

	if !plan.CustomSslKey.IsNull() || !plan.CustomSslKey.IsUnknown() {
		response.Data.Result.CustomSslKey = plan.CustomSslKey.ValueString()
		update = true
	}

	if !plan.CustomSslCrt.IsNull() || !plan.CustomSslCrt.IsUnknown() {
		response.Data.Result.CustomSslCrt = plan.CustomSslCrt.ValueString()
		update = true
	}
	if !plan.Forcessl.IsNull() || !plan.Forcessl.IsUnknown() {
		response.Data.Result.Forcessl = int(plan.Forcessl.ValueInt64())
		update = true
	}

	if !plan.ServiceHttp2.IsNull() || !plan.ServiceHttp2.IsUnknown() {
		response.Data.Result.ServiceHttp2 = int(plan.ServiceHttp2.ValueInt64())
		update = true
	}

	if !plan.GeoipMode.IsNull() || !plan.GeoipMode.IsUnknown() {
		response.Data.Result.GeoipMode = int(plan.GeoipMode.ValueInt64())
		update = true
	}

	if !plan.GeoipList.IsNull() || !plan.GeoipList.IsUnknown() {
		response.Data.Result.GeoipList = plan.GeoipList.ValueString()
		update = true
	}

	if !plan.GlobalWhitelistActive.IsNull() || !plan.GlobalWhitelistActive.IsUnknown() {
		response.Data.Result.GlobalWhitelistActive = int(plan.GlobalWhitelistActive.ValueInt64())
		update = true
	}

	if !plan.Http2https.IsNull() || !plan.Http2https.IsUnknown() {
		response.Data.Result.Http2https = int(plan.Http2https.ValueInt64())
		update = true
	}

	if !plan.Https2http.IsNull() || !plan.Https2http.IsUnknown() {
		response.Data.Result.Https2http = int(plan.Https2http.ValueInt64())
		update = true
	}

	if !plan.ProtectedIp.IsNull() || !plan.ProtectedIp.IsUnknown() {
		response.Data.Result.ProtectedIp = plan.ProtectedIp.ValueString()
		update = true
	}

	if !plan.Wwwredir.IsNull() || !plan.Wwwredir.IsUnknown() {
		response.Data.Result.Wwwredir = int(plan.Wwwredir.ValueInt64())
		update = true
	}

	if !plan.Cdn.IsNull() || !plan.Cdn.IsUnknown() {
		response.Data.Result.Cdn = int(plan.Cdn.ValueInt64())
		update = true
	}

	if !plan.CdnHost.IsNull() || !plan.CdnHost.IsUnknown() {
		response.Data.Result.CdnHost = plan.CdnHost.ValueString()
		update = true
	}

	if !plan.CdnProxyHost.IsNull() || !plan.CdnProxyHost.IsUnknown() {
		response.Data.Result.CdnProxyHost = plan.CdnProxyHost.ValueString()
		update = true
	}

	if update {
		respUpd, _, err := l7resource.Update(ctx, r.client, &response.Data.Result)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error Updating Servicepipe l7 resource",
				"Could not update l7 resource, unexpected error: "+err.Error()+"ResID"+strconv.Itoa(int(plan.L7ResourceID.ValueInt64())),
			)
			return
		}
		// Convert from the API data model to the Terraform data model
		plan = l7ItemToResourceModel(respUpd.Data.Result)
	}

	plan.OriginData = orig
	plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *l7resourceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state *l7resourceResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed order value from HashiCups
	response, _, err := l7resource.GetByID(ctx, r.client, int(state.L7ResourceID.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading HashiCups l7 resource",
			"Could not read HashiCups l7 resource ID "+strconv.Itoa(int(state.L7ResourceID.ValueInt64()))+": "+err.Error(),
		)
		return
	}

	orig := state.OriginData
	state = l7ItemToResourceModel(response.Data.Result)
	state.OriginData = orig

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *l7resourceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Get current plan
	var plan, state *l7resourceResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Update existing order
	opts := ResourceModelTol7Item(plan)
	opts.L7ResourceID = state.L7ResourceID.ValueInt64()

	jsonOpts, err := json.Marshal(opts)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Marshal req Servicepipe l7 resource",
			"Could not Marshal l7 resource, unexpected error: "+err.Error(),
		)
		return
	}

	update := false

	if !plan.L7ResourceName.Equal(state.L7ResourceName) {
		opts.L7ResourceName = plan.L7ResourceName.ValueString()
		update = true
	}

	if !plan.L7ResourceIsActive.Equal(state.L7ResourceIsActive) {
		opts.L7ResourceIsActive = int(plan.L7ResourceIsActive.ValueInt64())
		update = true
	}

	if !plan.L7ProtectionDisable.Equal(state.L7ProtectionDisable) {
		opts.L7ProtectionDisable = int(plan.L7ProtectionDisable.ValueInt64())
		update = true
	}

	if !plan.UseCustomSsl.Equal(state.UseCustomSsl) {
		opts.UseCustomSsl = int(plan.UseCustomSsl.ValueInt64())
		update = true
	}

	if !plan.UseLetsencryptSsl.Equal(state.UseLetsencryptSsl) {
		opts.UseLetsencryptSsl = int(plan.UseLetsencryptSsl.ValueInt64())
		update = true
	}

	if !plan.CustomSslKey.Equal(state.CustomSslKey) {
		opts.CustomSslKey = plan.CustomSslKey.ValueString()
		update = true
	}

	if !plan.CustomSslCrt.Equal(state.CustomSslCrt) {
		opts.CustomSslCrt = plan.CustomSslCrt.ValueString()
		update = true
	}
	if !plan.Forcessl.Equal(state.Forcessl) {
		opts.Forcessl = int(plan.Forcessl.ValueInt64())
		update = true
	}

	if !plan.ServiceHttp2.Equal(state.ServiceHttp2) {
		opts.ServiceHttp2 = int(plan.ServiceHttp2.ValueInt64())
		update = true
	}

	if !plan.GeoipMode.Equal(state.GeoipMode) {
		opts.GeoipMode = int(plan.GeoipMode.ValueInt64())
		update = true
	}

	if !plan.GeoipList.Equal(state.GeoipList) {
		opts.GeoipList = plan.GeoipList.ValueString()
		update = true
	}

	if !plan.GlobalWhitelistActive.Equal(state.GlobalWhitelistActive) {
		opts.GlobalWhitelistActive = int(plan.GlobalWhitelistActive.ValueInt64())
		update = true
	}

	if !plan.Http2https.Equal(state.Http2https) {
		opts.Http2https = int(plan.Http2https.ValueInt64())
		update = true
	}

	if !plan.Https2http.Equal(state.Https2http) {
		opts.Https2http = int(plan.Https2http.ValueInt64())
		update = true
	}

	if !plan.ProtectedIp.Equal(state.ProtectedIp) {
		opts.ProtectedIp = plan.ProtectedIp.ValueString()
		update = true
	}

	if !plan.Wwwredir.Equal(state.Wwwredir) {
		opts.Wwwredir = int(plan.Wwwredir.ValueInt64())
		update = true
	}

	if !plan.Cdn.Equal(state.Cdn) {
		opts.Cdn = int(plan.Cdn.ValueInt64())
		update = true
	}

	if !plan.CdnHost.Equal(state.CdnHost) {
		opts.CdnHost = plan.CdnHost.ValueString()
		update = true
	}

	if !plan.CdnProxyHost.Equal(state.CdnProxyHost) {
		opts.CdnProxyHost = plan.CdnProxyHost.ValueString()
		update = true
	}

	if !update {
		return
	}

	_, _, err = l7resource.Update(ctx, r.client, opts)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Servicepipe l7 resource",
			"Could not update l7 resource, unexpected error: "+err.Error()+"ResID"+strconv.Itoa(int(plan.L7ResourceID.ValueInt64()))+string(jsonOpts),
		)
		return
	}

	// Get refreshed l7resource value from Servicepipe
	response, _, err := l7resource.GetByID(ctx, r.client, int(state.L7ResourceID.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Servicepipe l7 resource",
			"Could not read Servicepipe l7 resource ID "+strconv.Itoa(int(state.L7ResourceID.ValueInt64()))+": "+err.Error(),
		)
		return
	}

	orig := plan.OriginData
	plan = l7ItemToResourceModel(response.Data.Result)

	plan.OriginData = orig
	plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))

	// Set refreshed state
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *l7resourceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state l7resourceResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	deleteOriginOpts := &l7resource.DeleteOpts{
		L7ResourceID: int(state.L7ResourceID.ValueInt64()),
	}

	// Delete existing resource
	result, _, err := l7resource.Delete(ctx, r.client, deleteOriginOpts)
	if err != nil || result.Data.Result != "ok" {
		resp.Diagnostics.AddError(
			"Error Deleting l7resource",
			"Could not delete l7resource, unexpected error: "+err.Error(),
		)
		return
	}
}

func ResourceModelTol7Item(model *l7resourceResourceModel) *l7resource.Item {
	return &l7resource.Item{
		// Assuming PartnerClientAccountId and other fields not in model are set elsewhere
		L7ResourceID:          model.L7ResourceID.ValueInt64(),
		L7ResourceName:        model.L7ResourceName.ValueString(),
		L7ResourceIsActive:    int(model.L7ResourceIsActive.ValueInt64()),
		L7ProtectionDisable:   int(model.L7ProtectionDisable.ValueInt64()),
		UseCustomSsl:          int(model.UseCustomSsl.ValueInt64()),
		UseLetsencryptSsl:     int(model.UseLetsencryptSsl.ValueInt64()),
		CustomSslKey:          model.CustomSslKey.ValueString(),
		CustomSslCrt:          model.CustomSslCrt.ValueString(),
		Forcessl:              int(model.Forcessl.ValueInt64()),
		ServiceHttp2:          int(model.ServiceHttp2.ValueInt64()),
		GeoipMode:             int(model.GeoipMode.ValueInt64()),
		GeoipList:             model.GeoipList.ValueString(),
		GlobalWhitelistActive: int(model.GlobalWhitelistActive.ValueInt64()),
		Http2https:            int(model.Http2https.ValueInt64()),
		Https2http:            int(model.Https2http.ValueInt64()),
		ProtectedIp:           model.ProtectedIp.ValueString(),
		Wwwredir:              int(model.Wwwredir.ValueInt64()),
		Cdn:                   int(model.Cdn.ValueInt64()),
		CdnHost:               model.CdnHost.ValueString(),
		CdnProxyHost:          model.CdnProxyHost.ValueString(),
		// OriginData:            model.OriginData.ValueString(),
	}
}

func l7ItemToResourceModel(item l7resource.Item) *l7resourceResourceModel {
	return &l7resourceResourceModel{
		L7ResourceID:          types.Int64Value(int64(item.L7ResourceID)),
		L7ResourceName:        types.StringValue(item.L7ResourceName),
		L7ResourceIsActive:    types.Int64Value(int64(item.L7ResourceIsActive)),
		L7ProtectionDisable:   types.Int64Value(int64(item.L7ProtectionDisable)),
		UseCustomSsl:          types.Int64Value(int64(item.UseCustomSsl)),
		UseLetsencryptSsl:     types.Int64Value(int64(item.UseLetsencryptSsl)),
		CustomSslKey:          types.StringValue(item.CustomSslKey),
		CustomSslCrt:          types.StringValue(item.CustomSslCrt),
		Forcessl:              types.Int64Value(int64(item.Forcessl)),
		ServiceHttp2:          types.Int64Value(int64(item.ServiceHttp2)),
		GeoipMode:             types.Int64Value(int64(item.GeoipMode)),
		GeoipList:             types.StringValue(item.GeoipList),
		GlobalWhitelistActive: types.Int64Value(int64(item.GlobalWhitelistActive)),
		Http2https:            types.Int64Value(int64(item.Http2https)),
		Https2http:            types.Int64Value(int64(item.Https2http)),
		ProtectedIp:           types.StringValue(item.ProtectedIp),
		Wwwredir:              types.Int64Value(int64(item.Wwwredir)),
		Cdn:                   types.Int64Value(int64(item.Cdn)),
		CdnHost:               types.StringValue(item.CdnHost),
		CdnProxyHost:          types.StringValue(item.CdnProxyHost),
		OriginData:            types.StringValue(item.OriginData),
	}
}
