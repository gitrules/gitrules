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
)

func RunServer(ctx context.Context, appServerAddr string, cfg *Config) {

	metricsRegistry := metrics.DefaultRegistry

	cc, err := githubapp.NewDefaultCachingClientCreator(
		cfg.Github,
		githubapp.WithClientUserAgent("gitrules-for-github/1.0.0"),
		githubapp.WithClientTimeout(3*time.Second),
		githubapp.WithClientCaching(false, func() httpcache.Cache { return httpcache.NewMemoryCache() }),
		githubapp.WithClientMiddleware(
			githubapp.ClientMetrics(metricsRegistry),
		),
	)
	must.NoError(ctx, err)

	prCommentHandler := &PRCommentHandler{
		ClientCreator: cc,
		preamble:      cfg.App.PullRequestPreamble,
	}

	webhookHandler := githubapp.NewDefaultEventDispatcher(cfg.Github, prCommentHandler)

	http.Handle(githubapp.DefaultWebhookRoute, webhookHandler)

	base.Infof("Starting GitRules for GitHub app server on %s ...", appServerAddr)
	must.NoError(ctx, http.ListenAndServe(appServerAddr, nil))
}
