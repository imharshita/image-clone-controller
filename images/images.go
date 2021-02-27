package images

import (
	"fmt"

	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/crane"
	"github.com/google/go-containerregistry/pkg/name"
)

var privateRegistry string = "backupregistry"

// func rename(name string) string {
// 	strings.Contains(name, "/")
// 	image := strings.Split(name, "/")

// 	image := strings.Split(name, ":")
// 	img, version := image[0], image[1]
// 	newName := privateRegistry + "/" + img + ":" + version
// 	return newName
// }

func retag(imgName string) (name.Tag, error) {
	tag, err := name.NewTag(imgName)
	if err != nil {
		return name.Tag{}, err
	}
	return tag, nil
}

func Process(imgName string) (string, error) {
	auth := authn.AuthConfig{
		Username: "backupregistry",
		Password: "mydockerimages",
	}
	authenticator := authn.FromConfig(auth)
	opt := crane.WithAuth(authenticator)
	img, err := crane.Pull(imgName, opt)
	fmt.Println("pulled image")
	if err != nil {
		return "", err
	}
	newName := "backupregistry/test:v1"
	//newName := rename(imgName)
	fmt.Println("rename", newName)
	tag, err := retag(newName)
	if err != nil {
		return "", err
	}

	if err := crane.Push(img, tag.String(), opt); err != nil {
		return "", err
	}
	fmt.Println("pushed image")
	return newName, nil
}
