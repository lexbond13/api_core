package logger

import (
	"github.com/lexbond13/api_core/config"
	"github.com/getsentry/raven-go"
	"github.com/getsentry/sentry-go"
	"net/http"
)

type Sentry struct {
	dsn string
	client *raven.Client
}

// init sentry for logging events and requests
func NewSentry(config *config.SentryConfig, isDebug bool) error {
	// init instance for listen panic
	err := sentry.Init(sentry.ClientOptions{
		Dsn: config.DSN,
		BeforeSend: func(event *sentry.Event, hint *sentry.EventHint) *sentry.Event {
			if hint.Context != nil {
				if _, ok := hint.Context.Value(sentry.RequestContextKey).(*http.Request); ok {
					// You have access to the original Request
					//fmt.Println(req)
				}
			}
			//fmt.Println(event)
			return event
		},
		Debug:            isDebug,
		AttachStacktrace: true,
	})

	if err != nil {
		return err
	}

	return nil
}
