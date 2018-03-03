package gi

import (
	"os"
	"gopkg.in/src-d/go-git.v4"
	"strings"
)

type Project interface {
	GetProjectName() *string
}

type project struct {
	*git.Repository
	name *string
}

// Instantiate Project from working directory, or the supplied path
func GetCurrentProject() Project {
	if len(os.Args) > 1 {
		_, err := os.Stat(os.Args[1])
		if !os.IsNotExist(err) {
			return GetProject(os.Args[1])
		}
	}
	dir, err := os.Getwd()
	if err != nil {
		return nil
	}
	return GetProject(dir)
}

// Instantiate Project from path
func GetProject(path string) Project {
	r, _ := git.PlainOpen(path)
	if r == nil {
		return nil
	}
	remotes, _ := r.Remotes()
	if len(remotes) == 0 {
		return nil
	}
	for _, url := range remotes[0].Config().URLs {
		if strings.Contains(url, "github") {
			name := GetNameFromUrl(url)
			if name != nil {
				return project{r, name}
			}
		}
	}
	return nil
}

func GetNameFromUrl(url string) *string {
	slash := strings.LastIndex(url, "/")
	if slash > 0 {
		url = url[slash+1:]
	}
	dot := strings.Index(url, ".")
	if dot > 0 {
		url = url[:dot]
	}
	return &url
}

func (p project) GetProjectName() *string {
	return p.name
}
