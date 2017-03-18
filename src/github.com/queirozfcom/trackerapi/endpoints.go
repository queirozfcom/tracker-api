package trackerapi

import (
	"context"
	"net/url"
	"strings"

	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/google/go-github/github"
)

// Endpoints collects all of the endpoints that compose a profile service. It's
// meant to be used as a helper struct, to collect all of the endpoints into a
// single parameter.
//
// In a server, it's useful for functions that need to operate on a per-endpoint
// basis. For example, you might pass an Endpoints to a function that produces
// an http.Handler, with each method (endpoint) wired up to a specific path. (It
// is probably a mistake in design to invoke the Service methods on the
// Endpoints struct in a server.)
//
// In a client, it's useful to collect individually constructed endpoints into a
// single type that implements the Service interface. For example, you might
// construct individual endpoints using transport/http.NewClient, combine them
// into an Endpoints, and return it to the caller as a Service.
type Endpoints struct {
	GetWatchedReposEndpoint endpoint.Endpoint
	//PostProfileEndpoint   endpoint.Endpoint
	//GetProfileEndpoint    endpoint.Endpoint
}

// MakeServerEndpoints returns an Endpoints struct where each endpoint invokes
// the corresponding method on the provided service. Useful in a profilesvc
// server.
func MakeServerEndpoints(s Service) Endpoints {
	return Endpoints{
		GetWatchedReposEndpoint: MakeGetWatchedReposEndpoint(s),
		//PostProfileEndpoint:   MakePostProfileEndpoint(s),
		//GetProfileEndpoint:    MakeGetProfileEndpoint(s),
	}
}

// MakeClientEndpoints returns an Endpoints struct where each endpoint invokes
// the corresponding method on the remote instance, via a transport/http.Client.
// Useful in a profilesvc client.
func MakeClientEndpoints(instance string) (Endpoints, error) {
	if !strings.HasPrefix(instance, "http") {
		instance = "http://" + instance
	}
	tgt, err := url.Parse(instance)
	if err != nil {
		return Endpoints{}, err
	}
	tgt.Path = ""

	options := []httptransport.ClientOption{}

	// Note that the request encoders need to modify the request URL, changing
	// the path and method. That's fine: we simply need to provide specific
	// encoders for each endpoint.

	return Endpoints{
		GetWatchedReposEndpoint: httptransport.NewClient("GET", tgt, encodeWatchedReposRequest, decodeWatchedReposResponse, options...).Endpoint(),
	}, nil
}

//// PostProfile implements Service. Primarily useful in a client.
//func (e Endpoints) PostProfile(ctx context.Context, p Profile) error {
//	request := postProfileRequest{Profile: p}
//	response, err := e.PostProfileEndpoint(ctx, request)
//	if err != nil {
//		return err
//	}
//	resp := response.(postProfileResponse)
//	return resp.Err
//}
//
//// GetProfile implements Service. Primarily useful in a client.
//func (e Endpoints) GetProfile(ctx context.Context, id string) (Profile, error) {
//	request := getProfileRequest{ID: id}
//	response, err := e.GetProfileEndpoint(ctx, request)
//	if err != nil {
//		return Profile{}, err
//	}
//	resp := response.(getProfileResponse)
//	return resp.Profile, resp.Err
//}
//
//// MakePostProfileEndpoint returns an endpoint via the passed service.
//// Primarily useful in a server.
//func MakePostProfileEndpoint(s Service) endpoint.Endpoint {
//	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
//		req := request.(postProfileRequest)
//		e := s.PostProfile(ctx, req.Profile)
//		return postProfileResponse{Err: e}, nil
//	}
//}
//
//// MakeGetProfileEndpoint returns an endpoint via the passed service.
//// Primarily useful in a server.
//func MakeGetProfileEndpoint(s Service) endpoint.Endpoint {
//	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
//		req := request.(getProfileRequest)
//		p, e := s.GetProfile(ctx, req.ID)
//		return getProfileResponse{Profile: p, Err: e}, nil
//	}
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

type postProfileRequest struct {
	Profile Profile
}

type postProfileResponse struct {
	Err error `json:"err,omitempty"`
}

func (r postProfileResponse) error() error { return r.Err }

type getProfileRequest struct {
	ID string
}

type getProfileResponse struct {
	Profile Profile `json:"profile,omitempty"`
	Err     error   `json:"err,omitempty"`
}

func (r getProfileResponse) error() error { return r.Err }
