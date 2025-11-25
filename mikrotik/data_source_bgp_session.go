package mikrotik

import (
	"context"

	"github.com/ddelnano/terraform-provider-mikrotik/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	tftypes "github.com/hashicorp/terraform-plugin-framework/types"
)

type bgpSession struct {
	client *client.Mikrotik
}

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &bgpSession{}
	_ datasource.DataSourceWithConfigure = &bgpSession{}
)

// NewBgpSessionDataSource is a helper function to simplify the provider implementation.
func NewBgpSessionDataSource() datasource.DataSource {
	return &bgpSession{}
}

func (d *bgpSession) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*client.Mikrotik)
}

// Metadata returns the data source type name.
func (d *bgpSession) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_bgp_session"
}

// Schema defines the schema for the data source.
func (d *bgpSession) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Reads MikroTik BGP session information. This is a read-only data source showing cached BGP session information including status, capabilities, and negotiated parameters. Even if the BGP session is not active anymore, the cache can still be stored for some time.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "Unique ID of the BGP session.",
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "Name of the BGP session to lookup.",
			},
			"established": schema.BoolAttribute{
				Computed:    true,
				Description: "Whether the BGP session is established.",
			},
			"remote_address": schema.StringAttribute{
				Computed:    true,
				Description: "Remote peer IP address.",
			},
			"remote_as": schema.Int64Attribute{
				Computed:    true,
				Description: "Remote peer AS number.",
			},
			"remote_id": schema.StringAttribute{
				Computed:    true,
				Description: "Remote peer router ID.",
			},
			"remote_capabilities": schema.StringAttribute{
				Computed:    true,
				Description: "Remote peer BGP capabilities (e.g., mp, rr, as4, gr).",
			},
			"remote_afi": schema.StringAttribute{
				Computed:    true,
				Description: "Remote peer address families (e.g., ip, ipv6, vpnv4).",
			},
			"remote_messages": schema.Int64Attribute{
				Computed:    true,
				Description: "Number of messages received from remote peer.",
			},
			"remote_bytes": schema.Int64Attribute{
				Computed:    true,
				Description: "Number of bytes received from remote peer.",
			},
			"remote_eor": schema.StringAttribute{
				Computed:    true,
				Description: "Remote peer End-of-RIB marker status.",
			},
			"remote_refused_cap_opt": schema.BoolAttribute{
				Computed:    true,
				Description: "Whether remote peer refused optional capability negotiation.",
			},
			"local_address": schema.StringAttribute{
				Computed:    true,
				Description: "Local IP address used for this session.",
			},
			"local_as": schema.Int64Attribute{
				Computed:    true,
				Description: "Local AS number.",
			},
			"local_id": schema.StringAttribute{
				Computed:    true,
				Description: "Local router ID.",
			},
			"local_capabilities": schema.StringAttribute{
				Computed:    true,
				Description: "Local BGP capabilities (e.g., mp, rr, as4, gr).",
			},
			"local_messages": schema.Int64Attribute{
				Computed:    true,
				Description: "Number of messages sent to remote peer.",
			},
			"local_bytes": schema.Int64Attribute{
				Computed:    true,
				Description: "Number of bytes sent to remote peer.",
			},
			"local_eor": schema.StringAttribute{
				Computed:    true,
				Description: "Local End-of-RIB marker status.",
			},
			"hold_time": schema.StringAttribute{
				Computed:    true,
				Description: "Negotiated hold time (e.g., 3m).",
			},
			"keepalive_time": schema.StringAttribute{
				Computed:    true,
				Description: "Negotiated keepalive time (e.g., 1m).",
			},
			"uptime": schema.StringAttribute{
				Computed:    true,
				Description: "Session uptime (e.g., 4s70ms).",
			},
			"output_procid": schema.Int64Attribute{
				Computed:    true,
				Description: "Output process ID.",
			},
			"output_keep_sent_attrs": schema.BoolAttribute{
				Computed:    true,
				Description: "Whether to keep sent attributes in memory.",
			},
			"output_last_notification": schema.StringAttribute{
				Computed:    true,
				Description: "Last BGP NOTIFICATION message sent (hex encoded).",
			},
			"input_procid": schema.Int64Attribute{
				Computed:    true,
				Description: "Input process ID.",
			},
			"input_limit_process_routes": schema.Int64Attribute{
				Computed:    true,
				Description: "Limit on the number of routes to process from this peer.",
			},
			"state": schema.StringAttribute{
				Computed:    true,
				Description: "BGP session state (e.g., established, idle, connect).",
			},
			"prefix_count": schema.Int64Attribute{
				Computed:    true,
				Description: "Number of prefixes received from this peer.",
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *bgpSession) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data bgpSessionModel
	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Find BGP session by name
	session, err := d.client.FindBgpSession(data.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading BGP Session",
			"Could not read BGP session name "+data.Name.ValueString()+": "+err.Error(),
		)
		return
	}

	// Map client model to Terraform model
	data.Id = tftypes.StringValue(session.Id)
	data.Name = tftypes.StringValue(session.Name)
	data.Established = tftypes.BoolValue(session.Established)
	data.RemoteAddress = tftypes.StringValue(session.RemoteAddress)
	data.RemoteAS = tftypes.Int64Value(int64(session.RemoteAS))
	data.RemoteID = tftypes.StringValue(session.RemoteID)
	data.RemoteCapabilities = tftypes.StringValue(session.RemoteCapabilities)
	data.RemoteAFI = tftypes.StringValue(session.RemoteAFI)
	data.RemoteMessages = tftypes.Int64Value(int64(session.RemoteMessages))
	data.RemoteBytes = tftypes.Int64Value(int64(session.RemoteBytes))
	data.RemoteEOR = tftypes.StringValue(session.RemoteEOR)
	data.RemoteRefusedCapOpt = tftypes.BoolValue(session.RemoteRefusedCapOpt)
	data.LocalAddress = tftypes.StringValue(session.LocalAddress)
	data.LocalAS = tftypes.Int64Value(int64(session.LocalAS))
	data.LocalID = tftypes.StringValue(session.LocalID)
	data.LocalCapabilities = tftypes.StringValue(session.LocalCapabilities)
	data.LocalMessages = tftypes.Int64Value(int64(session.LocalMessages))
	data.LocalBytes = tftypes.Int64Value(int64(session.LocalBytes))
	data.LocalEOR = tftypes.StringValue(session.LocalEOR)
	data.HoldTime = tftypes.StringValue(session.HoldTime)
	data.KeepaliveTime = tftypes.StringValue(session.KeepaliveTime)
	data.Uptime = tftypes.StringValue(session.Uptime)
	data.OutputProcID = tftypes.Int64Value(int64(session.OutputProcID))
	data.OutputKeepSentAttrs = tftypes.BoolValue(session.OutputKeepSentAttrs)
	data.OutputLastNotification = tftypes.StringValue(session.OutputLastNotification)
	data.InputProcID = tftypes.Int64Value(int64(session.InputProcID))
	data.InputLimitProcessRoutes = tftypes.Int64Value(int64(session.InputLimitProcessRoutes))
	data.State = tftypes.StringValue(session.State)
	data.PrefixCount = tftypes.Int64Value(int64(session.PrefixCount))

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

