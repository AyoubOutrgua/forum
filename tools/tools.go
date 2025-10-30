package tools

type Post struct {
	ID           int
	Title        string
	Description  string
	ImageUrl     string
	UserName     string
	CreationDate string
	Categories   []string
}

type Category struct {
	ID       int
	Category string
}

type ReactionStats struct {
	PostID        int
	LikesCount    int
	DislikesCount int
}

type PageData struct {
	Posts         []Post
	Categories    []Category
	IsLogin       IsLogin
	ReactionStats map[int]ReactionStats
	UserReactions map[int]int
	Comment       map[int][]Comment
	ConnectUserName string
}

type IsLogin struct {
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
