package grape

import (
	"net/url"
)

type Thing struct {
	Id   string
	Name string //fullname of a reddit thing (t1_, t2_, ...)
	Kind string //type
}

// Report reports a reddit Thing by user
// returns any errors recieved from reddit
func (t *Thing) Report(user *Redditor) error {
	if !user.IsLoggedIn() {
		return notLoggedInError
	}
	data := &url.Values{
		"id": {t.Name},
		"uh": {user.ModHash},
	}
	b, err := makePostRequest(Config.GetApiUrl("report"), data)
	if err != nil {
		return err
	}
	return parseSimpleErrorResponse(b)
}

//Hide hides a Thing for user, so that it won't appear on any requests
//returns any errors from reddit
func (t *Thing) Hide(user *Redditor) error {
	//id
	//ModHash
	data := &url.Values{
		"id": {t.Name},
		"uh": {user.ModHash},
	}
	b, err := makePostRequest(Config.GetApiUrl("hide"), data)
	if err != nil {
		return nil
	}
	return parseSimpleErrorResponse(b)
}

//Unhide undoes Hide to allow a Thing to turn up in user's requests
//returns any errors from reddit
func (t *Thing) Unhide(user *Redditor) error {
	//id
	//modhash
	data := &url.Values{
		"id": {t.Name},
		"uh": {user.ModHash},
	}
	b, err := makePostRequest(Config.GetApiUrl("unhide"), data)
	if err != nil {
		return nil
	}
	return parseSimpleErrorResponse(b)
}

// This might not apply to all objects
// func (t *Thing) Info() error {
// 	return nil
// }

//MarkNsfw marks a reddit thing as not safe for work for user
//returns errors from reddit
func (t *Thing) MarkNsfw(user *Redditor) error {
	//id
	//modhash
	data := &url.Values{
		"id": {t.Name},
		"uh": {user.ModHash},
	}
	b, err := makePostRequest(Config.GetApiUrl("marknsfw"), data)
	if err != nil {
		return nil
	}
	return parseSimpleErrorResponse(b)
}

//UnmarkNsfw unmarks a Thing as not safe for work
//returns errors recieved from reddit
func (t *Thing) UnmarkNsfw(user *Redditor) error {
	data := &url.Values{
		"id": {t.Name},
		"uh": {user.ModHash},
	}
	b, err := makePostRequest(Config.GetApiUrl("unmarknsfw"), data)
	if err != nil {
		return nil
	}
	return parseSimpleErrorResponse(b)
}
