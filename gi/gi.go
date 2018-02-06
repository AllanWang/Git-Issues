package gi

import (
	"fmt"
	"os"
	"gopkg.in/src-d/go-git.v4"
	"errors"
	"strings"
)

type Gi struct {
	git.Repository
}

var (
	ErrAmbiguousProjectName = errors.New("ambiguous project name; multiple remotes found")
	ErrNoRemote             = errors.New("no remote found")
)

func GetGitWd() *Gi {
	dir, err := os.Getwd()
	if err != nil {
		fmt.Printf("Could not open working directory")
		return nil
	}
	return GetGit(dir)
}

func GetGit(path string) *Gi {
	r, err := git.PlainOpen(path)
	if err == git.ErrRepositoryNotExists {
		fmt.Printf("No repo found at %s\n", path)
		return nil
	}
	if err != nil {
		fmt.Printf("An error occurred: %s\n", err.Error())
		return nil
	}
	return &Gi{*r}
}


func (g Gi) GetProjectName() (string, error) {
	remotes, err := g.Remotes()
	if err != nil {
		return "", err
	}
	switch len(remotes) {
	case 0:
		return "", ErrNoRemote
	case 1:
		remote := remotes[0].Config().URLs[0]
		return nameFromUrl(remote), nil
	default:
		return "", ErrAmbiguousProjectName
	}
}

func nameFromUrl(url string) string {
	slash := strings.LastIndex(url, "/")
	if slash > 0 {
		url = url[slash + 1:]
	}
	dot := strings.Index(url, ".")
	if dot > 0 {
		url = url[:dot]
	}
	return url
}
