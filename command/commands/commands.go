package commands

import (
	"github.com/stararks/dct/command/container"
	"github.com/stararks/dct/command/image"

	"github.com/spf13/cobra"
)

func AddCommands(cmd *cobra.Command) {
	cmd.AddCommand(
		image.CmdloadImageCommand(),
		image.CmdlistImageCommand(),
		image.CmdcleanImageCommand(),

		container.CmdcleanContainersCommand(),
	)
}
