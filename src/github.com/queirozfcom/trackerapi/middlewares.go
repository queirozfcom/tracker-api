package trackerapi

import (
	"context"
	"time"

	"github.com/go-kit/kit/log"
)

// Middleware describes a service (as opposed to endpoint) middleware.
type Middleware func(Service) Service

func LoggingMiddleware(logger log.Logger) Middleware {
	return func(next Service) Service {
		return &loggingMiddleware{
			next:   next,
			logger: logger,
		}
	}
}

type loggingMiddleware struct {
	next   Service
	logger log.Logger
}

func (mw loggingMiddleware) GetWatchedRepos(ctx context.Context, username string) ([]RepoInformation, error) {

	defer func(begin time.Time) {
		mw.logger.Log("method", "GetWatchedRepos", "username", username, "took", time.Since(begin))
	}(time.Now())

	return mw.next.GetWatchedRepos(ctx, username)
}

func (mw loggingMiddleware) GetStarredRepos(ctx context.Context, username string) ([]RepoInformation, error) {

	defer func(begin time.Time) {
		mw.logger.Log("method", "GetStarredRepos", "username", username, "took", time.Since(begin))
	}(time.Now())

	return mw.next.GetStarredRepos(ctx, username)
}
