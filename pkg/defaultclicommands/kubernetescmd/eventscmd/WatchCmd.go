package eventscmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/nativekubernetesoo"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	v1 "k8s.io/api/core/v1"
)

func NewWatchCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "watch",
		Short: "Watch K8s events",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := contextutils.GetVerbosityContextByCobraCmd(cmd)

			namespace, err := cmd.Flags().GetString("namespace")
			if err != nil {
				logging.LogGoErrorFatalWithTrace(err)
			}

			apiversion, err := cmd.Flags().GetString("api-version")
			if err != nil {
				logging.LogGoErrorFatalWithTrace(err)
			}

			kind, err := cmd.Flags().GetString("kind")
			if err != nil {
				logging.LogGoErrorFatalWithTrace(err)
			}

			name, err := cmd.Flags().GetString("name")
			if err != nil {
				logging.LogGoErrorFatalWithTrace(err)
			}

			err = nativekubernetesoo.WatchEvents(
				ctx,
				&kubernetesutils.WatchEventOptions{
					Namespace:                namespace,
					InvolvedObjectName:       name,
					InvolvedObjectAPIVersion: apiversion,
					InvolvedObjectKind:       kind,
				},
				onCreate,
				onUpdate,
				onDelete,
			)
			if err != nil {
				logging.LogGoErrorFatalWithTrace(err)
			}

			<-ctx.Done()
		},
	}

	cmd.PersistentFlags().String("namespace", "", "Only watch events in --namespace.")
	cmd.PersistentFlags().String("api-version", "", "Only watch events matching given --api-version. e.g apiextensions.k8s.io/v1")
	cmd.PersistentFlags().String("kind", "", "Only watch events matching given --kind.")
	cmd.PersistentFlags().String("name", "", "Only watch events matching given --name.")

	return cmd
}

func onCreate(event *v1.Event) {
	fmt.Println("create: " + nativekubernetesoo.EventToString(event))
}

func onUpdate(event *v1.Event) {
	fmt.Println("update: " + nativekubernetesoo.EventToString(event))
}

func onDelete(event *v1.Event) {
	fmt.Println("delete: " + nativekubernetesoo.EventToString(event))
}
