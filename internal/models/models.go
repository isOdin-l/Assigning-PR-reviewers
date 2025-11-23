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

type ResponseUser struct {
	UserId   string `json:"user_id"`
	Username string `json:"username"`
	TeamName string `json:"team_name"`
	IsActive bool   `json:"is_active"`
}

func ConvertToPostUserSetIsActive(user api.PostUserSetIsActive) *PostUserSetIsActive {
	return &PostUserSetIsActive{
		UserId:   user.UserId,
		IsActive: user.IsActive,
	}
}

func ConvertToResponseUser(user PostUserSetIsActive) *ResponseUser {
	return &ResponseUser{
		UserId:   user.UserId,
		IsActive: user.IsActive,
	}
}
