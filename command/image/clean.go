package image

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/stararks/dct/command/container"

	"github.com/spf13/cobra"
)

func CmdcleanImageCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cleani",
		Short: "Clean dangling images",
		Run: func(cmd *cobra.Command, args []string) {
			cleanDangling()
		},
	}
	return cmd
}

func cleanDangling() {
	args := []string{"images", "-q", "--filter", "dangling=true"}
	out, err := exec.Command("docker", args...).Output()
	if err != nil {
		fmt.Println("Get dangling images error")
		return
	}
	images := strings.Fields(strings.TrimSpace(strings.Replace(string(out), "\n", " ", -1)))

	if len(images) == 0 {
		return
	}

	ids := []string{"rmi"}
	ids = append(ids, images...)
	out, err = exec.Command("docker", ids...).Output()
	if err != nil {
		fmt.Println("Dangling images can't be cleaned. It is possible container running base on the images. Running docker native command to stop and clean the corresponding containers or use dtool cleanc subcommand to stop and clean ALL containers")
		fmt.Printf("Enter [y] to run cleancontainers default[No]:")
		input := bufio.NewScanner(os.Stdin)
		if input.Scan() && input.Text() == "y" {
			container.CleanContainers()
		}
	}
}
