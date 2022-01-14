package sentry

import (
	"fmt"
	"net/http"
	"time"

	"github.com/getsentry/sentry-go"
	goSentry "github.com/getsentry/sentry-go"
	sentryecho "github.com/getsentry/sentry-go/echo"
	logDump "github.com/sirupsen/logrus"
)

func New(config goSentry.ClientOptions) {
	goSentry.Init(goSentry.ClientOptions{
		Dsn:        config.Dsn,
		HTTPProxy:  config.HTTPProxy,
		HTTPSProxy: config.HTTPSProxy,
		TracesSampler: goSentry.TracesSamplerFunc(func(ctx goSentry.SamplingContext) goSentry.Sampled {
			return goSentry.SampledTrue
		}),
		BeforeSend: func(event *goSentry.Event, hint *goSentry.EventHint) *goSentry.Event {
			if hint.Context != nil {
				if req, ok := hint.Context.Value(goSentry.RequestContextKey).(*http.Request); ok {
					// You have access to the original Request
					logDump.Info(req)
				}
			}
			logDump.Info(event)
			return event
		},
		Debug:            true,
		AttachStacktrace: true,
	})
	defer goSentry.Flush(2 * time.Second)
}

func MiddlewareSentry() {
	if hub := sentryecho.GetHubFromContext(ctx); hub != nil {
		var (
			userId = fmt.Sprintf("%v", ctx.Get("RequestID"))
		)

		if ctx.Get("UserId") != nil {
			userId = fmt.Sprintf("%v", fmt.Sprintf("%v", ctx.Get("UserId")))
		}

		if hub := sentryecho.GetHubFromContext(ctx); hub != nil {

			hub := sentryecho.GetHubFromContext(ctx)

			hub.Scope().SetTransaction(fmt.Sprintf("%s", ctx.Path()))
			hub.Scope().SetUser(sentry.User{
				ID:        userId,
				IPAddress: m.GetLocalIP(),
			})
			hub.Scope().SetLevel(sentry.LevelError)
			hub.Scope().SetRequest(ctx.Request())
			sentry.Logger.SetFlags(time.Now().Minute())
			sentry.Logger.SetPrefix("[sentry SDK]")
		}
	}
}
