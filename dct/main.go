package main

import (
	"fmt"
	"os"

	"github.com/stararks/dct/command/commands"

	"github.com/spf13/cobra"
)

func main() {
	cmd := &cobra.Command{
		Use:   "dct",
		Short: "A docker container tool",
	}
	commands.AddCommands(cmd)
	if err := cmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
