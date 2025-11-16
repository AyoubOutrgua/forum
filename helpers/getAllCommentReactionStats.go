package helpers

import (
	"forum/database"
	"forum/tools"
)

// GetAllCommentReactionStats: Daba kat rajja3 (map, error)
func GetAllCommentReactionStats() (map[int]tools.CommentReactionStats, error) { // BEDDELNAHA
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
		return nil, err 
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
	
	return stats, nil 
}

func GetUserCommentReactions(userID int) (map[int]int, error) { 
	reactions := make(map[int]int)
	
	if userID == 0 {
		return reactions, nil 
	}
	
	rows, err := database.DataBase.Query(
		"SELECT commentId, reaction FROM commentReactions WHERE userId = ?",
		userID,
	)
	if err != nil {
		return nil, err 
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
	
	return reactions, nil 
}