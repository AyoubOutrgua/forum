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

// type PostCategory struct {
// 	PostID       int
// 	CategoryID int
// }

type Index struct {
	Posts      []Post
	Categories []Category
}
