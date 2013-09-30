package grape

import (
	"net/http"
	"strings"
	"sync"
	"time"
)

// minimum time between fetching from reddit
const expirationTime = time.Second * 30

// default cache size
const cacheSize = 25

// A cache entry represents the http response recieved from reddit
// along with the time it was retrieved
// if an expired entry is requested, it will be re-sent to reddit
// otherwise the cached response is sent.
type cacheEntry struct {
	resp      *http.Response
	retrieved time.Time
}

// RedditCache is a map of strings (urls) to their responses
// it is used to contain prefetch data from reddit, to minimise wasted time in requests
type redditCache struct {
	cache  map[string]*cacheEntry
	update chan struct{}
	*sync.RWMutex
}

// NewRedditCache returns a new initialised reddit cache
func newRedditCache() *redditCache {
	rc := new(redditCache)
	rc.cache = make(map[string]*cacheEntry, cacheSize)
	rc.update = make(chan struct{})
	rc.RWMutex = new(sync.RWMutex)
	return rc
}

// Update signals that the cache has changed by closing the update channel
func (r *redditCache) Update() {
	r.Lock()
	defer r.Unlock()
	close(r.update)
	r.update = make(chan struct{})
}

// GetUpdateChan gives the current update channel in the cache
func (r *redditCache) GetUpdateChan() chan struct{} {
	r.RLock()
	defer r.RUnlock()
	return r.update
}

// Set sets the key, value pair in the cache.
// If the url is an api call, the retrieval time is set as Now - expirationTime
// to prevent the blocking of api calls. The 30 second rule applies only to normal page retrieval.
func (r *redditCache) Set(key string, value *http.Response) {
	r.Lock()
	defer r.Unlock()
	if strings.Contains(key, "/api/") {
		r.cache[key] = &cacheEntry{value, time.Now().Add(-expirationTime)}
	} else {
		r.cache[key] = &cacheEntry{value, time.Now()}
	}
}

// Get returns the value for key, if it exists along with a bool denoting
// such existence
func (r *redditCache) Get(key string) (*http.Response, bool) {
	r.Lock()
	defer r.Unlock()
	entry, exists := r.cache[key]
	return entry.resp, exists
}

// IsExpired check a cache entry for expiration.
// if no entry exists for key exists, then the function returns true
// to ensure it is then fetched.
// If the request is to /api/login it will always return true
func (r *redditCache) IsExpired(key string) bool {
	r.Lock()
	defer r.Unlock()
	entry, exists := r.cache[key]
	if !exists {
		return true
	}
	if entry.retrieved.Add(expirationTime).After(time.Now()) {
		return false
	}
	return true
}
