package trackerapi

import (
	"context"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/google/go-github/github"
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

func (mw loggingMiddleware) GetWatchedRepos(ctx context.Context, username string) ([]github.Repository, error) {

	defer func(begin time.Time) {
		mw.logger.Log("method", "GetWatchedRepos", "username", username, "took", time.Since(begin))
	}(time.Now())

	return mw.next.GetWatchedRepos(ctx, username)
}

//func (mw loggingMiddleware) PostProfile(ctx context.Context, p Profile) (err error) {
//	defer func(begin time.Time) {
//		mw.logger.Log("method", "PostProfile", "id", p.ID, "took", time.Since(begin), "err", err)
//	}(time.Now())
//	return mw.next.PostProfile(ctx, p)
//}
//
//func (mw loggingMiddleware) GetProfile(ctx context.Context, id string) (p Profile, err error) {
//	defer func(begin time.Time) {
//		mw.logger.Log("method", "GetProfile", "id", id, "took", time.Since(begin), "err", err)
//	}(time.Now())
//	return mw.next.GetProfile(ctx, id)
//}
