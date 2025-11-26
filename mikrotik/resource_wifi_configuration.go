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

type wifiConfiguration struct {
	client *client.Mikrotik
}

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &wifiConfiguration{}
	_ resource.ResourceWithConfigure   = &wifiConfiguration{}
	_ resource.ResourceWithImportState = &wifiConfiguration{}
)

// NewWiFiConfigurationResource is a helper function to simplify the provider implementation.
func NewWiFiConfigurationResource() resource.Resource {
	return &wifiConfiguration{}
}

func (r *wifiConfiguration) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*client.Mikrotik)
}

// Metadata returns the resource type name.
func (r *wifiConfiguration) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_wifi_configuration"
}

// Schema defines the schema for the resource.
func (r *wifiConfiguration) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Creates a WiFi configuration profile for RouterOS v7 WiFi 6. Defines SSID, mode (AP/station), security, country settings, and WiFi 6 parameters.",
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
				Description: "Name of the configuration profile.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"ssid": schema.StringAttribute{
				Required:    true,
				Description: "SSID (network name) broadcast by this configuration.",
			},
			"mode": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("ap"),
				Description: "Operating mode: ap (access point), station (client), sniffer. Default: ap",
			},
			"hide_ssid": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
				Description: "Hide SSID in beacon frames. Default: false",
			},
			"security": schema.StringAttribute{
				Optional:    true,
				Description: "Reference to WiFi security profile. Use mikrotik_wifi_security resource.",
			},
			"country": schema.StringAttribute{
				Optional:    true,
				Description: "Country code for regulatory domain (e.g., US, GB, PL). Affects available channels and TX power.",
			},
			"installation": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("indoor"),
				Description: "Installation type: indoor, outdoor, any. Affects TX power limits. Default: indoor",
			},
			"tx_power": schema.Int64Attribute{
				Optional:    true,
				Description: "TX power in dBm. Leave empty for regulatory maximum.",
			},
			"tx_power_mode": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("default"),
				Description: "TX power mode: default, all-rates-fixed, manual-table. Default: default",
			},
			"he_guard_interval": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("long"),
				Description: "802.11ax Guard Interval: long (3.2µs), short (0.8µs). Default: long",
			},
			"he_frame_format": schema.StringAttribute{
				Optional:    true,
				Description: "802.11ax frame format: he-su, he-er-su, he-mu, he-tb (comma-separated for multiple).",
			},
			"distance": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("indoors"),
				Description: "Distance mode: indoors, dynamic. Default: indoors",
			},
			"max_station_count": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Default:     int64default.StaticInt64(2007),
				Description: "Maximum number of associated stations. Default: 2007",
			},
			"wps": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
				Description: "Enable WPS (WiFi Protected Setup). Default: false",
			},
			"wps_mode": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("disabled"),
				Description: "WPS mode: disabled, push-button, pin. Default: disabled",
			},
			"disabled": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
				Description: "Whether the configuration profile is disabled. Default: false",
			},
			"comment": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
				Description: "Comment for the configuration profile.",
			},
		},
	}
}

// Create creates the resource and sets the initial Terraform state.
func (r *wifiConfiguration) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var terraformModel wifiConfigurationModel
	var mikrotikModel client.WiFiConfiguration
	GenericCreateResource(&terraformModel, &mikrotikModel, r.client)(ctx, req, resp)
}

// Read refreshes the Terraform state with the latest data.
func (r *wifiConfiguration) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var terraformModel wifiConfigurationModel
	var mikrotikModel client.WiFiConfiguration
	GenericReadResource(&terraformModel, &mikrotikModel, r.client)(ctx, req, resp)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *wifiConfiguration) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var terraformModel wifiConfigurationModel
	var mikrotikModel client.WiFiConfiguration
	GenericUpdateResource(&terraformModel, &mikrotikModel, r.client)(ctx, req, resp)
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *wifiConfiguration) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var terraformModel wifiConfigurationModel
	var mikrotikModel client.WiFiConfiguration
	GenericDeleteResource(&terraformModel, &mikrotikModel, r.client)(ctx, req, resp)
}

func (r *wifiConfiguration) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}

type wifiConfigurationModel struct {
	Id               tftypes.String `tfsdk:"id"`
	Name             tftypes.String `tfsdk:"name"`
	SSID             tftypes.String `tfsdk:"ssid"`
	Mode             tftypes.String `tfsdk:"mode"`
	HideSSID         tftypes.Bool   `tfsdk:"hide_ssid"`
	Security         tftypes.String `tfsdk:"security"`
	Country          tftypes.String `tfsdk:"country"`
	Installation     tftypes.String `tfsdk:"installation"`
	TxPower          tftypes.Int64  `tfsdk:"tx_power"`
	TxPowerMode      tftypes.String `tfsdk:"tx_power_mode"`
	HEGuardInterval  tftypes.String `tfsdk:"he_guard_interval"`
	HEFrameFormat    tftypes.String `tfsdk:"he_frame_format"`
	Distance         tftypes.String `tfsdk:"distance"`
	MaxStationCount  tftypes.Int64  `tfsdk:"max_station_count"`
	WPS              tftypes.Bool   `tfsdk:"wps"`
	WPSMode          tftypes.String `tfsdk:"wps_mode"`
	Disabled         tftypes.Bool   `tfsdk:"disabled"`
	Comment          tftypes.String `tfsdk:"comment"`
}
