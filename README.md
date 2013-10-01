grape
======

grape is a Go wrapper for the Reddit API. It will provide functionality for all aspects of the API for use in the creation of bots. 

The current functionality is limited to the following:
 - Logging in
 - Retrieving information about any user
 - Retrieving information about the currently logged in user
 - Retrieving a subreddit (limited to front page only)
 - Retrieving sorted subreddits (by new, hot and top)
 - Retrieving comments from a link
 - Submitting a link from the currently logged in user to a subreddit
 - Submitting comments
 - Getting all mail (inbox as well as unread)
 - Cached responses from reddit and a priority system in preparation for prefetching of data

Immediate TODO:
 - Prefetching of extra information
 - Replying to private messages
 - Develop a way of incorporating OAuth

Example Code: 
```go
package main 

import (
  "github.com/Stringy/grape"
  "fmt"
)

func main() {
  user, err := grape.Login("username", "password", true)
  if err != nil {
    // handle error
  } 
  sub, err := grape.GetSubreddit("learnprogramming")
  if err != nil {
    // handle error
  }
  for i, item := range sub.Items {
    fmt.Println(item.Submission)
  }
}
```