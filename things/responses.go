package things

import ()

type UserResponse struct {
	Data Redditor
}

type RedditResponse struct {
	Data Subreddit
}

type jsonComment struct {
	Author      string
	Body        string
	ScoreHidden bool `json:"score_hidden"`
	Ups         int
	Downs       int
	Replies     struct {
		Data struct {
			Children []struct {
				Data jsonComment
			}
		}
	}
}

func commentFromJson(jComm jsonComment) Comment {
	comment := new(Comment)
	comment.Author = jComm.Author
	comment.Body = jComm.Body
	comment.ScoreHidden = jComm.ScoreHidden
	comment.Ups = jComm.Ups
	comment.Downs = jComm.Downs
	comment.Replies = make([]Comment, len(jComm.Replies.Data.Children))
	for i, jCommReply := range jComm.Replies.Data.Children {
		comment.Replies[i] = commentFromJson(jCommReply.Data)
	}
	return *comment
}

type CommentsResponse struct {
	Data struct {
		Children []struct {
			Data jsonComment `json:"data"`
		}
	}
}

type LoginResponse struct {
	Json struct {
		Errors [][]string
		Data   struct {
			ModHash string
			Cookie  string
		}
	}
}
