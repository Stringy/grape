package reddit

import (
	"encoding/json"
	"errors"
	"net/url"
	"strings"
)

type Thing struct {
	Id   string
	Name string //fullname of a reddit thing (t1_, t2_, ...)
	Kind string //type
}

//Report reports a reddit Thing by user
//returns any errors recieved from reddit
func (t *Thing) Report(user *Redditor) error {
	if !user.IsLoggedIn() {
		return NotLoggedInError
	}
	data := &url.Values{
		"id": {t.Name},
		"uh": {user.ModHash},
	}
	b, err := makePostRequest(ApiUrls["report"], data)
	if err != nil {
		return err
	}
	es := new(errorJson)
	err = json.Unmarshal(b, &es)
	if err != nil {
		return err
	}
	if len(es.Json.Errors) != 0 {
		return errors.New(strings.Join(es.Json.Errors[0], ", "))
	}
	return nil
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
	b, err := makePostRequest(ApiUrls["hide"], data)
	if err != nil {
		return nil
	}
	es := new(errorJson)
	err = json.Unmarshal(b, &es)
	if err != nil {
		return err
	}
	if len(es.Json.Errors) != 0 {
		return errors.New(strings.Join(es.Json.Errors[0], ", "))
	}
	return nil
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
	b, err := makePostRequest(ApiUrls["unhide"], data)
	if err != nil {
		return nil
	}
	es := new(errorJson)
	err = json.Unmarshal(b, &es)
	if err != nil {
		return err
	}
	if len(es.Json.Errors) != 0 {
		return errors.New(strings.Join(es.Json.Errors[0], ", "))
	}
	return nil
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
	b, err := makePostRequest(ApiUrls["marknsfw"], data)
	if err != nil {
		return nil
	}
	es := new(errorJson)
	err = json.Unmarshal(b, &es)
	if err != nil {
		return err
	}
	if len(es.Json.Errors) != 0 {
		return errors.New(strings.Join(es.Json.Errors[0], ", "))
	}

	return nil
}

//UnmarkNsfw unmarks a Thing as not safe for work
//returns errors recieved from reddit
func (t *Thing) UnmarkNsfw(user *Redditor) error {
	//id
	//modhash
	data := &url.Values{
		"id": {t.Name},
		"uh": {user.ModHash},
	}
	b, err := makePostRequest(ApiUrls["unmarknsfw"], data)
	if err != nil {
		return nil
	}
	es := new(errorJson)
	err = json.Unmarshal(b, &es)
	if err != nil {
		return err
	}
	if len(es.Json.Errors) != 0 {
		return errors.New(strings.Join(es.Json.Errors[0], ", "))
	}
	return nil
}
