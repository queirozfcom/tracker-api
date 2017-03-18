// Package client provides a  client based on a predefined Consul
// service name and relevant tags. Users must only provide the address of a
// Consul server.
package trackerapi

import (
	"io"
	"time"

	consulapi "github.com/hashicorp/consul/api"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/sd"
	"github.com/go-kit/kit/sd/consul"
	"github.com/go-kit/kit/sd/lb"

)

// New returns a service that's load-balanced over instances of  found
// in the provided Consul server. The mechanism of looking up 
// instances in Consul is hard-coded into the client.
func New(consulAddr string, logger log.Logger) (Service, error) {
	apiclient, err := consulapi.NewClient(&consulapi.Config{
		Address: consulAddr,
	})
	if err != nil {
		return nil, err
	}

	// As the implementer of , we declare and enforce these
	// parameters for all of the  consumers.
	var (
		consulService = "trackerapi"
		consulTags    = []string{"prod"}
		passingOnly   = true
		retryMax      = 3
		retryTimeout  = 500 * time.Millisecond
	)

	var (
		sdclient  = consul.NewClient(apiclient)
		endpoints Endpoints
	)
	{
		factory := factoryFor(MakePostProfileEndpoint)
		subscriber := consul.NewSubscriber(sdclient, factory, logger, consulService, consulTags, passingOnly)
		balancer := lb.NewRoundRobin(subscriber)
		retry := lb.Retry(retryMax, retryTimeout, balancer)
		endpoints.PostProfileEndpoint = retry
	}
	{
		factory := factoryFor(MakeGetProfileEndpoint)
		subscriber := consul.NewSubscriber(sdclient, factory, logger, consulService, consulTags, passingOnly)
		balancer := lb.NewRoundRobin(subscriber)
		retry := lb.Retry(retryMax, retryTimeout, balancer)
		endpoints.GetProfileEndpoint = retry
	}

	return endpoints, nil
}

func factoryFor(makeEndpoint func(Service) endpoint.Endpoint) sd.Factory {
	return func(instance string) (endpoint.Endpoint, io.Closer, error) {
		service, err := MakeClientEndpoints(instance)
		if err != nil {
			return nil, nil, err
		}
		return makeEndpoint(service), nil, nil
	}
}
