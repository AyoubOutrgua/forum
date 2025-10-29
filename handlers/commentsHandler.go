package handlers

import (
   
    "net/http"
    "strconv"
    "time"
    
    "forum/database"
)

func CreateCommentHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }
    
    // Check if user is logged in
    cookie, err := r.Cookie("session")
    if err != nil {
        http.Redirect(w, r, "/login", http.StatusSeeOther)
        return
    }
    
    // Get user ID
    userID := database.SelectUserID("SELECT id FROM users WHERE session = ?", cookie.Value)
    if userID == 0 {
        http.Redirect(w, r, "/login", http.StatusSeeOther)
        return
    }
    
    // Get form data
    commentText := r.FormValue("comment")
    postIDStr := r.FormValue("postId")
    
    if commentText == "" || postIDStr == "" {
        http.Error(w, "Missing data", http.StatusBadRequest)
        return
    }
    
    postID, err := strconv.Atoi(postIDStr)
    if err != nil {
        http.Error(w, "Invalid post ID", http.StatusBadRequest)
        return
    }
    
    // Insert comment
    query := `INSERT INTO comments (comment, postId, userId, creationDate) 
              VALUES (?, ?, ?, ?)`
    
    creationDate := time.Now().Format("2006-01-02 15:04:05")
    database.ExecuteData(query, commentText, postID, userID, creationDate)
    
    // Redirect back to home
    http.Redirect(w, r, "/", http.StatusSeeOther)
}