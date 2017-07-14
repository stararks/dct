package image

import (
	"fmt"
	"io/ioutil"
	"os/exec"
	"path"
	"path/filepath"
	"sync"

	"github.com/spf13/cobra"
)

var wg sync.WaitGroup
var sema = make(chan struct{}, 3)
var ch = make(chan string)
var recursive bool

func CmdloadImageCommand() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "loadi PATH",
		Short: "Load images from a directory",
		Run: func(cmd *cobra.Command, args []string) {
			dir := "."
			if len(args) != 0 {
				dir = args[0]
			}
			loadImages(dir)
		},
	}

	cmd.Flags().BoolVarP(&recursive, "recursive", "r", false, "load images recursively")

	return cmd
}

func loadImages(dir string) {
	dir = path.Clean(dir)
	entries, err := ioutil.ReadDir(dir)
	if err != nil {
		fmt.Printf("Load images from %s error: %s", dir, err)
	}
	for _, entry := range entries {
		abspath := filepath.Join(dir, entry.Name())
		if entry.IsDir() {
			if recursive {
				loadImages(abspath)
			}
		} else {
			wg.Add(1)
			go loadImage(abspath, &wg, ch)
		}
	}
	go func() {
		wg.Wait()
		close(ch)
	}()
	for out := range ch {
		fmt.Println(out)
	}

}

func loadImage(abspath string, wg *sync.WaitGroup, ch chan<- string) {
	defer wg.Done()
	sema <- struct{}{}
	defer func() { <-sema }()

	f := path.Base(abspath)

	cmd := exec.Command("docker", "load", "-i", abspath)
	if err := cmd.Run(); err != nil {
		ch <- fmt.Sprintf("Loading image %s ... error", f)
		return
	}
	ch <- fmt.Sprintf("Loading image %s ... done", f)
}
