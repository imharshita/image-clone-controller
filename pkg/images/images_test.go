package images_test

import (
	"testing"

	"github.com/imharshita/image-clone-controller/pkg/images"
)

func TestProcess(t *testing.T) {
	result, _ := images.Process("xy/nginx:1.14")
	if result != "backupregistry/nginx:1.14" {
		t.Errorf("Process(\"xy/nginx:1.14\") failed, expected %v, got %v", "backupregistry/nginx:1.14", result)
	}
}
