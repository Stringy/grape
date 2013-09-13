package client

import (
	"net/http"
	"sync"
)

//RedditCache is a map of strings (urls) to their responses
//it is used to contain prefetch data from reddit, to minimise wasted time in requests
type RedditCache struct {
	cache  map[string]*http.Response
	update chan struct{}
	*sync.RWMutex
}

//NewRedditCache returns a new initialised reddit cache
func NewRedditCache() *RedditCache {
	rc := new(RedditCache)
	rc.cache = make(map[string]*http.Response, cacheSize)
	rc.update = make(chan struct{})
	rc.RWMutex = new(sync.RWMutex)
	return rc
}

//Update signals that the cache has changed by closing the update channel
func (r *RedditCache) Update() {
	r.Lock()
	defer r.Unlock()
	close(r.update)
	r.update = make(chan struct{})
}

//GetUpdateChan gives the current update channel in the cache
func (r *RedditCache) GetUpdateChan() chan struct{} {
	r.RLock()
	defer r.RUnlock()
	return r.update
}

//Set sets the key, value pair in the cache
func (r *RedditCache) Set(key string, value *http.Response) {
	r.Lock()
	defer r.Unlock()
	r.cache[key] = value
}

//Get returns the value for key, if it exists along with a bool denoting
//such existence
func (r *RedditCache) Get(key string) (*http.Response, bool) {
	r.Lock()
	defer r.Unlock()
	resp, exists := r.cache[key]
	return resp, exists
}
