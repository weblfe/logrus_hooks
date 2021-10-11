package notify

import (
	"errors"
	"io"
	"net/http"
)

type (
	httpClient struct {
		url         string
		method      string
		contentType string
	}
)

func NewUrlClient(url string, method ...string) *httpClient {
	var client = new(httpClient)
	client.url = url
	if len(method) <= 0 {
		method = append(method, http.MethodGet)
	}
	client.method = method[0]
	return client
}

func (client *httpClient) Send(params map[string]string) error {
	if client.url == "" {
		return errors.New("http client miss request url")
	}
	var (
		err error
		res *http.Response
	)
	switch client.method {
	case http.MethodGet:
		res, err = http.Get(client.url)
	case http.MethodPost:
		res, err = http.Post(client.url, client.contentType, client.parseBody(params))
	case http.MethodPut:
		var httpClientImpl = http.Client{}
		req, errReq := http.NewRequest(http.MethodPut, client.url, client.parseBody(params))
		if errReq != nil {
			return errReq
		}
		req.Header.Set("Content-Type", client.contentType)
		res, err = httpClientImpl.Do(req)
	case http.MethodDelete:
		var httpClientImpl = http.Client{}
		req, errReq := http.NewRequest(http.MethodDelete, client.url, client.parseBody(params))
		if errReq != nil {
			return errReq
		}
		req.Header.Set("Content-Type", client.contentType)
		res, err = httpClientImpl.Do(req)
	default:
		err = errors.New("unSupport method")
		return err
	}
	if res.StatusCode != 200 {
		return errors.New("response error:" + res.Status)
	}
	return err
}

func (client *httpClient) parseBody(params map[string]string) io.Reader {
	switch client.method {
	case http.MethodGet:
	case http.MethodPost:
	case http.MethodPut:
	case http.MethodDelete:
	}
	return nil
}

func (client *httpClient) SetContentType(contentType string) *httpClient {
	if client.contentType == "" {
		client.contentType = contentType
	}
	return client
}
