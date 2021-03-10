package main

import (
	"bytes"
	"testing"
)

func Test_app(t *testing.T) {

	tests := []struct {
		name       string
		args       []string
		wantStdout string
		wantErr    bool
	}{
		// TODO: Add test cases.
		{
			"base",
			[]string{"srcimage", "--image", "us.gcr.io/sourcegraph-dev/buildkite-agent"},
			"",
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stdout := &bytes.Buffer{}
			err := app(tt.args, stdout)
			if (err != nil) != tt.wantErr {
				t.Errorf("app() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotStdout := stdout.String(); gotStdout != tt.wantStdout {
				t.Errorf("app() gotStdout = %v, want %v", gotStdout, tt.wantStdout)
			}
		})
	}
}
