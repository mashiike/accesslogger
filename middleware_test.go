package accesslogger_test

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/mashiike/accesslogger"
	"github.com/sebdah/goldie/v2"
)

func TestWrap(t *testing.T) {
	var past time.Duration
	base, _ := time.Parse(time.RFC3339, "2022-12-26T15:04:05+09:00")
	accesslogger.Clock = func() time.Time {
		return base.Add(past)
	}

	var combinedLogs bytes.Buffer
	var combinedDLogs bytes.Buffer
	var jsonLogs bytes.Buffer
	h := accesslogger.Wrap(
		http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			time.Sleep(10 * time.Millisecond)
			past += 10 * time.Millisecond
			fmt.Fprintf(w, "hoge")
		}),
		accesslogger.CombinedLogger(&combinedLogs),
		accesslogger.CombinedDLogger(&combinedDLogs),
		accesslogger.JSONLogger(&jsonLogs),
	)
	g := goldie.New(t)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	r.Header.Set("User-Agent", "go test client")
	h.ServeHTTP(w, r)
	r.Header.Set("Referer", "https://example.com")
	r.URL.Path = "/hoge"
	r.SetBasicAuth("hoge", "fuga")
	past = 1 * time.Second
	h.ServeHTTP(w, r)
	r.Header.Set("CloudFront-Viewer-Address", "222.222.333.333")
	h.ServeHTTP(w, r)

	g.Assert(t, "combined", combinedLogs.Bytes())
	g.Assert(t, "combinedD", combinedDLogs.Bytes())
	g.Assert(t, "json", jsonLogs.Bytes())
}
