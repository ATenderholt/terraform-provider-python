package provider

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

var _ provider.Provider = (*pythonPackageProvider)(nil)

func New() func() provider.Provider {
	return func() provider.Provider {
		return &pythonPackageProvider{}
	}
}

type pythonPackageProvider struct{}

type pythonPackageProviderModel struct {
	PipCommand    types.String `tfsdk:"pip_command"`
	PipExtraFlags types.String `tfsdk:"pip_extra_flags"`
}

func (p *pythonPackageProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"pip_command": schema.StringAttribute{
				Required:            true,
				Optional:            false,
				Description:         "Executable on path to install dependencies (e.g. pip3.10)",
				MarkdownDescription: "Executable on path to install dependencies (e.g. pip3.10)",
				Validators:          nil,
			},
			"pip_extra_flags": schema.StringAttribute{
				Required:            false,
				Optional:            true,
				Description:         "Extra flags to pass to pip when installing dependencies",
				MarkdownDescription: "Extra flags to pass to pip when installing dependencies",
				Validators:          nil,
			},
		},
		Description:         "Package Python & dependencies into single archive",
		MarkdownDescription: "Package Python & dependencies into single archive",
	}
}

func (p *pythonPackageProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var config pythonPackageProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.DataSourceData = config
}

func (p *pythonPackageProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "python_package"
}

func (p *pythonPackageProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewAwsLambdaDataSource,
	}
}

func (p *pythonPackageProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{}
}
