package helpers

import (
	"forum/database"
	"forum/tools"
	"net/http"
)

func GetAllCommentReactionStats(w http.ResponseWriter) map[int]tools.CommentReactionStats {
	stats := make(map[int]tools.CommentReactionStats)
	
	query := `
		SELECT 
			commentId,
			COALESCE(SUM(CASE WHEN reaction = 1 THEN 1 ELSE 0 END), 0) as likes,
			COALESCE(SUM(CASE WHEN reaction = -1 THEN 1 ELSE 0 END), 0) as dislikes
		FROM commentReactions
		GROUP BY commentId
	`
	
	rows, err := database.DataBase.Query(query)
	if err != nil {
		return stats
	}
	defer rows.Close()
	
	for rows.Next() {
		var stat tools.CommentReactionStats
		err := rows.Scan(&stat.CommentID, &stat.LikesCount, &stat.DislikesCount)
		if err != nil {
			continue
		}
		stats[stat.CommentID] = stat
	}
	
	return stats
}

func GetUserCommentReactions(w http.ResponseWriter, userID int) map[int]int {
	reactions := make(map[int]int)
	
	if userID == 0 {
		return reactions
	}
	
	rows, err := database.DataBase.Query(
		"SELECT commentId, reaction FROM commentReactions WHERE userId = ?",
		userID,
	)
	if err != nil {
		return reactions
	}
	defer rows.Close()
	
	for rows.Next() {
		var commentID, reaction int
		err := rows.Scan(&commentID, &reaction)
		if err != nil {
			continue
		}
		reactions[commentID] = reaction
	}
	
	return reactions
}