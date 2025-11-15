package store

import (
	"sync"

	"github.com/vyacheslavbytsko/Pull-Requests-Reviewers-Service/internal/api"
)

type Store struct {
	mu    sync.RWMutex
	Users map[string]*api.User
	Teams map[string]*api.Team
	PRs   map[string]*api.PullRequest
}

func NewStore() *Store {
	users := map[string]*api.User{
		"1": {UserId: "1", Username: "Alice", IsActive: true, TeamName: "Backend"},
		"2": {UserId: "2", Username: "Bob", IsActive: true, TeamName: "Backend"},
		"3": {UserId: "3", Username: "Charlie", IsActive: false, TeamName: "Frontend"},
	}
	teams := map[string]*api.Team{
		"Backend": {
			TeamName: "Backend",
			Members: []api.TeamMember{
				UserToTeamMember(users["1"]),
				UserToTeamMember(users["2"]),
			},
		},
		"Frontend": {
			TeamName: "Frontend",
			Members: []api.TeamMember{
				UserToTeamMember(users["3"]),
			},
		},
	}
	return &Store{
		Users: users,
		Teams: teams,
		PRs:   make(map[string]*api.PullRequest),
	}
}

func UserToTeamMember(u *api.User) api.TeamMember {
	return api.TeamMember{
		UserId:   u.UserId,
		Username: u.Username,
		IsActive: u.IsActive,
	}
}
