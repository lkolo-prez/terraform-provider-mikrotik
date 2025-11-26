package mikrotik

import (
	"context"

	"github.com/ddelnano/terraform-provider-mikrotik/client"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	tftypes "github.com/hashicorp/terraform-plugin-framework/types"
)

type wifiAccessList struct {
	client *client.Mikrotik
}

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &wifiAccessList{}
	_ resource.ResourceWithConfigure   = &wifiAccessList{}
	_ resource.ResourceWithImportState = &wifiAccessList{}
)

// NewWiFiAccessListResource is a helper function to simplify the provider implementation.
func NewWiFiAccessListResource() resource.Resource {
	return &wifiAccessList{}
}

func (r *wifiAccessList) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*client.Mikrotik)
}

// Metadata returns the resource type name.
func (r *wifiAccessList) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_wifi_access_list"
}

// Schema defines the schema for the resource.
func (r *wifiAccessList) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Creates a WiFi access list entry for RouterOS v7 WiFi 6. MAC-based access control with VLAN assignment, signal filtering, and client isolation.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Description: "Identifier of this resource assigned by RouterOS",
			},
			"mac_address": schema.StringAttribute{
				Required:    true,
				Description: "MAC address of the client (e.g., AA:BB:CC:DD:EE:FF).",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"action": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("accept"),
				Description: "Action to take: accept, reject, query-radius. Default: accept",
			},
			"interface": schema.StringAttribute{
				Optional:    true,
				Description: "Apply rule only on specified WiFi interface.",
			},
			"ssid_regexp": schema.StringAttribute{
				Optional:    true,
				Description: "Apply rule only on SSIDs matching regex pattern.",
			},
			"vlan_id": schema.Int64Attribute{
				Optional:    true,
				Description: "Assign client to specific VLAN ID (1-4094).",
			},
			"vlan_mode": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("disabled"),
				Description: "VLAN mode: disabled, use-tag, use-service-tag. Default: disabled",
			},
			"signal_range": schema.StringAttribute{
				Optional:    true,
				Description: "Accept only if signal within range (e.g., '-90..-60' for -90 to -60 dBm).",
			},
			"time": schema.StringAttribute{
				Optional:    true,
				Description: "Time when rule is active (e.g., '8h-17h').",
			},
			"client_isolation": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
				Description: "Enable client isolation for this MAC. Default: false",
			},
			"radius_accounting": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
				Description: "Enable RADIUS accounting for this MAC. Default: false",
			},
			"disabled": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
				Description: "Whether the access list entry is disabled. Default: false",
			},
			"comment": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
				Description: "Comment for the access list entry.",
			},
		},
	}
}

// Create creates the resource and sets the initial Terraform state.
func (r *wifiAccessList) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var terraformModel wifiAccessListModel
	var mikrotikModel client.WiFiAccessList
	GenericCreateResource(&terraformModel, &mikrotikModel, r.client)(ctx, req, resp)
}

// Read refreshes the Terraform state with the latest data.
func (r *wifiAccessList) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var terraformModel wifiAccessListModel
	var mikrotikModel client.WiFiAccessList
	GenericReadResource(&terraformModel, &mikrotikModel, r.client)(ctx, req, resp)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *wifiAccessList) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var terraformModel wifiAccessListModel
	var mikrotikModel client.WiFiAccessList
	GenericUpdateResource(&terraformModel, &mikrotikModel, r.client)(ctx, req, resp)
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *wifiAccessList) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var terraformModel wifiAccessListModel
	var mikrotikModel client.WiFiAccessList
	GenericDeleteResource(&terraformModel, &mikrotikModel, r.client)(ctx, req, resp)
}

func (r *wifiAccessList) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("mac_address"), req, resp)
}

type wifiAccessListModel struct {
	Id               tftypes.String `tfsdk:"id"`
	MacAddress       tftypes.String `tfsdk:"mac_address"`
	Action           tftypes.String `tfsdk:"action"`
	Interface        tftypes.String `tfsdk:"interface"`
	SSIDRegexp       tftypes.String `tfsdk:"ssid_regexp"`
	VlanID           tftypes.Int64  `tfsdk:"vlan_id"`
	VlanMode         tftypes.String `tfsdk:"vlan_mode"`
	SignalRange      tftypes.String `tfsdk:"signal_range"`
	Time             tftypes.String `tfsdk:"time"`
	ClientIsolation  tftypes.Bool   `tfsdk:"client_isolation"`
	RadiusAccounting tftypes.Bool   `tfsdk:"radius_accounting"`
	Disabled         tftypes.Bool   `tfsdk:"disabled"`
	Comment          tftypes.String `tfsdk:"comment"`
}
