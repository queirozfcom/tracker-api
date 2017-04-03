package trackerapi

import(
	"net/http"
	"strconv"
)

// DISCLAIMER this does not follow the contract for RoundTripper,
// this is just for experimentation purposes

type MyTransport struct {
	Transport http.RoundTripper
	MaxStale int //seconds
}

func (my *MyTransport) RoundTrip(req *http.Request) (*http.Response, error) {

	transport := my.Transport
	if transport == nil {
		transport = http.DefaultTransport
	}

	req2 := cloneRequest(req) // per RoundTripper contract

	req2.Header.Add("cache-control","max-stale="+strconv.Itoa(my.MaxStale))

	res, err := transport.RoundTrip(req2)

	return res, err
}

// cloneRequest returns a clone of the provided *http.Request.
// The clone is a shallow copy of the struct and its Header map.
func cloneRequest(r *http.Request) *http.Request {
	// shallow copy of the struct
	r2 := new(http.Request)
	*r2 = *r
	// deep copy of the Header
	r2.Header = make(http.Header, len(r.Header))
	for k, s := range r.Header {
		r2.Header[k] = append([]string(nil), s...)
	}
	return r2
}