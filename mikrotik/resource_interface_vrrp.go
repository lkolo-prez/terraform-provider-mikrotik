package mikrotik

import (
	"context"

	"github.com/ddelnano/terraform-provider-mikrotik/client"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
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

	tftypes "github.com/hashicorp/terraform-plugin-framework/types"
)

type interfaceVrrp struct {
	client *client.Mikrotik
}

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &interfaceVrrp{}
	_ resource.ResourceWithConfigure   = &interfaceVrrp{}
	_ resource.ResourceWithImportState = &interfaceVrrp{}
)

// NewInterfaceVrrpResource is a helper function to simplify the provider implementation.
func NewInterfaceVrrpResource() resource.Resource {
	return &interfaceVrrp{}
}

func (r *interfaceVrrp) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*client.Mikrotik)
}

func (r *interfaceVrrp) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_interface_vrrp"
}

func (r *interfaceVrrp) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Creates and manages a VRRP (Virtual Router Redundancy Protocol) interface for high availability setups.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The unique identifier of the VRRP interface.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "Name of the VRRP interface.",
				Required:    true,
			},
			"interface": schema.StringAttribute{
				Description: "Physical interface name that VRRP will run on.",
				Required:    true,
			},
			"vrid": schema.Int64Attribute{
				Description: "Virtual Router ID (1-255). Must be the same on both master and backup routers.",
				Required:    true,
				Validators: []validator.Int64{
					int64validator.Between(1, 255),
				},
			},
			"priority": schema.Int64Attribute{
				Description: "Priority of the router (1-254). Higher priority means higher chance to become master. Default: 100.",
				Optional:    true,
				Computed:    true,
				Default:     int64default.StaticInt64(100),
				Validators: []validator.Int64{
					int64validator.Between(1, 254),
				},
			},
			"version": schema.Int64Attribute{
				Description: "VRRP version (2 or 3). Default: 3. Version 3 supports both IPv4 and IPv6.",
				Optional:    true,
				Computed:    true,
				Default:     int64default.StaticInt64(3),
				Validators: []validator.Int64{
					int64validator.OneOf(2, 3),
				},
			},
			"authentication": schema.StringAttribute{
				Description: "Authentication type: 'none', 'simple', or 'ah' (Authentication Header). Default: none.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("none"),
				Validators: []validator.String{
					stringvalidator.OneOf("none", "simple", "ah"),
				},
			},
			"password": schema.StringAttribute{
				Description: "Password for authentication (if authentication is 'simple' or 'ah'). Marked as sensitive.",
				Optional:    true,
				Sensitive:   true,
			},
			"interval": schema.StringAttribute{
				Description: "Advertisement interval. Default: 1s. Format: integer followed by time unit (s, ms).",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("1s"),
			},
			"preemption_mode": schema.BoolAttribute{
				Description: "Enable preemption mode. If true, higher priority router will become master even if lower priority is currently master. Default: true.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(true),
			},
			"v3_protocol": schema.StringAttribute{
				Description: "Protocol for VRRP version 3: 'ipv4' or 'ipv6'. Only applicable when version=3.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("ipv4"),
				Validators: []validator.String{
					stringvalidator.OneOf("ipv4", "ipv6"),
				},
			},
			"on_backup": schema.StringAttribute{
				Description: "Script to execute when router transitions to BACKUP state.",
				Optional:    true,
			},
			"on_master": schema.StringAttribute{
				Description: "Script to execute when router transitions to MASTER state.",
				Optional:    true,
			},
			"disabled": schema.BoolAttribute{
				Description: "Whether the VRRP interface is disabled. Default: false.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"comment": schema.StringAttribute{
				Description: "Comment for the VRRP interface.",
				Optional:    true,
			},
			"running": schema.BoolAttribute{
				Description: "Whether the VRRP interface is currently running (computed).",
				Computed:    true,
			},
		},
	}
}

// interfaceVrrpModel describes the resource data model.
type interfaceVrrpModel struct {
	Id             tftypes.String `tfsdk:"id"`
	Name           tftypes.String `tfsdk:"name"`
	Interface      tftypes.String `tfsdk:"interface"`
	Vrid           tftypes.Int64  `tfsdk:"vrid"`
	Priority       tftypes.Int64  `tfsdk:"priority"`
	Version        tftypes.Int64  `tfsdk:"version"`
	Authentication tftypes.String `tfsdk:"authentication"`
	Password       tftypes.String `tfsdk:"password"`
	Interval       tftypes.String `tfsdk:"interval"`
	PreemptionMode tftypes.Bool   `tfsdk:"preemption_mode"`
	V3Protocol     tftypes.String `tfsdk:"v3_protocol"`
	OnBackup       tftypes.String `tfsdk:"on_backup"`
	OnMaster       tftypes.String `tfsdk:"on_master"`
	Disabled       tftypes.Bool   `tfsdk:"disabled"`
	Comment        tftypes.String `tfsdk:"comment"`
	Running        tftypes.Bool   `tfsdk:"running"`
}

func (r *interfaceVrrp) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var terraformModel interfaceVrrpModel
	var mikrotikModel client.InterfaceVrrp
	GenericCreateResource(&terraformModel, &mikrotikModel, r.client)(ctx, req, resp)
}

func (r *interfaceVrrp) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var terraformModel interfaceVrrpModel
	var mikrotikModel client.InterfaceVrrp
	GenericReadResource(&terraformModel, &mikrotikModel, r.client)(ctx, req, resp)
}

func (r *interfaceVrrp) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var terraformModel interfaceVrrpModel
	var mikrotikModel client.InterfaceVrrp
	GenericUpdateResource(&terraformModel, &mikrotikModel, r.client)(ctx, req, resp)
}

func (r *interfaceVrrp) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var terraformModel interfaceVrrpModel
	var mikrotikModel client.InterfaceVrrp
	GenericDeleteResource(&terraformModel, &mikrotikModel, r.client)(ctx, req, resp)
}

func (r *interfaceVrrp) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to name attribute
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}
