package main

import (
	"fmt"
	"os"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
)

func Hello() {
	fmt.Println("helo")
}

func main() {
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	fmt.Println(dir)

	r := GetGit(dir)
	if r == nil {
		panic("No git found")
	}
	t, _ := r.Tags()
	t.ForEach(func(reference *plumbing.Reference) error {
		fmt.Println(reference.Name())
		return nil
	})
}

func GetGit(path string) *git.Repository {
	r, err := git.PlainOpen(path)
	if err == git.ErrRepositoryNotExists {
		fmt.Printf("No repo found at %s\n", path)
		return nil
	}
	if err != nil {
		fmt.Printf("An error occurred: %s\n", err.Error())
		return nil
	}
	return r
}
