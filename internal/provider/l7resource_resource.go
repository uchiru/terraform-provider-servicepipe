package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/types"

	v1 "terraform-provider-servicepipe/internal/pkg/sdkv1"
	l7origin "terraform-provider-servicepipe/internal/pkg/sdkv1/l7origin"
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
	ServiceHTTP2          types.Int64  `tfsdk:"service_http2"`
	GeoipMode             types.Int64  `tfsdk:"geoip_mode"`
	GeoipList             types.String `tfsdk:"geoip_list"`
	GlobalWhitelistActive types.Int64  `tfsdk:"global_whitelist_active"`
	HTTP2https            types.Int64  `tfsdk:"http_2_https"`
	HTTPS2http            types.Int64  `tfsdk:"https_2_http"`
	ProtectedIp           types.String `tfsdk:"protected_ip"`
	Wwwredir              types.Int64  `tfsdk:"www_redir"`
	Cdn                   types.Int64  `tfsdk:"cdn"`
	CdnHost               types.String `tfsdk:"cdn_host"`
	CdnProxyHost          types.String `tfsdk:"cdn_proxy_host"`

	Origins []*l7originResourceModel `tfsdk:"origins"`

	LastUpdated types.String `tfsdk:"last_updated"`
}

type l7originResourceModel struct {
	L7ResourceID types.Int64  `tfsdk:"l7_resource_id"`
	ID           types.Int64  `tfsdk:"id"`
	Weight       types.Int64  `tfsdk:"weight"`
	Mode         types.String `tfsdk:"mode"`
	IP           types.String `tfsdk:"ip"`
	CreatedAt    types.Int64  `tfsdk:"created_at"`
	ModifiedAt   types.Int64  `tfsdk:"modified_at"`
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
			"last_updated": schema.StringAttribute{
				Computed: true,
			},
			"origins": schema.ListNestedAttribute{
				Required: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"l7_resource_id": schema.Int64Attribute{
							Computed: true,
						},
						"id": schema.Int64Attribute{
							Computed: true,
						},
						"weight": schema.Int64Attribute{
							Optional: true,
							Computed: true,
							Default:  int64default.StaticInt64(50),
						},
						"mode": schema.StringAttribute{
							Optional: true,
							Computed: true,
							Default:  stringdefault.StaticString(""),
						},
						"ip": schema.StringAttribute{
							Optional: true,
							Computed: true,
							Default:  stringdefault.StaticString(""),
						},
						"created_at": schema.Int64Attribute{
							Computed: true,
						},
						"modified_at": schema.Int64Attribute{
							Computed: true,
						},
					},
				},
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

	planOrigins := plan.Origins

	// Generate API request body from plan
	createOpts := &l7resource.CreateOpts{
		L7ResourceName: plan.L7ResourceName.ValueString(),
		OriginData:     plan.Origins[0].IP.ValueString(),
	}

	response, _, err := l7resource.Create(ctx, r.client, createOpts)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating l7resource",
			"Could not create l7resource, unexpected error: "+err.Error(),
		)
		return
	}

	updateOpts, update := CheckingL7resourcePlanAttrIsNull(*plan, &response.Data.Result)

	if update {
		respUpd, _, err := l7resource.Update(ctx, r.client, updateOpts)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error Updating Servicepipe l7 resource",
				"Could not update l7 resource, unexpected error: "+err.Error()+"ResID"+strconv.Itoa(int(plan.L7ResourceID.ValueInt64())),
			)
			return
		}

		results := hackSPSSLState(plan, respUpd)

		// Convert from the API data model to the Terraform data model
		plan = flatternL7ResourceModel(results.Data.Result)
	}

	var origins []*l7originResourceModel
	for _, v := range planOrigins {
		origin := expandL7OriginModel(v)

		item, ok := CheckExistingOriginByIP(ctx, r.client, origin.IP, response.Data.Result.L7ResourceID)
		if !ok {
			// Generate API request body from plan
			createOriginOpts := &l7origin.CreateOpts{
				L7ResourceID: response.Data.Result.L7ResourceID,
				IP:           origin.IP,
				Weight:       origin.Weight,
				Mode:         origin.Mode,
			}

			result, _, err := l7origin.Create(ctx, r.client, createOriginOpts)
			if err != nil {
				msg := fmt.Sprintf("%+v", result)
				resp.Diagnostics.AddError(
					"Error creating l7origin: "+msg,
					"Could not create l7origin, unexpected error: "+err.Error(),
				)
				return
			}
			item = &result.Data.Result
		}

		originOpts := expandL7OriginModel(v)
		originOpts.L7ResourceID = response.Data.Result.L7ResourceID
		originOpts.ID = item.ID
		updateOrig := false

		if !v.IP.IsNull() || !v.IP.IsUnknown() {
			originOpts.IP = v.IP.ValueString()
			updateOrig = true
		}

		if !v.Mode.IsNull() || !v.Mode.IsUnknown() {
			originOpts.Mode = v.Mode.ValueString()
			updateOrig = true
		}

		if !v.Weight.IsNull() || !v.Weight.IsUnknown() {
			originOpts.Weight = v.Weight.ValueInt64()
			updateOrig = true
		}

		if updateOrig {
			respUpd, _, err := l7origin.Update(ctx, r.client, originOpts)
			if err != nil {
				msg := fmt.Sprintf("%+v _ %+v", respUpd, originOpts)
				resp.Diagnostics.AddError(
					"Error Updating Servicepipe l7 origin: "+msg,
					"Could not update l7 origin, unexpected error: "+err.Error(),
				)
				return
			}
		}

		originResponse, _, err := l7origin.GetByID(ctx, r.client, int(response.Data.Result.L7ResourceID), int(item.ID))
		if err != nil {
			resp.Diagnostics.AddError(
				"Error Reading Servicepipe l7 origin",
				"Could not read Servicepipe l7 origin ID "+strconv.Itoa(int(item.ID))+": "+err.Error(),
			)
			return
		}

		originResponse.Data.Result.L7ResourceID = response.Data.Result.L7ResourceID
		origins = append(origins, flatternL7OriginModel(&originResponse.Data.Result))
	}
	plan.Origins = origins
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
	resourceResponse, _, err := l7resource.GetByID(ctx, r.client, int(state.L7ResourceID.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading HashiCups l7 resource",
			"Could not read HashiCups l7 resource ID "+strconv.Itoa(int(state.L7ResourceID.ValueInt64()))+": "+err.Error(),
		)
		return
	}

	var origins []*l7originResourceModel
	for _, v := range state.Origins {
		originResponse, _, err := l7origin.GetByID(ctx, r.client, int(state.L7ResourceID.ValueInt64()), int(v.ID.ValueInt64()))
		if err != nil {
			resp.Diagnostics.AddError(
				"Error Reading Servicepipe l7 origin",
				"Could not read Servicepipe l7 origin ID "+strconv.Itoa(int(v.ID.ValueInt64()))+": "+err.Error(),
			)
			return
		}

		originResponse.Data.Result.L7ResourceID = state.L7ResourceID.ValueInt64()
		origins = append(origins, flatternL7OriginModel(&originResponse.Data.Result))
	}

	state = flatternL7ResourceModel(resourceResponse.Data.Result)
	state.Origins = origins

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
	item := expandL7ResourceModel(plan)
	item.L7ResourceID = state.L7ResourceID.ValueInt64()

	jsonOpts, err := json.Marshal(item)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Marshal req Servicepipe l7 resource",
			"Could not Marshal l7 resource, unexpected error: "+err.Error(),
		)
		return
	}

	opts, update := CheckPlanVsState(plan, state, item)

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

	results := hackSPSSLState(plan, response)

	planOrigins := plan.Origins
	plan = flatternL7ResourceModel(results.Data.Result)

	for _, v := range planOrigins {
		origin := expandL7OriginModel(v)
		_, ok := CheckExistingOriginByIP(ctx, r.client, origin.IP, state.L7ResourceID.ValueInt64())
		if !ok {
			// Generate API request body from plan
			createOriginOpts := &l7origin.CreateOpts{
				L7ResourceID: state.L7ResourceID.ValueInt64(),
				IP:           origin.IP,
				Weight:       origin.Weight,
				Mode:         origin.Mode,
			}

			result, _, err := l7origin.Create(ctx, r.client, createOriginOpts)
			if err != nil {
				msg := fmt.Sprintf("%+v", result)
				resp.Diagnostics.AddError(
					"Error creating l7origin: "+msg,
					"Could not create l7origin, unexpected error: "+err.Error(),
				)
				return
			}

			result.Data.Result.L7ResourceID = state.L7ResourceID.ValueInt64()
			state.Origins = append(state.Origins, flatternL7OriginModel(&result.Data.Result))
		}
	}

	if len(state.Origins) != len(planOrigins) {
		for i, v := range state.Origins {
			_, ok := CheckPlanVsStateOrigin(planOrigins, v.IP.ValueString())
			if !ok {
				deleteOriginOpts := &l7origin.DeleteOpts{
					ID:           v.ID.ValueInt64(),
					L7ResourceID: state.L7ResourceID.ValueInt64(),
				}

				// Delete existing resource
				result, _, err := l7origin.Delete(ctx, r.client, deleteOriginOpts)
				if err != nil || result.Data.Result != "ok" {
					resp.Diagnostics.AddError(
						"Error Deleting l7origin",
						"Could not delete l7origin, unexpected error: "+err.Error(),
					)
					return
				}

				state.Origins = removeL7originFromState(state.Origins, i)
			}
		}
	}

	for _, s := range state.Origins {
		originOpts := expandL7OriginModel(s)
		originOpts.L7ResourceID = state.L7ResourceID.ValueInt64()

		updateOrig := false
		for _, p := range planOrigins {
			if p.IP == s.IP {
				if !p.Mode.Equal(s.Mode) {
					originOpts.Mode = p.Mode.ValueString()
					updateOrig = true
				}

				if !p.Weight.Equal(s.Weight) {
					originOpts.Weight = p.Weight.ValueInt64()
					updateOrig = true
				}

				if !p.IP.Equal(s.IP) {
					originOpts.IP = p.IP.ValueString()
					updateOrig = true
				}
			}
		}

		if !updateOrig {
			continue
		}

		_, _, err = l7origin.Update(ctx, r.client, originOpts)
		if err != nil {
			msg := fmt.Sprintf("%+v", originOpts)
			resp.Diagnostics.AddError(
				"Error Updating Servicepipe l7 origin: "+msg,
				"Could not update l7 origin, unexpected error: "+err.Error(),
			)
			return
		}
	}

	var origins []*l7originResourceModel
	for _, v := range state.Origins {
		originResponse, _, err := l7origin.GetByID(ctx, r.client, int(state.L7ResourceID.ValueInt64()), int(v.ID.ValueInt64()))
		if err != nil {
			resp.Diagnostics.AddError(
				"Error Reading Servicepipe l7 origin",
				"Could not read Servicepipe l7 origin ID "+strconv.Itoa(int(v.ID.ValueInt64()))+": "+err.Error(),
			)
			return
		}

		originResponse.Data.Result.L7ResourceID = state.L7ResourceID.ValueInt64()
		origins = append(origins, flatternL7OriginModel(&originResponse.Data.Result))
	}

	plan.Origins = origins
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

	result, _, err := l7resource.Delete(ctx, r.client, deleteOriginOpts)
	if err != nil || result.Data.Result != "ok" {
		resp.Diagnostics.AddError(
			"Error Deleting l7resource",
			"Could not delete l7resource, unexpected error: "+err.Error(),
		)
		return
	}
}

