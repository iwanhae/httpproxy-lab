package main

import (
	"net/http"
	"os"
	"time"

	proxy "github.com/go-httpproxy/httpproxy"
	"github.com/justinas/alice"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/hlog"
)

func main() {

	// Proxy
	prx, err := proxy.NewProxy()
	if err != nil {
		panic(err)
	}

	// Logger
	log := zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr}).With().Caller().Timestamp().Logger()

	h := alice.New().
		Append(hlog.NewHandler(log)).
		Append(hlog.AccessHandler(func(r *http.Request, status, size int, duration time.Duration) {
			hlog.FromRequest(r).Info().
				Str("method", r.Method).
				Stringer("url", r.URL).
				Int("status", status).
				Int("size", size).
				Dur("duration", duration).
				Str("host", r.Host).
				Msg("")
		})).
		Append(hlog.RemoteAddrHandler("ip")).
		Append(hlog.UserAgentHandler("user_agent")).
		Append(hlog.RefererHandler("referer")).
		Append(hlog.RequestIDHandler("req_id", "Request-Id")).
		Append(func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				prx.ServeHTTP(w, r)
				next.ServeHTTP(w, r)
			})
		}).
		Then(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			hlog.FromRequest(r).Info().
				Msg("new requests")
		}))

	log.Info().Str("addr", ":8080").Msg("listen")
	log.Fatal().Err(http.ListenAndServe(":8080", h))
}
