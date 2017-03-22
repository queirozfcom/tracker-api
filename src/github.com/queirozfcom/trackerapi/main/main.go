package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"../../trackerapi"
	"golang.org/x/oauth2"
	"github.com/go-kit/kit/log"
	"github.com/google/go-github/github"
	"github.com/gregjones/httpcache"
	"github.com/die-net/lrucache"
)

func main() {

	var transport http.RoundTripper

	var (
		httpAddr = flag.String("http.addr", ":8080", "HTTP listen address")
	)
	flag.Parse()

	tok, exists := os.LookupEnv("GITHUB_TOKEN")

	if ! exists {
		panic("need a github token")
	}

	// GitHub API authentication.
	transport = &oauth2.Transport{
		Source: oauth2.StaticTokenSource(&oauth2.Token{AccessToken: tok, TokenType: "Basic"}),
	}

	// Memory caching.
	transport = &httpcache.Transport{
		Transport:           transport,
		Cache:               lrucache.New(100000000, 21000),
		MarkCachedResponses: true,
	}

	httpClient := &http.Client{Transport: transport}

	githubClient := github.NewClient(httpClient)

	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}

	var s trackerapi.Service
	{
		s = trackerapi.NewInmemService(*githubClient)
		s = trackerapi.LoggingMiddleware(logger)(s)
	}

	var h http.Handler
	{
		h = trackerapi.MakeHTTPHandler(s, log.With(logger, "component", "HTTP"))
	}

	errs := make(chan error)
	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errs <- fmt.Errorf("%s", <-c)
	}()

	go func() {
		logger.Log("transport", "HTTP", "addr", *httpAddr)
		errs <- http.ListenAndServe(*httpAddr, h)
	}()

	logger.Log("exit", <-errs)
}
