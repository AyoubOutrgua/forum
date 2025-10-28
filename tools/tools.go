package tools

type Post struct {
	ID           int
	Title        string
	Description  string
	ImageUrl     string
	UserName     string
	CreationDate string
	Categories   []string
	Comments     []Comment  // Zid hadi
    CommentCount int
}

type Category struct {
	ID       int
	Category string
}

// type PostCategory struct {
// 	PostID       int
// 	CategoryID int
// }

type ReactionStats struct {
	PostID        int
	LikesCount    int
	DislikesCount int
}

type PageData struct {
	Posts      []Post
	Categories []Category
	IdLogin IdLogin
	ReactionStats map[int]ReactionStats 
	UserReactions map[int]int 
}

type IdLogin struct {
	LoggedIn bool
	UserID   int 
}

type Comment struct {
    ID           int
    CommentText  string
    PostID       int
    UserID       int
    UserName     string
    CreationDate string
}