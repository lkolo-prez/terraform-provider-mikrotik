package mikrotik

import (
	"context"

	"github.com/ddelnano/terraform-provider-mikrotik/client"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	tftypes "github.com/hashicorp/terraform-plugin-framework/types"
)

type bgpConnection struct {
	client *client.Mikrotik
}

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &bgpConnection{}
	_ resource.ResourceWithConfigure   = &bgpConnection{}
	_ resource.ResourceWithImportState = &bgpConnection{}
)

// NewBgpConnectionResource is a helper function to simplify the provider implementation.
func NewBgpConnectionResource() resource.Resource {
	return &bgpConnection{}
}

func (r *bgpConnection) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*client.Mikrotik)
}

func (r *bgpConnection) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_bgp_connection"
}

// Schema defines the schema for the resource.
func (r *bgpConnection) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages MikroTik BGP connection configuration (RouterOS v7+). The BGP connection menu defines BGP outgoing connections as well as acts as a template matcher for incoming BGP connections.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "Unique ID of this resource.",
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "Name of the BGP connection.",
			},
			"as": schema.Int64Attribute{
				Required:    true,
				Description: "Local AS number for this BGP connection.",
			},
			"instance": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
				Description: "BGP instance name. Required since RouterOS v7.20.",
			},
			"disabled": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
				Description: "Whether this BGP connection is disabled.",
			},
			"local_role": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("ebgp"),
				Description: "Role of this peer in the connection (ibgp, ebgp, rr-client).",
				Validators: []validator.String{
					stringvalidator.OneOf("ibgp", "ebgp", "rr-client"),
				},
			},
			"local_address": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
				Description: "Local IP address used for this BGP connection.",
			},
			"remote_address": schema.StringAttribute{
				Required:    true,
				Description: "Remote peer IP address or subnet (for listen mode).",
			},
			"remote_as": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Default:     int64default.StaticInt64(0),
				Description: "Remote peer AS number. If not specified, uses the same AS as local (iBGP).",
			},
			"remote_port": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Default:     int64default.StaticInt64(179),
				Description: "Remote BGP port (default 179).",
			},
			"listen": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
				Description: "Whether to listen for incoming connections on specified subnet. Should not be enabled in unsafe environments.",
			},
			"router_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
				Description: "BGP router ID for this connection.",
			},
			"nexthop_choice": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("default"),
				Description: "Nexthop selection method (default, force-self, propagate).",
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
			"connect_retry_time": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("2m"),
				Description: "Connect retry time for failed connections.",
			},
			"ttl": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("default"),
				Description: "TTL for BGP packets (default or specific value).",
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
			"address_families": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("ip"),
				Description: "Supported address families (ip, ipv6, l2vpn, vpnv4, vpnv6).",
			},
			"input_filter": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
				Description: "Input routing filter chain name.",
			},
			"input_accept_nlri": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
				Description: "Which NLRI types to accept from peer (unicast, multicast, both).",
			},
			"input_accept_communities": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
				Description: "Which BGP communities to accept (standard, extended, large, all).",
			},
			"output_filter": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
				Description: "Output routing filter chain name.",
			},
			"output_default_originate": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("never"),
				Description: "When to originate default route (never, always, if-installed).",
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
			"tcp_md5_key": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Sensitive:   true,
				Default:     stringdefault.StaticString(""),
				Description: "TCP MD5 authentication key.",
			},
			"use_mpls": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
				Description: "Whether to use MPLS for this connection.",
			},
			"vpnv4": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
				Description: "Whether to support VPNv4 address family.",
			},
			"vpnv6": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
				Description: "Whether to support VPNv6 address family.",
			},
			"vrf": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
				Description: "VRF instance name for this BGP connection.",
			},
			"route_distinguisher": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
				Description: "Route distinguisher for VPN routes (format: ASN:NN or IP:NN).",
			},
			"routing_table": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
				Description: "Routing table name for this BGP connection.",
			},
			"comment": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
				Description: "Comment for this BGP connection.",
			},
			"templates": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
				Description: "BGP template names to apply (comma-separated).",
			},
		},
	}
}

