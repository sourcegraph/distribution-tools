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
	flag.BoolVarP(&help, "help", "h", false, "print help")

	flag.Parse()

	if help || image == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}
	token := os.Getenv("IMAGE_FETCH_TOKEN")
	if token == "" {
		fmt.Println("Use the env var IMAGE_FETCH_TOKEN to provide the registry access token ")
		os.Exit(1)
	}

	dbg := os.Getenv("DEBUG_MODE")
	if dbg != "" {
		fetchtag.DebugMode = true
	}

	result, err := app(image, token)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
	fmt.Println(result)
}
func app(image, token string) (string, error) {
	imgRef, err := fetchtag.Transform(image)
	if err != nil {
		return "", err
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
			return "", err
		}
		imgRef.Version = tag
		imgRef.Sha256 = digest
	} else if imgRef.Sha256 == "" {
		// fetch only digest (faster)
		digest, err := fetcher.FetchImageDigest(imgRef.Version)
		if err != nil {
			return "", err
		}
		imgRef.Sha256 = digest
	}
	return fmt.Sprintf("%s/%s:%s@%s", imgRef.Registry, imgRef.Name, imgRef.Version, imgRef.Sha256), nil
}
