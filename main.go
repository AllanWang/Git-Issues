package main

import (
	"./gi"
	"os"
	"fmt"
)

func main() {
	r := gi.GetGitWd()
	if r == nil {
		fmt.Println("Could not get git instance")
		os.Exit(-1)
	}
	name, _ := r.GetProjectName()
	fmt.Printf("Welcome to %s\n", name)
	head, _ := r.Head()
	fmt.Println(head.String())
}
