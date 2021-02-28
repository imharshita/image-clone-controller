package images

import (
	"log"
	"os"
	"strings"

	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/crane"
	"github.com/google/go-containerregistry/pkg/name"
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

func retag(imgName string) (name.Tag, error) {
	tag, err := name.NewTag(imgName)
	if err != nil {
		return name.Tag{}, err
	}
	return tag, nil
}

// Process public image to retag and push to private registry
func Process(imgName string) (string, error) {
	username := os.Getenv("USERNAME")
	password := os.Getenv("PASSWORD")
	auth := authn.AuthConfig{
		Username: username,
		Password: password,
	}
	authenticator := authn.FromConfig(auth)
	opt := crane.WithAuth(authenticator)
	img, err := crane.Pull(imgName)
	if err != nil {
		log.Fatal(err)
		return "", err
	}
	newName := rename(imgName)
	tag, err := retag(newName)
	if err != nil {
		log.Fatal(err)
		return "", err
	}
	if err := crane.Push(img, tag.String(), opt); err != nil {
		log.Fatal(err)
		return "", err
	}
	return newName, nil
}
