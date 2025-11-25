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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type ospfInterfaceTemplateV7 struct {
	Area               types.String `tfsdk:"area"`
	Networks           types.List   `tfsdk:"networks"`
	Interfaces         types.List   `tfsdk:"interfaces"`
	Type               types.String `tfsdk:"type"`
	Disabled           types.Bool   `tfsdk:"disabled"`
	Comment            types.String `tfsdk:"comment"`
	Cost               types.Int64  `tfsdk:"cost"`
	Priority           types.Int64  `tfsdk:"priority"`
	Passive            types.Bool   `tfsdk:"passive"`
	Auth               types.String `tfsdk:"auth"`
	AuthKey            types.String `tfsdk:"auth_key"`
	AuthId             types.Int64  `tfsdk:"auth_id"`
	HelloInterval      types.String `tfsdk:"hello_interval"`
	DeadInterval       types.String `tfsdk:"dead_interval"`
	RetransmitInterval types.String `tfsdk:"retransmit_interval"`
	TransmitDelay      types.String `tfsdk:"transmit_delay"`
	WaitTime           types.String `tfsdk:"wait_time"`
	VlinkTransitArea   types.String `tfsdk:"vlink_transit_area"`
	VlinkNeighborId    types.String `tfsdk:"vlink_neighbor_id"`
	// Computed
	Id      types.String `tfsdk:"id"`
	Dynamic types.Bool   `tfsdk:"dynamic"`
	Invalid types.Bool   `tfsdk:"invalid"`
}

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &ospfInterfaceTemplateV7Resource{}
	_ resource.ResourceWithConfigure   = &ospfInterfaceTemplateV7Resource{}
	_ resource.ResourceWithImportState = &ospfInterfaceTemplateV7Resource{}
)

// NewOspfInterfaceTemplateV7Resource is a helper function to simplify the provider implementation.
func NewOspfInterfaceTemplateV7Resource() resource.Resource {
	return &ospfInterfaceTemplateV7Resource{}
}

type ospfInterfaceTemplateV7Resource struct {
	client *client.Mikrotik
}

func (r *ospfInterfaceTemplateV7Resource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*client.Mikrotik)
}

func (r *ospfInterfaceTemplateV7Resource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_ospf_interface_template_v7"
}

// Schema defines the schema for the resource.
func (r *ospfInterfaceTemplateV7Resource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages OSPF interface template configuration (RouterOS v7). Templates define how OSPF operates on matched interfaces/networks.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "Unique ID of the OSPF interface template.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"area": schema.StringAttribute{
				Required:    true,
				Description: "Name of the OSPF area this template belongs to.",
			},
			"networks": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Description: "List of network prefixes (CIDR format) to match interfaces. Example: ['192.168.1.0/24', '10.0.0.0/8'].",
			},
			"interfaces": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Description: "List of interface names to match directly. Example: ['ether1', 'bridge1'].",
			},
			"type": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("broadcast"),
				Description: "Network type: 'broadcast' (Ethernet), 'ptp' (point-to-point), 'ptmp' (point-to-multipoint), 'nbma' (non-broadcast), 'ptmp-broadcast', 'virtual-link'. Default is 'broadcast'.",
				Validators: []validator.String{
					stringvalidator.OneOf("broadcast", "ptp", "ptmp", "nbma", "ptmp-broadcast", "virtual-link"),
				},
			},
			"disabled": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
				Description: "Whether this template is disabled. Default is false.",
			},
			"comment": schema.StringAttribute{
				Optional:    true,
				Description: "Comment for this interface template.",
			},
			"cost": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Default:     int64default.StaticInt64(10),
				Description: "Interface cost (metric). Lower is better. Default is 10.",
			},
			"priority": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Default:     int64default.StaticInt64(128),
				Description: "Router priority for DR/BDR election (0-255). Higher wins. Priority 0 means never become DR/BDR. Default is 128.",
			},
			"passive": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
				Description: "If true, interface does not send/receive OSPF packets (only advertised). Default is false.",
			},
			"auth": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("none"),
				Description: "Authentication type: 'none', 'simple' (plaintext), 'md5', 'sha1', 'sha256', 'sha384', 'sha512'. Default is 'none'.",
				Validators: []validator.String{
					stringvalidator.OneOf("none", "simple", "md5", "sha1", "sha256", "sha384", "sha512"),
				},
			},
			"auth_key": schema.StringAttribute{
				Optional:    true,
				Sensitive:   true,
				Description: "Authentication key/password. Required when auth is not 'none'.",
			},
			"auth_id": schema.Int64Attribute{
				Optional:    true,
				Description: "Authentication key ID (1-255). Used for MD5 and SHA authentication to support key rollover.",
			},
			"hello_interval": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("10s"),
				Description: "Hello packet interval. RouterOS time format (e.g., '10s', '30s'). Default is '10s'.",
			},
			"dead_interval": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("40s"),
				Description: "Dead interval (neighbor considered down after no hello). Should be 4x hello_interval. Default is '40s'.",
			},
			"retransmit_interval": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("5s"),
				Description: "LSA retransmission interval. Default is '5s'.",
			},
			"transmit_delay": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("1s"),
				Description: "LSA transmission delay (link-state age increment). Default is '1s'.",
			},
			"wait_time": schema.StringAttribute{
				Optional:    true,
				Description: "Time to wait before electing DR/BDR on interface startup. Default is dead_interval.",
			},
			"vlink_transit_area": schema.StringAttribute{
				Optional:    true,
				Description: "Transit area name for virtual link. Required when type='virtual-link'.",
			},
			"vlink_neighbor_id": schema.StringAttribute{
				Optional:    true,
				Description: "Remote router ID for virtual link (format: 1.1.1.1). Required when type='virtual-link'.",
			},
			"dynamic": schema.BoolAttribute{
				Computed:    true,
				Description: "Whether this template was dynamically created.",
			},
			"invalid": schema.BoolAttribute{
				Computed:    true,
				Description: "Whether this template configuration is invalid.",
			},
		},
	}
}

// Create creates the resource and sets the initial Terraform state.
func (r *ospfInterfaceTemplateV7Resource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var terraformModel ospfInterfaceTemplateV7
	var mikrotikModel client.OspfInterfaceTemplateV7
	GenericCreateResource(&terraformModel, &mikrotikModel, r.client)(ctx, req, resp)
}

// Read refreshes the Terraform state with the latest data.
func (r *ospfInterfaceTemplateV7Resource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var terraformModel ospfInterfaceTemplateV7
	var mikrotikModel client.OspfInterfaceTemplateV7
	GenericReadResource(&terraformModel, &mikrotikModel, r.client)(ctx, req, resp)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *ospfInterfaceTemplateV7Resource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var terraformModel ospfInterfaceTemplateV7
	var mikrotikModel client.OspfInterfaceTemplateV7
	GenericUpdateResource(&terraformModel, &mikrotikModel, r.client)(ctx, req, resp)
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *ospfInterfaceTemplateV7Resource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var terraformModel ospfInterfaceTemplateV7
	var mikrotikModel client.OspfInterfaceTemplateV7
	GenericDeleteResource(&terraformModel, &mikrotikModel, r.client)(ctx, req, resp)
}

func (r *ospfInterfaceTemplateV7Resource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Import by ID (interface templates don't have unique name)
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
