reddit
======

reddit is a Go wrapper for the Reddit API. It will provide functionality for all aspects of the API for use in the creation of bots. 

The current functionality is limited to the following:
 - Logging in
 - Retrieving information about any user
 - Retrieving information about the currently logged in user
 - Retrieving a subreddit (limited to front page only)
 - Retrieving comments from a link
 - Submitting a link from the currently logged in user to a subreddit

Immediate TODO:
 - Submitting comments
 - User account controls (mail/deletion/creation etc)
 - Captcha

Example Code: 
```go
package main 

import (
  "github.com/Stringy/reddit"
)

func main() {
  user, err := reddit.Login("username", "password", true)
  if err != nil {
    // handle error
  } 
  err = user.SubmitLink(
      "learnprogramming", 
      "Amazing Search Engine",
	    "",
	    "www.google.com",
      true
      )
  if err != nil {
    // handle error
  }
}
```