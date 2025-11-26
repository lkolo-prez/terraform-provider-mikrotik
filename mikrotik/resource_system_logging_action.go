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

	tftypes "github.com/hashicorp/terraform-plugin-framework/types"
)

type systemLoggingAction struct {
	client *client.Mikrotik
}

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &systemLoggingAction{}
	_ resource.ResourceWithConfigure   = &systemLoggingAction{}
	_ resource.ResourceWithImportState = &systemLoggingAction{}
)

// NewSystemLoggingActionResource is a helper function to simplify the provider implementation.
func NewSystemLoggingActionResource() resource.Resource {
	return &systemLoggingAction{}
}

func (r *systemLoggingAction) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*client.Mikrotik)
}

func (r *systemLoggingAction) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_system_logging_action"
}

func (r *systemLoggingAction) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Creates and manages a system logging action that defines where and how logs are stored or transmitted.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The unique identifier of the logging action.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "Name of the logging action. Used by logging rules to reference this action.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"target": schema.StringAttribute{
				Description: "Log target type: 'disk', 'echo', 'email', 'memory', 'remote'.",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.OneOf("disk", "echo", "email", "memory", "remote"),
				},
			},
			"remote": schema.StringAttribute{
				Description: "Remote syslog server address and port. Format: 'ip:port' or 'hostname:port'. Example: '192.168.1.10:514'. Required when target='remote'.",
				Optional:    true,
			},
			"remote_port": schema.StringAttribute{
				Description: "Remote syslog port (deprecated, use 'remote' with 'ip:port' format instead). Default: 514.",
				Optional:    true,
			},
			"bsd_syslog": schema.BoolAttribute{
				Description: "Use BSD syslog format for remote logging. Default: false.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"syslog_facility": schema.StringAttribute{
				Description: "Syslog facility code: 'kern', 'user', 'mail', 'daemon', 'auth', 'syslog', 'lpr', 'news', 'uucp', 'cron', 'authpriv', 'ftp', 'ntp', 'security', 'console', 'local0-local7'. Default: 'daemon'.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("daemon"),
				Validators: []validator.String{
					stringvalidator.OneOf(
						"kern", "user", "mail", "daemon", "auth", "syslog", "lpr", "news",
						"uucp", "cron", "authpriv", "ftp", "ntp", "security", "console",
						"local0", "local1", "local2", "local3", "local4", "local5", "local6", "local7",
					),
				},
			},
			"src_address": schema.StringAttribute{
				Description: "Source IP address for remote syslog connections. Uses router's egress interface IP by default.",
				Optional:    true,
			},
			"memory": schema.StringAttribute{
				Description: "Number of log lines to keep in memory. Default: 100. Only applicable when target='memory'.",
				Optional:    true,
			},
			"disk_file_name": schema.StringAttribute{
				Description: "Log file name for disk target. Example: 'logs'. Only applicable when target='disk'.",
				Optional:    true,
			},
			"disk_file_count": schema.StringAttribute{
				Description: "Number of log files for rotation. Default: 2. Only applicable when target='disk'.",
				Optional:    true,
			},
			"disk_lines_per_file": schema.StringAttribute{
				Description: "Maximum lines per log file before rotation. Default: 1000. Only applicable when target='disk'.",
				Optional:    true,
			},
			"remember": schema.BoolAttribute{
				Description: "Remember logs on disk even after reboot (for memory target). Default: false.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
		},
	}
}

// systemLoggingActionModel maps the data source schema data.
type systemLoggingActionModel struct {
	Id                tftypes.String `tfsdk:"id"`
	Name              tftypes.String `tfsdk:"name"`
	Target            tftypes.String `tfsdk:"target"`
	Remote            tftypes.String `tfsdk:"remote"`
	RemotePort        tftypes.String `tfsdk:"remote_port"`
	BsdSyslog         tftypes.Bool   `tfsdk:"bsd_syslog"`
	SyslogFacility    tftypes.String `tfsdk:"syslog_facility"`
	SrcAddress        tftypes.String `tfsdk:"src_address"`
	Memory            tftypes.String `tfsdk:"memory"`
	DiskFileName      tftypes.String `tfsdk:"disk_file_name"`
	DiskFileCount     tftypes.String `tfsdk:"disk_file_count"`
	DiskLinesPerFile  tftypes.String `tfsdk:"disk_lines_per_file"`
	Remember          tftypes.Bool   `tfsdk:"remember"`
}

