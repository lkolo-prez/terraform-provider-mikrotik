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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"

	tftypes "github.com/hashicorp/terraform-plugin-framework/types"
)

type snmpCommunity struct {
	client *client.Mikrotik
}

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &snmpCommunity{}
	_ resource.ResourceWithConfigure   = &snmpCommunity{}
	_ resource.ResourceWithImportState = &snmpCommunity{}
)

// NewSnmpCommunityResource is a helper function to simplify the provider implementation.
func NewSnmpCommunityResource() resource.Resource {
	return &snmpCommunity{}
}

func (r *snmpCommunity) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*client.Mikrotik)
}

func (r *snmpCommunity) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_snmp_community"
}

func (r *snmpCommunity) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages SNMP communities for network monitoring access control.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The ID of the SNMP community.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "SNMP community name. Required.",
				Required:    true,
			},
			"security": schema.StringAttribute{
				Description: "Security level: 'none' (no restrictions), 'authorized' (IP restriction), 'private' (encrypted). Default: 'none'.",
				Optional:    true,
				Computed:    true,
				Validators: []validator.String{
					stringvalidator.OneOf("none", "authorized", "private"),
				},
			},
			"read_access": schema.BoolAttribute{
				Description: "Allow read access. Default: true.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(true),
			},
			"write_access": schema.BoolAttribute{
				Description: "Allow write access (SET operations). Default: false.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"address": schema.StringAttribute{
				Description: "Allowed source IP address or network (CIDR). Empty = all addresses. Example: '192.168.1.0/24'.",
				Optional:    true,
			},
			"disabled": schema.BoolAttribute{
				Description: "Whether the community is disabled. Default: false.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
		},
	}
}

// snmpCommunityModel maps the resource schema data.
type snmpCommunityModel struct {
	Id          tftypes.String `tfsdk:"id"`
	Name        tftypes.String `tfsdk:"name"`
	Security    tftypes.String `tfsdk:"security"`
	ReadAccess  tftypes.Bool   `tfsdk:"read_access"`
	WriteAccess tftypes.Bool   `tfsdk:"write_access"`
	Address     tftypes.String `tfsdk:"address"`
	Disabled    tftypes.Bool   `tfsdk:"disabled"`
}

// Create creates the resource and sets the initial Terraform state.
func (r *snmpCommunity) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data snmpCommunityModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	community := &client.SnmpCommunity{
		Name:        data.Name.ValueString(),
		Security:    data.Security.ValueString(),
		ReadAccess:  data.ReadAccess.ValueBool(),
		WriteAccess: data.WriteAccess.ValueBool(),
		Address:     data.Address.ValueString(),
		Disabled:    data.Disabled.ValueBool(),
	}

	created, err := r.client.AddSnmpCommunity(community)
	if err != nil {
		resp.Diagnostics.AddError("Error creating SNMP community", err.Error())
		return
	}

	// Map response to state
	data.Id = tftypes.StringValue(created.Id)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Read refreshes the Terraform state with the latest data.
func (r *snmpCommunity) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data snmpCommunityModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	community, err := r.client.FindSnmpCommunity(data.Id.ValueString())
	if err != nil {
		if client.IsNotFoundError(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading SNMP community", err.Error())
		return
	}

	// Map response to state
	data.Name = tftypes.StringValue(community.Name)
	data.Security = tftypes.StringValue(community.Security)
	data.ReadAccess = tftypes.BoolValue(community.ReadAccess)
	data.WriteAccess = tftypes.BoolValue(community.WriteAccess)
	data.Address = tftypes.StringValue(community.Address)
	data.Disabled = tftypes.BoolValue(community.Disabled)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *snmpCommunity) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data snmpCommunityModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	community := &client.SnmpCommunity{
		Id:          data.Id.ValueString(),
		Name:        data.Name.ValueString(),
		Security:    data.Security.ValueString(),
		ReadAccess:  data.ReadAccess.ValueBool(),
		WriteAccess: data.WriteAccess.ValueBool(),
		Address:     data.Address.ValueString(),
		Disabled:    data.Disabled.ValueBool(),
	}

	_, err := r.client.UpdateSnmpCommunity(community)
	if err != nil {
		resp.Diagnostics.AddError("Error updating SNMP community", err.Error())
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *snmpCommunity) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data snmpCommunityModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteSnmpCommunity(data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error deleting SNMP community", err.Error())
		return
	}
}

func (r *snmpCommunity) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