// Create creates the resource and sets the initial Terraform state.
func (r *bgpConnection) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var terraformModel bgpConnectionModel
	var mikrotikModel client.BgpConnection
	GenericCreateResource(&terraformModel, &mikrotikModel, r.client)(ctx, req, resp)
}

// Read refreshes the Terraform state with the latest data.
func (r *bgpConnection) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var terraformModel bgpConnectionModel
	var mikrotikModel client.BgpConnection
	GenericReadResource(&terraformModel, &mikrotikModel, r.client)(ctx, req, resp)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *bgpConnection) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var terraformModel bgpConnectionModel
	var mikrotikModel client.BgpConnection
	GenericUpdateResource(&terraformModel, &mikrotikModel, r.client)(ctx, req, resp)
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *bgpConnection) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var terraformModel bgpConnectionModel
	var mikrotikModel client.BgpConnection
	GenericDeleteResource(&terraformModel, &mikrotikModel, r.client)(ctx, req, resp)
}

func (r *bgpConnection) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}

type bgpConnectionModel struct {
	Id                     tftypes.String `tfsdk:"id"`
	Name                   tftypes.String `tfsdk:"name"`
	AS                     tftypes.Int64  `tfsdk:"as"`
	Instance               tftypes.String `tfsdk:"instance"`
	Disabled               tftypes.Bool   `tfsdk:"disabled"`
	LocalRole              tftypes.String `tfsdk:"local_role"`
	LocalAddress           tftypes.String `tfsdk:"local_address"`
	RemoteAddress          tftypes.String `tfsdk:"remote_address"`
	RemoteAS               tftypes.Int64  `tfsdk:"remote_as"`
	RemotePort             tftypes.Int64  `tfsdk:"remote_port"`
	Listen                 tftypes.Bool   `tfsdk:"listen"`
	RouterID               tftypes.String `tfsdk:"router_id"`
	NexthopChoice          tftypes.String `tfsdk:"nexthop_choice"`
	HoldTime               tftypes.String `tfsdk:"hold_time"`
	KeepaliveTime          tftypes.String `tfsdk:"keepalive_time"`
	ConnectRetryTime       tftypes.String `tfsdk:"connect_retry_time"`
	TTL                    tftypes.String `tfsdk:"ttl"`
	Multihop               tftypes.Bool   `tfsdk:"multihop"`
	UseBFD                 tftypes.Bool   `tfsdk:"use_bfd"`
	AddressFamily          tftypes.String `tfsdk:"address_families"`
	InputFilter            tftypes.String `tfsdk:"input_filter"`
	InputAcceptNLRI        tftypes.String `tfsdk:"input_accept_nlri"`
	InputAcceptCommunities tftypes.String `tfsdk:"input_accept_communities"`
	OutputFilter           tftypes.String `tfsdk:"output_filter"`
	OutputDefaultOriginate tftypes.String `tfsdk:"output_default_originate"`
	OutputNetwork          tftypes.String `tfsdk:"output_network"`
	OutputRedistribute     tftypes.String `tfsdk:"output_redistribute"`
	TCPMd5Key              tftypes.String `tfsdk:"tcp_md5_key"`
	UseMPLS                tftypes.Bool   `tfsdk:"use_mpls"`
	VPNV4                  tftypes.Bool   `tfsdk:"vpnv4"`
	VPNV6                  tftypes.Bool   `tfsdk:"vpnv6"`
	VRF                    tftypes.String `tfsdk:"vrf"`
	RouteDistinguisher     tftypes.String `tfsdk:"route_distinguisher"`
	RoutingTable           tftypes.String `tfsdk:"routing_table"`
	Comment                tftypes.String `tfsdk:"comment"`
	Templates              tftypes.String `tfsdk:"templates"`
}
