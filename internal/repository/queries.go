package repository

import (
	"fmt"

	"github.com/isOdin-l/Assigning-PR-reviewers/internal/database"
	"github.com/isOdin-l/Assigning-PR-reviewers/internal/models"
)

// Users
func getColumnsUserIsReviewer() []string {
	return []string{"id", "author_id", "name", "status"}
}

func getJoinUserIsReviewer() string {
	return fmt.Sprintf("%s USING (pr_id)", database.PrReviewsTable)
}

func getUpdateUserIsActiveSuffix() string {
	return "RETURNING user_name, team_name"
}

// Team
func getSelectColumnsTeam() []string {
	return []string{"id", "user_name", "is_active"}
}

func getInsertColumnsTeamMember() []string {
	return []string{"id", "team_name", "user_name", "is_active"}
}

// PullRequest
func getInsertColumnsPr() []string {
	return []string{"id", "author_id", "name"}
}

func getInsertValuesPr(pullRequest *models.PullRequestCreate) []string {
	return []string{pullRequest.PullRequestId, pullRequest.AuthorId, pullRequest.PullRequestName}
}

func getInsertColumnsPrReviewer() []string {
	return []string{"pr_id", "reviewer_id"}
}

func getSelectColumnsPrMerge() []string {
	return []string{"author_id", "name", "status", "merged_at"}
}

func getSelectColumnsPr() []string {
	return []string{"id", "author_id", "name", "status"}
}
