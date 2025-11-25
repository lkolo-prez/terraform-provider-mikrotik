package mikrotik

import (
	"context"

	"github.com/ddelnano/terraform-provider-mikrotik/client"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	tftypes "github.com/hashicorp/terraform-plugin-framework/types"
)

type bgpTemplate struct {
	client *client.Mikrotik
}

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &bgpTemplate{}
	_ resource.ResourceWithConfigure   = &bgpTemplate{}
	_ resource.ResourceWithImportState = &bgpTemplate{}
)

// NewBgpTemplateResource is a helper function to simplify the provider implementation.
func NewBgpTemplateResource() resource.Resource {
	return &bgpTemplate{}
}

func (r *bgpTemplate) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*client.Mikrotik)
}

func (r *bgpTemplate) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_bgp_template"
}

// Schema defines the schema for the resource.
func (r *bgpTemplate) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages MikroTik BGP template configuration (RouterOS v7+). Templates allow configuration reuse across multiple BGP connections. The template contains all BGP protocol-related configuration options.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "Unique ID of this resource.",
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "Name of the BGP template.",
			},
			"as": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Default:     int64default.StaticInt64(0),
				Description: "Default AS number for connections using this template.",
			},
			"disabled": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
				Description: "Whether this BGP template is disabled.",
			},
			"router_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
				Description: "BGP router ID for connections using this template.",
			},
			"address_families": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("ip"),
				Description: "Supported address families (ip, ipv6, l2vpn, vpnv4, vpnv6).",
			},
			"capabilities": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
				Description: "BGP capabilities to advertise (e.g., mp, rr, as4, gr).",
			},
			"as_override": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
				Description: "Whether to override remote AS number with local AS in updates.",
			},
			"cisco": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
				Description: "Whether to use Cisco-compatible BGP route-refresh.",
			},
			"remove_private_as": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
				Description: "Whether to remove private AS numbers from AS_PATH.",
			},
			"comment": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
				Description: "Comment for this BGP template.",
			},
			"connect_retry_time": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("2m"),
				Description: "Connect retry time for failed connections.",
			},
			"hold_time": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("3m"),
				Description: "BGP hold time (e.g., 3m).",
			},
			"keepalive_time": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("1m"),
				Description: "BGP keepalive time (e.g., 1m).",
			},
			"input_affixes": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
				Description: "Input prefix list names.",
			},
			"input_filter": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
				Description: "Input routing filter chain name.",
			},
			"input_limit": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Default:     int64default.StaticInt64(0),
				Description: "Maximum number of prefixes to accept (0 = unlimited).",
			},
			"input_accept_communities": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
				Description: "Which BGP communities to accept (standard, extended, large, all).",
			},
			"input_accept_nlri": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
				Description: "Which NLRI types to accept from peer (unicast, multicast, both).",
			},
			"input_accept_originated": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
				Description: "Whether to accept routes originated by this router.",
			},
			"input_ignore_as_path_len": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
				Description: "Whether to ignore AS_PATH length in best path selection.",
			},
			"input_limit_process_routes_ipv4": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Default:     int64default.StaticInt64(0),
				Description: "Limit on the number of IPv4 routes to process from this peer.",
			},
			"input_limit_process_routes_ipv6": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Default:     int64default.StaticInt64(0),
				Description: "Limit on the number of IPv6 routes to process from this peer.",
			},
			"multihop": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
				Description: "Whether to allow multihop BGP connections.",
			},
			"use_bfd": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
				Description: "Whether to use BFD for fast failure detection.",
			},
			"ttl": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("default"),
				Description: "TTL for BGP packets (default or specific value).",
			},
			"nexthop_choice": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("default"),
				Description: "Nexthop selection method (default, force-self, propagate).",
			},
			"output_affixes": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
				Description: "Output prefix list names.",
			},
			"output_default_originate": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("never"),
				Description: "When to originate default route (never, always, if-installed).",
			},
			"output_filter": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
				Description: "Output routing filter chain name.",
			},
			"output_filter_chain": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
				Description: "Output filter chain for advanced filtering.",
			},
			"output_keepalive_time": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
				Description: "Custom keepalive time for output.",
			},
			"output_network": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
				Description: "Networks to advertise via BGP.",
			},
			"output_redistribute": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
				Description: "Which routes to redistribute into BGP (connected, static, ospf, rip).",
			},
			"passive": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
				Description: "Whether to wait for remote peer to initiate connection.",
			},
			"route_reflect": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
				Description: "Whether this router is a route reflector for this template.",
			},
			"graceful_restart": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
				Description: "Graceful restart configuration (yes, no, restart-time, stale-time).",
			},
		},
	}
}

