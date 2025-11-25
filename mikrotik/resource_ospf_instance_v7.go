package mikrotik

import (
	"context"

	"github.com/ddelnano/terraform-provider-mikrotik/client"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type ospfInstanceV7 struct {
	Name                  types.String `tfsdk:"name"`
	Version               types.String `tfsdk:"version"`
	RouterId              types.String `tfsdk:"router_id"`
	DomainId              types.String `tfsdk:"domain_id"`
	Disabled              types.Bool   `tfsdk:"disabled"`
	Comment               types.String `tfsdk:"comment"`
	Vrf                   types.String `tfsdk:"vrf"`
	RoutingTable          types.String `tfsdk:"routing_table"`
	RedistributeConnected types.Bool   `tfsdk:"redistribute_connected"`
	RedistributeStatic    types.Bool   `tfsdk:"redistribute_static"`
	RedistributeBgp       types.Bool   `tfsdk:"redistribute_bgp"`
	RedistributeRip       types.Bool   `tfsdk:"redistribute_rip"`
	RedistributeOspf      types.Bool   `tfsdk:"redistribute_ospf"`
	OriginateDefault      types.String `tfsdk:"originate_default"`
	InFilterChain         types.String `tfsdk:"in_filter_chain"`
	OutFilterChain        types.String `tfsdk:"out_filter_chain"`
	// Computed
	Id            types.String `tfsdk:"id"`
	RoutingMarks  types.String `tfsdk:"routing_marks"`
	Dynamic       types.Bool   `tfsdk:"dynamic"`
	Invalid       types.Bool   `tfsdk:"invalid"`
}

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &ospfInstanceV7Resource{}
	_ resource.ResourceWithConfigure   = &ospfInstanceV7Resource{}
	_ resource.ResourceWithImportState = &ospfInstanceV7Resource{}
)

// NewOspfInstanceV7Resource is a helper function to simplify the provider implementation.
func NewOspfInstanceV7Resource() resource.Resource {
	return &ospfInstanceV7Resource{}
}

type ospfInstanceV7Resource struct {
	client *client.Mikrotik
}

func (r *ospfInstanceV7Resource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*client.Mikrotik)
}

func (r *ospfInstanceV7Resource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_ospf_instance_v7"
}

// Schema defines the schema for the resource.
func (r *ospfInstanceV7Resource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages OSPF v2/v3 instance (RouterOS v7 unified OSPF).",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "Unique ID of the OSPF instance.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "Name of the OSPF instance.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"version": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("2"),
				Description: "OSPF version: '2' for IPv4 or '3' for IPv6. Default is '2'.",
				Validators: []validator.String{
					stringvalidator.OneOf("2", "3"),
				},
			},
			"router_id": schema.StringAttribute{
				Optional:    true,
				Description: "Router ID in IPv4 address format (e.g., 1.1.1.1). If not set, RouterOS auto-selects one.",
			},
			"domain_id": schema.StringAttribute{
				Optional:    true,
				Description: "OSPF domain ID for multi-instance support (RFC 6549). Used to distinguish multiple OSPF instances.",
			},
			"disabled": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
				Description: "Whether the OSPF instance is disabled. Default is false.",
			},
			"comment": schema.StringAttribute{
				Optional:    true,
				Description: "Comment for the OSPF instance.",
			},
			"vrf": schema.StringAttribute{
				Optional:    true,
				Description: "VRF name to associate this OSPF instance with. Requires VRF to be configured.",
			},
			"routing_table": schema.StringAttribute{
				Optional:    true,
				Description: "Routing table name to use for this OSPF instance. Leave empty for main table.",
			},
			"redistribute_connected": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
				Description: "Redistribute connected routes into OSPF. Default is false.",
			},
			"redistribute_static": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
				Description: "Redistribute static routes into OSPF. Default is false.",
			},
			"redistribute_bgp": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
				Description: "Redistribute BGP routes into OSPF. Default is false.",
			},
			"redistribute_rip": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
				Description: "Redistribute RIP routes into OSPF. Default is false.",
			},
			"redistribute_ospf": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
				Description: "Redistribute other OSPF instance routes. Default is false.",
			},
			"originate_default": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("never"),
				Description: "Control default route (0.0.0.0/0) origination: 'never' (default), 'always', or 'if-installed'.",
				Validators: []validator.String{
					stringvalidator.OneOf("never", "always", "if-installed"),
				},
			},
			"in_filter_chain": schema.StringAttribute{
				Optional:    true,
				Description: "Name of the routing filter chain for incoming routes.",
			},
			"out_filter_chain": schema.StringAttribute{
				Optional:    true,
				Description: "Name of the routing filter chain for outgoing routes.",
			},
			"routing_marks": schema.StringAttribute{
				Computed:    true,
				Description: "Routing marks associated with this instance.",
			},
			"dynamic": schema.BoolAttribute{
				Computed:    true,
				Description: "Whether this instance was dynamically created.",
			},
			"invalid": schema.BoolAttribute{
				Computed:    true,
				Description: "Whether this instance configuration is invalid.",
			},
		},
	}
}

// Create creates the resource and sets the initial Terraform state.
func (r *ospfInstanceV7Resource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var terraformModel ospfInstanceV7
	var mikrotikModel client.OspfInstanceV7
	GenericCreateResource(&terraformModel, &mikrotikModel, r.client)(ctx, req, resp)
}

// Read refreshes the Terraform state with the latest data.
func (r *ospfInstanceV7Resource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var terraformModel ospfInstanceV7
	var mikrotikModel client.OspfInstanceV7
	GenericReadResource(&terraformModel, &mikrotikModel, r.client)(ctx, req, resp)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *ospfInstanceV7Resource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var terraformModel ospfInstanceV7
	var mikrotikModel client.OspfInstanceV7
	GenericUpdateResource(&terraformModel, &mikrotikModel, r.client)(ctx, req, resp)
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *ospfInstanceV7Resource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var terraformModel ospfInstanceV7
	var mikrotikModel client.OspfInstanceV7
	GenericDeleteResource(&terraformModel, &mikrotikModel, r.client)(ctx, req, resp)
}

func (r *ospfInstanceV7Resource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Import by name
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}
