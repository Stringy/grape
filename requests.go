package reddit

import (
	"bytes"
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

	resps = make(chan *http.Response)
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
	responseCache = new(RedditCache)
	responseCache.cache = make(map[string]*http.Response, cacheSize)
	responseCache.RWMutex = new(sync.RWMutex)
	lf, err := os.Create("reddit.log")
	if err != nil {
		panic(err)
	}
	log.SetOutput(lf)
	go cacheResponses()
	go makeRequests()
}

type RedditCache struct {
	cache map[string]*http.Response
	*sync.RWMutex
}

func (r *RedditCache) Set(key string, value *http.Response) {
	r.Lock()
	defer r.Unlock()
	r.cache[key] = value
}

func (r *RedditCache) Get(key string) *http.Response {
	r.Lock()
	defer r.Unlock()
	return r.cache[key]
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

// func getPostJsonBytes(link string, data *url.Values) ([]byte, error) {
// 	u, err := url.Parse(link)
// 	if err != nil {
// 		return nil, err
// 	}
// 	post_form.URL = u
// 	post_form.Host = u.Host
// 	content_len := len(data.Encode())
// 	post_form.ContentLength = int64(content_len)
// 	post_form.Body = &ClosingBuffer{bytes.NewBufferString(data.Encode())}
// 	resp, err := client.Do(post_form)
// 	if err != nil {
// 		return nil, err
// 	}
// 	//	fmt.Println(resp)
// 	if resp.StatusCode != http.StatusOK {
// 		return nil, errors.New(
// 			fmt.Sprintf("http: unexpected status code from request: %d", resp.StatusCode))
// 	}
// 	if len(client.Jar.Cookies(actual_url)) == 0 {
// 		client.Jar.SetCookies(actual_url, resp.Cookies())
// 	}
// 	buf := new(bytes.Buffer)
// 	_, err = io.Copy(buf, resp.Body)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return buf.Bytes(), nil
// }

// func getJsonBytes(link string) ([]byte, error) {
// 	u, err := url.Parse(link)
// 	if err != nil {
// 		return nil, err
// 	}
// 	get.URL = u
// 	get.Host = u.Host
// 	resp, err := client.Do(get)
// 	if err != nil {
// 		return nil, err
// 	}
// 	if resp.StatusCode != http.StatusOK {
// 		return nil, errors.New(
// 			fmt.Sprintf("http: unexpected status code from request: %d", resp.StatusCode))
// 	}
// 	buf := new(bytes.Buffer)
// 	_, err = io.Copy(buf, resp.Body)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return buf.Bytes(), nil
// }

func makePostRequest(link string, data *url.Values) ([]byte, error) {
	req, err := http.NewRequest("POST", link, &ClosingBuffer{bytes.NewBufferString(data.Encode())})
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", UserAgent)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	log.Printf("Sending POST request to scheduler\n")
	priorities[0] <- req
	return nil, nil
}

func makeGetRequest(url string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", UserAgent)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	log.Printf("Sendinng GET request to scheduler")
	priorities[0] <- req
	return nil, nil
}

func cacheResponses() {
	for {
		select {
		case resp := <-resps:
			if resp != nil {
				u := resp.Request.URL.String()
				responseCache.Set(u, resp)
			}
		default:
		}
	}
}

func makeRequests() {
	schedule := time.Tick(2 * time.Second)
	log.Printf("Entering For\n")
	var req *http.Request
	for {
		select {
		case <-schedule:
			log.Println("Recieved from ticker")
			for i := 0; i < len(priorities); i++ {
				select {
				case req = <-priorities[i]:
					log.Printf("Making request of priority: %d\n", i)
					go doRequest(req)
					goto CONTINUE
				default:
				}
			}
			log.Println("Trying first possible job")
			select {
			case req = <-priorities[0]: //high
				log.Printf("Making high priority request to %v\n", req.URL)
				go doRequest(req)
				goto CONTINUE
			case req = <-priorities[1]: //medium
				log.Printf("Making medium priority request to %v\n", req.URL)
				go doRequest(req)
				goto CONTINUE
			case req = <-priorities[2]: //low
				log.Printf("Making low priority request to %v\n", req.URL)
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
	log.Printf("Doing request: %v", req.URL)
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("error in response from %v\n", req.URL)
	}
	resps <- resp
}
