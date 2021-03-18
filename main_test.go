package main

import (
	"os"
	"testing"
)

func Test_app(t *testing.T) {
	token := os.Getenv("IMAGE_FETCH_TOKEN")
	if token == "" {
		t.Fatal("must set token")
	}

	tests := []struct {
		name    string
		image   string
		token   string
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			"base",
			"cron",
			token,
			"us.gcr.io/sourcegraph-dev/cron:e3de51c21e9b069235f382ca3d9063943af31d39@sha256:e17738da3fbb408544fe8a25106733f2e8cb28a5d6e70818d3db9b9fc8cac3ce",
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := app(tt.image, tt.token)
			if (err != nil) != tt.wantErr {
				t.Errorf("app() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("app() got = %v, want %v", got, tt.want)
			}
		})
	}
}
