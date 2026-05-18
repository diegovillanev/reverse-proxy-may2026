// Package proxy
package proxy

import (
	"bytes"
	"container/list"
	"io"
	"log/slog"
	"net"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

var HopHeaders = []string{
	"Connection",
	"Keep-Alive",
	"Proxy-Authenticate",
	"Proxy-Authorization",
	"Te",
	"Trailer",
	"Transfer-Encoding",
	"Upgrade",
}

func StripHopByHop(h http.Header) {
	if conn := h.Get("Connection"); conn != "" {
		for _, name := range strings.Split(conn, ",") {
			if name = strings.TrimSpace(name); name != "" {
				h.Del(name)
			}
		}
	}
	for _, name := range HopHeaders {
		h.Del(name)
	}
}

// ----------------------------------------------------------------------
// LRU Cache, better than the TTL Map cache
// ----------------------------------------------------------------------

type Entry struct {
	Key     string
	Status  int
	Header  http.Header
	Body    []byte
	Expires time.Time
}

type LRUCache struct {
	Mu      sync.Mutex
	TTL     time.Duration
	MaxSize int
	Entries map[string]*list.Element
	Order   *list.List
}

func NewCache(ttl time.Duration, maxSize int) *LRUCache {
	return &LRUCache{
		TTL:     ttl,
		MaxSize: maxSize,
		Entries: make(map[string]*list.Element),
		Order:   list.New(),
	}
}

func (c *LRUCache) Get(key string) (*Entry, bool) {
	c.Mu.Lock()
	defer c.Mu.Unlock()
	elem, ok := c.Entries[key]

	if !ok {
		return nil, false
	}
	entry := elem.Value.(*Entry)
	if time.Now().After(entry.Expires) {
		c.remove(elem)
		return nil, false
	}
	c.Order.MoveToFront(elem)
	return entry, true
}

func (c *LRUCache) Put(key string, newEntry *Entry) {
	newEntry.Expires = time.Now().Add(c.TTL)
	c.Mu.Lock()
	defer c.Mu.Unlock()
	if elem, ok := c.Entries[key]; ok {
		c.Order.MoveToFront(elem)
		elem.Value = newEntry // to update the TTL
		return
	}
	if c.Order.Len() >= c.MaxSize {
		c.remove(c.Order.Back())
	}
	elem := c.Order.PushFront(newEntry)
	c.Entries[key] = elem
}

func (c *LRUCache) remove(elem *list.Element) {
	entry := elem.Value.(*Entry)
	println("Evicting:", entry.Key, ", Len before:", c.Order.Len())
	delete(c.Entries, entry.Key)
	c.Order.Remove(elem)
	println("Len after:", c.Order.Len())
}

func CacheKey(r *http.Request) string {
	return r.Method + " " + r.URL.RequestURI()
}

func IsCacheable(r *http.Request, resp *http.Response) bool {
	if r.Method != http.MethodGet && r.Method != http.MethodHead {
		return false
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return false
	}
	cc := resp.Header.Get("Cache-Control")
	if strings.Contains(cc, "no-store") || strings.Contains(cc, "no-cache") {
		return false
	}
	return true
}

func WriteEntry(w http.ResponseWriter, e *Entry, status string) {
	for k, vs := range e.Header {
		for _, v := range vs {
			w.Header().Add(k, v)
		}
	}
	w.Header().Set("X-Cache", status)
	w.WriteHeader(e.Status)
	_, _ = io.Copy(w, bytes.NewReader(e.Body))
}

// ----------------------------------------------------------------------
// Proxy.
// ----------------------------------------------------------------------

type ReverseProxy struct {
	Upstream *url.URL
	Client   *http.Client
	Logger   *slog.Logger
	Cache    *LRUCache
}

func (p *ReverseProxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Create a Cache Key to look for a cached request first
	key := CacheKey(r)

	// We look up into the Cache Map using a Mutex to protect against race conditions
	// If we find one, we send it and log accordingly.
	if hit, ok := p.Cache.Get(key); ok {
		WriteEntry(w, hit, "HIT")
		p.Logger.LogAttrs(
			r.Context(),
			slog.LevelInfo,
			"Cache HIT ",
			slog.String("method", r.Method),
			slog.String("uri", r.URL.RequestURI()),
		)
		return
	}

	outURL := *p.Upstream
	outURL.Path = r.URL.Path
	outURL.RawQuery = r.URL.RawQuery

	outReq, err := http.NewRequestWithContext(r.Context(), r.Method, outURL.String(), r.Body)
	if err != nil {
		p.Logger.Error("Bad Gateway", "error", err.Error(), "status", http.StatusBadGateway)
		http.Error(w, "Bad Gateway: "+err.Error(), http.StatusBadGateway)
		return
	}

	outReq.Header = r.Header.Clone()
	StripHopByHop(outReq.Header)

	outReq.Host = p.Upstream.Host
	if clientIP, _, splitErr := net.SplitHostPort(r.RemoteAddr); splitErr == nil {
		if prior := outReq.Header.Get("X-Forwarded-For"); prior != "" {
			clientIP = prior + ", " + clientIP
		}
		outReq.Header.Set("X-Forwarded-For", clientIP)
	}

	resp, err := p.Client.Do(outReq)
	if err != nil {
		p.Logger.Error("Bad Gateway", "error", err.Error(), "status", http.StatusBadGateway)
		http.Error(w, "Bad Gateway: "+err.Error(), http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		p.Logger.Error("Bad Gateway", "error", err.Error(), "status", http.StatusBadGateway)
		http.Error(w, "Bad Gateway: "+err.Error(), http.StatusBadGateway)
	}

	// This seems a far better alternative than the for loop
	hdrs := resp.Header.Clone()
	StripHopByHop(hdrs)
	e := &Entry{Key: key, Status: resp.StatusCode, Header: hdrs, Body: body}

	if IsCacheable(r, resp) {
		p.Cache.Put(key, e)
	}

	WriteEntry(w, e, "MISS")
	p.Logger.LogAttrs(
		r.Context(),
		slog.LevelInfo,
		"Cache MISS",
		slog.String("method", r.Method),
		slog.String("uri", r.URL.RequestURI()),
		slog.Int("status-code", resp.StatusCode),
		slog.Int("body-size", len(body)),
	)
}
