package models

import "github.com/isOdin-l/Assigning-PR-reviewers/pkg/api"

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

func ConvertToGetTeamGetParams(team *api.GetTeamGetParams) *GetTeamGetParams {
	return &GetTeamGetParams{
		TeamName: team.TeamName,
	}
}
