package python

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ datasource.DataSource = (*awsLambdaDataSource)(nil)

func NewAwsLambdaDataSource() datasource.DataSource {
	return &awsLambdaDataSource{}
}

type awsLambdaDataSource struct{}

type awsLambdaDataSourceModel struct {
	Id                       types.String `tfsdk:"id"`
	SourceDir                types.String `tfsdk:"source_dir"`
	ArchivePath              types.String `tfsdk:"archive_path"`
	ArchiveBase64Sha256      types.String `tfsdk:"archive_base64sha256"`
	DependenciesPath         types.String `tfsdk:"dependencies_path"`
	DependenciesBase64Sha256 types.String `tfsdk:"dependencies_base64sha256"`
}

func (d *awsLambdaDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_aws_lambda"
}

func (d *awsLambdaDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"source_dir": schema.StringAttribute{
				Description:         "Directory containing Python code, including requirements.txt to install dependencies",
				MarkdownDescription: "Directory containing Python code, including requirements.txt to install dependencies",
				Required:            true,
			},
			"archive_path": schema.StringAttribute{
				Description:         "Path for resulting ZIP file containing Python code",
				MarkdownDescription: "Path for resulting ZIP file containing Python code",
				Required:            true,
			},
			"archive_base64sha256": schema.StringAttribute{
				Description:         "Base64 Encoded SHA256 checksum of ZIP file containing Python code",
				MarkdownDescription: "Base64 Encoded SHA256 checksum of ZIP file containing Python code",
				Computed:            true,
			},
			"dependencies_path": schema.StringAttribute{
				Description:         "Path for resulting ZIP file containing dependencies",
				MarkdownDescription: "Path for resulting ZIP file containing dependencies",
				Optional:            true,
			},
			"dependencies_base64sha256": schema.StringAttribute{
				Description:         "Base64 Encoded SHA256 checksum of ZIP file containing dependencies",
				MarkdownDescription: "Base64 Encoded SHA256 checksum of ZIP file containing dependencies",
				Computed:            true,
			},
		},
	}
}

func (d *awsLambdaDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data awsLambdaDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Create file
	archivePath := data.ArchivePath.ValueString()
	a := NewArchiver(archivePath)
	err := a.Open()
	if err != nil {
		resp.Diagnostics.AddError(
			"unable to open archive",
			fmt.Sprintf("unable to open archive '%s': %v", archivePath, err),
		)
		return
	}
	defer a.Close()

	// TODO: add excludes
	err = a.ArchiveDir(data.SourceDir.ValueString(), "", nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"unable to create archive",
			fmt.Sprintf("unable to create archive '%s': %v", archivePath, err),
		)
		return
	}
	a.Close()

	checksum, err := Checksum(archivePath)
	if err != nil {
		resp.Diagnostics.AddError(
			"unable to checksum archive",
			fmt.Sprintf("unable to checksum archive '%s': %v", archivePath, err),
		)
		return
	}

	// Example data value setting
	data.Id = data.SourceDir
	data.ArchiveBase64Sha256 = types.StringValue(checksum)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
