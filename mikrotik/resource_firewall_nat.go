package mikrotik

import (
	"context"

	"github.com/ddelnano/terraform-provider-mikrotik/client"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"

	tftypes "github.com/hashicorp/terraform-plugin-framework/types"
)

type firewallNat struct {
	client *client.Mikrotik
}

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &firewallNat{}
	_ resource.ResourceWithConfigure   = &firewallNat{}
	_ resource.ResourceWithImportState = &firewallNat{}
)

// NewFirewallNatResource is a helper function to simplify the provider implementation.
func NewFirewallNatResource() resource.Resource {
	return &firewallNat{}
}

func (r *firewallNat) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*client.Mikrotik)
}

func (r *firewallNat) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_firewall_nat"
}

func (r *firewallNat) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Creates and manages a firewall NAT rule for network address translation.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The unique identifier of the NAT rule.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"chain": schema.StringAttribute{
				Description: "NAT chain: 'srcnat' (source NAT) or 'dstnat' (destination NAT).",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.OneOf("srcnat", "dstnat"),
				},
			},
			"action": schema.StringAttribute{
				Description: "NAT action: 'masquerade', 'dst-nat', 'src-nat', 'netmap', 'redirect', 'accept', 'passthrough', 'jump', 'return'.",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.OneOf(
						"masquerade", "dst-nat", "src-nat", "netmap", "redirect",
						"accept", "passthrough", "jump", "return", "same",
					),
				},
			},
			"disabled": schema.BoolAttribute{
				Description: "Whether the NAT rule is disabled. Default: false.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"comment": schema.StringAttribute{
				Description: "Comment for the NAT rule.",
				Optional:    true,
			},

			// Matching criteria
			"src_address": schema.StringAttribute{
				Description: "Source IP address or network (CIDR). Example: '192.168.1.0/24'.",
				Optional:    true,
			},
			"dst_address": schema.StringAttribute{
				Description: "Destination IP address or network (CIDR). Example: '10.0.0.0/8'.",
				Optional:    true,
			},
			"src_address_list": schema.StringAttribute{
				Description: "Source address list name to match.",
				Optional:    true,
			},
			"dst_address_list": schema.StringAttribute{
				Description: "Destination address list name to match.",
				Optional:    true,
			},
			"protocol": schema.StringAttribute{
				Description: "IP protocol: 'tcp', 'udp', 'icmp', etc. or protocol number.",
				Optional:    true,
			},
			"src_port": schema.StringAttribute{
				Description: "Source port(s). Example: '80' or '80-90' or '80,443'.",
				Optional:    true,
			},
			"dst_port": schema.StringAttribute{
				Description: "Destination port(s). Example: '80' or '80-90' or '80,443'.",
				Optional:    true,
			},
			"in_interface": schema.StringAttribute{
				Description: "Input interface name to match.",
				Optional:    true,
			},
			"out_interface": schema.StringAttribute{
				Description: "Output interface name to match.",
				Optional:    true,
			},
			"in_interface_list": schema.StringAttribute{
				Description: "Input interface list name to match.",
				Optional:    true,
			},
			"out_interface_list": schema.StringAttribute{
				Description: "Output interface list name to match.",
				Optional:    true,
			},

			// Connection tracking
			"connection_state": schema.StringAttribute{
				Description: "Connection state to match. Example: 'new,established,related'.",
				Optional:    true,
			},
			"connection_nat_state": schema.StringAttribute{
				Description: "Connection NAT state to match: 'srcnat', 'dstnat'.",
				Optional:    true,
			},
			"connection_mark": schema.StringAttribute{
				Description: "Connection mark to match.",
				Optional:    true,
			},
			"packet_mark": schema.StringAttribute{
				Description: "Packet mark to match.",
				Optional:    true,
			},
			"routing_mark": schema.StringAttribute{
				Description: "Routing mark to match.",
				Optional:    true,
			},

			// NAT action parameters
			"to_addresses": schema.StringAttribute{
				Description: "NAT to IP address(es). Required for 'dst-nat' and 'src-nat'. Example: '192.168.1.100' or '192.168.1.1-192.168.1.10'.",
				Optional:    true,
			},
			"to_ports": schema.StringAttribute{
				Description: "NAT to port(s). Example: '8080' or '8000-9000'.",
				Optional:    true,
			},

			// Logging
			"log": schema.BoolAttribute{
				Description: "Enable logging for this rule. Default: false.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"log_prefix": schema.StringAttribute{
				Description: "Log message prefix.",
				Optional:    true,
			},

			// Advanced matching
			"icmp_options": schema.StringAttribute{
				Description: "ICMP options to match. Example: '0:0' for echo-reply.",
				Optional:    true,
			},
			"limit": schema.StringAttribute{
				Description: "Rate limit for matching packets. Example: '5/1m,10:packet'.",
				Optional:    true,
			},
			"time": schema.StringAttribute{
				Description: "Time range to match. Example: '08:00-17:00,mon,tue,wed,thu,fri'.",
				Optional:    true,
			},
			"random": schema.StringAttribute{
				Description: "Random match probability (1-99).",
				Optional:    true,
			},
			"hotspot": schema.StringAttribute{
				Description: "HotSpot authentication status: 'none', 'auth', 'http', 'https'.",
				Optional:    true,
			},
			"content_type": schema.StringAttribute{
				Description: "HTTP content type to match.",
				Optional:    true,
			},
			"layer7_protocol": schema.StringAttribute{
				Description: "Layer 7 protocol name to match.",
				Optional:    true,
			},
			"psd": schema.StringAttribute{
				Description: "Port scan detection parameters.",
				Optional:    true,
			},
			"tcp_flags": schema.StringAttribute{
				Description: "TCP flags to match. Example: 'syn,!ack,!rst'.",
				Optional:    true,
			},
			"tcp_mss": schema.StringAttribute{
				Description: "TCP MSS value to match. Example: '500-1500'.",
				Optional:    true,
			},
			"dst_limit": schema.StringAttribute{
				Description: "Destination address limit. Example: '1/1h,5,dst-address/1d'.",
				Optional:    true,
			},
			"address_list": schema.StringAttribute{
				Description: "Add matching source address to specified address list.",
				Optional:    true,
			},
			"address_list_timeout": schema.StringAttribute{
				Description: "Timeout for address list entry. Example: '1d'.",
				Optional:    true,
			},
			"packet_size": schema.StringAttribute{
				Description: "Packet size to match. Example: '64-128'.",
				Optional:    true,
			},
			"src_address_type": schema.StringAttribute{
				Description: "Source address type to match.",
				Optional:    true,
			},
			"dst_address_type": schema.StringAttribute{
				Description: "Destination address type to match.",
				Optional:    true,
			},

			// Computed fields
			"bytes": schema.Int64Attribute{
				Description: "Number of bytes matched by this rule (computed).",
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"packets": schema.Int64Attribute{
				Description: "Number of packets matched by this rule (computed).",
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"dynamic": schema.BoolAttribute{
				Description: "Whether the rule is dynamic (computed).",
				Computed:    true,
			},
			"invalid": schema.BoolAttribute{
				Description: "Whether the rule is invalid (computed).",
				Computed:    true,
			},
		},
	}
}