func expandL7ResourceModel(model *l7resourceResourceModel) *l7resource.Item {
	return &l7resource.Item{
		L7ResourceID:          model.L7ResourceID.ValueInt64(),
		L7ResourceName:        model.L7ResourceName.ValueString(),
		L7ResourceIsActive:    int(model.L7ResourceIsActive.ValueInt64()),
		L7ProtectionDisable:   int(model.L7ProtectionDisable.ValueInt64()),
		UseCustomSsl:          int(model.UseCustomSsl.ValueInt64()),
		UseLetsencryptSsl:     int(model.UseLetsencryptSsl.ValueInt64()),
		CustomSslKey:          model.CustomSslKey.ValueString(),
		CustomSslCrt:          model.CustomSslCrt.ValueString(),
		Forcessl:              int(model.Forcessl.ValueInt64()),
		ServiceHTTP2:          int(model.ServiceHTTP2.ValueInt64()),
		GeoipMode:             int(model.GeoipMode.ValueInt64()),
		GeoipList:             model.GeoipList.ValueString(),
		GlobalWhitelistActive: int(model.GlobalWhitelistActive.ValueInt64()),
		HTTP2https:            int(model.HTTP2https.ValueInt64()),
		HTTPS2http:            int(model.HTTPS2http.ValueInt64()),
		ProtectedIp:           model.ProtectedIp.ValueString(),
		Wwwredir:              int(model.Wwwredir.ValueInt64()),
		Cdn:                   int(model.Cdn.ValueInt64()),
		CdnHost:               model.CdnHost.ValueString(),
		CdnProxyHost:          model.CdnProxyHost.ValueString(),
	}
}

