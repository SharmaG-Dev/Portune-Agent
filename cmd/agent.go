package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/SharmaG-Dev/Portune-Agent/internal/agent"
	"github.com/spf13/cobra"
)

var (
	serverURL string

	agentToken string

	localTarget string

	requestedSubdomain string

	tunnelName string
)

var agentCmd = &cobra.Command{
	Use:   "agent",
	Short: "Start Portune tunnel agent",

	RunE: func(
		cmd *cobra.Command,
		args []string,
	) error {

		if agentToken == "" {

			agentToken =
				os.Getenv(
					"PORTUNE_AGENT_TOKEN",
				)
		}

		if agentToken == "" {

			return fmt.Errorf(
				"agent token is required: use --token or PORTUNE_AGENT_TOKEN",
			)
		}

		if serverURL == "" {

			return fmt.Errorf(
				"Portune server URL is required",
			)
		}

		if localTarget == "" {

			return fmt.Errorf(
				"local target is required",
			)
		}

		ctx,
			cancel :=
			signal.NotifyContext(
				context.Background(),

				os.Interrupt,

				syscall.SIGTERM,
			)

		defer cancel()

		portuneAgent :=
			agent.New(
				agent.Config{
					ServerURL: serverURL,

					Token: agentToken,

					LocalTarget: localTarget,

					Subdomain: requestedSubdomain,

					TunnelName: tunnelName,
				},
			)

		return portuneAgent.Run(ctx)
	},
}

func init() {

	rootCmd.AddCommand(
		agentCmd,
	)

	agentCmd.Flags().
		StringVar(
			&serverURL,

			"server",

			"http://localhost:3000/tunnel",

			"Portune tunnel server URL",
		)

	agentCmd.Flags().
		StringVar(
			&agentToken,

			"token",

			"",

			"Portune agent authentication token",
		)

	agentCmd.Flags().
		StringVar(
			&localTarget,

			"target",

			"http://localhost:4000",

			"Local application URL",
		)

	agentCmd.Flags().
		StringVar(
			&requestedSubdomain,

			"subdomain",

			"",

			"Requested public subdomain",
		)

	agentCmd.Flags().
		StringVar(
			&tunnelName,

			"name",

			"",

			"Tunnel name",
		)
}
