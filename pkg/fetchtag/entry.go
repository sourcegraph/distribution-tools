package fetchtag

import (
	"fmt"
	"strings"

	"github.com/sourcegraph/distribution-tools/pkg/dockerimg"
)

func Transform(image string) (dockerimg.ImageReference, error) {

	println("looking for image", image)

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
		println("no tag")
	} else {
		// naming
		imgRef.Version = tagged.Tag()
	}

	path := dockerimg.Path(named)
	// expect sourcegraph-dev/$imageName
	if strings.Contains(path, "/") {
		imgRef.Name = path
	} else {
		println("no provided project")
		imgRef.Name = fmt.Sprintf("%s/%s", defaultProject, path)
	}

	d := dockerimg.Domain(named)
	if d == "" {
		println("no domain")
		d = defaultRegistry
	}
	imgRef.Registry = d

	//if dockerimg.DomainIsNotHostName(d) {
	//	imgRef.Name = fmt.Sprintf("%s/%s", d, path)
	//	imgRef.Registry = ""
	//}

	// ?
	imgRef.Key = imgRef.Name

	// we only allow sha256... do we need both crypto libs ?
	if digested, ok := r.(dockerimg.Digested); ok {
		imgRef.Sha256 = strings.TrimPrefix(digested.Digest().String(),
			"sha256:")
	}

	return imgRef, nil

}