func flatternL7ResourceModel(item l7resource.Item) *l7resourceResourceModel {
	return &l7resourceResourceModel{
		L7ResourceID:          types.Int64Value(item.L7ResourceID),
		L7ResourceName:        types.StringValue(item.L7ResourceName),
		L7ResourceIsActive:    types.Int64Value(int64(item.L7ResourceIsActive)),
		L7ProtectionDisable:   types.Int64Value(int64(item.L7ProtectionDisable)),
		UseCustomSsl:          types.Int64Value(int64(item.UseCustomSsl)),
		UseLetsencryptSsl:     types.Int64Value(int64(item.UseLetsencryptSsl)),
		CustomSslKey:          types.StringValue(item.CustomSslKey),
		CustomSslCrt:          types.StringValue(item.CustomSslCrt),
		Forcessl:              types.Int64Value(int64(item.Forcessl)),
		ServiceHTTP2:          types.Int64Value(int64(item.ServiceHTTP2)),
		GeoipMode:             types.Int64Value(int64(item.GeoipMode)),
		GeoipList:             types.StringValue(item.GeoipList),
		GlobalWhitelistActive: types.Int64Value(int64(item.GlobalWhitelistActive)),
		HTTP2https:            types.Int64Value(int64(item.HTTP2https)),
		HTTPS2http:            types.Int64Value(int64(item.HTTPS2http)),
		ProtectedIp:           types.StringValue(item.ProtectedIp),
		Wwwredir:              types.Int64Value(int64(item.Wwwredir)),
		Cdn:                   types.Int64Value(int64(item.Cdn)),
		CdnHost:               types.StringValue(item.CdnHost),
		CdnProxyHost:          types.StringValue(item.CdnProxyHost),
	}
}

