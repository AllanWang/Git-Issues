package gi

import (
	"github.com/rivo/tview"
	"github.com/gdamore/tcell"
	"github.com/google/go-github/github"
	"os/exec"
	"os"
	"fmt"
)

type Tui interface {
	Gi
	Log(msg string, a ... interface{})
	Fetch() []*github.Issue
	Run()
	Stop()
}

type tuiBase struct {
	Gi
	*tview.Application
	logs *tview.TextView
}

type registration struct {
	tuiBase
	info  *tview.TextView
	link  *tview.Button
	input *tview.InputField
}

type issueList struct {
	tuiBase
	issues  *tview.List
	content *tview.List
}

func createTuiBase(config *Config) tuiBase {
	return tuiBase{GetGi(config), tview.NewApplication(), tview.NewTextView(),}
}

func (tui tuiBase) eventCapture(event *tcell.EventKey) *tcell.EventKey {
	switch event.Key() {
	case tcell.KeyRune:
		switch event.Rune() {
		case 'q':
			tui.Stop()
			os.Exit(0)
		default:
			tui.Log(string(event.Rune()))
		}
	default:
		tui.Log(tcell.KeyNames[event.Key()])
	}

	return event
}

func getGithubToken(base tuiBase) {
	tui := registration{base, tview.NewTextView(),
		tview.NewButton("Link"), tview.NewInputField()}
	tui.link.SetSelectedFunc(func() {
		exec.Command("xdg-open", "https://github.com/settings/tokens/new")
	})
	tui.info.SetWrap(true).SetWordWrap(true)

	if base.GetConfig().GithubToken == "" {
		tui.info.SetText("No token found; please paste a new one below")
	} else {
		tui.info.SetText("Invalid token found; please paste a new one below")
	}
	tui.input.SetText("Token").SetAcceptanceFunc(func(textToCheck string, lastChar rune) bool {
		if base.SetGithubToken(textToCheck) {
			tui.Stop()
		}
		return true
	})
	flex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(tui.info, 0, 1, false).
		AddItem(tui.input, 0, 1, false).
		AddItem(tui.link, 0, 1, true).
		AddItem(tui.logs, 1, 1, false)
	tui.SetRoot(flex, true).SetInputCapture(tui.eventCapture).Run()
}

func CreateIssueListView() Tui {
	base := createTuiBase(GetConfig())
	if base.GetClient() == nil {
		getGithubToken(base)
	}
	tui := issueList{base, tview.NewList(), tview.NewList()}
	tui.issues.SetBorder(true).SetTitle("Issues")
	tui.content.SetBorder(true).SetTitle("Content")
	tui.logs.SetTextColor(tcell.ColorWhite).SetBorder(false).SetBorderPadding(0, 0, 1, 1)

	flex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(tview.NewFlex().SetDirection(tview.FlexColumn).
		AddItem(tui.issues, 0, 1, true).
		AddItem(tui.content, 0, 3, false), 0, 1, false).
		AddItem(tui.logs, 1, 1, false)

	tui.SetRoot(flex, true).SetInputCapture(tui.eventCapture).SetFocus(flex)
	return tui
}

/*
 * -------------------------------------
 * Base
 * -------------------------------------
 */

func (tui tuiBase) Log(format string, a ... interface{}) {
	if len(a) == 0 {
		tui.logs.SetText(format)
	} else {
		tui.logs.SetText(fmt.Sprintf(format, a))
	}
}

func (tui tuiBase) Run() {
	err := tui.Application.Run()
	if err != nil {
		panic(err)
	}
}

/*
 * -------------------------------------
 * IssueList
 * -------------------------------------
 */

func (tui issueList) Fetch() []*github.Issue {
	project, err := tui.GetProjectName()
	if err != nil {
		tui.Log("Error: %s", err)
		return []*github.Issue{}
	}
	client := tui.GetClient()
	if client == nil {
		tui.Log("Nil client")
		return []*github.Issue{}
	}
	issues, _, err := client.Issues.
		ListByRepo(tui.Ctx(), *tui.GetUserName(), project, nil)
	if err != nil {
		tui.Log("Error %s: %s", project, err)
		return []*github.Issue{}
	}
	tui.issues.Clear()
	for i, issue := range issues {
		var key rune
		if i < 10 {
			key = rune(i)
		} else {
			key = 0
		}
		tui.issues.AddItem(*issue.Title, *issue.User.Name, key, func() {
			tui.showIssue(issue)
		})
	}
	tui.Log("Done %d", len(issues))
	return issues
}

func(tui issueList) showIssue(issue *github.Issue) {

}