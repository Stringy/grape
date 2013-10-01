package grape

import (
	"errors"
	"strings"
)

type response interface {
	hasErrors() bool
	getError() error
}

type userResponse struct {
	Data Redditor
	errorJson
}

func (u *userResponse) hasErrors() bool {
	return len(u.errorJson.Json.Errors) > 0
}

func (u *userResponse) getError() error {
	return errors.New(strings.Join(u.errorJson.Json.Errors[0], ", "))
}

type redditResponse struct {
	Data Subreddit
	errorJson
}

func (u *redditResponse) hasErrors() bool {
	return len(u.errorJson.Json.Errors) > 0
}

func (u *redditResponse) getError() error {
	return errors.New(strings.Join(u.errorJson.Json.Errors[0], ", "))
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
	errorJson
}

func (u *commentsResponse) hasErrors() bool {
	return len(u.errorJson.Json.Errors) > 0
}

func (u *commentsResponse) getError() error {
	return errors.New(strings.Join(u.errorJson.Json.Errors[0], ", "))
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

func (u *loginResponse) hasErrors() bool {
	return len(u.Json.Errors) > 0
}

func (u *loginResponse) getError() error {
	return errors.New(strings.Join(u.Json.Errors[0], ", "))
}

type messageResponse struct {
	Data struct {
		Children []struct {
			Msg Message `json:"data"`
		}
	}
	errorJson
}

func (u *messageResponse) hasErrors() bool {
	return len(u.errorJson.Json.Errors) > 0
}

func (u *messageResponse) getError() error {
	return errors.New(strings.Join(u.errorJson.Json.Errors[0], ", "))
}
