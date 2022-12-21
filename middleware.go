package accesslogger

import (
	"net/http"
)

func New(loggers ...Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			l := NewAccessLog(r)
			responseWriter := &ResponseWriter{
				ResponseWriter: w,
			}
			defer func() {
				err := recover()
				l = l.WriteResponseInfo(responseWriter)
				for _, logger := range loggers {
					logger.WriteAccessLog(l)
				}
				if err != nil {
					panic(err)
				}
			}()
			next.ServeHTTP(responseWriter, r)
		})
	}
}

func Wrap(next http.Handler, loggers ...Logger) http.Handler {
	m := New(loggers...)
	return m(next)
}
