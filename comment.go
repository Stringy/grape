package reddit

type Comment struct {
	Author      string
	Body        string
	ScoreHidden bool
	Ups         int
	Downs       int
	Replies     []Comment
}
