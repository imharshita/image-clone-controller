package images

import (
	"os"
	"strings"

	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1/remote"
)

var registry string = os.Getenv("REGISTRY")

func rename(name string) string {
	var img string
	image := strings.Split(name, "/")
	if len(image) == 2 {
		img = image[1]
	} else {
		img = image[0]
	}
	newName := registry + "/" + img
	return newName
}

// Process public image to retag and push to private registry
func Process(imgName string) (string, error) {
	ref, err := name.ParseReference(imgName)
	if err != nil {
		return "", err
	}
	username := os.Getenv("USERNAME")
	password := os.Getenv("PASSWORD")
	if len(username) == 0 && len(password) == 0 {
		return "", err
	}
	auth := authn.AuthConfig{
		Username: username,
		Password: password,
	}
	authenticator := authn.FromConfig(auth)
	opt := remote.WithAuth(authenticator)
	img, err := remote.Image(ref, opt)
	if err != nil {
		return "", err
	}
	newName := rename(imgName)
	newRef, err := name.ParseReference(newName)
	if err != nil {
		return "", err
	}
	if err := remote.Write(newRef, img, opt); err != nil {
		return "", err
	}
	return newName, nil
}
