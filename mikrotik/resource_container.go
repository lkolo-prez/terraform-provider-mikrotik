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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type container struct {
	client *client.Mikrotik
}

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &container{}
	_ resource.ResourceWithConfigure   = &container{}
	_ resource.ResourceWithImportState = &container{}
)

// NewContainerResource is a helper function to simplify the provider implementation.
func NewContainerResource() resource.Resource {
	return &container{}
}

func (r *container) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*client.Mikrotik)
}

// Metadata returns the resource type name.
func (r *container) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_container"
}

// Schema defines the schema for the resource.
func (r *container) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages OCI containers on RouterOS v7.4+. " +
			"Run containerized applications like Pi-hole, Prometheus, or custom services. " +
			"Compatible with Docker Hub, GCR, Quay, and other OCI registries.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Description: "Unique ID of this container instance",
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "Container name - must be unique",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"remote_image": schema.StringAttribute{
				Optional:    true,
				Description: "Container image name from registry (e.g., pihole/pihole:latest). Use with registry_url configured in container_config.",
				Validators: []validator.String{
					stringvalidator.ConflictsWith(path.MatchRoot("file")),
				},
			},
			"tag": schema.StringAttribute{
				Computed:    true,
				Description: "Image tag (read-only)",
			},
			"digest": schema.StringAttribute{
				Computed:    true,
				Description: "Image digest (read-only)",
			},
			"file": schema.StringAttribute{
				Optional:    true,
				Description: "Path to container tarball file for offline import (e.g., disk1/pihole.tar)",
				Validators: []validator.String{
					stringvalidator.ConflictsWith(path.MatchRoot("remote_image")),
				},
			},
			"interface": schema.StringAttribute{
				Required:    true,
				Description: "Veth interface name for container networking (must be pre-configured)",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"root_dir": schema.StringAttribute{
				Required:    true,
				Description: "Directory path for container root filesystem. MUST be on external storage (e.g., disk1/containers/pihole)",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"cmd": schema.StringAttribute{
				Optional:    true,
				Description: "Command arguments to pass to container entrypoint",
			},
			"entrypoint": schema.StringAttribute{
				Optional:    true,
				Description: "Override container entrypoint (e.g., /bin/sh)",
			},
			"workdir": schema.StringAttribute{
				Optional:    true,
				Description: "Working directory for cmd/entrypoint",
			},
			"mounts": schema.StringAttribute{
				Optional:    true,
				Description: "Comma-separated list of mount names (e.g., 'MOUNT_PIHOLE_ETC,MOUNT_PIHOLE_DNSMASQ'). Mounts must be pre-configured in /container/mounts",
			},
			"envlist": schema.StringAttribute{
				Optional:    true,
				Description: "Environment variable list name. Environment variables must be pre-configured in /container/envs with matching list name",
			},
			"dns": schema.StringAttribute{
				Optional:    true,
				Description: "Custom DNS servers for container (e.g., '1.1.1.1,8.8.8.8')",
			},
			"domain_name": schema.StringAttribute{
				Optional:    true,
				Description: "Domain name for container",
			},
			"hostname": schema.StringAttribute{
				Optional:    true,
				Description: "Hostname for container identification",
			},
			"logging": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
				Description: "Enable logging to RouterOS logs (yes/no)",
			},
			"start_on_boot": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
				Description: "Auto-start container on device boot",
			},
			"auto_restart_interval": schema.StringAttribute{
				Optional:    true,
				Description: "Interval for auto-restart on failure (e.g., '10s', '1m')",
			},
			"stop_signal": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Default:     int64default.StaticInt64(15),
				Description: "Linux signal number to send if container doesn't stop after 10 seconds. Default: 15 (SIGTERM)",
			},
			"devices": schema.StringAttribute{
				Optional:    true,
				Description: "Pass-through physical devices to container (e.g., '/dev/ttyUSB0')",
			},
			"cpu_list": schema.StringAttribute{
				Optional:    true,
				Description: "CPU cores container can use (e.g., '0-1' for cores 0 and 1)",
			},
			"user": schema.StringAttribute{
				Optional:    true,
				Description: "User and group to run container as (e.g., '1000:1000' or 'www-data'). Set to '0:0' for root privileges.",
			},
			"memory_high": schema.Int64Attribute{
				Optional:    true,
				Description: "RAM usage limit in bytes for this container (overrides global limit)",
			},
			"status": schema.StringAttribute{
				Computed:    true,
				Description: "Container status (stopped, running, extracting, etc.)",
			},
			"comment": schema.StringAttribute{
				Optional:    true,
				Description: "Comment for documentation",
			},
		},
	}
}

