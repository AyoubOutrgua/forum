package helpers

import (
	"net/http"

	"forum/database"
	"forum/tools"
)

func GetAllReactionStats(w http.ResponseWriter) map[int]tools.ReactionStats {
	stats := make(map[int]tools.ReactionStats)

	query := `
		SELECT 
			postId,
			COALESCE(SUM(CASE WHEN reaction = 1 THEN 1 ELSE 0 END), 0) as likes,
			COALESCE(SUM(CASE WHEN reaction = -1 THEN 1 ELSE 0 END), 0) as dislikes
		FROM postReactions
		GROUP BY postId
	`

	rows, err := database.DataBase.Query(query)
	if err != nil {
		return stats
	}
	defer rows.Close()

	for rows.Next() {
		var stat tools.ReactionStats
		if err := rows.Scan(&stat.PostID, &stat.LikesCount, &stat.DislikesCount); err != nil {
			continue
		}
		stats[stat.PostID] = stat
	}

	return stats
}

func GetUserPostReactions(w http.ResponseWriter, userID int) map[int]int {
	reactions := make(map[int]int)

	if userID == 0 {
		return reactions
	}

	rows, err := database.DataBase.Query(
		"SELECT postId, reaction FROM postReactions WHERE userId = ?",
		userID,
	)
	if err != nil {
		return reactions
	}
	defer rows.Close()

	for rows.Next() {
		var postID, reaction int
		if err := rows.Scan(&postID, &reaction); err != nil {
			continue
		}
		reactions[postID] = reaction
	}

	return reactions
}
