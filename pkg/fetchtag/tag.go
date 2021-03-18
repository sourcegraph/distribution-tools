package fetchtag

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"time"
)

// Largely copied from https://github.com/slimsag/update-docker-tags/blob/201509b910e3a78948ec8951c7b7fb33dd711931/update-docker-tags.go
// Thank you @slimsag

const (
	// TODO Get these from config
	defaultRegistry = "gcr.io"
	defaultProject  = "sourcegraph-dev"
)

type Fetcher struct {
	Name     string // ${repo}/${image-Name}
	Registry string

	AuthToken string
}

// Effectively the same as:
//
// 	$ export token=$(curl -s "https://auth.docker.io/token?service=registry.docker.io&scope=repository:sourcegraph/server:pull" | jq -r .token)
//
func fetchAuthToken(repositoryName string) (string, error) {
	resp, err := http.Get(fmt.Sprintf("https://auth.docker.io/token?service=registry.docker.io&scope=Fetcher:%s:pull", repositoryName))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	result := struct {
		Token string
	}{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return "", err
	}
	return result.Token, nil
}

// Effectively the same as:
//
// 	$ curl -H "Authorization: Bearer $token" https://index.docker.io/v2/sourcegraph/server/tags/list
// or curl -H "Authorization: Bearer $token" https://us.gcr.io/v2/sourcegraph-dev/chrome/tags/list
// see https://docs.docker.com/registry/spec/api/#listing-image-tags
func (r *Fetcher) fetchAllTags() ([]string, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("https://%s/v2/%s/tags/list", r.Registry, r.Name), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+r.AuthToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	result := struct {
		Tags []string
	}{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}
	return result.Tags, nil
}

type tagListResp struct {
	Manifest map[string]manifest
	Name     string
	Tags     []string
}

type manifest struct {
	ImageSizeBytes string
	LayerId        string
	MediaType      string
	Tag            []string
	TimeCreatedMs  string
	TimeUploadedMs string
}

// TODO: ptr return
func (r *Fetcher) imageTagList() (tagListResp, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("https://%s/v2/%s/tags/list", r.Registry, r.Name), nil)
	if err != nil {
		return tagListResp{}, err
	}
	req.Header.Set("Authorization", "Bearer "+r.AuthToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return tagListResp{}, err
	}
	defer resp.Body.Close()

	result := tagListResp{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	return result, err
}

// Effectively the same as:
//
//  $ curl -s -D - -H "Authorization: Bearer $token" -H "Accept: application/vnd.docker.distribution.manifest.v2+json" https://index.docker.io/v2/sourcegraph/server/manifests/3.12.1 | grep Docker-Content-Digest
//
func (r *Fetcher) FetchImageDigest(tag string) (string, error) {
	req, err := http.NewRequest("GET",
		fmt.Sprintf("https://%s/v2/%s/manifests/%s", r.Registry, r.Name, tag), nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", r.AuthToken))
	req.Header.Set("Accept", "application/vnd.docker.distribution.manifest.v2+json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	return resp.Header.Get("Docker-Content-Digest"), nil
}

func unixTimeMs() int64 {
	now := time.Now()
	nanos := now.UnixNano()

	return nanos / 1000000
}

func (r *Fetcher) FindLatestImageByTime() (tag, digest string, err error) {
	tagResp, err := r.imageTagList()
	if err != nil {
		return "", "", err
	}
	now := unixTimeMs()

	minDelta := int64(math.MaxInt64)
	var delta int64
	for k, v := range tagResp.Manifest {
		timeUploadedMs, err := strconv.ParseInt(v.TimeUploadedMs, 10, 64)
		if err != nil {
			return "", "", err
		}

		delta = now - timeUploadedMs

		if delta < minDelta {
			if len(v.Tag) > 0 {
				// just grab any tag
				tag = v.Tag[0]
				minDelta = delta
				digest = k
			} else {
				println("possible latest image lacks a tag")
				continue
			}
		}
	}
	return tag, digest, err
}
