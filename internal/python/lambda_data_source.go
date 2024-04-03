package python

import (
	"context"

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

	// Read API call logic

	// Example data value setting
	data.Id = types.StringValue("example-id")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
