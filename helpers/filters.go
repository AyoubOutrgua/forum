package helpers

import (
	"database/sql"
	"fmt"
	"time"
)

type PostView struct {
	ID           int
	Title        string
	Description  string
	ImageUrl     string
	CreationDate string
	Category     string
	UserName     string
	UserID       int
}

func formatTimeNull(s sql.NullString) string {
	if !s.Valid || s.String == "" {
		return ""
	}
	layouts := []string{
		time.RFC3339,
		"2006-01-02 15:04:05",
		"2006-01-02",
	}
	for _, l := range layouts {
		if t, err := time.Parse(l, s.String); err == nil {
			return t.Format("Jan 02 2006 15:04")
		}
	}
	return s.String
}

func PostsByAuthor(db *sql.DB, authorID int) ([]PostView, error) {
	const q = `
        SELECT p.id,
               p.title,
               COALESCE(p.content, p.description, '') AS description,
               COALESCE(p.image_url, p.image, '') AS image_url,
               COALESCE(STRFTIME('%Y-%m-%d %H:%M:%S', p.created_at), p.created_at) as created_at,
               COALESCE(c.category, '') as category,
               COALESCE(u.userName, '') as userName,
               p.user_id
        FROM posts p
        LEFT JOIN categories c ON p.category_id = c.id
        LEFT JOIN users u ON p.user_id = u.id
        WHERE p.user_id = ?
        ORDER BY p.created_at DESC
    `
	rows, err := db.Query(q, authorID)
	if err != nil {
		return nil, fmt.Errorf("PostsByAuthor query: %w", err)
	}
	defer rows.Close()

	var out []PostView
	for rows.Next() {
		var p PostView
		var created sql.NullString
		var image sql.NullString
		var category sql.NullString
		var userName sql.NullString
		if err := rows.Scan(&p.ID, &p.Title, &p.Description, &image, &created, &category, &userName, &p.UserID); err != nil {
			return nil, fmt.Errorf("PostsByAuthor scan: %w", err)
		}
		if image.Valid {
			p.ImageUrl = image.String
		}
		p.CreationDate = formatTimeNull(created)
		if category.Valid {
			p.Category = category.String
		}
		if userName.Valid {
			p.UserName = userName.String
		}
		out = append(out, p)
	}
	return out, rows.Err()
}

func PostsByCategory(db *sql.DB, categoryID int) ([]PostView, error) {
	const q = `
        SELECT p.id,
               p.title,
               COALESCE(p.content, p.description, '') AS description,
               COALESCE(p.image_url, p.image, '') AS image_url,
               COALESCE(STRFTIME('%Y-%m-%d %H:%M:%S', p.created_at), p.created_at) as created_at,
               COALESCE(c.category, '') as category,
               COALESCE(u.userName, '') as userName,
               p.user_id
        FROM posts p
        LEFT JOIN categories c ON p.category_id = c.id
        LEFT JOIN users u ON p.user_id = u.id
        WHERE p.category_id = ?
        ORDER BY p.created_at DESC
    `
	rows, err := db.Query(q, categoryID)
	if err != nil {
		return nil, fmt.Errorf("PostsByCategory query: %w", err)
	}
	defer rows.Close()

	var out []PostView
	for rows.Next() {
		var p PostView
		var created sql.NullString
		var image sql.NullString
		var category sql.NullString
		var userName sql.NullString
		if err := rows.Scan(&p.ID, &p.Title, &p.Description, &image, &created, &category, &userName, &p.UserID); err != nil {
			return nil, fmt.Errorf("PostsByCategory scan: %w", err)
		}
		if image.Valid {
			p.ImageUrl = image.String
		}
		p.CreationDate = formatTimeNull(created)
		if category.Valid {
			p.Category = category.String
		}
		if userName.Valid {
			p.UserName = userName.String
		}
		out = append(out, p)
	}
	return out, rows.Err()
}

func LikedPostsByUser(db *sql.DB, userID int) ([]PostView, error) {
	const q = `
        SELECT p.id,
               p.title,
               COALESCE(p.content, p.description, '') AS description,
               COALESCE(p.image_url, p.image, '') AS image_url,
               COALESCE(STRFTIME('%Y-%m-%d %H:%M:%S', p.created_at), p.created_at) as created_at,
               COALESCE(c.category, '') as category,
               COALESCE(u.userName, '') as userName,
               p.user_id
        FROM posts p
        INNER JOIN reactions r ON r.post_id = p.id
        LEFT JOIN categories c ON p.category_id = c.id
        LEFT JOIN users u ON p.user_id = u.id
        WHERE r.user_id = ? AND (r.type = 'like' OR r.value = 1)
        ORDER BY r.created_at DESC
    `
	rows, err := db.Query(q, userID)
	if err != nil {
		return nil, fmt.Errorf("LikedPostsByUser query: %w", err)
	}
	defer rows.Close()

	var out []PostView
	for rows.Next() {
		var p PostView
		var created sql.NullString
		var image sql.NullString
		var category sql.NullString
		var userName sql.NullString
		if err := rows.Scan(&p.ID, &p.Title, &p.Description, &image, &created, &category, &userName, &p.UserID); err != nil {
			return nil, fmt.Errorf("LikedPostsByUser scan: %w", err)
		}
		if image.Valid {
			p.ImageUrl = image.String
		}
		p.CreationDate = formatTimeNull(created)
		if category.Valid {
			p.Category = category.String
		}
		if userName.Valid {
			p.UserName = userName.String
		}
		out = append(out, p)
	}
	return out, rows.Err()
}
