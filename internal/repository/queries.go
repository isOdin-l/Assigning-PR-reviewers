package repository

import (
	"fmt"

	"github.com/isOdin-l/Assigning-PR-reviewers/internal/database"
)

// Users
func getColumnsUserIsReviewer() []string {
	return []string{"id", "author_id", "name", "status"}
}

func getJoinUserIsReviewer() string {
	return fmt.Sprintf("%s ON id = pr_id", database.PrReviewsTable)
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

func getInsertColumnsPrReviewer() []string {
	return []string{"pr_id", "reviewer_id"}
}

func getSelectColumnsPrMerge() []string {
	return []string{"author_id", "name", "status", "merged_at"}
}

func getSelectColumnsPr() []string {
	return []string{"id", "author_id", "name", "status"}
}
