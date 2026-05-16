// Package proxy
package proxy

import (
	"io"
	"log/slog"
	"net"
	"net/http"
	"net/url"
	"strings"
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

type ReverseProxy struct {
	Upstream *url.URL
	Client   *http.Client
	Logger   *slog.Logger
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

func (p *ReverseProxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	outURL := *p.Upstream
	outURL.Path = r.URL.Path
	outURL.RawQuery = r.URL.RawQuery

	outReq, err := http.NewRequestWithContext(r.Context(), r.Method, outURL.String(), r.Body)
	if err != nil {
		p.Logger.Error("Bad Gateway", "error", err.Error(), "status", http.StatusBadGateway)
		http.Error(w, "Bad Gateway: "+err.Error(), http.StatusBadGateway)
		return
	}

	for k, vs := range r.Header {
		for _, v := range vs {
			outReq.Header.Add(k, v)
		}
	}
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

	for k, vs := range resp.Header {
		for _, v := range vs {
			w.Header().Add(k, v)
		}
	}
	StripHopByHop(w.Header())

	w.WriteHeader(resp.StatusCode)
	if _, err := io.Copy(w, resp.Body); err != nil {
		p.Logger.Error("Copy error after status sent", "error", err)
	}
}
