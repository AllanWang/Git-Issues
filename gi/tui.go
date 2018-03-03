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
	Log(msg string)
	Fetch() []*github.Issue
	Run()
	Stop()
}

type tuiBase struct {
	Gi
	*tview.Application
	logs *tview.InputField
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
	content *tview.TextView
}

func createTuiBase(config *Config) tuiBase {
	return tuiBase{GetGi(config), tview.NewApplication(), tview.NewInputField(),}
}

func (tui tuiBase) eventCapture(event *tcell.EventKey) *tcell.EventKey {
	if event.Key() == tcell.KeyRune {
		if !tui.logs.HasFocus() {
			switch event.Rune() {
			case 'i':
				tui.SetFocus(tui.logs)
				tui.logs.SetText("")
				//tui.logs.Clear()
				tui.Draw()
				return nil
			case 'q':
				tui.Stop()
			default:
				tui.Log(string(event.Rune()))
			}
		}
	} else {
		switch event.Key() {
		case tcell.KeyCtrlC, tcell.KeyCtrlZ:
			tui.Stop()
		default:
			if !tui.logs.HasFocus() {
				tui.Log(tcell.KeyNames[event.Key()])
			}
		}
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
			tui.Application.Stop()
		}
		return true
	})
	flex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(tui.info, 0, 1, false).
		AddItem(tui.input, 0, 1, false).
		AddItem(tui.link, 0, 1, true).
		AddItem(tui.logs, 1, 1, false)
	flex.SetBackgroundColor(tcell.ColorBlack)
	tui.SetRoot(flex, true).SetInputCapture(tui.eventCapture).Run()
}

func theme(box *tview.Box) *tview.Box {
	box.SetBackgroundColor(tcell.ColorBlack).
		SetBorderColor(tcell.ColorGray).
		SetTitleColor(tcell.ColorWhite)
	return box
}

func themeText(view *tview.TextView) *tview.TextView {
	view.SetTextColor(tcell.ColorWhite).
		SetTitleColor(tcell.ColorWhite)
	return view
}

func themeInput(view *tview.InputField) *tview.InputField {
	view.SetFieldTextColor(tcell.ColorWhite).
		SetBackgroundColor(tcell.ColorBlack).
		SetTitleColor(tcell.ColorWhite)
	return view
}

func themeList(l *tview.List) *tview.List {
	l.SetMainTextColor(tcell.ColorWhite).
		SetSecondaryTextColor(tcell.ColorGray).
		SetSelectedBackgroundColor(tcell.ColorDarkGray).
		SetSelectedTextColor(tcell.ColorWhite).
		SetTitleColor(tcell.ColorWhite).
		SetBackgroundColor(tcell.ColorBlack)
	return l
}

func CreateIssueListView() Tui {
	base := createTuiBase(GetConfig())
	if base.GetClient() == nil {
		getGithubToken(base)
	}
	tui := issueList{base, tview.NewList(), tview.NewTextView()}
	tui.issues.ShowSecondaryText(false).SetBorder(true).SetTitle("Issues").
		SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyLeft:
			return nil
		case tcell.KeyRight:
			tui.SetFocus(tui.content)
			return nil
		}
		return event
	})
	tui.content.SetScrollable(true).SetBorder(true).
		SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyLeft:
			tui.SetFocus(tui.issues)
			return nil
		case tcell.KeyRight:
			return nil
		}
		return event
	})
	tui.logs.SetBorder(false).SetBorderPadding(0, 0, 1, 1).
		SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEsc:
			tui.SetFocus(tui.issues)
			return nil
		}
		return event
	})

	themeList(tui.issues)
	themeText(tui.content)
	themeInput(tui.logs)

	flex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(tview.NewFlex().SetDirection(tview.FlexColumn).
		AddItem(tui.issues, 0, 1, true).
		AddItem(tui.content, 0, 3, false), 0, 1, false).
		AddItem(tui.logs, 1, 1, false)

	theme(flex.Box)

	tui.SetRoot(flex, true).SetInputCapture(tui.eventCapture).SetFocus(tui.issues)
	return tui
}

/*
 * -------------------------------------
 * Base
 * -------------------------------------
 */

func (tui tuiBase) Log(format string) {
	tui.logs.SetText(format)
}

func (tui tuiBase) Run() {
	err := tui.Application.Run()
	if err != nil {
		panic(err)
	}
}

func (tui tuiBase) Stop() {
	tui.Application.Stop()
	os.Exit(0)
}

/*
 * -------------------------------------
 * IssueList
 * -------------------------------------
 */

func (tui issueList) Fetch() []*github.Issue {
	project := tui.GetProjectName()
	if project == nil {
		tui.Log("No project found")
		return []*github.Issue{}
	}
	client := tui.GetClient()
	if client == nil {
		tui.Log("Nil client")
		return []*github.Issue{}
	}
	issues, _, err := client.Issues.
		ListByRepo(tui.Ctx(), *tui.GetUserName(), *project, nil)
	if err != nil {
		tui.Log(fmt.Sprintf("Error %s: %s", *project, err))
		return []*github.Issue{}
	}
	tui.issues.Clear()
	if len(issues) > 0 {
		for _, issue := range issues {
			tui.issues.AddItem(*issue.Title, "", 0, nil)
		}
	} else {
		tui.issues.AddItem("No issues found", "", 0, nil)
	}
	tui.issues.SetChangedFunc(func(i int, _ string, _ string, _ rune) {
		tui.showIssue(issues[i])
	})
	tui.Log(fmt.Sprintf("Done %d %s", len(issues), *tui.GetProjectName()))
	return issues
}

func (tui issueList) showIssue(issue *github.Issue) {
	tui.content.SetText(*issue.Body)
}
