package helpers

import (
	"net/http"

	"forum/database"
	"forum/tools"
)

func GetAllCategories(w http.ResponseWriter) []tools.Category {
	categoriesQuery := `
			SELECT c.id, c.category
			FROM categories AS c;
			`
	categories, errSelect := database.SelectAllCategories(categoriesQuery)
	if errSelect != nil {
		Errorhandler(w, "Status Internal Server Error", http.StatusInternalServerError)
		return nil
	}
	return categories
}
