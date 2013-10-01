package grape

type userResponse struct {
	Data Redditor
}

type redditResponse struct {
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

func (jc *jsonComment) toComment() Comment {
	comment := new(Comment)
	comment.Author = jc.Author
	comment.Body = jc.Body
	comment.ScoreHidden = jc.ScoreHidden
	comment.Ups = jc.Ups
	comment.Downs = jc.Downs
	comment.Replies = make([]Comment, len(jc.Replies.Data.Children))
	for i, jcReply := range jc.Replies.Data.Children {
		comment.Replies[i] = jcReply.Data.toComment()
	}
	return *comment
}

type commentsResponse struct {
	Data struct {
		Children []struct {
			Data jsonComment `json:"data"`
		}
	}
}

type loginResponse struct {
	Json struct {
		Errors [][]string
		Data   struct {
			ModHash string
			Cookie  string
		}
	}
}

type messageResponse struct {
	Data struct {
		Children []struct {
			Msg Message `json:"data"`
		}
	}
}
