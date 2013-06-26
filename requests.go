package reddit

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sync"
)

var (
	get       *http.Request
	post      *http.Request
	post_form *http.Request
	client    *http.Client
)

func init() {
	jar := NewJar()
	client = &http.Client{nil, nil, jar}
	get, err := http.NewRequest("GET", "", nil)
	if err != nil {
		panic(err)
	}
	post, err = http.NewRequest("POST", "", nil)
	if err != nil {
		panic(err)
	}
	post_form, err = http.NewRequest("POST", "", nil)
	if err != nil {
		panic(err)
	}
	post_form.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	get.Header.Set("User-Agent", UserAgent)
	post.Header.Set("User-Agent", UserAgent)
	post_form.Header.Set("User-Agent", UserAgent)
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
	get.Header.Set("User-Agent", UserAgent)
	post.Header.Set("User-Agent", UserAgent)
	post_form.Header.Set("User-Agent", UserAgent)
}

func getPostJsonBytes(u *url.URL, data *url.Values) ([]byte, error) {
	post_form.URL = u
	post_form.Body = &ClosingBuffer{bytes.NewBufferString(data.Encode())}
	resp, err := client.Do(post_form)
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(
			fmt.Sprintf("http: unexpected status code from request: %d", resp.StatusCode))
	}
	if err != nil {
		return nil, err
	}
	client.Jar.SetCookies(actual_url, resp.Cookies())
	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, resp.Body)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func getJsonBytes(u *url.URL) ([]byte, error) {
	get.URL = u
	resp, err := client.Do(get)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(
			fmt.Sprintf("http: unexpected status code from request: %d", resp.StatusCode))
	}
	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, resp.Body)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
