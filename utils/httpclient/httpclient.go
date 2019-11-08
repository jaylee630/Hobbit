package httpclient

import (
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/pkg/errors"
)

type (
	HttpClient struct {
		cli    http.Client
		req    *http.Request
		resp   *http.Response
		header map[string]string
		err    error
	}

	HttpResp struct {
		StatusCode int
		Body       []byte
	}
)

func NewHTTPClient() *HttpClient {
	return &HttpClient{
		header: make(map[string]string),
	}
}

func (h *HttpClient) Request() *http.Request {
	return h.req
}

func (h *HttpClient) Response() *http.Response {
	return h.resp
}

func (h *HttpClient) WriteHeader(k, v string) *HttpClient {

	h.header[k] = v
	return h
}

func (h *HttpClient) doWriteHeader() {

	if h.req == nil {
		return
	}

	for k, v := range h.header {
		h.req.Header.Set(k, v)
	}
}

func (h *HttpClient) Timeout(second int) *HttpClient {
	h.cli.Timeout = time.Duration(second) * time.Second

	return h
}

func (h *HttpClient) Get(url string) *HttpClient {

	return h.Do("GET", url, nil)
}

func (h *HttpClient) Post(url string, body string) *HttpClient {

	return h.Do("POST", url, strings.NewReader(body))

}

func (h *HttpClient) Do(method, url string, body io.Reader) *HttpClient {

	h.req, h.err = http.NewRequest(strings.ToUpper(method), url, body)
	h.doWriteHeader()
	h.resp, h.err = h.cli.Do(h.req)

	return h
}

func (h *HttpClient) Result(resp *HttpResp) *HttpClient {

	if h.err != nil {
		return h
	}

	if h.resp == nil {
		h.err = errors.New("nil response, please do a request first. ")
		return h
	}

	defer h.resp.Body.Close()
	body, err := ioutil.ReadAll(h.resp.Body)

	if err != nil {
		h.err = errors.New("nil response, please do a request first. ")
		return h
	}

	resp.StatusCode = h.resp.StatusCode
	resp.Body = body

	return h
}

func (h *HttpClient) Error() error {
	return h.err
}
