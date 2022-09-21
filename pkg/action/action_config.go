package action

import (
	"fmt"
	"github.com/sethvargo/go-githubactions"
	"os"
	"strings"
)

type Config struct {
	token *string
	repo  *string
	owner *string

	types         *[]string
	scopes        *[]string
	headerPattern *string
}

func NewActionConfig(action *githubactions.Action) (*Config, error) {
	ownerRepoSlug := os.Getenv("GITHUB_REPOSITORY")
	ownerRepo := strings.Split(ownerRepoSlug, "/")
	if len(ownerRepo) != 2 {
		return nil, fmt.Errorf("GITHUB_REPOSITORY is not found")
	}
	owner, repo := ownerRepo[0], ownerRepo[1]

	token := os.Getenv("GITHUB_TOKEN")

	types := strings.Fields(strings.TrimSpace(action.GetInput("types")))
	scopes := strings.Fields(strings.TrimSpace(action.GetInput("scopes")))
	headerPattern := action.GetInput("headerPattern")

	return &Config{
		token:         &token,
		repo:          &repo,
		owner:         &owner,
		types:         &types,
		scopes:        &scopes,
		headerPattern: &headerPattern,
	}, nil
}

func (ac *Config) GetOwner() string {
	if ac == nil || ac.owner == nil {
		return ""
	}
	return *ac.owner
}

func (ac *Config) GetRepo() string {
	if ac == nil || ac.repo == nil {
		return ""
	}
	return *ac.repo
}

func (ac *Config) GetToken() string {
	if ac == nil || ac.token == nil {
		return ""
	}
	return *ac.token
}

func (ac *Config) GetTypes() []string {
	if ac == nil || ac.types == nil {
		return nil
	}
	return *ac.types
}

func (ac *Config) GetScopes() []string {
	if ac == nil || ac.scopes == nil {
		return nil
	}
	return *ac.scopes
}

func (ac *Config) GetHeaderPattern() string {
	if ac == nil || ac.headerPattern == nil {
		return ""
	}
	return *ac.headerPattern
}
