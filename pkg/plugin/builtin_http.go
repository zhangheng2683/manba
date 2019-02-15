package plugin

import (
	"bytes"
	"io/ioutil"
	"net/http"
)

// HTTPResult result
type HTTPResult struct {
	rsp  *http.Response
	err  error
	body string
}

func newHTTPResult(rsp *http.Response, err error) *HTTPResult {
	if rsp != nil {
		defer rsp.Body.Close()

		data, err := ioutil.ReadAll(rsp.Body)
		if err != nil {
			return &HTTPResult{
				err: err,
			}
		}

		return &HTTPResult{
			err:  err,
			body: string(data),
			rsp:  rsp,
		}
	}

	return &HTTPResult{
		err: err,
		rsp: rsp,
	}
}

// HasError returns true if has a error
func (res *HTTPResult) HasError() bool {
	return res.err != nil
}

// Error returns error
func (res *HTTPResult) Error() string {
	if res.err != nil {
		return res.err.Error()
	}

	return ""
}

// StatusCode returns status code
func (res *HTTPResult) StatusCode() int {
	if res.HasError() {
		return 0
	}

	return res.rsp.StatusCode
}

// Header returns http response header
func (res *HTTPResult) Header() map[string][]string {
	headers := make(map[string][]string)
	if res.HasError() {
		return headers
	}

	for key, values := range res.rsp.Header {
		headers[key] = values
	}

	return headers
}

// Cookie returns http response cookie
func (res *HTTPResult) Cookie() []*http.Cookie {
	if res.HasError() {
		return nil
	}

	return res.rsp.Cookies()
}

// Body returns http response body
func (res *HTTPResult) Body() string {
	if res.HasError() {
		return ""
	}

	return res.body
}

// HTTP http support
type HTTP struct {
}

// Get go get
func (h *HTTP) Get(url string) *HTTPResult {
	rsp, err := http.DefaultClient.Get(url)
	return newHTTPResult(rsp, err)
}

// Post do post
func (h *HTTP) Post(url string, body string, header map[string][]string) *HTTPResult {
	return h.do("POST", url, body, header)
}

// PostJSON do post
func (h *HTTP) PostJSON(url string, body string, header map[string][]string) *HTTPResult {
	header["Content-Type"] = []string{"application/json"}
	return h.Post(url, body, header)
}

// Put do put
func (h *HTTP) Put(url string, body string, header map[string][]string) *HTTPResult {
	return h.do("PUT", url, body, header)
}

// PutJSON do put json
func (h *HTTP) PutJSON(url string, body string, header map[string][]string) *HTTPResult {
	header["Content-Type"] = []string{"application/json"}
	return h.Put(url, body, header)
}

// Delete do delete
func (h *HTTP) Delete(url string, body string, header map[string][]string) *HTTPResult {
	return h.do("DELETE", url, body, header)
}

// DeleteJSON do delete json
func (h *HTTP) DeleteJSON(url string, body string, header map[string][]string) *HTTPResult {
	header["Content-Type"] = []string{"application/json"}
	return h.Delete(url, body, header)
}

func (h *HTTP) do(method string, url string, body string, header map[string][]string) *HTTPResult {
	r := bytes.NewReader([]byte(body))
	req, err := http.NewRequest(method, url, r)
	if err != nil {
		return newHTTPResult(nil, err)
	}

	for key, values := range header {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}

	rsp, err := http.DefaultClient.Do(req)
	return newHTTPResult(rsp, err)
}
