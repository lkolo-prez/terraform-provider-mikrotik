package mikrotik

import (
	"context"

	"github.com/ddelnano/terraform-provider-mikrotik/client"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type containerConfig struct {
	client *client.Mikrotik
}

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &containerConfig{}
	_ resource.ResourceWithConfigure   = &containerConfig{}
	_ resource.ResourceWithImportState = &containerConfig{}
)

// NewContainerConfigResource is a helper function to simplify the provider implementation.
func NewContainerConfigResource() resource.Resource {
	return &containerConfig{}
}

func (r *containerConfig) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*client.Mikrotik)
}

// Metadata returns the resource type name.
func (r *containerConfig) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_container_config"
}

// Schema defines the schema for the resource.
func (r *containerConfig) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages global container configuration on RouterOS v7.4+. " +
			"This is a singleton resource - there is only one container configuration per device. " +
			"Configure registry URL, temporary directory, and memory limits for containers.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Description: "Unique ID - always 'container_config' for singleton",
			},
			"registry_url": schema.StringAttribute{
				Optional:    true,
				Description: "External registry URL from where containers will be downloaded. Default: https://registry-1.docker.io",
			},
			"tmpdir": schema.StringAttribute{
				Optional:    true,
				Description: "Container extraction directory. Should point to external storage (e.g., disk1/containers/tmp)",
			},
			"memory_high": schema.Int64Attribute{
				Optional: true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
				Description: "RAM usage limit in bytes for all containers (soft limit). " +
					"When exceeded, processes are throttled and put under reclaim pressure. " +
					"Example: 256000000 (256MB). Set to 0 for unlimited.",
			},
			"username": schema.StringAttribute{
				Optional:    true,
				Sensitive:   true,
				Description: "Username for registry authentication (RouterOS 7.8+)",
			},
			"password": schema.StringAttribute{
				Optional:    true,
				Sensitive:   true,
				Description: "Password for registry authentication (RouterOS 7.8+)",
			},
		},
	}
}

// containerConfigModel maps the resource schema data.
type containerConfigModel struct {
	Id          types.String `tfsdk:"id"`
	RegistryUrl types.String `tfsdk:"registry_url"`
	Tmpdir      types.String `tfsdk:"tmpdir"`
	MemoryHigh  types.Int64  `tfsdk:"memory_high"`
	Username    types.String `tfsdk:"username"`
	Password    types.String `tfsdk:"password"`
}

// Create creates the resource and sets the initial Terraform state.
func (r *containerConfig) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan containerConfigModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	config := &client.ContainerConfig{
		RegistryUrl: plan.RegistryUrl.ValueString(),
		Tmpdir:      plan.Tmpdir.ValueString(),
		MemoryHigh:  int(plan.MemoryHigh.ValueInt64()),
		Username:    plan.Username.ValueString(),
		Password:    plan.Password.ValueString(),
	}

	updatedConfig, err := r.client.UpdateContainerConfig(config)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Configuring Container Settings",
			"Could not update container config: "+err.Error(),
		)
		return
	}

	// Singleton resource always has the same ID
	plan.Id = types.StringValue("container_config")
	plan.RegistryUrl = types.StringValue(updatedConfig.RegistryUrl)
	plan.Tmpdir = types.StringValue(updatedConfig.Tmpdir)
	plan.MemoryHigh = types.Int64Value(int64(updatedConfig.MemoryHigh))
	// Don't read back sensitive values
	if plan.Username.IsNull() {
		plan.Username = types.StringValue("")
	}
	if plan.Password.IsNull() {
		plan.Password = types.StringValue("")
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

// Read refreshes the Terraform state with the latest data.
func (r *containerConfig) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state containerConfigModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	config, err := r.client.GetContainerConfig()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Container Config",
			"Could not read container config: "+err.Error(),
		)
		return
	}

	state.Id = types.StringValue("container_config")
	state.RegistryUrl = types.StringValue(config.RegistryUrl)
	state.Tmpdir = types.StringValue(config.Tmpdir)
	state.MemoryHigh = types.Int64Value(int64(config.MemoryHigh))
	// Keep existing sensitive values from state

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *containerConfig) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan containerConfigModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	config := &client.ContainerConfig{
		RegistryUrl: plan.RegistryUrl.ValueString(),
		Tmpdir:      plan.Tmpdir.ValueString(),
		MemoryHigh:  int(plan.MemoryHigh.ValueInt64()),
		Username:    plan.Username.ValueString(),
		Password:    plan.Password.ValueString(),
	}

	updatedConfig, err := r.client.UpdateContainerConfig(config)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Container Config",
			"Could not update container config: "+err.Error(),
		)
		return
	}

	plan.RegistryUrl = types.StringValue(updatedConfig.RegistryUrl)
	plan.Tmpdir = types.StringValue(updatedConfig.Tmpdir)
	plan.MemoryHigh = types.Int64Value(int64(updatedConfig.MemoryHigh))

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

// Delete deletes the resource - for singleton, this resets to defaults
func (r *containerConfig) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// For singleton resource, deletion means resetting to defaults
	config := &client.ContainerConfig{
		RegistryUrl: "https://registry-1.docker.io",
		Tmpdir:      "",
		MemoryHigh:  0,
		Username:    "",
		Password:    "",
	}

	_, err := r.client.UpdateContainerConfig(config)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Resetting Container Config",
			"Could not reset container config to defaults: "+err.Error(),
		)
		return
	}

	// Resource is removed from state
}

// ImportState imports the resource into Terraform state.
func (r *containerConfig) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Singleton resource - always use the same ID
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
	resp.State.SetAttribute(ctx, path.Root("id"), "container_config")
}
