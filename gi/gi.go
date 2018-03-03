package gi

import (
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
	"context"
)

type Gi interface {
	Project
	SetGithubToken(token string) bool
	GetUserName() *string
	GetClient() *github.Client
	GetConfig() *Config
	Ctx() context.Context
}

type gi struct {
	Project
	user         *github.User
	config       *Config
	githubClient *github.Client
	ctx      context.Context
}

func GetGi(config *Config) Gi {
	gi := gi{GetCurrentProject(), nil, config,
		nil, context.Background()}
	gi.updateGithubClient()
	return gi
}

func (gi *gi) updateGithubClient() bool {
	token := gi.config.GithubToken
	if token == "" {
		return false
	}
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(gi.ctx, ts)

	client := github.NewClient(tc)
	user, _, _ := client.Users.Get(gi.ctx, "")
	if user == nil {
		return false
	}
	gi.user = user
	gi.githubClient = client
	return true
}

func (gi gi) SetGithubToken(token string) bool {
	gi.config.GithubToken = token
	gi.config.Save()
	return gi.updateGithubClient()
}

func (gi gi) GetUserName() *string {
	return gi.user.Login
}

func (gi gi) GetClient() *github.Client {
	return gi.githubClient
}

func (gi gi) GetConfig() *Config {
	return gi.config
}

func (gi gi) Ctx() context.Context {
	return gi.ctx
}

func (gi gi) GetUser() *github.User {
	user, _, _ := gi.githubClient.Users.Get(gi.ctx, "")
	return user
}
