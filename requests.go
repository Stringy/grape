package reddit

import (
	"bytes"
	"net/http"
	"net/url"
	"sync"
	"time"
)

type priority chan *http.Request

var (
	priorities    []priority
	responseCache RedditCache
)

const (
	dontPrefetch = iota
	commentReq
	listingReq
	userReq

	cacheSize = 25
)

func init() {
	priorities = make([]priority, 3)
	responseCache = *new(RedditCache)
	responseCache.cache = make(map[string]*http.Response, cacheSize)
}

type RedditCache struct {
	cache map[string]*http.Response
	*sync.RWMutex
}

func (r *RedditCache) Get(url string) (*http.Response, error) {
	return nil, nil
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
	priorities[0] <- req
	return nil, nil
}

func cacheRequests(resps chan *http.Response) {
	for {
		select {
		case resp := <-resps:
			// cache response for retrieval
			_ = resp
		default:
		}
	}
}

func makeRequests(resps chan *http.Response, errs chan error) {
	schedule := time.Tick(2 * time.Second)
	for _ = range schedule {
		for i := range priorities {
			select {
			case req := <-priorities[i]:
				resp, err := client.Do(req)
				if err != nil {
					panic(err)
				}
				resps <- resp
			default:
			}
		}
		select {
		case req := <-priorities[0]: //high
			resp, err := client.Do(req)
			if err != nil {
				panic(err)
			}
			resps <- resp
		case req := <-priorities[1]: //medium
			resp, err := client.Do(req)
			if err != nil {
				panic(err)
			}
			resps <- resp
		case req := <-priorities[2]: //low
			resp, err := client.Do(req)
			if err != nil {
				panic(err)
			}
			resps <- resp
		}
	}
}
