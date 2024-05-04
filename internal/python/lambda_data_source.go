package python

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var _ datasource.DataSource = (*awsLambdaDataSource)(nil)

func NewAwsLambdaDataSource() datasource.DataSource {
	return &awsLambdaDataSource{}
}

type awsLambdaDataSource struct {
	pipExecutor PipExecutor
}

type awsLambdaDataSourceModel struct {
	Id                       types.String `tfsdk:"id"`
	SourceDir                types.String `tfsdk:"source_dir"`
	ArchivePath              types.String `tfsdk:"archive_path"`
	ArchiveBase64Sha256      types.String `tfsdk:"archive_base64sha256"`
	DependenciesPath         types.String `tfsdk:"dependencies_path"`
	DependenciesBase64Sha256 types.String `tfsdk:"dependencies_base64sha256"`
	ExtraArgs                types.String `tfsdk:"extra_args"`
}

func (d *awsLambdaDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	model, ok := req.ProviderData.(pythonPackageProviderModel)
	if !ok {
		resp.Diagnostics.AddError("Unexpected DataSource Configure Type",
			fmt.Sprintf("Expected pythonPackageProviderModel, but got %T", req.ProviderData))
		return
	}

	d.pipExecutor = NewPipExecutor(model.PipCommand.ValueString())
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
			"extra_args": schema.StringAttribute{
				Description:         "Additional arguments for pip install",
				MarkdownDescription: "Additional arguments for `pip install`",
				Optional:            true,
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

	depPath, err := d.installDependencies(ctx, data)
	if err != nil {
		resp.Diagnostics.AddError(
			"unable to install dependencies",
			fmt.Sprintf("unable to install dependencies: %v", err),
		)
		return
	}

	var depChecksum string
	if depPath != "" {
		depChecksum, err = d.packageDependencies(ctx, data, depPath)
	}

	if err != nil {
		resp.Diagnostics.AddError(
			"unable to package dependencies",
			fmt.Sprintf("unable to package dependencies: %v", err),
		)
		return
	}

	// Example data value setting
	data.Id = data.SourceDir
	data.ArchiveBase64Sha256 = types.StringValue(checksum)
	data.DependenciesBase64Sha256 = types.StringValue(depChecksum)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (d *awsLambdaDataSource) installDependencies(ctx context.Context, data awsLambdaDataSourceModel) (string, error) {
	sourceDir := data.SourceDir.ValueString()
	reqPath := filepath.Join(sourceDir, "requirements.txt")
	_, err := os.Stat(reqPath)
	if os.IsNotExist(err) {
		tflog.Debug(ctx, "No requirements.txt found in sourceDir, skipping running pip", map[string]interface{}{
			"sourceDir": sourceDir,
		})
		return "", nil
	}

	tempDir := os.TempDir()
	now := time.Now()
	installPath := filepath.Join(tempDir, "terraform", now.Format("20060102-150405.00000"))
	tflog.Debug(ctx, "installing python dependencies", map[string]interface{}{
		"sourceDir":   sourceDir,
		"installPath": installPath,
	})

	var extraArgs []string
	if !data.ExtraArgs.IsNull() {
		extraArgs = strings.Split(data.ExtraArgs.ValueString(), " ")
	}

	err = d.pipExecutor.Install(ctx, reqPath, installPath, extraArgs...)
	if err != nil {
		tflog.Error(ctx, "unable to install dependencies via pip", map[string]interface{}{
			"sourceDir": sourceDir,
			"error":     err,
		})
		return "", fmt.Errorf("unable to install dependencies via pip: %w", err)
	}

	return installPath, nil
}

func (d *awsLambdaDataSource) packageDependencies(ctx context.Context, data awsLambdaDataSourceModel, installPath string) (string, error) {
	version, err := d.pipExecutor.GetPythonVersion(ctx)
	if err != nil {
		tflog.Error(ctx, "unable to determine python version from pip", map[string]interface{}{
			"error": err,
		})
		return "", fmt.Errorf("unable to determine python version from pip: %w", err)
	}

	archivePath := data.DependenciesPath.ValueString()
	a := NewArchiver(archivePath)
	err = a.Open()
	if err != nil {
		tflog.Error(ctx, "unable to open archiver for dependencies", map[string]interface{}{
			"dependenciesPath": archivePath,
			"error":            err,
		})
		return "", fmt.Errorf("unable open archiver for dependencies: %w", err)
	}
	defer a.Close()

	root := filepath.Join("/python", "lib", "python"+version, "site-packages")
	err = a.ArchiveDir(installPath, root, []string{"*.pyc", "**/*.pyc"})
	if err != nil {
		tflog.Error(ctx, "unable to archive python dependencies", map[string]interface{}{
			"dependenciesPath": archivePath,
			"error":            err,
		})
		return "", fmt.Errorf("unable to archive python dependencies: %w", err)
	}
	a.Close()

	checksum, err := Checksum(archivePath)
	if err != nil {
		tflog.Error(ctx, "unable to checksum python dependencies", map[string]interface{}{
			"dependenciesPath": archivePath,
			"error":            err,
		})
		return "", fmt.Errorf("unable to checksum python dependencies: %w", err)
	}

	return checksum, nil
}
