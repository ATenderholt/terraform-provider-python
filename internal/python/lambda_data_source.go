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
	Id                 types.String `tfsdk:"id"`
	SourceDir          types.String `tfsdk:"source_dir"`
	OutputPath         types.String `tfsdk:"output_path"`
	OutputBase64Sha256 types.String `tfsdk:"output_base64sha256"`
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
			"output_path": schema.StringAttribute{
				Description:         "Path for resulting ZIP file containing Python code and its dependencies",
				MarkdownDescription: "Path for resulting ZIP file containing Python code and its dependencies",
				Required:            true,
			},
			"output_base64sha256": schema.StringAttribute{
				Description: "Base64 Encoded SHA256 checksum of output file",
				Computed:    true,
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
	outputPath := data.OutputPath.ValueString()
	a := NewArchiver(outputPath)
	err := a.Open()
	if err != nil {
		resp.Diagnostics.AddError(
			"unable to open archive",
			fmt.Sprintf("unable to open archive '%s': %v", outputPath, err),
		)
		return
	}
	defer a.Close()

	// TODO: add excludes
	err = a.ArchiveDir(data.SourceDir.ValueString(), "", nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"unable to create archive",
			fmt.Sprintf("unable to create archive '%s': %v", outputPath, err),
		)
		return
	}
	a.Close()

	checksum, err := Checksum(outputPath)
	if err != nil {
		resp.Diagnostics.AddError(
			"unable to checksum archive",
			fmt.Sprintf("unable to checksum archive '%s': %v", outputPath, err),
		)
		return
	}

	// Example data value setting
	data.Id = data.SourceDir
	data.OutputBase64Sha256 = types.StringValue(checksum)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
