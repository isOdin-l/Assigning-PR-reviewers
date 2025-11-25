package models

import (
	"time"

	"github.com/isOdin-l/Assigning-PR-reviewers/pkg/api"
)

// TODO: расположить коверторы по блокам, разделённым этими комментариями:
// API => INTERNAL
// API <= INTERNAL

func convertToTeamMember(teamMembers []api.TeamMember) *[]TeamMember {
	members := make([]TeamMember, len(teamMembers), cap(teamMembers))
	for idx, v := range teamMembers {
		members[idx] = TeamMember{UserId: v.UserId, Username: v.Username, IsActive: v.IsActive}
	}
	return &members
}
func ConvertToTeam(team *api.Team) *Team {
	return &Team{
		TeamName: team.TeamName,
		Members:  *convertToTeamMember(team.Members),
	}
}

func convertToApiTeamMember(teamMembers []TeamMember) *[]api.TeamMember {
	members := make([]api.TeamMember, len(teamMembers), cap(teamMembers))
	for idx, v := range teamMembers {
		members[idx] = api.TeamMember{UserId: v.UserId, Username: v.Username, IsActive: v.IsActive}
	}
	return &members
}
func ConvertToApiTeam(team *Team) *api.Team {
	return &api.Team{
		TeamName: team.TeamName,
		Members:  *convertToApiTeamMember(team.Members),
	}
}

func ConvertToPostUserSetIsActive(user *api.PostUserSetIsActive) *PostUserSetIsActive {
	return &PostUserSetIsActive{
		UserId:   user.UserId,
		IsActive: user.IsActive,
	}
}

func ConvertToUser(user *PostUserSetIsActive) *User {
	return &User{
		UserId:   user.UserId,
		IsActive: user.IsActive,
	}
}

func ConvertToGetTeamParams(team *api.GetTeamParams) *GetTeamParams {
	return &GetTeamParams{
		TeamName: team.TeamName,
	}
}

func ConvertToPullRequestCreate(pullRequest *api.PullRequestCreate) *PullRequestCreate {
	return &PullRequestCreate{
		AuthorId:        pullRequest.AuthorId,
		PullRequestId:   pullRequest.PullRequestId,
		PullRequestName: pullRequest.PullRequestName,
	}
}

func ConvertToPullRequestMerge(pullRequest *api.PullRequestMerge) *PullRequestMerge {
	return &PullRequestMerge{
		PullRequestId: pullRequest.PullRequestId,
	}
}

func ConvertToPullRequestReassign(pullRequest *api.PullRequestReassign) *PullRequestReassign {
	return &PullRequestReassign{
		OldUserId:     pullRequest.OldUserId,
		PullRequestId: pullRequest.PullRequestId,
	}
}

func ConvertToApiPullRequestCreate(pullRequest PullRequest) *api.ResponsePullRequestCreate {
	return &api.ResponsePullRequestCreate{
		PR: struct {
			PullRequestId     string                "json:\"pull_request_id\""
			PullRequestName   string                "json:\"pull_request_name\""
			AuthorId          string                "json:\"author_id\""
			Status            api.PullRequestStatus "json:\"status\""
			AssignedReviewers []string              "json:\"assigned_reviewers\""
		}{
			PullRequestId:     pullRequest.PullRequestId,
			PullRequestName:   pullRequest.PullRequestName,
			AuthorId:          pullRequest.AuthorId,
			Status:            api.PullRequestStatus(pullRequest.Status),
			AssignedReviewers: pullRequest.AssignedReviewers,
		},
	}
}

