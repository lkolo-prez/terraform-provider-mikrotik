package mikrotik

import (
	"context"

	"github.com/ddelnano/terraform-provider-mikrotik/client"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type ospfAreaV7 struct {
	Name            types.String `tfsdk:"name"`
	AreaId          types.String `tfsdk:"area_id"`
	Instance        types.String `tfsdk:"instance"`
	Type            types.String `tfsdk:"type"`
	Disabled        types.Bool   `tfsdk:"disabled"`
	Comment         types.String `tfsdk:"comment"`
	DefaultCost     types.Int64  `tfsdk:"default_cost"`
	NoSummaries     types.Bool   `tfsdk:"no_summaries"`
	NssaTranslator  types.String `tfsdk:"nssa_translator"`
	NssaPropagation types.Bool   `tfsdk:"nssa_propagation"`
	// Computed
	Id      types.String `tfsdk:"id"`
	Dynamic types.Bool   `tfsdk:"dynamic"`
	Invalid types.Bool   `tfsdk:"invalid"`
}

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &ospfAreaV7Resource{}
	_ resource.ResourceWithConfigure   = &ospfAreaV7Resource{}
	_ resource.ResourceWithImportState = &ospfAreaV7Resource{}
)

// NewOspfAreaV7Resource is a helper function to simplify the provider implementation.
func NewOspfAreaV7Resource() resource.Resource {
	return &ospfAreaV7Resource{}
}

type ospfAreaV7Resource struct {
	client *client.Mikrotik
}

func (r *ospfAreaV7Resource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*client.Mikrotik)
}

func (r *ospfAreaV7Resource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_ospf_area_v7"
}

// Schema defines the schema for the resource.
func (r *ospfAreaV7Resource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages OSPF area configuration (RouterOS v7). Areas organize OSPF routing domains hierarchically.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "Unique ID of the OSPF area.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "Name of the OSPF area.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"area_id": schema.StringAttribute{
				Required:    true,
				Description: "Area ID in IPv4 address format (e.g., 0.0.0.0 for backbone, 1.1.1.1 for area 1).",
			},
			"instance": schema.StringAttribute{
				Required:    true,
				Description: "Name of the OSPF instance this area belongs to.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"type": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("default"),
				Description: "Area type: 'default' (standard), 'stub', or 'nssa' (Not-So-Stubby Area). Default is 'default'.",
				Validators: []validator.String{
					stringvalidator.OneOf("default", "stub", "nssa"),
				},
			},
			"disabled": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
				Description: "Whether the area is disabled. Default is false.",
			},
			"comment": schema.StringAttribute{
				Optional:    true,
				Description: "Comment for the OSPF area.",
			},
			"default_cost": schema.Int64Attribute{
				Optional:    true,
				Description: "Cost of the default route injected into stub area. Only for stub/nssa areas.",
			},
			"no_summaries": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
				Description: "If true, creates a totally stubby area (no Type 3 LSAs). Only for stub areas. Default is false.",
			},
			"nssa_translator": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("candidate"),
				Description: "NSSA ABR translator role: 'yes' (always), 'no' (never), 'candidate' (elect). Only for NSSA areas. Default is 'candidate'.",
				Validators: []validator.String{
					stringvalidator.OneOf("yes", "no", "candidate"),
				},
			},
			"nssa_propagation": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(true),
				Description: "Whether to propagate NSSA Type-7 LSAs. Only for NSSA areas. Default is true.",
			},
			"dynamic": schema.BoolAttribute{
				Computed:    true,
				Description: "Whether this area was dynamically created.",
			},
			"invalid": schema.BoolAttribute{
				Computed:    true,
				Description: "Whether this area configuration is invalid.",
			},
		},
	}
}

// Create creates the resource and sets the initial Terraform state.
func (r *ospfAreaV7Resource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var terraformModel ospfAreaV7
	var mikrotikModel client.OspfAreaV7
	GenericCreateResource(&terraformModel, &mikrotikModel, r.client)(ctx, req, resp)
}

// Read refreshes the Terraform state with the latest data.
func (r *ospfAreaV7Resource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var terraformModel ospfAreaV7
	var mikrotikModel client.OspfAreaV7
	GenericReadResource(&terraformModel, &mikrotikModel, r.client)(ctx, req, resp)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *ospfAreaV7Resource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var terraformModel ospfAreaV7
	var mikrotikModel client.OspfAreaV7
	GenericUpdateResource(&terraformModel, &mikrotikModel, r.client)(ctx, req, resp)
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *ospfAreaV7Resource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var terraformModel ospfAreaV7
	var mikrotikModel client.OspfAreaV7
	GenericDeleteResource(&terraformModel, &mikrotikModel, r.client)(ctx, req, resp)
}

func (r *ospfAreaV7Resource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Import by name
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}