func expandL7OriginModel(item *l7originResourceModel) *l7origin.Item {
	return &l7origin.Item{
		L7ResourceID: item.L7ResourceID.ValueInt64(),
		ID:           item.ID.ValueInt64(),
		Weight:       item.Weight.ValueInt64(),
		Mode:         item.Mode.ValueString(),
		IP:           item.IP.ValueString(),
		CreatedAt:    item.CreatedAt.ValueInt64(),
		ModifiedAt:   item.ModifiedAt.ValueInt64(),
	}
}

func flatternL7OriginModel(item *l7origin.Item) *l7originResourceModel {
	return &l7originResourceModel{
		L7ResourceID: types.Int64Value(item.L7ResourceID),
		ID:           types.Int64Value(item.ID),
		IP:           types.StringValue(item.IP),
		Weight:       types.Int64Value(item.Weight),
		Mode:         types.StringValue(item.Mode),
		CreatedAt:    types.Int64Value(item.CreatedAt),
		ModifiedAt:   types.Int64Value(item.ModifiedAt),
	}
}

func CheckExistingOriginByIP(ctx context.Context, client *v1.Client, ip string, resourceID int64) (*l7origin.Item, bool) {
	listOpts := &l7origin.ListOpts{
		L7ResourceID: resourceID,
	}

	origins, _, err := l7origin.List(ctx, client, listOpts)
	if err != nil {
		log.Fatal(err)
	}

	for _, item := range origins {
		if item.IP == ip {
			return item, true
		}
	}

	return &l7origin.Item{}, false
}

func CheckPlanVsStateOrigin(origins []*l7originResourceModel, ip string) (*l7origin.Item, bool) {
	for _, item := range origins {
		if item.IP.ValueString() == ip {
			return expandL7OriginModel(item), true
		}
	}

	return &l7origin.Item{}, false
}

