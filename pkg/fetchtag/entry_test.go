package fetchtag

import (
	"reflect"
	"testing"

	"github.com/sourcegraph/distribution-tools/pkg/dockerimg"
)

func TestTransform(t *testing.T) {

	tests := []struct {
		name    string
		image   string
		want    dockerimg.ImageReference
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			"base",
			"frontend",
			dockerimg.ImageReference{
				Registry: defaultRegistry,
				Name:     defaultProject + "/" + "frontend",
				Version:  "",
				Sha256:   "",
				Key:      defaultProject + "/" + "frontend",
			},
			false,
		},
		{
			"base 2",
			"repo-updater",
			dockerimg.ImageReference{
				Registry: defaultRegistry,
				Name:     defaultProject + "/" + "repo-updater",
				Version:  "",
				Sha256:   "",
				Key:      defaultProject + "/" + "repo-updater",
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Transform(tt.image)
			if (err != nil) != tt.wantErr {
				t.Errorf("Transform() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Transform() got = %v, want %v", got, tt.want)
			}
		})
	}
}