func ConvertToApiPullRequestMerge(pullRequest *PullRequest) *api.ResponsePullRequestMerge {
	return &api.ResponsePullRequestMerge{
		PR: struct {
			PullRequestId     string                "json:\"pull_request_id\""
			PullRequestName   string                "json:\"pull_request_name\""
			AuthorId          string                "json:\"author_id\""
			Status            api.PullRequestStatus "json:\"status\""
			AssignedReviewers []string              "json:\"assigned_reviewers\""
			MergedAt          *time.Time            "json:\"mergedAt\""
		}{
			PullRequestId:     pullRequest.PullRequestId,
			PullRequestName:   pullRequest.PullRequestName,
			AuthorId:          pullRequest.AuthorId,
			Status:            api.PullRequestStatus(pullRequest.Status),
			AssignedReviewers: pullRequest.AssignedReviewers,
			MergedAt:          pullRequest.MergedAt,
		},
	}
}

func ConvertToApiPullRequestReassign(pullRequest *PullRequest, replacedByUserId string) *api.ResponsePullRequestReassign {
	return &api.ResponsePullRequestReassign{
		PR: struct {
			PullRequestId     string                "json:\"pull_request_id\""
			PullRequestName   string                "json:\"pull_request_name\""
			AuthorId          string                "json:\"author_id\""
			Status            api.PullRequestStatus "json:\"status\""
			AssignedReviewers []string              "json:\"assigned_reviewers\""
		}{
			PullRequestId:     pullRequest.PullRequestId,
			PullRequestName:   pullRequest.PullRequestName,
			AuthorId:          pullRequest.AuthorId,
			Status:            api.PullRequestStatus(pullRequest.Status),
			AssignedReviewers: pullRequest.AssignedReviewers,
		},
		ReplacedByUserId: replacedByUserId,
	}
}

func ConvertToApiResponseTeam(team *Team) *api.ResponseTeam {
	return &api.ResponseTeam{
		Team: struct {
			Members  []api.TeamMember "json:\"members\""
			TeamName string           "json:\"team_name\""
		}{
			TeamName: team.TeamName,
			Members:  *convertToApiTeamMember(team.Members),
		},
	}
}

func ConvertToApiResponseSetUserActive(user *User) *api.ResponseSetUserActive {
	return &api.ResponseSetUserActive{
		User: struct {
			UserId   string "json:\"user_id\""
			UserName string "json:\"username\""
			TeamName string "json:\"team_name\""
			IsActive bool   "json:\"is_active\""
		}{
			UserId:   user.UserId,
			UserName: user.Username,
			TeamName: user.TeamName,
			IsActive: user.IsActive,
		},
	}
}

func ConvertToApiResponseGetPRsWhereUserIsReviewer(userPrs *PRsWhereUserIsReviewer) *api.ResponseGetPRsWhereUserIsReviewer {
	prs := make([]struct {
		PullRequestId   string                "json:\"pull_request_id\""
		PullRequestName string                "json:\"pull_request_name\""
		AuthorId        string                "json:\"author_id\""
		Status          api.PullRequestStatus "json:\"status\""
	}, 0)

	result := &api.ResponseGetPRsWhereUserIsReviewer{UserId: userPrs.User_id}

	for _, pr := range userPrs.PullRequests {
		prs = append(prs, struct {
			PullRequestId   string                "json:\"pull_request_id\""
			PullRequestName string                "json:\"pull_request_name\""
			AuthorId        string                "json:\"author_id\""
			Status          api.PullRequestStatus "json:\"status\""
		}{
			PullRequestId:   pr.PullRequestId,
			PullRequestName: pr.PullRequestName,
			AuthorId:        pr.AuthorId,
			Status:          api.PullRequestStatus(pr.Status),
		})
	}
	result.PR = prs

	return result
}

func ConvertToStringUser(reviewersId *[]UserId) *[]string {
	result := make([]string, 0)
	for _, v := range *reviewersId {
		result = append(result, v.UserId)
	}
	return &result
}

func ConvertToStringReviewer(reviewersId *[]ReviewerId) *[]string {
	result := make([]string, 0)
	for _, v := range *reviewersId {
		result = append(result, v.ReviewerId)
	}
	return &result
}
