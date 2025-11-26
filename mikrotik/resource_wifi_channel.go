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

type wifiChannel struct {
	client *client.Mikrotik
}

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &wifiChannel{}
	_ resource.ResourceWithConfigure   = &wifiChannel{}
	_ resource.ResourceWithImportState = &wifiChannel{}
)

// NewWiFiChannelResource is a helper function to simplify the provider implementation.
func NewWiFiChannelResource() resource.Resource {
	return &wifiChannel{}
}

func (r *wifiChannel) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*client.Mikrotik)
}

// Metadata returns the resource type name.
func (r *wifiChannel) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_wifi_channel"
}

// Schema defines the schema for the resource.
func (r *wifiChannel) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Creates a WiFi channel profile for RouterOS v7 WiFi 6. Configures band, frequency, channel width, and DFS settings.",
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
				Description: "Name of the channel profile.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"band": schema.StringAttribute{
				Optional:    true,
				Description: "WiFi band: 2ghz-ax (WiFi 6 2.4GHz), 5ghz-ax (WiFi 6 5GHz), 5ghz-ac (WiFi 5), 6ghz-ax (WiFi 6E).",
			},
			"frequency": schema.StringAttribute{
				Optional:    true,
				Description: "Center frequency in MHz (e.g., 2412, 5180). Use 'auto' for automatic selection.",
			},
			"width": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("20mhz"),
				Description: "Channel width: 20mhz, 40mhz, 80mhz, 160mhz, 20/40mhz, 20/40/80mhz. Default: 20mhz",
			},
			"secondary_frequency": schema.StringAttribute{
				Optional:    true,
				Description: "Secondary frequency for 40MHz+ channels (e.g., above, below, or specific frequency).",
			},
			"skip_dfs_channels": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
				Description: "Skip DFS (radar) channels when scanning. Default: false",
			},
			"reuse_dfs_channels": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
				Description: "Reuse DFS channels without waiting for CAC. Default: false",
			},
			"control_channel_position": schema.StringAttribute{
				Optional:    true,
				Description: "Control channel position for 40MHz+: lower, upper. Only for 40/80/160 MHz.",
			},
			"disabled": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
				Description: "Whether the channel profile is disabled. Default: false",
			},
			"comment": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
				Description: "Comment for the channel profile.",
			},
		},
	}
}

// Create creates the resource and sets the initial Terraform state.
func (r *wifiChannel) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var terraformModel wifiChannelModel
	var mikrotikModel client.WiFiChannel
	GenericCreateResource(&terraformModel, &mikrotikModel, r.client)(ctx, req, resp)
}

// Read refreshes the Terraform state with the latest data.
func (r *wifiChannel) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var terraformModel wifiChannelModel
	var mikrotikModel client.WiFiChannel
	GenericReadResource(&terraformModel, &mikrotikModel, r.client)(ctx, req, resp)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *wifiChannel) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var terraformModel wifiChannelModel
	var mikrotikModel client.WiFiChannel
	GenericUpdateResource(&terraformModel, &mikrotikModel, r.client)(ctx, req, resp)
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *wifiChannel) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var terraformModel wifiChannelModel
	var mikrotikModel client.WiFiChannel
	GenericDeleteResource(&terraformModel, &mikrotikModel, r.client)(ctx, req, resp)
}

func (r *wifiChannel) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}

type wifiChannelModel struct {
	Id                     tftypes.String `tfsdk:"id"`
	Name                   tftypes.String `tfsdk:"name"`
	Band                   tftypes.String `tfsdk:"band"`
	Frequency              tftypes.String `tfsdk:"frequency"`
	Width                  tftypes.String `tfsdk:"width"`
	SecondaryFrequency     tftypes.String `tfsdk:"secondary_frequency"`
	SkipDFSChannels        tftypes.Bool   `tfsdk:"skip_dfs_channels"`
	ReuseDFSChannels       tftypes.Bool   `tfsdk:"reuse_dfs_channels"`
	ControlChannelPosition tftypes.String `tfsdk:"control_channel_position"`
	Disabled               tftypes.Bool   `tfsdk:"disabled"`
	Comment                tftypes.String `tfsdk:"comment"`
}
