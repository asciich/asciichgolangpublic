package miniocmdoptions

import (
	"context"

	"github.com/minio/madmin-go/v3"
	"github.com/minio/minio-go/v7"
	"github.com/spf13/cobra"
)

type MinioCmdOptions struct {
	OverrideUse    string
	GetEndpoint    func(context.Context, *cobra.Command) string
	GetClient      func(context.Context, *cobra.Command) *minio.Client
	GetAdminClient func(context.Context, *cobra.Command) *madmin.AdminClient
}
