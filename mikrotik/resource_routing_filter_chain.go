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

type routingFilterChain struct {
	client *client.Mikrotik
}

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &routingFilterChain{}
	_ resource.ResourceWithConfigure   = &routingFilterChain{}
	_ resource.ResourceWithImportState = &routingFilterChain{}
)

// NewRoutingFilterChainResource is a helper function to simplify the provider implementation.
func NewRoutingFilterChainResource() resource.Resource {
	return &routingFilterChain{}
}

func (r *routingFilterChain) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*client.Mikrotik)
}

// Metadata returns the resource type name.
func (r *routingFilterChain) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_routing_filter_chain"
}

// Schema defines the schema for the resource.
func (r *routingFilterChain) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a RouterOS v7 routing filter chain. " +
			"Chains group multiple filter rules together and can be referenced by BGP connections, OSPF instances, and other routing protocols.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Description: "Unique ID of this resource.",
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "Name of the filter chain. Must be unique.",
			},
			"dynamic": schema.BoolAttribute{
				Optional: true,
				Computed: true,
				Default:  booldefault.StaticBool(false),
				Description: "Whether the chain is dynamic. Dynamic chains can be modified by the system.",
			},
			"disabled": schema.BoolAttribute{
				Optional: true,
				Computed: true,
				Default:  booldefault.StaticBool(false),
				Description: "Whether the chain is disabled. Disabled chains are not evaluated.",
			},
			"comment": schema.StringAttribute{
				Optional:    true,
				Description: "Comment for the routing filter chain.",
			},
		},
	}
}

// routingFilterChainModel maps the resource schema data.
type routingFilterChainModel struct {
	Id      types.String `tfsdk:"id"`
	Name    types.String `tfsdk:"name"`
	Dynamic types.Bool   `tfsdk:"dynamic"`
	Disabled types.Bool  `tfsdk:"disabled"`
	Comment types.String `tfsdk:"comment"`
}

// Create creates the resource and sets the initial Terraform state.
func (r *routingFilterChain) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan routingFilterChainModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	chain := &client.RoutingFilterChain{
		Name:         plan.Name.ValueString(),
		DynamicChain: plan.Dynamic.ValueBool(),
		Disabled:     plan.Disabled.ValueBool(),
		Comment:      plan.Comment.ValueString(),
	}

	created, err := r.client.AddRoutingFilterChain(chain)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating routing filter chain",
			"Could not create routing filter chain: "+err.Error(),
		)
		return
	}

	// Map response to state
	plan.Id = types.StringValue(created.Id)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

// Read refreshes the Terraform state with the latest data.
func (r *routingFilterChain) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state routingFilterChainModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	chain, err := r.client.FindRoutingFilterChain(state.Name.ValueString())
	if err != nil {
		if client.IsNotFoundError(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error reading routing filter chain",
			"Could not read routing filter chain "+state.Name.ValueString()+": "+err.Error(),
		)
		return
	}

	// Map response to state
	state.Id = types.StringValue(chain.Id)
	state.Name = types.StringValue(chain.Name)
	state.Dynamic = types.BoolValue(chain.DynamicChain)
	state.Disabled = types.BoolValue(chain.Disabled)
	state.Comment = types.StringValue(chain.Comment)

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *routingFilterChain) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan routingFilterChainModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	chain := &client.RoutingFilterChain{
		Id:           plan.Id.ValueString(),
		Name:         plan.Name.ValueString(),
		DynamicChain: plan.Dynamic.ValueBool(),
		Disabled:     plan.Disabled.ValueBool(),
		Comment:      plan.Comment.ValueString(),
	}

	_, err := r.client.UpdateRoutingFilterChain(chain)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating routing filter chain",
			"Could not update routing filter chain "+plan.Name.ValueString()+": "+err.Error(),
		)
		return
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *routingFilterChain) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state routingFilterChainModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteRoutingFilterChain(state.Name.ValueString())
	if err != nil {
		if client.IsNotFoundError(err) {
			return
		}
		resp.Diagnostics.AddError(
			"Error deleting routing filter chain",
			"Could not delete routing filter chain "+state.Name.ValueString()+": "+err.Error(),
		)
		return
	}
}

// ImportState imports the resource state.
func (r *routingFilterChain) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Import by name
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}
