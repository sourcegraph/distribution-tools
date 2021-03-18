package main

import (
	"fmt"
	"os"

	"github.com/sourcegraph/distribution-tools/pkg/fetchtag"

	flag "github.com/spf13/pflag"
)

func main() {

	var image string
	var help bool
	flag.StringVarP(&image, "image", "i", "", "image to update, in registry/repo/tag format")
	flag.BoolVarP(&help, "help", "h", false, "use IMAGE_FETCH_TOKEN env var to provide registry access")

	flag.Parse()

	token := os.Getenv("IMAGE_FETCH_TOKEN")
	if token == "" {
		fmt.Println("Use the env var IMAGE_FETCH_TOKEN to provide the registry access token ")
		os.Exit(1)
	}

	if help || image == "" {
		flag.PrintDefaults()
		os.Exit(0)
	}

	if err := app(image, token); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}

}

func app(image, token string) error {

	imgRef, err := fetchtag.Transform(image)
	if err != nil {
		return err
	}

	fetcher := fetchtag.Fetcher{
		Name:      imgRef.Name,
		Registry:  imgRef.Registry,
		AuthToken: token,
	}
	// missing tag
	if imgRef.Version == "" {
		// create repo struct from imageRef
		tag, digest, err := fetcher.FindLatestImageByTime()
		if err != nil {
			return err
		}
		imgRef.Version = tag
		imgRef.Sha256 = digest
	} else if imgRef.Sha256 == "" {
		// fetch only digest (faster)
		digest, err := fetcher.FetchImageDigest(imgRef.Version)
		if err != nil {
			return err
		}
		imgRef.Sha256 = digest
	}
	fmt.Printf("%s/%s:%s@%s", imgRef.Registry, imgRef.Name, imgRef.Version, imgRef.Sha256)

	return nil
}
