package repository

import (
	"fmt"

	"github.com/isOdin-l/Assigning-PR-reviewers/internal/database"
)

func getColumnsUserIsReviewer() []string {
	return []string{"id", "author_id", "name", "status"}
}

func getJoinUserIsReviewer() string {
	return fmt.Sprintf("%s USING (pr_id)", database.PrReviewsTable)
}

func getUpdateUserIsActiveSuffix() string {
	return "RETURNING user_name, team_name"
}
