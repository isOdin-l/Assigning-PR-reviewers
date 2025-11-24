package models

import "github.com/isOdin-l/Assigning-PR-reviewers/pkg/api"

type ResponsePRsWhereUserIsReviewer struct {
	User_id      string                 `json:"user_id"`
	PullRequests []api.PullRequestShort `json:"pr"`
}

type PostUserSetIsActive struct {
	UserId   string
	IsActive bool
}

type User struct {
	UserId   string
	Username string
	TeamName string
	IsActive bool
}

type GetTeamGetParams struct {
	TeamName string
}

type Team struct {
	Members  []TeamMember
	TeamName string
}

type TeamMember struct {
	UserId   string
	Username string
	IsActive bool
}
