package main

import (
	"github.com/allanwang/git-issues/gi"
	"fmt"
)

func main() {
	tui := gi.CreateIssueListView()
	data := tui.Fetch()
	for _, issue := range data {
		fmt.Println(issue.Body)
	}
	tui.Run()
}
