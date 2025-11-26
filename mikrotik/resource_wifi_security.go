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

type wifiSecurity struct {
	client *client.Mikrotik
}

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &wifiSecurity{}
	_ resource.ResourceWithConfigure   = &wifiSecurity{}
	_ resource.ResourceWithImportState = &wifiSecurity{}
)

// NewWiFiSecurityResource is a helper function to simplify the provider implementation.
func NewWiFiSecurityResource() resource.Resource {
	return &wifiSecurity{}
}

func (r *wifiSecurity) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*client.Mikrotik)
}

// Metadata returns the resource type name.
func (r *wifiSecurity) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_wifi_security"
}

// Schema defines the schema for the resource.
func (r *wifiSecurity) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Creates a WiFi security profile for RouterOS v7 WiFi 6. Supports WPA2, WPA3, enterprise (EAP/RADIUS), and advanced features like PMF and Fast Transition.",
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
				Description: "Name of the security profile.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			// Authentication
			"authentication_types": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("wpa2-psk,wpa3-psk"),
				Description: "Authentication types: wpa2-psk, wpa3-psk, wpa2-eap, wpa3-eap (comma-separated). Default: wpa2-psk,wpa3-psk",
			},
			"encryption": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("ccmp"),
				Description: "Encryption ciphers: ccmp (AES), gcmp, ccmp-256, gcmp-256 (comma-separated). Default: ccmp",
			},
			"passphrase": schema.StringAttribute{
				Optional:    true,
				Sensitive:   true,
				Description: "WPA passphrase (8-63 characters). Required for PSK authentication.",
			},
			
			// Protected Management Frames (PMF) - required for WPA3
			"management_protection": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("disabled"),
				Description: "Protected Management Frames (PMF/802.11w): disabled, optional, required. Required for WPA3.",
			},
			"management_protection_key": schema.StringAttribute{
				Optional:    true,
				Sensitive:   true,
				Description: "PMF key for management frame encryption.",
			},
			
			// Enterprise (EAP/RADIUS)
			"eap_methods": schema.StringAttribute{
				Optional:    true,
				Description: "EAP methods for enterprise: eap-tls, eap-ttls, peap (comma-separated). For wpa2-eap/wpa3-eap.",
			},
			"eap_radius_server": schema.StringAttribute{
				Optional:    true,
				Description: "RADIUS server IP address for EAP authentication.",
			},
			"eap_radius_secret": schema.StringAttribute{
				Optional:    true,
				Sensitive:   true,
				Description: "RADIUS shared secret.",
			},
			"eap_radius_port": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Default:     int64default.StaticInt64(1812),
				Description: "RADIUS server port. Default: 1812",
			},
			"eap_radius_accounting": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
				Description: "Enable RADIUS accounting. Default: false",
			},
			"eap_radius_accounting_port": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Default:     int64default.StaticInt64(1813),
				Description: "RADIUS accounting port. Default: 1813",
			},
			"eap_tls_certificate": schema.StringAttribute{
				Optional:    true,
				Description: "TLS certificate name for EAP-TLS authentication.",
			},
			
			// WPA3 SAE (Simultaneous Authentication of Equals)
			"sae_pwe": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("both"),
				Description: "SAE Password Element method for WPA3: hunting-and-pecking, hash-to-element, both. Default: both",
			},
			"sae_groups": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("19,20,21"),
				Description: "SAE elliptic curve groups: 19, 20, 21 (comma-separated). Default: 19,20,21",
			},
			
			// Fast Transition (802.11r)
			"ft": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
				Description: "Enable Fast Transition (802.11r) for fast roaming. Default: false",
			},
			"ft_over_ds": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
				Description: "Fast Transition over Distribution System. Default: false",
			},
			"ft_preserve_vlan": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
				Description: "Preserve VLAN assignment during FT roaming. Default: false",
			},
			
			// Group rekey
			"group_rekey": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("1h"),
				Description: "Group key rekey interval (e.g., 1h, 30m). Default: 1h",
			},
			
			"disabled": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
				Description: "Whether the security profile is disabled. Default: false",
			},
			"comment": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
				Description: "Comment for the security profile.",
			},
		},
	}
}

// Create creates the resource and sets the initial Terraform state.
func (r *wifiSecurity) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var terraformModel wifiSecurityModel
	var mikrotikModel client.WiFiSecurity
	GenericCreateResource(&terraformModel, &mikrotikModel, r.client)(ctx, req, resp)
}

// Read refreshes the Terraform state with the latest data.
func (r *wifiSecurity) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var terraformModel wifiSecurityModel
	var mikrotikModel client.WiFiSecurity
	GenericReadResource(&terraformModel, &mikrotikModel, r.client)(ctx, req, resp)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *wifiSecurity) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var terraformModel wifiSecurityModel
	var mikrotikModel client.WiFiSecurity
	GenericUpdateResource(&terraformModel, &mikrotikModel, r.client)(ctx, req, resp)
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *wifiSecurity) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var terraformModel wifiSecurityModel
	var mikrotikModel client.WiFiSecurity
	GenericDeleteResource(&terraformModel, &mikrotikModel, r.client)(ctx, req, resp)
}

func (r *wifiSecurity) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}

type wifiSecurityModel struct {
	Id                         tftypes.String `tfsdk:"id"`
	Name                       tftypes.String `tfsdk:"name"`
	AuthenticationTypes        tftypes.String `tfsdk:"authentication_types"`
	Encryption                 tftypes.String `tfsdk:"encryption"`
	Passphrase                 tftypes.String `tfsdk:"passphrase"`
	ManagementProtection       tftypes.String `tfsdk:"management_protection"`
	ManagementProtectionKey    tftypes.String `tfsdk:"management_protection_key"`
	EapMethods                 tftypes.String `tfsdk:"eap_methods"`
	EapRadiusServer            tftypes.String `tfsdk:"eap_radius_server"`
	EapRadiusSecret            tftypes.String `tfsdk:"eap_radius_secret"`
	EapRadiusPort              tftypes.Int64  `tfsdk:"eap_radius_port"`
	EapRadiusAccounting        tftypes.Bool   `tfsdk:"eap_radius_accounting"`
	EapRadiusAccountingPort    tftypes.Int64  `tfsdk:"eap_radius_accounting_port"`
	EapTlsCertificate          tftypes.String `tfsdk:"eap_tls_certificate"`
	SAE_PWE                    tftypes.String `tfsdk:"sae_pwe"`
	SAE_Groups                 tftypes.String `tfsdk:"sae_groups"`
	FT                         tftypes.Bool   `tfsdk:"ft"`
	FTOverDS                   tftypes.Bool   `tfsdk:"ft_over_ds"`
	FTPreserveVlan             tftypes.Bool   `tfsdk:"ft_preserve_vlan"`
	GroupRekey                 tftypes.String `tfsdk:"group_rekey"`
	Disabled                   tftypes.Bool   `tfsdk:"disabled"`
	Comment                    tftypes.String `tfsdk:"comment"`
}
