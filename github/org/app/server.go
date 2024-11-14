package app

import (
	"context"
	"net/http"
	"time"

	"github.com/gitrules/gitrules/lib/base"
	"github.com/gitrules/gitrules/lib/must"
	"github.com/gregjones/httpcache"
	"github.com/palantir/go-githubapp/githubapp"
	"github.com/rcrowley/go-metrics"
	"github.com/rs/zerolog"
)

func RunServer(ctx context.Context, addr string, cfg *Config) {

	metricsRegistry := metrics.DefaultRegistry

	cc, err := githubapp.NewDefaultCachingClientCreator(
		cfg.Github,
		githubapp.WithClientUserAgent("gitrules-for-organizations/1.0.0"),
		githubapp.WithClientTimeout(3*time.Second),
		githubapp.WithClientCaching(false, func() httpcache.Cache { return httpcache.NewMemoryCache() }),
		githubapp.WithClientMiddleware(
			// add logger from this context to request context
			func(next http.RoundTripper) http.RoundTripper {
				return roundTripperFunc(
					func(r *http.Request) (*http.Response, error) {
						logger := zerolog.Ctx(ctx)
						return next.RoundTrip(r.WithContext(logger.WithContext(r.Context())))
					},
				)
			},
			// log http headers
			// func(next http.RoundTripper) http.RoundTripper {
			// 	return roundTripperFunc(
			// 		func(r *http.Request) (*http.Response, error) {
			// 			base.Debugf("request_header: %v", r.Header)
			// 			res, err := next.RoundTrip(r)
			// 			base.Debugf("response_header: %v", res.Header)
			// 			return res, err
			// 		},
			// 	)
			// },
			githubapp.ClientMetrics(metricsRegistry),
			// githubapp.ClientLogging(
			// 	zerolog.DebugLevel,
			// 	githubapp.LogRequestBody(".*"),
			// 	githubapp.LogResponseBody(".*"),
			// ),
		),
	)
	must.NoError(ctx, err)

	webhookHandler := githubapp.NewDefaultEventDispatcher(
		cfg.Github,
		&InstallationHandler{
			ClientCreator: cc,
			DeployRelease: cfg.App.DeployRelease,
		},
	)

	http.Handle(githubapp.DefaultWebhookRoute, webhookHandler)

	base.Infof("Starting GitRules for Organizations cron ...")
	go func() {
		updateCron(ctx, cc)
	}()

	base.Infof("Starting GitRules for Organizations app server on %s ...", addr)
	must.NoError(ctx, http.ListenAndServe(addr, nil))
}

type roundTripperFunc func(*http.Request) (*http.Response, error)

func (fn roundTripperFunc) RoundTrip(r *http.Request) (*http.Response, error) {
	return fn(r)
}
