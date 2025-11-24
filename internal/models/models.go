package models

import (
	"time"

	"github.com/isOdin-l/Assigning-PR-reviewers/pkg/api"
)

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

type GetTeamParams struct {
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

type PullRequest struct {
	PullRequestId     string            `db:"id"`
	AuthorId          string            `db:"author_id"`
	PullRequestName   string            `db:"name"`
	Status            PullRequestStatus `db:"status"`
	AssignedReviewers []string
	CreatedAt         *time.Time
	MergedAt          *time.Time
}
type PullRequestShort struct {
	AuthorId        string
	PullRequestId   string
	PullRequestName string
	Status          PullRequestStatus
}
type PullRequestStatus string

type PullRequestCreate struct {
	AuthorId        string
	PullRequestId   string
	PullRequestName string
}
type PullRequestMerge struct {
	PullRequestId string
}
type PullRequestReassign struct {
	OldUserId     string
	PullRequestId string
}

const (
	PullRequestStatusMERGED PullRequestStatus = "MERGED"
	PullRequestStatusOPEN   PullRequestStatus = "OPEN"
)
