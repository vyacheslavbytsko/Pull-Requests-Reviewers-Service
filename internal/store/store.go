package store

import (
	"sync"

	"github.com/vyacheslavbytsko/Pull-Requests-Reviewers-Service/internal/api"
)

type Store struct {
	Mu    sync.RWMutex
	Users map[string]*api.User
	Teams map[string]*api.Team
	PRs   map[string]*api.PullRequest
}

func NewStore() *Store {
	return &Store{
		Users: make(map[string]*api.User),
		Teams: make(map[string]*api.Team),
		PRs:   make(map[string]*api.PullRequest),
	}
}
