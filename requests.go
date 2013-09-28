package reddit

import (
	"bytes"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

var (
	priorities    []chan *http.Request // priorities channel
	responseCache *redditCache         // cache of prefetched results

	resps       = make(chan *http.Response) // responses to be cached
	cacheUpdate = make(chan struct{})       // signal chan for signifying updated cache

	client *http.Client // http client for making requests
)

const (
	// prefetching type enumeration
	// used for future prefetching logic
	dontPrefetch = iota
	commentReq
	listingReq
	userReq

	// default cache size
	cacheSize = 25
)

func init() {
	priorities = []chan *http.Request{
		make(chan *http.Request),
		make(chan *http.Request),
		make(chan *http.Request),
	}
	client = &http.Client{nil, nil, newJar()}
	responseCache = newRedditCache()
	go cacheResponses()
	go makeRequests()
}

//Jar is an implementation of a CookieJar for use in the http client
type jar struct {
	sync.Mutex
	cookies map[string][]*http.Cookie
}

// NewJar creates and returns a new Cookie Jar
func newJar() *jar {
	jar := new(jar)
	jar.cookies = make(map[string][]*http.Cookie)
	return jar
}

// SetCookies handles the receipt of the cookies in a reply for the
// given URL.  It may or may not choose to save the cookies, depending
// on the jar's policy and implementation.
func (jar *jar) SetCookies(u *url.URL, cookies []*http.Cookie) {
	jar.Lock()
	jar.cookies[u.Host] = cookies
	jar.Unlock()
}

// Cookies returns the cookies to send in a request for the given URL.
// It is up to the implementation to honor the standard cookie use
// restrictions such as in RFC 6265.
func (jar *jar) Cookies(u *url.URL) []*http.Cookie {
	return jar.cookies[u.Host]
}

// MakePostRequest adds a post request to the request schedule and waits for the
// existence of a response in the cache
func makePostRequest(url string, data *url.Values) ([]byte, error) {
	req, err := http.NewRequest("POST", url, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", config.UserAgent)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	//req.Header.Set("Content-Length", fmt.Sprintf("%d", len(data.Encode())))
	//debug.Println(req)
	priorities[0] <- req
	cache := responseCache.GetUpdateChan()
	for {
		select {
		case _, ok := <-cache:
			if !ok {
				resp, exists := responseCache.Get(req.URL.String())
				if exists {
					debug.Printf("cache retrieved desired response for %v", req.URL)
					buf := new(bytes.Buffer)
					_, err := io.Copy(buf, resp.Body)
					if err != nil {
						return nil, err
					}
					return buf.Bytes(), nil
				}
			}
		default:
		}
	}
	return nil, nil
}

// MakeGetRequest adds a get request to the request schedule and waits for the
// existence of a response in the cache
func makeGetRequest(url string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", config.UserAgent)
	//req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	priorities[0] <- req
	cache := responseCache.GetUpdateChan()
	for {
		select {
		case _, ok := <-cache:
			if !ok {
				resp, exists := responseCache.Get(req.URL.String())
				if exists {
					debug.Printf("retrieved desired response for %v from cache.", req.URL)
					buf := new(bytes.Buffer)
					_, err := io.Copy(buf, resp.Body)
					if err != nil {
						return nil, err
					}
					return buf.Bytes(), nil
				}
			}
		}
	}
	return nil, nil
}

// cacheResponses is run as a go routine upon startup. it waits for a response
// and adds it to the cache.
func cacheResponses() {
	for {
		select {
		case resp := <-resps:
			if resp != nil {
				u := resp.Request.URL.String()
				debug.Printf("caching response from %s\n", u)
				responseCache.Set(u, resp)
				responseCache.Update()
			}
		}
	}
}

//makeRequests is run as a go routine. It checks for a job every two seconds to
//conform to the reddit API. It then takes a job from the priority channels and
//starts a routine to process it.
func makeRequests() {
	schedule := time.Tick(2 * time.Second)
	var req *http.Request
	for {
		select {
		case <-schedule:
			//check for jobs in order of priority
			for i := 0; i < len(priorities); i++ {
				select {
				case req = <-priorities[i]:
					go doRequest(req)
				default:
				}
			}
			//take first available job
			select {
			case req = <-priorities[0]: //high
				go doRequest(req)
			case req = <-priorities[1]: //medium
				go doRequest(req)
			case req = <-priorities[2]: //low
				go doRequest(req)
			}
		}
	}
}

// doRequest is called in a new go routine to send the request to reddit
// it then updates the client cookies and sends the response to be cached for
// retrieval
func doRequest(req *http.Request) {
	debug.Printf("client doing request: %v", req.URL)
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("error in response from %v\n\t%v\n", req.URL, err)
	}
	//	debug.Println(resp)
	if len(resp.Cookies()) != 0 {
		client.Jar.SetCookies(reddit_url, resp.Cookies())
	}
	resps <- resp
}
