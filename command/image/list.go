package image

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/stararks/dct/repo"

	"github.com/spf13/cobra"
)

var registry string
var defaultPort = "5000"

var protocal = "http://"

func CmdlistImageCommand() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "listi REGISTRY",
		Short: "List image from a registry",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				fmt.Println("Miss registry")
				return
			}
			colon := strings.LastIndex(args[0], ":")
			if colon == -1 {
				registry = args[0]
			} else {
				registry = args[0][:colon]
				defaultPort = args[0][colon+1:]
			}
			listImages(registry)
		},
	}
	return cmd
}

func listImages(registry string) {
	var image repo.Images
	var catalog repo.Repository
	var err error

	url := fmt.Sprintf("%s:%s/v2/", registry, defaultPort)

	if err = get(url+"_catalog", &catalog); err != nil {
		fmt.Println(err)
		return
	}

	for _, item := range catalog.Repos {
		if err = get(url+*item+"/tags/list", &image); err != nil {
			fmt.Println(err)
			continue
		}
		for _, item := range image.ImangeTags {
			fmt.Printf("%s:%s/%s:%s\n", registry, defaultPort, image.ImageName, *item)
		}
	}
}

func get(url string, result interface{}) error {
	resp, err := http.Get(protocal + url)
	if err != nil {
		log.Fatal(err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		protocal = "https://"
		tr := &http.Transport{
			MaxIdleConns:       10,
			IdleConnTimeout:    30 * time.Second,
			DisableCompression: true,

			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		client := &http.Client{Transport: tr}
		resp, err = client.Get(protocal + url)
		if err != nil {
			log.Fatal(err)
			return err
		}

	}
	if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
		return err
	}
	return nil
}