// Create creates the resource and sets the initial Terraform state.
func (r *bgpTemplate) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var terraformModel bgpTemplateModel
	var mikrotikModel client.BgpTemplate
	GenericCreateResource(&terraformModel, &mikrotikModel, r.client)(ctx, req, resp)
}

// Read refreshes the Terraform state with the latest data.
func (r *bgpTemplate) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var terraformModel bgpTemplateModel
	var mikrotikModel client.BgpTemplate
	GenericReadResource(&terraformModel, &mikrotikModel, r.client)(ctx, req, resp)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *bgpTemplate) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var terraformModel bgpTemplateModel
	var mikrotikModel client.BgpTemplate
	GenericUpdateResource(&terraformModel, &mikrotikModel, r.client)(ctx, req, resp)
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *bgpTemplate) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var terraformModel bgpTemplateModel
	var mikrotikModel client.BgpTemplate
	GenericDeleteResource(&terraformModel, &mikrotikModel, r.client)(ctx, req, resp)
}

func (r *bgpTemplate) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}

type bgpTemplateModel struct {
	Id                              tftypes.String `tfsdk:"id"`
	Name                            tftypes.String `tfsdk:"name"`
	AS                              tftypes.Int64  `tfsdk:"as"`
	Disabled                        tftypes.Bool   `tfsdk:"disabled"`
	RouterID                        tftypes.String `tfsdk:"router_id"`
	AddressFamily                   tftypes.String `tfsdk:"address_families"`
	Capabilities                    tftypes.String `tfsdk:"capabilities"`
	AsOverride                      tftypes.Bool   `tfsdk:"as_override"`
	Cisco                           tftypes.Bool   `tfsdk:"cisco"`
	RemovePrivateAS                 tftypes.Bool   `tfsdk:"remove_private_as"`
	Comment                         tftypes.String `tfsdk:"comment"`
	ConnectRetryTime                tftypes.String `tfsdk:"connect_retry_time"`
	HoldTime                        tftypes.String `tfsdk:"hold_time"`
	KeepaliveTime                   tftypes.String `tfsdk:"keepalive_time"`
	InputAffixFilters               tftypes.String `tfsdk:"input_affixes"`
	InputFilter                     tftypes.String `tfsdk:"input_filter"`
	InputLimit                      tftypes.Int64  `tfsdk:"input_limit"`
	InputAcceptCommunities          tftypes.String `tfsdk:"input_accept_communities"`
	InputAcceptNLRI                 tftypes.String `tfsdk:"input_accept_nlri"`
	InputAcceptOriginated           tftypes.Bool   `tfsdk:"input_accept_originated"`
	InputIgnoreAsPathLen            tftypes.Bool   `tfsdk:"input_ignore_as_path_len"`
	InputLimitProcessRoutesIPv4     tftypes.Int64  `tfsdk:"input_limit_process_routes_ipv4"`
	InputLimitProcessRoutesIPv6     tftypes.Int64  `tfsdk:"input_limit_process_routes_ipv6"`
	Multihop                        tftypes.Bool   `tfsdk:"multihop"`
	UseBFD                          tftypes.Bool   `tfsdk:"use_bfd"`
	TTL                             tftypes.String `tfsdk:"ttl"`
	NexthopChoice                   tftypes.String `tfsdk:"nexthop_choice"`
	OutputAffixFilters              tftypes.String `tfsdk:"output_affixes"`
	OutputDefaultOriginate          tftypes.String `tfsdk:"output_default_originate"`
	OutputFilter                    tftypes.String `tfsdk:"output_filter"`
	OutputFilterChain               tftypes.String `tfsdk:"output_filter_chain"`
	OutputKeepaliveTime             tftypes.String `tfsdk:"output_keepalive_time"`
	OutputNetwork                   tftypes.String `tfsdk:"output_network"`
	OutputRedistribute              tftypes.String `tfsdk:"output_redistribute"`
	Passive                         tftypes.Bool   `tfsdk:"passive"`
	RouteReflect                    tftypes.Bool   `tfsdk:"route_reflect"`
	GracefulRestart                 tftypes.String `tfsdk:"graceful_restart"`
}
