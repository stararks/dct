package container

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

func CmdcleanContainersCommand() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "cleanc",
		Short: "Clean ALL stopped and running containers",
		Run: func(cmd *cobra.Command, args []string) {
			CleanContainers()
		},
	}
	return cmd
}

func operateContainer(op, c string) ([]string, error) {
	switch op {
	case "ps":
		args := []string{op, "-qa"}
		out, err := exec.Command("docker", args...).Output()
		if err != nil {
			return nil, fmt.Errorf("Get container error")
		}
		containers := strings.Fields(strings.TrimSpace(strings.Replace(string(out), "\n", " ", -1)))
		return containers, nil

	case "inspect":
		args := []string{op, "--format=\"{{ .Name }}\"", c}
		out, err := exec.Command("docker", args...).Output()
		if err != nil {
			return nil, err
		}
		return []string{string(out[1 : len(string(out))-1])}, nil

	case "stop", "rm":
		_, err := exec.Command("docker", op, c).Output()
		return nil, err
	default:
		return nil, fmt.Errorf("Wrong docker command:%s", op)
	}
}

func CleanContainers() {
	containers, err := operateContainer("ps", "")
	if err != nil {
		os.Exit(1)
	}

	for _, c := range containers {
		cn, err := operateContainer("inspect", c)
		if err != nil {
			fmt.Println(err)
			continue
		}

		fmt.Printf("Stoping %s ...", cn[0])
		_, err = operateContainer("stop", c)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Printf("\rStopping %s ... done\n", cn[0])

	}

	time.Sleep(1 * time.Second) // retrieve containers in case automatically removed container exits
	containers, err = operateContainer("ps", "")
	if err != nil {
		os.Exit(1)
	}

	for _, c := range containers {
		cn, err := operateContainer("inspect", c)
		if err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Printf("Removing %s ...", cn[0])

		_, err = operateContainer("rm", c)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Printf("\rRemoving %s ... done\n", cn[0])

	}
}
