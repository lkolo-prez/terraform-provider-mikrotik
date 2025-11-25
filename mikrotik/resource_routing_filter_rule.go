package mikrotik

import (
	"context"

	"github.com/ddelnano/terraform-provider-mikrotik/client"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type routingFilterRule struct {
	client *client.Mikrotik
}

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &routingFilterRule{}
	_ resource.ResourceWithConfigure   = &routingFilterRule{}
	_ resource.ResourceWithImportState = &routingFilterRule{}
)

// NewRoutingFilterRuleResource is a helper function to simplify the provider implementation.
func NewRoutingFilterRuleResource() resource.Resource {
	return &routingFilterRule{}
}

func (r *routingFilterRule) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*client.Mikrotik)
}

// Metadata returns the resource type name.
func (r *routingFilterRule) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_routing_filter_rule"
}

// Schema defines the schema for the resource.
func (r *routingFilterRule) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a RouterOS v7 routing filter rule. " +
			"Routing filters in v7 use a new rule-based syntax for filtering and manipulating routes from BGP, OSPF, and other routing protocols.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Description: "Unique ID of this resource.",
			},
			"chain": schema.StringAttribute{
				Required:    true,
				Description: "Name of the filter chain this rule belongs to. Multiple rules can share the same chain.",
			},
			"rule": schema.StringAttribute{
				Required: true,
				Description: "Filter rule expression in RouterOS v7 syntax. " +
					"Examples: " +
					"`if (dst == 0.0.0.0/0) { reject }` - Deny default route. " +
					"`if (dst in 10.0.0.0/8 && bgp-communities includes 65001:100) { accept }` - Accept prefixes with community. " +
					"`if (bgp-communities includes 65001:200) { set bgp-local-pref 200; accept }` - Set local-pref and accept. " +
					"Supports prefix matching (dst, dst-len), BGP attributes (bgp-as-path, bgp-communities, bgp-local-pref, bgp-med), " +
					"OSPF attributes (ospf-type, ospf-tag), and actions (accept, reject, jump, return, set).",
			},
			"disabled": schema.BoolAttribute{
				Optional: true,
				Computed: true,
				Default:  booldefault.StaticBool(false),
				Description: "Whether the rule is disabled. Disabled rules are not evaluated.",
			},
			"comment": schema.StringAttribute{
				Optional:    true,
				Description: "Comment for the routing filter rule.",
			},
			"invalid": schema.BoolAttribute{
				Computed:    true,
				Description: "Whether the rule syntax is invalid. Read-only.",
			},
			"dynamic": schema.BoolAttribute{
				Computed:    true,
				Description: "Whether the rule is dynamically created. Read-only.",
			},
		},
	}
}

// routingFilterRuleModel maps the resource schema data.
type routingFilterRuleModel struct {
	Id       types.String `tfsdk:"id"`
	Chain    types.String `tfsdk:"chain"`
	Rule     types.String `tfsdk:"rule"`
	Disabled types.Bool   `tfsdk:"disabled"`
	Comment  types.String `tfsdk:"comment"`
	Invalid  types.Bool   `tfsdk:"invalid"`
	Dynamic  types.Bool   `tfsdk:"dynamic"`
}

// Create creates the resource and sets the initial Terraform state.
func (r *routingFilterRule) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan routingFilterRuleModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	rule := &client.RoutingFilterRule{
		Chain:    plan.Chain.ValueString(),
		Rule:     plan.Rule.ValueString(),
		Disabled: plan.Disabled.ValueBool(),
		Comment:  plan.Comment.ValueString(),
	}

	created, err := r.client.AddRoutingFilterRule(rule)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating routing filter rule",
			"Could not create routing filter rule: "+err.Error(),
		)
		return
	}

	// Map response to state
	plan.Id = types.StringValue(created.Id)
	plan.Invalid = types.BoolValue(created.Invalid)
	plan.Dynamic = types.BoolValue(created.Dynamic)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

// Read refreshes the Terraform state with the latest data.
func (r *routingFilterRule) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state routingFilterRuleModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	rule, err := r.client.FindRoutingFilterRuleById(state.Id.ValueString())
	if err != nil {
		if client.IsNotFoundError(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error reading routing filter rule",
			"Could not read routing filter rule "+state.Id.ValueString()+": "+err.Error(),
		)
		return
	}

	// Map response to state
	state.Chain = types.StringValue(rule.Chain)
	state.Rule = types.StringValue(rule.Rule)
	state.Disabled = types.BoolValue(rule.Disabled)
	state.Comment = types.StringValue(rule.Comment)
	state.Invalid = types.BoolValue(rule.Invalid)
	state.Dynamic = types.BoolValue(rule.Dynamic)

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *routingFilterRule) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan routingFilterRuleModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	rule := &client.RoutingFilterRule{
		Id:       plan.Id.ValueString(),
		Chain:    plan.Chain.ValueString(),
		Rule:     plan.Rule.ValueString(),
		Disabled: plan.Disabled.ValueBool(),
		Comment:  plan.Comment.ValueString(),
	}

	updated, err := r.client.UpdateRoutingFilterRule(rule)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating routing filter rule",
			"Could not update routing filter rule "+plan.Id.ValueString()+": "+err.Error(),
		)
		return
	}

	// Map response to state
	plan.Invalid = types.BoolValue(updated.Invalid)
	plan.Dynamic = types.BoolValue(updated.Dynamic)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *routingFilterRule) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state routingFilterRuleModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteRoutingFilterRule(state.Id.ValueString())
	if err != nil {
		if client.IsNotFoundError(err) {
			return
		}
		resp.Diagnostics.AddError(
			"Error deleting routing filter rule",
			"Could not delete routing filter rule "+state.Id.ValueString()+": "+err.Error(),
		)
		return
	}
}

// ImportState imports the resource state.
func (r *routingFilterRule) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Import by ID
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