// Create creates the resource and sets the initial Terraform state.
func (r *systemLoggingAction) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data systemLoggingActionModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	action := &client.SystemLoggingAction{
		Name:             data.Name.ValueString(),
		Target:           data.Target.ValueString(),
		Remote:           data.Remote.ValueString(),
		RemotePort:       data.RemotePort.ValueString(),
		BsdSyslog:        data.BsdSyslog.ValueBool(),
		SyslogFacility:   data.SyslogFacility.ValueString(),
		SrcAddress:       data.SrcAddress.ValueString(),
		Memory:           data.Memory.ValueString(),
		DiskFileName:     data.DiskFileName.ValueString(),
		DiskFileCount:    data.DiskFileCount.ValueString(),
		DiskLinesPerFile: data.DiskLinesPerFile.ValueString(),
		Remember:         data.Remember.ValueBool(),
	}

	created, err := r.client.AddSystemLoggingAction(action)
	if err != nil {
		resp.Diagnostics.AddError("Error creating system logging action", err.Error())
		return
	}

	// Map response to state
	data.Id = tftypes.StringValue(created.Id)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Read refreshes the Terraform state with the latest data.
func (r *systemLoggingAction) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data systemLoggingActionModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	action, err := r.client.FindSystemLoggingAction(data.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading system logging action", err.Error())
		return
	}

	// Map response to state
	data.Id = tftypes.StringValue(action.Id)
	data.Target = tftypes.StringValue(action.Target)
	data.Remote = tftypes.StringValue(action.Remote)
	data.RemotePort = tftypes.StringValue(action.RemotePort)
	data.BsdSyslog = tftypes.BoolValue(action.BsdSyslog)
	data.SyslogFacility = tftypes.StringValue(action.SyslogFacility)
	data.SrcAddress = tftypes.StringValue(action.SrcAddress)
	data.Memory = tftypes.StringValue(action.Memory)
	data.DiskFileName = tftypes.StringValue(action.DiskFileName)
	data.DiskFileCount = tftypes.StringValue(action.DiskFileCount)
	data.DiskLinesPerFile = tftypes.StringValue(action.DiskLinesPerFile)
	data.Remember = tftypes.BoolValue(action.Remember)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *systemLoggingAction) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data systemLoggingActionModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	action := &client.SystemLoggingAction{
		Id:               data.Id.ValueString(),
		Name:             data.Name.ValueString(),
		Target:           data.Target.ValueString(),
		Remote:           data.Remote.ValueString(),
		RemotePort:       data.RemotePort.ValueString(),
		BsdSyslog:        data.BsdSyslog.ValueBool(),
		SyslogFacility:   data.SyslogFacility.ValueString(),
		SrcAddress:       data.SrcAddress.ValueString(),
		Memory:           data.Memory.ValueString(),
		DiskFileName:     data.DiskFileName.ValueString(),
		DiskFileCount:    data.DiskFileCount.ValueString(),
		DiskLinesPerFile: data.DiskLinesPerFile.ValueString(),
		Remember:         data.Remember.ValueBool(),
	}

	updated, err := r.client.UpdateSystemLoggingAction(action)
	if err != nil {
		resp.Diagnostics.AddError("Error updating system logging action", err.Error())
		return
	}

	// Map response to state
	data.Id = tftypes.StringValue(updated.Id)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *systemLoggingAction) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data systemLoggingActionModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteSystemLoggingAction(data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error deleting system logging action", err.Error())
		return
	}
}

// ImportState imports the resource state.
func (r *systemLoggingAction) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Import by name
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}
