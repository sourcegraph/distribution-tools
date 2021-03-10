package fetchtag

import (
	"fmt"
	"os"
	"testing"
)

func TestFetcher_imageTagList(t *testing.T) {

	token := os.Getenv("IMAGE_FETCH_TOKEN")
	if token == "" {
		t.Fatal("must set token")
	}

	type fields struct {
		Name      string
		Registry  string
		AuthToken string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			"base",
			fields{"sourcegraph-dev/frontend",
				"us.gcr.io",
				token},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Fetcher{
				Name:      tt.fields.Name,
				Registry:  tt.fields.Registry,
				AuthToken: tt.fields.AuthToken,
			}
			got, err := r.imageTagList()
			if (err != nil) != tt.wantErr {
				t.Errorf("imageTagList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			fmt.Println(got)
			//if !reflect.DeepEqual(got, tt.want) {
			//	t.Errorf("imageTagList() got = %v, want %v", got, tt.want)
			//}
		})
	}
}

func TestFetcher_FindLatestImage(t *testing.T) {
	token := os.Getenv("IMAGE_FETCH_TOKEN")
	if token == "" {
		t.Fatal("must set token")
	}

	type fields struct {
		Name      string
		Registry  string
		AuthToken string
	}
	tests := []struct {
		name       string
		fields     fields
		wantTag    string
		wantDigest string
		wantErr    bool
	}{
		// TODO: Add test cases.
		{
			"base",
			fields{"sourcegraph-dev/cron",
				"us.gcr.io",
				token},
			"e3de51c21e9b069235f382ca3d9063943af31d39",
			"sha256:e17738da3fbb408544fe8a25106733f2e8cb28a5d6e70818d3db9b9fc8cac3ce",
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Fetcher{
				Name:      tt.fields.Name,
				Registry:  tt.fields.Registry,
				AuthToken: tt.fields.AuthToken,
			}
			gotTag, gotDigest, err := r.FindLatestImageByTime()
			if (err != nil) != tt.wantErr {
				t.Errorf("FindLatestImage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotTag != tt.wantTag {
				t.Errorf("FindLatestImage() gotTag = %v, want %v", gotTag, tt.wantTag)
			}
			if gotDigest != tt.wantDigest {
				t.Errorf("FindLatestImage() gotDigest = %v, want %v", gotDigest, tt.wantDigest)
			}
		})
	}
}
