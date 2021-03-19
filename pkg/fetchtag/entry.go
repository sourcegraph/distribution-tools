package fetchtag

import (
	"fmt"
	"strings"

	"github.com/sourcegraph/distribution-tools/pkg/dockerimg"
)

var DebugMode = false

func Transform(image string) (dockerimg.ImageReference, error) {
	debugPrint("looking for image " + image)

	imgRef := dockerimg.ImageReference{}

	r, err := dockerimg.Parse(image)
	if err != nil {
		return dockerimg.ImageReference{}, err
	}
	named, ok := r.(dockerimg.Named)
	if !ok {
		return dockerimg.ImageReference{}, fmt.Errorf("no name")
	}

	tagged, ok := r.(dockerimg.Tagged)
	if !ok {
		debugPrint("no tag")
	} else {
		// naming
		imgRef.Version = tagged.Tag()
	}

	path := dockerimg.Path(named)
	// expect sourcegraph-dev/$imageName
	if strings.Contains(path, "/") {
		imgRef.Name = path
	} else {
		debugPrint("no provided project")
		imgRef.Name = fmt.Sprintf("%s/%s", defaultProject, path)
	}

	d := dockerimg.Domain(named)
	if d == "" {
		debugPrint("no domain")
		d = defaultRegistry
	}
	imgRef.Registry = d

	imgRef.Key = imgRef.Name

	// we only allow sha256... do we need both crypto libs ?
	if digested, ok := r.(dockerimg.Digested); ok {
		imgRef.Sha256 = strings.TrimPrefix(digested.Digest().String(),
			"sha256:")
	}
	return imgRef, nil
}

func debugPrint(s string) {
	if DebugMode {
		println(s)
	}
}
