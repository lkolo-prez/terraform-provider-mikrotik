package mikrotik

import (
	"context"
	"fmt"
	"time"

	"github.com/ddelnano/terraform-provider-mikrotik/client"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"

	tftypes "github.com/hashicorp/terraform-plugin-framework/types"
)

type systemBackup struct {
	client *client.Mikrotik
}

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource              = &systemBackup{}
	_ resource.ResourceWithConfigure = &systemBackup{}
)

// NewSystemBackupResource is a helper function to simplify the provider implementation.
func NewSystemBackupResource() resource.Resource {
	return &systemBackup{}
}

func (r *systemBackup) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*client.Mikrotik)
}

func (r *systemBackup) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_system_backup"
}

func (r *systemBackup) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Creates and manages RouterOS system backups. Note: This resource creates backups but doesn't download them. Use with file upload/download mechanisms for remote storage.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The backup filename (includes .backup extension).",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "Backup name (without .backup extension). Auto-generated if not specified.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplaceIfConfigured(),
				},
			},
			"password": schema.StringAttribute{
				Description: "Backup encryption password. Required for encrypted backups.",
				Optional:    true,
				Sensitive:   true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"dont_encrypt": schema.BoolAttribute{
				Description: "Disable encryption (not recommended). Default: false.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"size": schema.StringAttribute{
				Description: "Backup file size (computed after creation).",
				Computed:    true,
			},
			"creation_time": schema.StringAttribute{
				Description: "Backup creation timestamp (computed).",
				Computed:    true,
			},
		},
	}
}

// systemBackupModel maps the resource schema data.
type systemBackupModel struct {
	Id           tftypes.String `tfsdk:"id"`
	Name         tftypes.String `tfsdk:"name"`
	Password     tftypes.String `tfsdk:"password"`
	DontEncrypt  tftypes.Bool   `tfsdk:"dont_encrypt"`
	Size         tftypes.String `tfsdk:"size"`
	CreationTime tftypes.String `tfsdk:"creation_time"`
}

// Create creates the resource and sets the initial Terraform state.
func (r *systemBackup) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data systemBackupModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate backup name if not provided
	backupName := data.Name.ValueString()
	if backupName == "" {
		backupName = fmt.Sprintf("terraform-backup-%s", time.Now().Format("20060102-150405"))
		data.Name = tftypes.StringValue(backupName)
	}

	// Create backup
	params := &client.SystemBackupSave{
		Name:        backupName,
		Password:    data.Password.ValueString(),
		DontEncrypt: data.DontEncrypt.ValueBool(),
	}

	err := r.client.SaveSystemBackup(params)
	if err != nil {
		resp.Diagnostics.AddError("Error creating system backup", err.Error())
		return
	}

	// Wait a moment for backup to complete
	time.Sleep(2 * time.Second)

	// Verify backup was created and get details
	backupFile := backupName + ".backup"
	backup, err := r.client.FindSystemBackup(backupFile)
	if err != nil {
		resp.Diagnostics.AddError("Error verifying backup creation", err.Error())
		return
	}

	// Map response to state
	data.Id = tftypes.StringValue(backup.Name)
	data.Size = tftypes.StringValue(backup.Size)
	data.CreationTime = tftypes.StringValue(backup.Creation)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Read refreshes the Terraform state with the latest data.
func (r *systemBackup) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data systemBackupModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	backup, err := r.client.FindSystemBackup(data.Id.ValueString())
	if err != nil {
		if client.IsNotFoundError(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading system backup", err.Error())
		return
	}

	// Map response to state
	data.Size = tftypes.StringValue(backup.Size)
	data.CreationTime = tftypes.StringValue(backup.Creation)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *systemBackup) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Backups are immutable - changes require replacement
	// This should not be called due to RequiresReplace modifiers
	resp.Diagnostics.AddError(
		"Update not supported",
		"System backups are immutable. Changes require resource replacement.",
	)
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *systemBackup) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data systemBackupModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteSystemBackup(data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error deleting system backup", err.Error())
		return
	}
}