func removeL7originFromState(slice []*l7originResourceModel, s int) []*l7originResourceModel {
	return append(slice[:s], slice[s+1:]...)
}

func hackSPSSLState(plan *l7resourceResourceModel, l7res *l7resource.Data) *l7resource.Data {
	// Hack - servicepipe api doesn't support ssl cert params in response
	if !plan.CustomSslKey.IsNull() || !plan.CustomSslKey.IsUnknown() {
		l7res.Data.Result.CustomSslKey = plan.CustomSslKey.ValueString()
	}

	if !plan.CustomSslCrt.IsNull() || !plan.CustomSslCrt.IsUnknown() {
		l7res.Data.Result.CustomSslCrt = plan.CustomSslCrt.ValueString()
	}

	return l7res
}

func CheckPlanVsState(plan *l7resourceResourceModel, state *l7resourceResourceModel, item *l7resource.Item) (*l7resource.Item, bool) {
	update := false

	if !plan.L7ResourceName.Equal(state.L7ResourceName) {
		item.L7ResourceName = plan.L7ResourceName.ValueString()
		update = true
	}

	if !plan.L7ResourceIsActive.Equal(state.L7ResourceIsActive) {
		item.L7ResourceIsActive = int(plan.L7ResourceIsActive.ValueInt64())
		update = true
	}

	if !plan.L7ProtectionDisable.Equal(state.L7ProtectionDisable) {
		item.L7ProtectionDisable = int(plan.L7ProtectionDisable.ValueInt64())
		update = true
	}

	if !plan.UseCustomSsl.Equal(state.UseCustomSsl) {
		item.UseCustomSsl = int(plan.UseCustomSsl.ValueInt64())
		update = true
	}

	if !plan.UseLetsencryptSsl.Equal(state.UseLetsencryptSsl) {
		item.UseLetsencryptSsl = int(plan.UseLetsencryptSsl.ValueInt64())
		update = true
	}

	if !plan.CustomSslKey.Equal(state.CustomSslKey) {
		item.CustomSslKey = plan.CustomSslKey.ValueString()
		update = true
	}

	if !plan.CustomSslCrt.Equal(state.CustomSslCrt) {
		item.CustomSslCrt = plan.CustomSslCrt.ValueString()
		update = true
	}
	if !plan.Forcessl.Equal(state.Forcessl) {
		item.Forcessl = int(plan.Forcessl.ValueInt64())
		update = true
	}

	if !plan.ServiceHTTP2.Equal(state.ServiceHTTP2) {
		item.ServiceHTTP2 = int(plan.ServiceHTTP2.ValueInt64())
		update = true
	}

	if !plan.GeoipMode.Equal(state.GeoipMode) {
		item.GeoipMode = int(plan.GeoipMode.ValueInt64())
		update = true
	}

	if !plan.GeoipList.Equal(state.GeoipList) {
		item.GeoipList = plan.GeoipList.ValueString()
		update = true
	}

	if !plan.GlobalWhitelistActive.Equal(state.GlobalWhitelistActive) {
		item.GlobalWhitelistActive = int(plan.GlobalWhitelistActive.ValueInt64())
		update = true
	}

	if !plan.HTTP2https.Equal(state.HTTP2https) {
		item.HTTP2https = int(plan.HTTP2https.ValueInt64())
		update = true
	}

	if !plan.HTTPS2http.Equal(state.HTTPS2http) {
		item.HTTPS2http = int(plan.HTTPS2http.ValueInt64())
		update = true
	}

	if !plan.ProtectedIp.Equal(state.ProtectedIp) {
		item.ProtectedIp = plan.ProtectedIp.ValueString()
		update = true
	}

	if !plan.Wwwredir.Equal(state.Wwwredir) {
		item.Wwwredir = int(plan.Wwwredir.ValueInt64())
		update = true
	}

	if !plan.Cdn.Equal(state.Cdn) {
		item.Cdn = int(plan.Cdn.ValueInt64())
		update = true
	}

	if !plan.CdnHost.Equal(state.CdnHost) {
		item.CdnHost = plan.CdnHost.ValueString()
		update = true
	}

	if !plan.CdnProxyHost.Equal(state.CdnProxyHost) {
		item.CdnProxyHost = plan.CdnProxyHost.ValueString()
		update = true
	}

	return item, update
}