// containerModel maps the resource schema data.
type containerModel struct {
	Id                  types.String `tfsdk:"id"`
	Name                types.String `tfsdk:"name"`
	RemoteImage         types.String `tfsdk:"remote_image"`
	Tag                 types.String `tfsdk:"tag"`
	Digest              types.String `tfsdk:"digest"`
	File                types.String `tfsdk:"file"`
	Interface           types.String `tfsdk:"interface"`
	RootDir             types.String `tfsdk:"root_dir"`
	Cmd                 types.String `tfsdk:"cmd"`
	Entrypoint          types.String `tfsdk:"entrypoint"`
	Workdir             types.String `tfsdk:"workdir"`
	Mounts              types.String `tfsdk:"mounts"`
	Envlist             types.String `tfsdk:"envlist"`
	Dns                 types.String `tfsdk:"dns"`
	DomainName          types.String `tfsdk:"domain_name"`
	Hostname            types.String `tfsdk:"hostname"`
	Logging             types.Bool   `tfsdk:"logging"`
	StartOnBoot         types.Bool   `tfsdk:"start_on_boot"`
	AutoRestartInterval types.String `tfsdk:"auto_restart_interval"`
	StopSignal          types.Int64  `tfsdk:"stop_signal"`
	Devices             types.String `tfsdk:"devices"`
	CpuList             types.String `tfsdk:"cpu_list"`
	User                types.String `tfsdk:"user"`
	MemoryHigh          types.Int64  `tfsdk:"memory_high"`
	Status              types.String `tfsdk:"status"`
	Comment             types.String `tfsdk:"comment"`
}

// Create creates the resource and sets the initial Terraform state.
func (r *container) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan containerModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	con := &client.Container{
		Name:                plan.Name.ValueString(),
		RemoteImage:         plan.RemoteImage.ValueString(),
		File:                plan.File.ValueString(),
		Interface:           plan.Interface.ValueString(),
		RootDir:             plan.RootDir.ValueString(),
		Cmd:                 plan.Cmd.ValueString(),
		Entrypoint:          plan.Entrypoint.ValueString(),
		Workdir:             plan.Workdir.ValueString(),
		Mounts:              plan.Mounts.ValueString(),
		Envlist:             plan.Envlist.ValueString(),
		Dns:                 plan.Dns.ValueString(),
		DomainName:          plan.DomainName.ValueString(),
		Hostname:            plan.Hostname.ValueString(),
		Logging:             plan.Logging.ValueBool(),
		StartOnBoot:         plan.StartOnBoot.ValueBool(),
		AutoRestartInterval: plan.AutoRestartInterval.ValueString(),
		StopSignal:          int(plan.StopSignal.ValueInt64()),
		Devices:             plan.Devices.ValueString(),
		CpuList:             plan.CpuList.ValueString(),
		User:                plan.User.ValueString(),
		MemoryHigh:          int(plan.MemoryHigh.ValueInt64()),
		Comment:             plan.Comment.ValueString(),
	}

	created, err := r.client.AddContainer(con)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Container",
			"Could not create container: "+err.Error(),
		)
		return
	}

	plan.Id = types.StringValue(created.Id)
	plan.Tag = types.StringValue(created.Tag)
	plan.Digest = types.StringValue(created.Digest)
	plan.Status = types.StringValue(created.Status)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

