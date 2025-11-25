package models

import (
	"time"
)

type GetTeamParams struct {
	TeamName string
}

type Team struct {
	Members  []TeamMember
	TeamName string
}

type TeamMember struct {
	UserId   string `db:"id"`
	Username string `db:"user_name"`
	IsActive bool   `db:"is_active"`
}

// User
type User struct {
	UserId   string
	Username string
	TeamName string
	IsActive bool
}
type PRsWhereUserIsReviewer struct {
	User_id      string             `json:"user_id"`
	PullRequests []PullRequestShort `json:"pr"`
}

type PostUserSetIsActive struct {
	UserId   string
	IsActive bool
}

type ReviewerId struct {
	ReviewerId string `db:"reviewer_id"`
}

type UserId struct {
	UserId string `db:"id"`
}

// PullRequest
type PullRequestStatus string

const (
	PullRequestStatusMERGED PullRequestStatus = "MERGED"
	PullRequestStatusOPEN   PullRequestStatus = "OPEN"
)

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
	PullRequestId   string            `db:"id"`
	AuthorId        string            `db:"author_id"`
	PullRequestName string            `db:"name"`
	Status          PullRequestStatus `db:"status"`
}

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