type bgpSessionModel struct {
	Id                       tftypes.String `tfsdk:"id"`
	Name                     tftypes.String `tfsdk:"name"`
	Established              tftypes.Bool   `tfsdk:"established"`
	RemoteAddress            tftypes.String `tfsdk:"remote_address"`
	RemoteAS                 tftypes.Int64  `tfsdk:"remote_as"`
	RemoteID                 tftypes.String `tfsdk:"remote_id"`
	RemoteCapabilities       tftypes.String `tfsdk:"remote_capabilities"`
	RemoteAFI                tftypes.String `tfsdk:"remote_afi"`
	RemoteMessages           tftypes.Int64  `tfsdk:"remote_messages"`
	RemoteBytes              tftypes.Int64  `tfsdk:"remote_bytes"`
	RemoteEOR                tftypes.String `tfsdk:"remote_eor"`
	RemoteRefusedCapOpt      tftypes.Bool   `tfsdk:"remote_refused_cap_opt"`
	LocalAddress             tftypes.String `tfsdk:"local_address"`
	LocalAS                  tftypes.Int64  `tfsdk:"local_as"`
	LocalID                  tftypes.String `tfsdk:"local_id"`
	LocalCapabilities        tftypes.String `tfsdk:"local_capabilities"`
	LocalMessages            tftypes.Int64  `tfsdk:"local_messages"`
	LocalBytes               tftypes.Int64  `tfsdk:"local_bytes"`
	LocalEOR                 tftypes.String `tfsdk:"local_eor"`
	HoldTime                 tftypes.String `tfsdk:"hold_time"`
	KeepaliveTime            tftypes.String `tfsdk:"keepalive_time"`
	Uptime                   tftypes.String `tfsdk:"uptime"`
	OutputProcID             tftypes.Int64  `tfsdk:"output_procid"`
	OutputKeepSentAttrs      tftypes.Bool   `tfsdk:"output_keep_sent_attrs"`
	OutputLastNotification   tftypes.String `tfsdk:"output_last_notification"`
	InputProcID              tftypes.Int64  `tfsdk:"input_procid"`
	InputLimitProcessRoutes  tftypes.Int64  `tfsdk:"input_limit_process_routes"`
	State                    tftypes.String `tfsdk:"state"`
	PrefixCount              tftypes.Int64  `tfsdk:"prefix_count"`
}
