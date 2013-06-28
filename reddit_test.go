package reddit

import (
	_ "fmt"
	"testing"
)

// func TestGetSubreddit(t *testing.T) {
// 	sub, err := GetSubreddit("reddit_test0")
// 	if err != nil {
// 		t.Errorf("Error from subreddit retrieval: %v", err)
// 		t.FailNow()
// 	}
// 	if len(sub.Items) != 8 {
// 		t.Errorf("Unexpected number of items from subreddit: %v", len(sub.Items))
// 		t.Fail()
// 	}
// 	if sub.Name != "reddit_test0" {
// 		t.Errorf("Queried incorrect subreddit: %v", sub.Name)
// 		t.Fail()
// 	}
// 	sub, err = GetSubreddit("not_a_subreddit")
// 	if err == nil {
// 		t.Errorf("Expected error from GetSubreddit")
// 		t.Fail()
// 	}
// }

// func TestLogin(t *testing.T) {
// 	redditor, err := Login("reddit", "password", true)
// 	if err != nil {
// 		t.Errorf("Unexpected error from login: %v", err)
// 		t.FailNow()
// 	}
// 	if redditor.ModHash == "" {
// 		t.Errorf("nil modhash returned from login")
// 		t.Fail()
// 	}
// 	if len(client.Jar.Cookies(actual_url)) == 0 {
// 		t.Errorf("nil cookies from login")
// 		t.Fail()
// 	}
// }

// func TestGetRedditor(t *testing.T) {
// 	redditor, err := GetRedditor("reddit")
// 	if err != nil {
// 		t.Errorf("Unexpected error returned from login: %v", err)
// 		t.FailNow()
// 	}
// 	if !redditor.IsMod {
// 		t.Errorf("Expected mod status for user 'reddit'")
// 		t.Fail()
// 	}
// 	if redditor.Name != "reddit" {
// 		t.Errorf("Unexpected redditor name %s", redditor.Name)
// 		t.Fail()
// 	}
// 	if redditor.CKarma != 0 || redditor.LKarma != 3 {
// 		t.Errorf(
// 			"Unexpected karma results\n\tC: %d\n\tL: %d",
// 			redditor.CKarma,
// 			redditor.LKarma)
// 		t.Fail()
// 	}
// }

// func TestSubmitLink(t *testing.T) {
// 	user, err := Login("stringy", "test", true)
// 	if err != nil {
// 		t.Errorf("%v", err)
// 		t.Fail()
// 	}

// 	// capt, err := user.getCaptcha()
// 	// if err != nil {
// 	// 	t.Errorf("%v", err)
// 	// 	t.Fail()
// 	// }
// 	//	t.Log(capt)
// 	err = user.SubmitLink(
// 		"reddit_test0",
// 		"This is a test",
// 		"",
// 		"http://www.google.com",
// 		KindLink)

// 	if err != nil {
// 		t.Errorf("%v", err)
// 		t.Fail()
// 	}
// }

func TestSubmitComment(t *testing.T) {
	user, err := Login("stringy", "test", true)
	if err != nil {
		t.Errorf("%v", err)
		t.FailNow()
	}
	sub, err := GetSubreddit("reddit_test0")
	if err != nil {
		t.Errorf("%v", err)
		t.FailNow()
	}
	comment := new(Comment)
	comment.Author = user.Name
	comment.Body = "This is a test comment"
	err = sub.Items[0].PostComment(user, comment)
	if err != nil {
		t.Errorf("%v", err)
		t.FailNow()
	}
}

// func TestDeleteAccount(t *testing.T) {
// 	u, err := Login("Stringy", "test", true)
// 	if err != nil {
// 		t.Errorf("%v", err)
// 		t.Fail()
// 	}
// 	err = u.DeleteAccount("test")
// 	if err != nil {
// 		t.Errorf("%v", err)
// 		t.Fail()
// 	}
// }