// Read refreshes the Terraform state with the latest data.
func (r *container) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state containerModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	con, err := r.client.FindContainerById(state.Id.ValueString())
	if err != nil {
		if client.IsNotFoundError(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error Reading Container",
			"Could not read container ID "+state.Id.ValueString()+": "+err.Error(),
		)
		return
	}

	state.Name = types.StringValue(con.Name)
	state.RemoteImage = types.StringValue(con.RemoteImage)
	state.Tag = types.StringValue(con.Tag)
	state.Digest = types.StringValue(con.Digest)
	state.File = types.StringValue(con.File)
	state.Interface = types.StringValue(con.Interface)
	state.RootDir = types.StringValue(con.RootDir)
	state.Cmd = types.StringValue(con.Cmd)
	state.Entrypoint = types.StringValue(con.Entrypoint)
	state.Workdir = types.StringValue(con.Workdir)
	state.Mounts = types.StringValue(con.Mounts)
	state.Envlist = types.StringValue(con.Envlist)
	state.Dns = types.StringValue(con.Dns)
	state.DomainName = types.StringValue(con.DomainName)
	state.Hostname = types.StringValue(con.Hostname)
	state.Logging = types.BoolValue(con.Logging)
	state.StartOnBoot = types.BoolValue(con.StartOnBoot)
	state.AutoRestartInterval = types.StringValue(con.AutoRestartInterval)
	state.StopSignal = types.Int64Value(int64(con.StopSignal))
	state.Devices = types.StringValue(con.Devices)
	state.CpuList = types.StringValue(con.CpuList)
	state.User = types.StringValue(con.User)
	if con.MemoryHigh > 0 {
		state.MemoryHigh = types.Int64Value(int64(con.MemoryHigh))
	}
	state.Status = types.StringValue(con.Status)
	state.Comment = types.StringValue(con.Comment)

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *container) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan containerModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	con := &client.Container{
		Id:                  plan.Id.ValueString(),
		Name:                plan.Name.ValueString(),
		RemoteImage:         plan.RemoteImage.ValueString(),
		File:                plan.File.ValueString(),
		Interface:           plan.Interface.ValueString(),
		RootDir:             plan.RootDir.ValueString(),
		Cmd:                 plan.Cmd.ValueString(),
		Entrypoint:          plan.Entrypoint.ValueString(),
		Workdir:             plan.Workdir.ValueString(),
		Mounts:              plan.Mounts.ValueString(),
		Envlist:             plan.Envlist.ValueString(),
		Dns:                 plan.Dns.ValueString(),
		DomainName:          plan.DomainName.ValueString(),
		Hostname:            plan.Hostname.ValueString(),
		Logging:             plan.Logging.ValueBool(),
		StartOnBoot:         plan.StartOnBoot.ValueBool(),
		AutoRestartInterval: plan.AutoRestartInterval.ValueString(),
		StopSignal:          int(plan.StopSignal.ValueInt64()),
		Devices:             plan.Devices.ValueString(),
		CpuList:             plan.CpuList.ValueString(),
		User:                plan.User.ValueString(),
		MemoryHigh:          int(plan.MemoryHigh.ValueInt64()),
		Comment:             plan.Comment.ValueString(),
	}

	updated, err := r.client.UpdateContainer(con)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Container",
			"Could not update container: "+err.Error(),
		)
		return
	}

	plan.Tag = types.StringValue(updated.Tag)
	plan.Digest = types.StringValue(updated.Digest)
	plan.Status = types.StringValue(updated.Status)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *container) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state containerModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Stop container before deletion if it's running
	err := r.client.StopContainer(state.Name.ValueString())
	if err != nil {
		// Don't fail if container is already stopped
		resp.Diagnostics.AddWarning(
			"Error Stopping Container",
			"Could not stop container before deletion (may already be stopped): "+err.Error(),
		)
	}

	err = r.client.DeleteContainer(state.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Container",
			"Could not delete container: "+err.Error(),
		)
		return
	}
}

// ImportState imports the resource into Terraform state.
func (r *container) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Import by name
	con, err := r.client.FindContainer(req.ID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Importing Container",
			"Could not find container: "+err.Error(),
		)
		return
	}

	resp.State.SetAttribute(ctx, path.Root("id"), con.Id)
	resp.State.SetAttribute(ctx, path.Root("name"), con.Name)
}
