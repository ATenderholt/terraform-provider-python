package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

func LogDebug(ctx context.Context, msg string, args map[string]interface{}) {
	tflog.Debug(ctx, msg, args)
	fmt.Printf("[DEBUG] %s %v\n", msg, args)
}

func LogError(ctx context.Context, msg string, args map[string]interface{}) {
	tflog.Error(ctx, msg, args)
	fmt.Printf("[ERROR] %s %v\n", msg, args)
}
