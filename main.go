package main

import (
	"github.com/allanwang/git-issues/gi"
)

func main() {
	tui := gi.CreateIssueListView()
	tui.Fetch()
	tui.Run()
}
