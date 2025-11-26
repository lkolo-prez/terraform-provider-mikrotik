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

type interfaceWiFi struct {
	client *client.Mikrotik
}

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &interfaceWiFi{}
	_ resource.ResourceWithConfigure   = &interfaceWiFi{}
	_ resource.ResourceWithImportState = &interfaceWiFi{}
)

// NewInterfaceWiFiResource is a helper function to simplify the provider implementation.
func NewInterfaceWiFiResource() resource.Resource {
	return &interfaceWiFi{}
}

func (r *interfaceWiFi) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*client.Mikrotik)
}

// Metadata returns the resource type name.
func (r *interfaceWiFi) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_interface_wifi"
}

// Schema defines the schema for the resource.
func (r *interfaceWiFi) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Creates a WiFi 6 (802.11ax) interface for RouterOS v7. This is the new WiFi stack replacing legacy /interface/wireless.",
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
				Description: "Name of the WiFi interface.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"configuration": schema.StringAttribute{
				Optional:    true,
				Description: "Reference to WiFi configuration profile (SSID, mode, etc). Use mikrotik_wifi_configuration resource.",
			},
			"datapath": schema.StringAttribute{
				Optional:    true,
				Description: "Reference to WiFi datapath profile (bridge, VLAN). Use mikrotik_wifi_datapath resource.",
			},
			"channel": schema.StringAttribute{
				Optional:    true,
				Description: "Reference to WiFi channel profile (band, frequency). Use mikrotik_wifi_channel resource.",
			},
			"master_interface": schema.StringAttribute{
				Optional:    true,
				Description: "Master interface for virtual APs. Leave empty for physical interface, set to physical interface name for guest networks.",
			},
			"mac_address": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "MAC address of the interface. Automatically assigned if not specified.",
			},
			"mtu": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Default:     int64default.StaticInt64(1500),
				Description: "Maximum Transmission Unit. Default is 1500.",
			},
			"arp": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("enabled"),
				Description: "ARP mode: disabled, enabled, proxy-arp, reply-only. Default is enabled.",
			},
			"arp_timeout": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("auto"),
				Description: "ARP timeout. Default is auto.",
			},
			"disabled": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
				Description: "Whether the interface is disabled. Default is false.",
			},
			"comment": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
				Description: "Comment for the WiFi interface.",
			},
			// Computed runtime fields
			"running": schema.BoolAttribute{
				Computed:    true,
				Description: "Whether the interface is currently running (read-only).",
			},
			"radio": schema.StringAttribute{
				Computed:    true,
				Description: "Physical radio this interface is using (read-only).",
			},
		},
	}
}

// Create creates the resource and sets the initial Terraform state.
func (r *interfaceWiFi) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var terraformModel interfaceWiFiModel
	var mikrotikModel client.InterfaceWiFi
	GenericCreateResource(&terraformModel, &mikrotikModel, r.client)(ctx, req, resp)
}

// Read refreshes the Terraform state with the latest data.
func (r *interfaceWiFi) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var terraformModel interfaceWiFiModel
	var mikrotikModel client.InterfaceWiFi
	GenericReadResource(&terraformModel, &mikrotikModel, r.client)(ctx, req, resp)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *interfaceWiFi) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var terraformModel interfaceWiFiModel
	var mikrotikModel client.InterfaceWiFi
	GenericUpdateResource(&terraformModel, &mikrotikModel, r.client)(ctx, req, resp)
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *interfaceWiFi) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var terraformModel interfaceWiFiModel
	var mikrotikModel client.InterfaceWiFi
	GenericDeleteResource(&terraformModel, &mikrotikModel, r.client)(ctx, req, resp)
}

func (r *interfaceWiFi) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}

type interfaceWiFiModel struct {
	Id              tftypes.String `tfsdk:"id"`
	Name            tftypes.String `tfsdk:"name"`
	Configuration   tftypes.String `tfsdk:"configuration"`
	Datapath        tftypes.String `tfsdk:"datapath"`
	Channel         tftypes.String `tfsdk:"channel"`
	MasterInterface tftypes.String `tfsdk:"master_interface"`
	MacAddress      tftypes.String `tfsdk:"mac_address"`
	Mtu             tftypes.Int64  `tfsdk:"mtu"`
	Arp             tftypes.String `tfsdk:"arp"`
	ArpTimeout      tftypes.String `tfsdk:"arp_timeout"`
	Disabled        tftypes.Bool   `tfsdk:"disabled"`
	Comment         tftypes.String `tfsdk:"comment"`
	// Computed
	Running tftypes.Bool   `tfsdk:"running"`
	Radio   tftypes.String `tfsdk:"radio"`
}
