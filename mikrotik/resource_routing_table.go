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

type routingTable struct {
	client *client.Mikrotik
}

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &routingTable{}
	_ resource.ResourceWithConfigure   = &routingTable{}
	_ resource.ResourceWithImportState = &routingTable{}
)

// NewRoutingTableResource is a helper function to simplify the provider implementation.
func NewRoutingTableResource() resource.Resource {
	return &routingTable{}
}

func (r *routingTable) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*client.Mikrotik)
}

// Metadata returns the resource type name.
func (r *routingTable) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_routing_table"
}

// Schema defines the schema for the resource.
func (r *routingTable) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Creates a MikroTik Routing Table / VRF for RouterOS v7. RouterOS v7 introduces proper VRF (Virtual Routing and Forwarding) support with multiple routing tables, essential for enterprise deployments with multi-tenancy, MPLS L3VPN, and route isolation scenarios.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Description: "Unique ID of this resource.",
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "Name of the routing table. Used as VRF identifier in BGP, OSPF, and other routing protocols. This name is referenced by other resources (e.g., `mikrotik_bgp_instance_v7.vrf`).",
			},
			"fib": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
				Description: "Forwarding Information Base (FIB) table selection. Specifies which FIB table to use for this routing table. If not specified, the system will use the default FIB.",
			},
			"disabled": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
				Description: "Whether the routing table is disabled. When disabled, routes in this table will not be used for forwarding decisions.",
			},
			"comment": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
				Description: "Descriptive comment for this routing table/VRF.",
			},
		},
	}
}

// Create creates the resource and sets the initial Terraform state.
func (r *routingTable) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var terraformModel routingTableModel
	var mikrotikModel client.RoutingTable
	GenericCreateResource(&terraformModel, &mikrotikModel, r.client)(ctx, req, resp)
}

// Read refreshes the Terraform state with the latest data.
func (r *routingTable) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var terraformModel routingTableModel
	var mikrotikModel client.RoutingTable
	GenericReadResource(&terraformModel, &mikrotikModel, r.client)(ctx, req, resp)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *routingTable) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var terraformModel routingTableModel
	var mikrotikModel client.RoutingTable
	GenericUpdateResource(&terraformModel, &mikrotikModel, r.client)(ctx, req, resp)
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *routingTable) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var terraformModel routingTableModel
	var mikrotikModel client.RoutingTable
	GenericDeleteResource(&terraformModel, &mikrotikModel, r.client)(ctx, req, resp)
}

func (r *routingTable) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}

type routingTableModel struct {
	Id       tftypes.String `tfsdk:"id"`
	Name     tftypes.String `tfsdk:"name"`
	FIB      tftypes.String `tfsdk:"fib"`
	Disabled tftypes.Bool   `tfsdk:"disabled"`
	Comment  tftypes.String `tfsdk:"comment"`
}
