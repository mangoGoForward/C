package action

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/apache/pulsar-test-infra/docbot/pkg/logger"
	"github.com/google/go-github/v45/github"
	"golang.org/x/oauth2"
	"regexp"
)

const (
	TitleUnmatchPattern = `Please follow [Pulsar Pull Request Naming Convention Guide]
(https://docs.google.com/document/d/1d8Pw6ZbWk-_pCKdOmdvx9rnhPiyuxwq60_TrD68d7BA/edit#) to write a PR title.`

	openedActionType = "opened"
	editedActionType = "edited"
)

type Action struct {
	config *Config

	globalContext context.Context
	client        *github.Client

	prNumber int
}

func NewAction(ctx context.Context, ac *Config) *Action {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: ac.GetToken()},
	)

	tc := oauth2.NewClient(ctx, ts)

	return &Action{
		config:        ac,
		globalContext: ctx,
		client:        github.NewClient(tc),
	}
}

func (a *Action) Run(prNumber int, actionType string) error {
	a.prNumber = prNumber

	switch actionType {
	case openedActionType, editedActionType:
		return a.checkPRTitle()
	}
	return nil
}

func (a *Action) checkPRTitle() error {
	pr, _, err := a.client.PullRequests.Get(a.globalContext, a.config.GetOwner(), a.config.GetRepo(), a.prNumber)
	if err != nil {
		return fmt.Errorf("get PR: %v", err)
	}
	title := pr.Title
	logger.Infof("The PR's title: %v\n", *title)

	re := regexp.MustCompile(a.config.GetHeaderPattern())
	matched := re.FindSubmatch([]byte(*title))
	if len(matched) == 4 {
		titleType := bytes.NewBuffer(matched[1]).String()
		if !existInArr(titleType, a.config.GetTypes()) {
			return errors.New(TitleUnmatchPattern)
		}

		titleScope := bytes.NewBuffer(matched[2]).String()
		if !existInArr(titleScope, a.config.GetScopes()) {
			return errors.New(TitleUnmatchPattern)
		}
		return nil
	}
	return errors.New(TitleUnmatchPattern)
}

func existInArr(target string, origin []string) bool {
	titleMatched := false
	for _, val := range origin {
		if target == val {
			titleMatched = true
			break
		}
	}
	return titleMatched
}
