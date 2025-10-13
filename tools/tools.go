package tools

type Post struct {
	ID           int
	Title        string
	Description  string
	ImageUrl     string
	UserName     string
	CreationDate string
}

type Category struct {
	ID       int
	Category string
}

type Insex struct {
	Posts      []Post
	Categories []Category
}
