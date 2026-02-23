package miniocmd

import (
	"context"

	"github.com/minio/minio-go/v7"
	"github.com/spf13/cobra"
	"github.com/asciich/asciichgolangpublic/pkg/defaultclicommands/storagecmd/s3cmd/miniocmd/miniocmdoptions"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/mustutils"
	"github.com/asciich/asciichgolangpublic/pkg/storage/s3/nativeminioclient"
	"github.com/asciich/asciichgolangpublic/pkg/storage/s3/s3options"
)

// For the default and generic commands set the options argument to nil.
// The options can be used to define a complete command tree for a specific minio server.
func NewMinioCmd(options *miniocmdoptions.MinioCmdOptions) *cobra.Command {
	if options == nil {
		options = &miniocmdoptions.MinioCmdOptions{}
	}

	if options.GetEndpoint == nil {
		options.GetEndpoint = defaultGetEndpoint
	}

	if options.GetClient == nil {
		options.GetClient = defaultGetClient
	}

	use := "minio"
	if options.OverrideUse != "" {
		use = options.OverrideUse
	}

	cmd := &cobra.Command{
		Use:   use,
		Short: "Minio (S3 compatible server) related commands.",
	}

	cmd.AddCommand(
		NewListBuckets(options),
		NewListObjectsCmd(options),
	)

	cmd.PersistentFlags().String("endpoint", "", "The minio endpoint/ server to use.")

	return cmd
}

func defaultGetClient(ctx context.Context, cmd *cobra.Command) *minio.Client {
	endpoint := defaultGetEndpoint(ctx, cmd)
	return mustutils.Must(nativeminioclient.NewClientFromEnvVars(ctx, endpoint, &s3options.NewS3ClientOptions{}))
}

func defaultGetEndpoint(ctx context.Context, cmd *cobra.Command) string {
	endpoint, err := cmd.Flags().GetString("endpoint")
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	logging.LogInfoByCtxf(ctx, "Minio endpoint is '%s'.", endpoint)

	return endpoint
}
