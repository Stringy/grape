package reddit

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"sync"
	"time"
)

type priority chan *http.Request

var (
	priorities    []chan *http.Request
	responseCache *RedditCache

	resps       = make(chan *http.Response)
	cacheUpdate = make(chan struct{})
)

const (
	dontPrefetch = iota
	commentReq
	listingReq
	userReq

	cacheSize = 25
)

func init() {
	priorities = []chan *http.Request{
		make(chan *http.Request),
		make(chan *http.Request),
		make(chan *http.Request),
	}
	responseCache = NewRedditCache()
	lf, err := os.Create("reddit.log")
	if err != nil {
		panic(err)
	}
	log.SetOutput(lf)
	go cacheResponses()
	go makeRequests()
}

type RedditCache struct {
	cache  map[string]*http.Response
	update chan struct{}
	*sync.RWMutex
}

func NewRedditCache() *RedditCache {
	rc := new(RedditCache)
	rc.cache = make(map[string]*http.Response, cacheSize)
	rc.update = make(chan struct{})
	rc.RWMutex = new(sync.RWMutex)
	return rc
}

func (r *RedditCache) Update() {
	r.Lock()
	defer r.Unlock()
	close(r.update)
	r.update = make(chan struct{})
}

func (r *RedditCache) GetUpdateChan() chan struct{} {
	r.RLock()
	defer r.RUnlock()
	return r.update
}

func (r *RedditCache) Set(key string, value *http.Response) {
	r.Lock()
	defer r.Unlock()
	r.cache[key] = value
}

func (r *RedditCache) Get(key string) (*http.Response, bool) {
	r.Lock()
	defer r.Unlock()
	resp, exists := r.cache[key]
	return resp, exists
}

type ClosingBuffer struct {
	*bytes.Buffer
}

func (cb *ClosingBuffer) Close() error {
	return nil
}

type Jar struct {
	sync.Mutex
	cookies map[string][]*http.Cookie
}

func NewJar() *Jar {
	jar := new(Jar)
	jar.cookies = make(map[string][]*http.Cookie)
	return jar
}

// SetCookies handles the receipt of the cookies in a reply for the
// given URL.  It may or may not choose to save the cookies, depending
// on the jar's policy and implementation.
func (jar *Jar) SetCookies(u *url.URL, cookies []*http.Cookie) {
	jar.Lock()
	jar.cookies[u.Host] = cookies
	jar.Unlock()
}

// Cookies returns the cookies to send in a request for the given URL.
// It is up to the implementation to honor the standard cookie use
// restrictions such as in RFC 6265.
func (jar *Jar) Cookies(u *url.URL) []*http.Cookie {
	return jar.cookies[u.Host]
}

func SetUserAgent(ua string) {
	UserAgent = ua
}

func makePostRequest(link string, data *url.Values) ([]byte, error) {
	req, err := http.NewRequest("POST", link, &ClosingBuffer{bytes.NewBufferString(data.Encode())})
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", UserAgent)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	priorities[0] <- req
	cache := responseCache.GetUpdateChan()
	for {
		select {
		case _, ok := <-cache:
			if !ok {
				resp, exists := responseCache.Get(req.URL.String())
				if exists {
					log.Printf("[CACHE] retrieved desired response for %v", req.URL)
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

func makeGetRequest(url string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", UserAgent)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	priorities[0] <- req
	cache := responseCache.GetUpdateChan()
	for {
		select {
		case _, ok := <-cache:
			if !ok {
				resp, exists := responseCache.Get(req.URL.String())
				if exists {
					log.Printf("[CACHE] retrieved desired response for %v", req.URL)
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

func cacheResponses() {
	for {
		select {
		case resp := <-resps:
			if resp != nil {
				u := resp.Request.URL.String()
				log.Printf("[CACHE] caching response from %s\n", u)
				responseCache.Set(u, resp)
				responseCache.Update()
			}
		default:
		}
	}
}

func makeRequests() {
	schedule := time.Tick(2 * time.Second)
	var req *http.Request
	for {
		select {
		case <-schedule:
			for i := 0; i < len(priorities); i++ {
				select {
				case req = <-priorities[i]:
					go doRequest(req)
					goto CONTINUE
				default:
				}
			}
			select {
			case req = <-priorities[0]: //high
				go doRequest(req)
				goto CONTINUE
			case req = <-priorities[1]: //medium
				go doRequest(req)
				goto CONTINUE
			case req = <-priorities[2]: //low
				go doRequest(req)
				goto CONTINUE
			default:
			}
		default:
		}
	CONTINUE:
	}
}

func doRequest(req *http.Request) {
	log.Printf("[CLIENT] doing request: %v", req.URL)
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("error in response from %v\n", req.URL)
	}
	resps <- resp
}
