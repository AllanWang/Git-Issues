package gi

import (
	"fmt"
	"os"
	"gopkg.in/src-d/go-git.v4"
	"errors"
	"strings"
	"gopkg.in/src-d/go-billy.v4/helper/chroot"
	"gopkg.in/src-d/go-git.v4/storage/filesystem"
	"path"
)

type Project interface {
	GetProjectName() (string, error)
	GetRepository() *git.Repository
}

type project struct {
	*git.Repository
}

var (
	ErrAmbiguousProjectName = errors.New("ambiguous project name; multiple remotes found")
	ErrNoRemote             = errors.New("no remote found")
)

// Instantiate Project from working directory
func GetProjectWd() Project {
	dir, err := os.Getwd()
	if err != nil {
		fmt.Printf("Could not open working directory")
		return nil
	}
	return GetProject(dir)
}

// Instantiate Project from path
func GetProject(path string) Project {
	r, err := git.PlainOpen(path)
	if err == git.ErrRepositoryNotExists {
		fmt.Printf("No repo found at %s\n", path)
		return nil
	}
	if err != nil {
		fmt.Printf("An error occurred: %s\n", err.Error())
		return nil
	}
	return project{r}
}

func (g project) GetRepository() *git.Repository {
	return g.Repository
}

// Attempts to retrieve the git parent directory
// Falls back to prettified remote name
func (g project) GetProjectName() (string, error) {
	// Try to grab the repository Storer
	s, ok := g.Storer.(*filesystem.Storage)
	if !ok {
		return g.getRemoteName()
	}

	// Try to get the underlying billy.Filesystem
	fs, ok := s.Filesystem().(*chroot.ChrootHelper)
	if !ok {
		return g.getRemoteName()
	}
	name := path.Base(path.Dir(fs.Root()))
	return name, nil
}

// Attempts to retrieve a remote name
func (g project) getRemoteName() (string, error) {
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

// Trim url to get last segment before extensions
func nameFromUrl(url string) string {
	slash := strings.LastIndex(url, "/")
	if slash > 0 {
		url = url[slash+1:]
	}
	dot := strings.Index(url, ".")
	if dot > 0 {
		url = url[:dot]
	}
	return url
}
