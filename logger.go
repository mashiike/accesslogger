package accesslogger

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type AccessLog struct {
	RequestHeader  http.Header `json:"-"`
	ResponseHeader http.Header `json:"-"`
	RemoteAddr     string      `json:"remote_addr,omitempty"`
	AccessedAt     time.Time   `json:"accessed_at,omitempty"`
	UserAgent      string      `json:"user_agent,omitempty"`
	Referer        string      `json:"referer,omitempty"`
	BasicAuthUser  string      `json:"basic_auth_user,omitempty"`
	Request        string      `json:"request,omitempty"`
	StatusCode     int         `json:"status_code,omitempty"`
	BodyByteSent   int         `json:"body_byte_sent,omitempty"`
	FirstSentAt    time.Time   `json:"first_sent_at,omitempty"`
	LastSentAt     time.Time   `json:"last_sent_at,omitempty"`
	FirstSentTime  int64       `json:"first_sent_time,omitempty"`
	ResponseTime   int64       `json:"response_time,omitempty"`
}

func NewAccessLog(r *http.Request) *AccessLog {
	user, _, _ := r.BasicAuth()
	l := &AccessLog{
		RequestHeader: r.Header,
		RemoteAddr:    coalesce(r.Header.Get("CloudFront-Viewer-Address"), r.RemoteAddr),
		AccessedAt:    Clock(),
		UserAgent:     r.UserAgent(),
		Referer:       r.Referer(),
		BasicAuthUser: user,
		Request:       fmt.Sprintf("%s %s %s", r.Method, r.URL.Path, r.Proto),
	}
	return l
}

func (l *AccessLog) WriteResponseInfo(w *ResponseWriter) *AccessLog {
	l.StatusCode = w.StatusCode
	l.BodyByteSent = w.BodyByteSent
	l.FirstSentAt = w.FirstWriteTime
	l.LastSentAt = w.LastWriteTime
	l.FirstSentTime = l.FirstSentAt.Sub(l.AccessedAt).Microseconds()
	l.ResponseTime = l.LastSentAt.Sub(l.AccessedAt).Microseconds()
	l.ResponseHeader = w.Header()
	return l
}

type Logger interface {
	WriteAccessLog(*AccessLog)
}

type FormatLogger struct {
	io.Writer
	LogFormat
}

func (logger FormatLogger) WriteAccessLog(l *AccessLog) {
	fmt.Fprintln(logger.Writer, logger.LogFormat(l))
}

type LogFormat func(*AccessLog) string

func CombinedLogFormat(l *AccessLog) string {
	return fmt.Sprintf(
		`%s %s %s [%s] "%s" %s %s "%s" "%s"`,
		coalesce(l.RemoteAddr, "-"),
		"-",
		coalesce(l.BasicAuthUser, "-"),
		l.AccessedAt.Format("02/Jan/2006:15:04:05 -0700"),
		coalesce(l.Request, "-"),
		coalesce(emptyif(fmt.Sprintf("%d", l.StatusCode), "0"), "-"),
		coalesce(fmt.Sprintf("%d", l.BodyByteSent), "0"),
		coalesce(l.Referer, "-"),
		coalesce(l.UserAgent, "-"),
	)
}

func CombinedLogger(w io.Writer) FormatLogger {
	return FormatLogger{
		Writer:    w,
		LogFormat: CombinedLogFormat,
	}
}

func CombinedDLogFormat(l *AccessLog) string {
	return fmt.Sprintf(
		`%s %s`,
		CombinedLogFormat(l),
		coalesce(fmt.Sprintf("%d", l.ResponseTime), "0"),
	)
}

func CombinedDLogger(w io.Writer) FormatLogger {
	return FormatLogger{
		Writer:    w,
		LogFormat: CombinedDLogFormat,
	}
}

func JSONLogFormat(l *AccessLog) string {
	bs, err := json.Marshal(l)
	if err != nil {
		return `{}`
	}
	return string(bs)
}

func JSONLogger(w io.Writer) FormatLogger {
	return FormatLogger{
		Writer:    w,
		LogFormat: JSONLogFormat,
	}
}
