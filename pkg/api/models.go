package api

type GetUsersGetReview struct {
	UserId string `form:"user_id"`
}

type PullRequestShort struct {
	AuthorId        string            `json:"author_id" form:"author_id"`
	PullRequestId   string            `json:"pull_request_id" form:"pull_request_id"`
	PullRequestName string            `json:"pull_request_name" form:"pull_request_name"`
	Status          PullRequestStatus `json:"status" form:"status"`
}
type PullRequestStatus string

type PostUserSetIsActive struct {
	UserId   string `json:"user_id" form:"user_id" `
	IsActive bool   `json:"is_active" form:"is_active"`
}

const (
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
