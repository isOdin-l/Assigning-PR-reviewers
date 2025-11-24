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
