package trackerapi

import(
	"net/http"
)

type MyClient struct{
  *http.Client
}

func (my *MyClient)Do(req *http.Request) (*http.Response, error){

	req.Header.Set("cache-control","max-stale=3600")

	return my.Client.Do(req)

}