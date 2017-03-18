package trackerapi

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/google/go-github/github"
)

type Endpoints struct {
	GetWatchedReposEndpoint endpoint.Endpoint
}

// MakeServerEndpoints returns an Endpoints struct where each endpoint invokes
// the corresponding method on the provided service. Useful in a profilesvc
// server.
func MakeServerEndpoints(s Service) Endpoints {
	return Endpoints{
		GetWatchedReposEndpoint: MakeGetWatchedReposEndpoint(s),
	}
}
//
//// MakeClientEndpoints returns an Endpoints struct where each endpoint invokes
//// the corresponding method on the remote instance, via a transport/http.Client.
//// Useful in a profilesvc client.
//func MakeClientEndpoints(instance string) (Endpoints, error) {
//	if !strings.HasPrefix(instance, "http") {
//		instance = "http://" + instance
//	}
//	tgt, err := url.Parse(instance)
//	if err != nil {
//		return Endpoints{}, err
//	}
//	tgt.Path = ""
//
//	options := []httptransport.ClientOption{}
//
//	// Note that the request encoders need to modify the request URL, changing
//	// the path and method. That's fine: we simply need to provide specific
//	// encoders for each endpoint.
//
//	return Endpoints{
//		GetWatchedReposEndpoint: httptransport.NewClient("GET", tgt, encodeWatchedReposRequest, decodeWatchedReposResponse, options...).Endpoint(),
//	}, nil
//}

func (e Endpoints) GetWatchedRepos(ctx context.Context, username string) ([]github.Repository, error) {
	request := getWatchedReposRequest{Username: username}
	response, err := e.GetWatchedReposEndpoint(ctx, request)
	if err != nil {
		return []github.Repository{}, err
	}

	resp := response.(getWatchedReposResponse)

	return resp.Repos, resp.Err
}

func MakeGetWatchedReposEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(getWatchedReposRequest)
		repos, e := s.GetWatchedRepos(ctx, req.Username)
		return getWatchedReposResponse{Repos: repos, Err: e}, nil
	}
}

// We have two options to return errors from the business logic.
//
// We could return the error via the endpoint itself. That makes certain things
// a little bit easier, like providing non-200 HTTP responses to the client. But
// Go kit assumes that endpoint errors are (or may be treated as)
// transport-domain errors. For example, an endpoint error will count against a
// circuit breaker error count.
//
// Therefore, it's often better to return service (business logic) errors in the
// response object. This means we have to do a bit more work in the HTTP
// response encoder to detect e.g. a not-found error and provide a proper HTTP
// status code. That work is done with the errorer interface, in transport.go.
// Response types that may contain business-logic errors implement that
// interface.

type getWatchedReposRequest struct {
	Username string
}

type getWatchedReposResponse struct {
	Repos []github.Repository `json:"repos,omitempty"`
	Err   error `json:"err,omitempty"`
}

func (r getWatchedReposResponse) error() error { return r.Err }
