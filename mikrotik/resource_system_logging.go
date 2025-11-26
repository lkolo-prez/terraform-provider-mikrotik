package mikrotik

import (
	"context"

	"github.com/ddelnano/terraform-provider-mikrotik/client"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"

	tftypes "github.com/hashicorp/terraform-plugin-framework/types"
)

type systemLogging struct {
	client *client.Mikrotik
}

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &systemLogging{}
	_ resource.ResourceWithConfigure   = &systemLogging{}
	_ resource.ResourceWithImportState = &systemLogging{}
)

// NewSystemLoggingResource is a helper function to simplify the provider implementation.
func NewSystemLoggingResource() resource.Resource {
	return &systemLogging{}
}

func (r *systemLogging) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*client.Mikrotik)
}

func (r *systemLogging) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_system_logging"
}

func (r *systemLogging) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Creates and manages a system logging rule that routes log topics to specific actions.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The unique identifier of the logging rule.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"topics": schema.StringAttribute{
				Description: "Comma-separated list of log topics to match. Examples: 'firewall,info', 'system,error,critical', 'bgp'. Common topics: account, system, firewall, wireless, web, script, critical, error, warning, info, debug.",
				Required:    true,
			},
			"action": schema.StringAttribute{
				Description: "Reference to a logging action (destination) by name. Must match an existing mikrotik_system_logging_action resource.",
				Required:    true,
			},
			"prefix": schema.StringAttribute{
				Description: "Optional prefix to add to log messages. Useful for filtering or categorizing logs.",
				Optional:    true,
			},
			"disabled": schema.BoolAttribute{
				Description: "Whether the logging rule is disabled. Default: false.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
		},
	}
}

// systemLoggingModel maps the data source schema data.
type systemLoggingModel struct {
	Id       tftypes.String `tfsdk:"id"`
	Topics   tftypes.String `tfsdk:"topics"`
	Action   tftypes.String `tfsdk:"action"`
	Prefix   tftypes.String `tfsdk:"prefix"`
	Disabled tftypes.Bool   `tfsdk:"disabled"`
}

// Create creates the resource and sets the initial Terraform state.
func (r *systemLogging) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data systemLoggingModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	logging := &client.SystemLogging{
		Topics:   data.Topics.ValueString(),
		Action:   data.Action.ValueString(),
		Prefix:   data.Prefix.ValueString(),
		Disabled: data.Disabled.ValueBool(),
	}

	created, err := r.client.AddSystemLogging(logging)
	if err != nil {
		resp.Diagnostics.AddError("Error creating system logging rule", err.Error())
		return
	}

	// Map response to state
	data.Id = tftypes.StringValue(created.Id)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Read refreshes the Terraform state with the latest data.
func (r *systemLogging) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data systemLoggingModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	logging, err := r.client.FindSystemLogging(data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading system logging rule", err.Error())
		return
	}

	// Map response to state
	data.Topics = tftypes.StringValue(logging.Topics)
	data.Action = tftypes.StringValue(logging.Action)
	data.Prefix = tftypes.StringValue(logging.Prefix)
	data.Disabled = tftypes.BoolValue(logging.Disabled)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *systemLogging) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data systemLoggingModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	logging := &client.SystemLogging{
		Id:       data.Id.ValueString(),
		Topics:   data.Topics.ValueString(),
		Action:   data.Action.ValueString(),
		Prefix:   data.Prefix.ValueString(),
		Disabled: data.Disabled.ValueBool(),
	}

	updated, err := r.client.UpdateSystemLogging(logging)
	if err != nil {
		resp.Diagnostics.AddError("Error updating system logging rule", err.Error())
		return
	}

	// Map response to state
	data.Id = tftypes.StringValue(updated.Id)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *systemLogging) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data systemLoggingModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteSystemLogging(data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error deleting system logging rule", err.Error())
		return
	}
}

// ImportState imports the resource state.
func (r *systemLogging) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Use the ID directly
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
