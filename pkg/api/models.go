package api

import "time"

// === Request Api Models ===
type GetUsersGetReview struct {
	UserId string `form:"user_id"`
}

type PostUserSetIsActive struct {
	UserId   string `json:"user_id" form:"user_id" `
	IsActive bool   `json:"is_active" form:"is_active"`
}

// Error
const (
	SERVERERROR ErrorResponseErrorCode = "SERVER_ERROR"
	NOCANDIDATE ErrorResponseErrorCode = "NO_CANDIDATE"
	NOTASSIGNED ErrorResponseErrorCode = "NOT_ASSIGNED"
	NOTFOUND    ErrorResponseErrorCode = "NOT_FOUND"
	PREXISTS    ErrorResponseErrorCode = "PR_EXISTS"
	PRMERGED    ErrorResponseErrorCode = "PR_MERGED"
	TEAMEXISTS  ErrorResponseErrorCode = "TEAM_EXISTS"
)

type ErrorResponse struct {
	Error struct {
		Code    ErrorResponseErrorCode `json:"code"`
		Message string                 `json:"message"`
	} `json:"error"`
}
type ErrorResponseErrorCode string

// Team
type GetTeamParams struct {
	TeamName string `form:"team_name" json:"team_name"`
}

type Team struct {
	Members  []TeamMember `json:"members"`
	TeamName string       `json:"team_name"`
}

type TeamMember struct {
	IsActive bool   `json:"is_active"`
	UserId   string `json:"user_id"`
	Username string `json:"username"`
}

// PullRequest
type PullRequestShort struct {
	AuthorId        string            `json:"author_id" form:"author_id"`
	PullRequestId   string            `json:"pull_request_id" form:"pull_request_id"`
	PullRequestName string            `json:"pull_request_name" form:"pull_request_name"`
	Status          PullRequestStatus `json:"status" form:"status"`
}
type PullRequestStatus string

type PullRequestCreate struct {
	AuthorId        string `json:"author_id"`
	PullRequestId   string `json:"pull_request_id"`
	PullRequestName string `json:"pull_request_name"`
}

type PullRequestMerge struct {
	PullRequestId string `json:"pull_request_id"`
}

type PullRequestReassign struct {
	OldUserId     string `json:"old_user_id"`
	PullRequestId string `json:"pull_request_id"`
}

// === Response Api Models ===

type ResponseTeam struct {
	Team struct {
		Members  []TeamMember `json:"members"`
		TeamName string       `json:"team_name"`
	} `json:"team"`
}

type ResponseSetUserActive struct {
	User struct {
		UserId   string `json:"user_id"`
		UserName string `json:"username"`
		TeamName string `json:"team_name"`
		IsActive bool   `json:"is_active"`
	} `json:"user"`
}

type ResponseGetPRsWhereUserIsReviewer struct {
	UserId string `json:"user_id"`
	PR     []struct {
		PullRequestId   string            `json:"pull_request_id"`
		PullRequestName string            `json:"pull_request_name"`
		AuthorId        string            `json:"author_id"`
		Status          PullRequestStatus `json:"status"`
	} `json:"pull_requests"`
}

type ResponsePullRequestCreate struct {
	PR struct {
		PullRequestId     string            `json:"pull_request_id"`
		PullRequestName   string            `json:"pull_request_name"`
		AuthorId          string            `json:"author_id"`
		Status            PullRequestStatus `json:"status"`
		AssignedReviewers []string          `json:"assigned_reviewers"`
	} `json:"pr"`
}

type ResponsePullRequestMerge struct {
	PR struct {
		PullRequestId     string            `json:"pull_request_id"`
		PullRequestName   string            `json:"pull_request_name"`
		AuthorId          string            `json:"author_id"`
		Status            PullRequestStatus `json:"status"`
		AssignedReviewers []string          `json:"assigned_reviewers"`
		MergedAt          *time.Time        `json:"mergedAt"`
	} `json:"pr"`
}

type ResponsePullRequestReassign struct {
	PR struct {
		PullRequestId     string            `json:"pull_request_id"`
		PullRequestName   string            `json:"pull_request_name"`
		AuthorId          string            `json:"author_id"`
		Status            PullRequestStatus `json:"status"`
		AssignedReviewers []string          `json:"assigned_reviewers"`
	} `json:"pr"`
	ReplacedByUserId string `json:"replaced_by"`
}
