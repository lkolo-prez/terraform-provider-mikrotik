package mikrotik

import (
	"context"

	"github.com/ddelnano/terraform-provider-mikrotik/client"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"

	tftypes "github.com/hashicorp/terraform-plugin-framework/types"
)

type bgpInstanceV7 struct {
	client *client.Mikrotik
}

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &bgpInstanceV7{}
	_ resource.ResourceWithConfigure   = &bgpInstanceV7{}
	_ resource.ResourceWithImportState = &bgpInstanceV7{}
)

// NewBgpInstanceV7Resource is a helper function to simplify the provider implementation.
func NewBgpInstanceV7Resource() resource.Resource {
	return &bgpInstanceV7{}
}

func (r *bgpInstanceV7) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*client.Mikrotik)
}

// Metadata returns the resource type name.
func (r *bgpInstanceV7) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_bgp_instance_v7"
}

// Schema defines the schema for the resource.
func (s *bgpInstanceV7) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Creates a MikroTik BGP Instance for RouterOS v7.20+. Starting from RouterOS v7.20, BGP instances are explicitly defined instead of auto-detecting based on router-ids. BGP routing instance is necessary for best path route selection and other instance-dependent features like VPN and EVPN.",

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
				Description: "Name of the BGP instance. Required field since RouterOS v7.20.",
			},
			"as": schema.Int64Attribute{
				Required:    true,
				Description: "The 32-bit BGP autonomous system number. Must be a value within 0 to 4294967295.",
			},
			"router_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
				Description: "BGP Router ID for this instance. If not specified or set to 0.0.0.0, BGP will automatically select one of router's IP addresses.",
			},
			"client_to_client_reflection": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(true),
				Description: "In case this instance is a route reflector: whether to redistribute routes learned from one routing reflection client to other clients.",
			},
			"cluster_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
				Description: "In case this instance is a route reflector: cluster ID of the route reflector cluster this instance belongs to.",
			},
			"confederation": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Default:     int64default.StaticInt64(0),
				Description: "In case of BGP confederations: autonomous system number that identifies the confederation as a whole.",
			},
			"ignore_as_path_len": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
				Description: "Whether to ignore AS_PATH attribute in BGP route selection algorithm.",
			},
			"out_filter": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
				Description: "Output routing filter chain used by all BGP peers belonging to this instance.",
			},
			"routing_table": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
				Description: "Name of routing table this BGP instance operates on.",
			},
			"redistribute_connected": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
				Description: "If enabled, this BGP instance will redistribute connected routes.",
			},
			"redistribute_ospf": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
				Description: "If enabled, this BGP instance will redistribute routes learned by OSPF.",
			},
			"redistribute_other_bgp": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
				Description: "If enabled, this BGP instance will redistribute routes learned by other BGP instances.",
			},
			"redistribute_rip": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
				Description: "If enabled, this BGP instance will redistribute routes learned by RIP.",
			},
			"redistribute_static": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
				Description: "If enabled, the router will redistribute static routes added to its routing database.",
			},
			"disabled": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
				Description: "Whether instance is disabled.",
			},
			"comment": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
				Description: "Comment for this BGP instance.",
			},
			"vrf": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
				Description: "VRF (Virtual Routing and Forwarding) instance for VPN/EVPN support.",
			},
		},
	}
}

// Create creates the resource and sets the initial Terraform state.
func (r *bgpInstanceV7) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var terraformModel bgpInstanceV7Model
	var mikrotikModel client.BgpInstanceV7
	GenericCreateResource(&terraformModel, &mikrotikModel, r.client)(ctx, req, resp)
}

// Read refreshes the Terraform state with the latest data.
func (r *bgpInstanceV7) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var terraformModel bgpInstanceV7Model
	var mikrotikModel client.BgpInstanceV7
	GenericReadResource(&terraformModel, &mikrotikModel, r.client)(ctx, req, resp)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *bgpInstanceV7) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var terraformModel bgpInstanceV7Model
	var mikrotikModel client.BgpInstanceV7
	GenericUpdateResource(&terraformModel, &mikrotikModel, r.client)(ctx, req, resp)
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *bgpInstanceV7) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var terraformModel bgpInstanceV7Model
	var mikrotikModel client.BgpInstanceV7
	GenericDeleteResource(&terraformModel, &mikrotikModel, r.client)(ctx, req, resp)
}

func (r *bgpInstanceV7) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}

type bgpInstanceV7Model struct {
	Id                       tftypes.String `tfsdk:"id"`
	Name                     tftypes.String `tfsdk:"name"`
	AS                       tftypes.Int64  `tfsdk:"as"`
	RouterID                 tftypes.String `tfsdk:"router_id"`
	ClientToClientReflection tftypes.Bool   `tfsdk:"client_to_client_reflection"`
	ClusterID                tftypes.String `tfsdk:"cluster_id"`
	Confederation            tftypes.Int64  `tfsdk:"confederation"`
	IgnoreAsPathLen          tftypes.Bool   `tfsdk:"ignore_as_path_len"`
	OutFilter                tftypes.String `tfsdk:"out_filter"`
	RoutingTable             tftypes.String `tfsdk:"routing_table"`
	RedistributeConnected    tftypes.Bool   `tfsdk:"redistribute_connected"`
	RedistributeOspf         tftypes.Bool   `tfsdk:"redistribute_ospf"`
	RedistributeOtherBgp     tftypes.Bool   `tfsdk:"redistribute_other_bgp"`
	RedistributeRip          tftypes.Bool   `tfsdk:"redistribute_rip"`
	RedistributeStatic       tftypes.Bool   `tfsdk:"redistribute_static"`
	Disabled                 tftypes.Bool   `tfsdk:"disabled"`
	Comment                  tftypes.String `tfsdk:"comment"`
	VRF                      tftypes.String `tfsdk:"vrf"`
}
