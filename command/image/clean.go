package image

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

// CmdcleanImageCommand clean image
func CmdcleanImageCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cleani [reserve-file]",
		Short: "delete images",
		Run: func(cmd *cobra.Command, args []string) {
			var reserve []string
			images := getCacheImages([]string{"images", "-q"})
			if images == nil {
				return
			}
			stat, _ := os.Stdin.Stat()
			if stat.Mode()&os.ModeNamedPipe != 0 {
				stdinContent, _ := ioutil.ReadAll(os.Stdin)
				reserve = strings.Fields(string(stdinContent))
			}
			if len(args) == 1 {
				fileContent, err := getReserveImageFromFile(args[0])
				if err != nil {
					fmt.Println("Get reserved images error")
					return
				}
				reserve = append(reserve, fileContent...)
			}
			if len(reserve) != 0 {
				for _, r := range reserve {
					for i, v := range images {
						if v == r {
							images = append(images[:i], images[i+1:]...)
						}
					}
				}
				cleanImage(images)
			}
			cleanDangling()
			return
		},
	}
	return cmd
}

func cleanImage(images []string) {
	for _, image := range images {
		_, err := exec.Command("docker", "rmi", "--force", image).Output()
		if err != nil {
			fmt.Printf("Image (ID %s) can't be deleted.\n", image)
			// fmt.Printf("Enter [y] to clean containers default[N/No]:")
			// input := bufio.NewScanner(os.Stdin)
			// if input.Scan() && input.Text() == "y" {
			// 	container.CleanContainers()
			// }
		}
	}

}

func cleanDangling() {
	images := getCacheImages([]string{"images", "-q", "--filter", "dangling=true"})
	if images != nil {
		fmt.Println("Dangling image will be removed from local cache")
		cleanImage(images)
	}
}

func getReserveImageFromFile(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var items []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		items = append(items, scanner.Text())
	}
	return items, scanner.Err()
}

func getCacheImages(args []string) []string {
	out, err := exec.Command("docker", args...).Output()
	if err != nil {
		fmt.Println("Get images error")
		return nil
	}
	images := strings.Fields(strings.TrimSpace(strings.Replace(string(out), "\n", " ", -1)))

	if len(images) == 0 {
		return nil
	}

	encountered := map[string]bool{}
	result := []string{}
	for v := range images {
		if encountered[images[v]] != true {
			encountered[images[v]] = true
			result = append(result, images[v])
		}
	}

	return result
}
