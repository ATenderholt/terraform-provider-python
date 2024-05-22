package python

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

var _ provider.Provider = (*pythonProvider)(nil)

func New() func() provider.Provider {
	return func() provider.Provider {
		return &pythonProvider{}
	}
}

type pythonProvider struct{}

type pythonProviderModel struct {
	PipCommand types.String `tfsdk:"pip_command"`
}

func (p *pythonProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"pip_command": schema.StringAttribute{
				Required:            true,
				Optional:            false,
				Description:         "Executable on path to install dependencies (e.g. pip3.10)",
				MarkdownDescription: "Executable on path to install dependencies (e.g. pip3.10)",
				Validators:          nil,
			},
		},
		Description:         "Package Python & dependencies into archives suitable for Cloud serverless",
		MarkdownDescription: "Package Python & dependencies into archives suitable for Cloud serverless",
	}
}

func (p *pythonProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var config pythonProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.ResourceData = config
	resp.DataSourceData = config
}

func (p *pythonProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "python"
}

func (p *pythonProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewAwsLambdaDataSource,
	}
}

func (p *pythonProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{}
}