func CheckingL7resourcePlanAttrIsNull(plan l7resourceResourceModel, item *l7resource.Item) (*l7resource.Item, bool) {
	update := false

	if !plan.L7ResourceName.IsNull() || !plan.L7ResourceName.IsUnknown() {
		item.L7ResourceName = plan.L7ResourceName.ValueString()
		update = true
	}

	if !plan.L7ResourceIsActive.IsNull() || !plan.L7ResourceIsActive.IsUnknown() {
		item.L7ResourceIsActive = int(plan.L7ResourceIsActive.ValueInt64())
		update = true
	}

	if !plan.L7ProtectionDisable.IsNull() || !plan.L7ProtectionDisable.IsUnknown() {
		item.L7ProtectionDisable = int(plan.L7ProtectionDisable.ValueInt64())
		update = true
	}

	if !plan.UseCustomSsl.IsNull() || !plan.UseCustomSsl.IsUnknown() {
		item.UseCustomSsl = int(plan.UseCustomSsl.ValueInt64())
		update = true
	}

	if !plan.UseLetsencryptSsl.IsNull() || !plan.UseLetsencryptSsl.IsUnknown() {
		item.UseLetsencryptSsl = int(plan.UseLetsencryptSsl.ValueInt64())
		update = true
	}

	if !plan.CustomSslKey.IsNull() || !plan.CustomSslKey.IsUnknown() {
		item.CustomSslKey = plan.CustomSslKey.ValueString()
		update = true
	}

	if !plan.CustomSslCrt.IsNull() || !plan.CustomSslCrt.IsUnknown() {
		item.CustomSslCrt = plan.CustomSslCrt.ValueString()
		update = true
	}

	if !plan.Forcessl.IsNull() || !plan.Forcessl.IsUnknown() {
		item.Forcessl = int(plan.Forcessl.ValueInt64())
		update = true
	}

	if !plan.ServiceHTTP2.IsNull() || !plan.ServiceHTTP2.IsUnknown() {
		item.ServiceHTTP2 = int(plan.ServiceHTTP2.ValueInt64())
		update = true
	}

	if !plan.GeoipMode.IsNull() || !plan.GeoipMode.IsUnknown() {
		item.GeoipMode = int(plan.GeoipMode.ValueInt64())
		update = true
	}

	if !plan.GeoipList.IsNull() || !plan.GeoipList.IsUnknown() {
		item.GeoipList = plan.GeoipList.ValueString()
		update = true
	}

	if !plan.GlobalWhitelistActive.IsNull() || !plan.GlobalWhitelistActive.IsUnknown() {
		item.GlobalWhitelistActive = int(plan.GlobalWhitelistActive.ValueInt64())
		update = true
	}

	if !plan.HTTP2https.IsNull() || !plan.HTTP2https.IsUnknown() {
		item.HTTP2https = int(plan.HTTP2https.ValueInt64())
		update = true
	}

	if !plan.HTTPS2http.IsNull() || !plan.HTTPS2http.IsUnknown() {
		item.HTTPS2http = int(plan.HTTPS2http.ValueInt64())
		update = true
	}

	if !plan.ProtectedIp.IsNull() || !plan.ProtectedIp.IsUnknown() {
		item.ProtectedIp = plan.ProtectedIp.ValueString()
		update = true
	}

	if !plan.Wwwredir.IsNull() || !plan.Wwwredir.IsUnknown() {
		item.Wwwredir = int(plan.Wwwredir.ValueInt64())
		update = true
	}

	if !plan.Cdn.IsNull() || !plan.Cdn.IsUnknown() {
		item.Cdn = int(plan.Cdn.ValueInt64())
		update = true
	}

	if !plan.CdnHost.IsNull() || !plan.CdnHost.IsUnknown() {
		item.CdnHost = plan.CdnHost.ValueString()
		update = true
	}

	if !plan.CdnProxyHost.IsNull() || !plan.CdnProxyHost.IsUnknown() {
		item.CdnProxyHost = plan.CdnProxyHost.ValueString()
		update = true
	}

	return item, update
}
