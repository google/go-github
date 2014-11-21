package github

import (
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"
)

func TestObeyRateLimit(t *testing.T) {
	setup()
	defer teardown()

	limit := 60
	s := &rateLimitServer{limit, limit, time.Now()}
	mux.Handle("/", s)

	client.wait = func(d time.Duration) {
		t.Logf("waiting %s", d)
		s.Remaining++
	}

	client.ObeyRateLimit = true
	for i := 0; i < limit*2; i++ {
		t.Logf("-- request %d", i)
		if _, _, err := client.Octocat("foo"); err != nil {
			t.Errorf("unexpected error: %v", err)
			break
		}
	}

	// Turn off automatic rate limit enforcement and make too many requests
	// to demonstrate that the alternative is over-quota errors.
	client.ObeyRateLimit = false
	var err error
	var resp *Response
	for i := 0; i < limit*2; i++ {
		_, resp, err = client.Octocat("foo")
		if err != nil {
			break
		}
	}
	if resp.StatusCode != http.StatusForbidden {
		t.Errorf("expected to be rate limited")
	}
}

type rateLimitServer struct {
	Limit, Remaining int
	Reset            time.Time
}

func (rl *rateLimitServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(headerRateLimit, fmt.Sprintf("%d", rl.Limit))
	w.Header().Set(headerRateRemaining, fmt.Sprintf("%d", rl.Remaining))

	tokensNeeded := rl.Limit - rl.Remaining
	tokensPerSecond := rl.Limit / 60 / 60
	timeUntilReset := time.Duration(tokensNeeded*tokensPerSecond) * time.Second
	rl.Reset = time.Now().Add(timeUntilReset)
	w.Header().Set(headerRateReset, fmt.Sprintf("%d", rl.Reset))

	if r.URL.Path != "/rate_limit" {
		if rl.Remaining <= 0 {
			http.Error(w, "Over quota", http.StatusForbidden)
		}
		rl.Remaining--
	}
	io.WriteString(w, "{}")
}
