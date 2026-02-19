package tailscalecmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func NewDockerClientInstructionsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "docker-client-instructions",
		Short: "Show the instructions how to run a local tailscale client in docker.",
		Run: func(cmd *cobra.Command, args []string) {
			instructions := `
Start tailscale in container:
=============================

Start the tailscale docker container with providing a SOCKS5 proxy:
    docker run --name tailscale-client --rm -d tailscale/tailscale tailscaled --tun=userspace-networking --socks5-server=localhost:8080

Connect using your preauth key to your headscale instance:
    docker exec -it tailscale-client tailscale up --login-server=https://your-headscale.example.com --auth-key="${YOUR_PREAUTH_KEY}"

Show the tailscale status
    docker exec -it tailscale-client tailscale status

Send data from one client to another using userspace-networking:
    On the receiver:
		tailscale serve --bg --tcp 8888 tcp://localhost:8888 ; nc -l -p 8888 > out.data
	On the sender:
		cat tosend.data | tailscale nc receiverip 8888
	Close receiving on the sender:
		tailscale serve --tcp=8888 off

Disconnect:
    docker exec -it tailscale-client tailscale down

Remove the container:
	docker rm --force tailscale-client
`
			fmt.Print(instructions)
		},
	}

	return cmd
}