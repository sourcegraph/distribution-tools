package dockerimg

import (
	// must be imported for the go-digest pkg to work
	// these are the supported images digests we can use
	_ "crypto/sha256"
)

type ImageReference struct {
	Registry string
	Name     string
	Version  string
	Sha256   string
	Key      string
}
