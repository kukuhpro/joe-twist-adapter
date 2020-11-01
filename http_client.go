package twist

import (
	"net/http"
)

type clientAPI interface {
	Send(request *http.Request) (*http.Response, error)
}

type httpClient struct {
}

func (h *httpClient) Send(request *http.Request) (*http.Response, error) {
	var c http.Client
	res, err := c.Do(request)
	if err != nil {
		return nil, err
	}
	return res, err
}