// firewallNatModel describes the resource data model.
type firewallNatModel struct {
	Id       tftypes.String `tfsdk:"id"`
	Chain    tftypes.String `tfsdk:"chain"`
	Action   tftypes.String `tfsdk:"action"`
	Disabled tftypes.Bool   `tfsdk:"disabled"`
	Comment  tftypes.String `tfsdk:"comment"`

	// Matching criteria
	SrcAddress       tftypes.String `tfsdk:"src_address"`
	DstAddress       tftypes.String `tfsdk:"dst_address"`
	SrcAddressList   tftypes.String `tfsdk:"src_address_list"`
	DstAddressList   tftypes.String `tfsdk:"dst_address_list"`
	Protocol         tftypes.String `tfsdk:"protocol"`
	SrcPort          tftypes.String `tfsdk:"src_port"`
	DstPort          tftypes.String `tfsdk:"dst_port"`
	InInterface      tftypes.String `tfsdk:"in_interface"`
	OutInterface     tftypes.String `tfsdk:"out_interface"`
	InInterfaceList  tftypes.String `tfsdk:"in_interface_list"`
	OutInterfaceList tftypes.String `tfsdk:"out_interface_list"`

	// Connection tracking
	ConnectionState    tftypes.String `tfsdk:"connection_state"`
	ConnectionNatState tftypes.String `tfsdk:"connection_nat_state"`
	ConnectionMark     tftypes.String `tfsdk:"connection_mark"`
	PacketMark         tftypes.String `tfsdk:"packet_mark"`
	RoutingMark        tftypes.String `tfsdk:"routing_mark"`

	// NAT action parameters
	ToAddresses tftypes.String `tfsdk:"to_addresses"`
	ToPorts     tftypes.String `tfsdk:"to_ports"`

	// Logging
	Log       tftypes.Bool   `tfsdk:"log"`
	LogPrefix tftypes.String `tfsdk:"log_prefix"`

	// Advanced matching
	IcmpOptions        tftypes.String `tfsdk:"icmp_options"`
	Limit              tftypes.String `tfsdk:"limit"`
	Time               tftypes.String `tfsdk:"time"`
	Random             tftypes.String `tfsdk:"random"`
	HotspotAuth        tftypes.String `tfsdk:"hotspot"`
	ContentType        tftypes.String `tfsdk:"content_type"`
	Layer7Protocol     tftypes.String `tfsdk:"layer7_protocol"`
	Psd                tftypes.String `tfsdk:"psd"`
	TcpFlags           tftypes.String `tfsdk:"tcp_flags"`
	TcpMss             tftypes.String `tfsdk:"tcp_mss"`
	DstLimit           tftypes.String `tfsdk:"dst_limit"`
	AddressList        tftypes.String `tfsdk:"address_list"`
	AddressListTimeout tftypes.String `tfsdk:"address_list_timeout"`
	PacketSize         tftypes.String `tfsdk:"packet_size"`
	SrcAddressType     tftypes.String `tfsdk:"src_address_type"`
	DstAddressType     tftypes.String `tfsdk:"dst_address_type"`

	// Computed fields
	Bytes   tftypes.Int64 `tfsdk:"bytes"`
	Packets tftypes.Int64 `tfsdk:"packets"`
	Dynamic tftypes.Bool  `tfsdk:"dynamic"`
	Invalid tftypes.Bool  `tfsdk:"invalid"`
}

func (r *firewallNat) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var terraformModel firewallNatModel
	var mikrotikModel client.FirewallNat
	GenericCreateResource(&terraformModel, &mikrotikModel, r.client)(ctx, req, resp)
}

func (r *firewallNat) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var terraformModel firewallNatModel
	var mikrotikModel client.FirewallNat
	GenericReadResource(&terraformModel, &mikrotikModel, r.client)(ctx, req, resp)
}

func (r *firewallNat) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var terraformModel firewallNatModel
	var mikrotikModel client.FirewallNat
	GenericUpdateResource(&terraformModel, &mikrotikModel, r.client)(ctx, req, resp)
}

func (r *firewallNat) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var terraformModel firewallNatModel
	var mikrotikModel client.FirewallNat
	GenericDeleteResource(&terraformModel, &mikrotikModel, r.client)(ctx, req, resp)
}

func (r *firewallNat) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
