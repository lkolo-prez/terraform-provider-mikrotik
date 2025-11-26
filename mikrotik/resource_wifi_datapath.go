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

type wifiDatapath struct {
	client *client.Mikrotik
}

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &wifiDatapath{}
	_ resource.ResourceWithConfigure   = &wifiDatapath{}
	_ resource.ResourceWithImportState = &wifiDatapath{}
)

// NewWiFiDatapathResource is a helper function to simplify the provider implementation.
func NewWiFiDatapathResource() resource.Resource {
	return &wifiDatapath{}
}

func (r *wifiDatapath) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*client.Mikrotik)
}

// Metadata returns the resource type name.
func (r *wifiDatapath) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_wifi_datapath"
}

// Schema defines the schema for the resource.
func (r *wifiDatapath) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Creates a WiFi datapath profile for RouterOS v7 WiFi 6. Configures bridge integration, VLAN tagging, client isolation, and Layer 2 parameters.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Description: "Identifier of this resource assigned by RouterOS",
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "Name of the datapath profile.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"bridge": schema.StringAttribute{
				Optional:    true,
				Description: "Bridge interface to add WiFi clients to (e.g., bridge1).",
			},
			"bridge_horizon": schema.Int64Attribute{
				Optional:    true,
				Description: "Bridge horizon for split-horizon bridging (0-429496729).",
			},
			"bridge_cost": schema.Int64Attribute{
				Optional:    true,
				Description: "Bridge path cost for STP/RSTP (0-4294967295).",
			},
			"vlan_id": schema.Int64Attribute{
				Optional:    true,
				Description: "VLAN ID to tag WiFi traffic (1-4094).",
			},
			"vlan_mode": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("disabled"),
				Description: "VLAN mode: disabled, use-tag, use-service-tag. Default: disabled",
			},
			"client_isolation": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
				Description: "Prevent clients from communicating with each other (AP isolation). Default: false",
			},
			"client_to_client_forwarding": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(true),
				Description: "Allow client-to-client forwarding within same AP. Default: true",
			},
			"arp": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("enabled"),
				Description: "ARP mode: disabled, enabled, proxy-arp, reply-only. Default: enabled",
			},
			"arp_timeout": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("auto"),
				Description: "ARP timeout. Default: auto",
			},
			"interface_list": schema.StringAttribute{
				Optional:    true,
				Description: "Add interface to specified interface list.",
			},
			"l2mtu": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Default:     int64default.StaticInt64(1598),
				Description: "Layer 2 MTU. Default: 1598",
			},
			"mtu": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Default:     int64default.StaticInt64(1500),
				Description: "Maximum Transmission Unit. Default: 1500",
			},
			"disabled": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
				Description: "Whether the datapath profile is disabled. Default: false",
			},
			"comment": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
				Description: "Comment for the datapath profile.",
			},
		},
	}
}

// Create creates the resource and sets the initial Terraform state.
func (r *wifiDatapath) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var terraformModel wifiDatapathModel
	var mikrotikModel client.WiFiDatapath
	GenericCreateResource(&terraformModel, &mikrotikModel, r.client)(ctx, req, resp)
}

// Read refreshes the Terraform state with the latest data.
func (r *wifiDatapath) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var terraformModel wifiDatapathModel
	var mikrotikModel client.WiFiDatapath
	GenericReadResource(&terraformModel, &mikrotikModel, r.client)(ctx, req, resp)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *wifiDatapath) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var terraformModel wifiDatapathModel
	var mikrotikModel client.WiFiDatapath
	GenericUpdateResource(&terraformModel, &mikrotikModel, r.client)(ctx, req, resp)
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *wifiDatapath) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var terraformModel wifiDatapathModel
	var mikrotikModel client.WiFiDatapath
	GenericDeleteResource(&terraformModel, &mikrotikModel, r.client)(ctx, req, resp)
}

func (r *wifiDatapath) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}

type wifiDatapathModel struct {
	Id                        tftypes.String `tfsdk:"id"`
	Name                      tftypes.String `tfsdk:"name"`
	Bridge                    tftypes.String `tfsdk:"bridge"`
	BridgeHorizon             tftypes.Int64  `tfsdk:"bridge_horizon"`
	BridgeCost                tftypes.Int64  `tfsdk:"bridge_cost"`
	VlanID                    tftypes.Int64  `tfsdk:"vlan_id"`
	VlanMode                  tftypes.String `tfsdk:"vlan_mode"`
	ClientIsolation           tftypes.Bool   `tfsdk:"client_isolation"`
	ClientToClientForwarding  tftypes.Bool   `tfsdk:"client_to_client_forwarding"`
	ARP                       tftypes.String `tfsdk:"arp"`
	ARPTimeout                tftypes.String `tfsdk:"arp_timeout"`
	InterfaceList             tftypes.String `tfsdk:"interface_list"`
	L2MTU                     tftypes.Int64  `tfsdk:"l2mtu"`
	MTU                       tftypes.Int64  `tfsdk:"mtu"`
	Disabled                  tftypes.Bool   `tfsdk:"disabled"`
	Comment                   tftypes.String `tfsdk:"comment"`
}
