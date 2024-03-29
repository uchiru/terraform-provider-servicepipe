package sdkv1

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"time"
)

const (
	// AppVersion is a version of the application.
	appVersion = "0.1.0"

	// AppName is a global application name.
	appName = "sdkv1-servicepipe"

	// userAgent contains a basic user agent that will be used in queries.
	userAgent = appName + "/" + appVersion

	// defaultEndpoint represents default endpoint for ServicePipe API v1.
	defaultEndpoint = "https://api.servicepipe.ru/api/v1"

	// defaultHTTPTimeout represents the default timeout (in seconds) for HTTP
	// requests.
	defaultHTTPTimeout = 120

	// defaultDialTimeout represents the default timeout (in seconds) for HTTP
	// connection establishments.
	defaultDialTimeout = 60

	// defaultKeepalive represents the default keep-alive period for an active
	// network connection.
	defaultKeepaliveTimeout = 60

	// defaultMaxIdleConns represents the maximum number of idle (keep-alive)
	// connections.
	defaultMaxIdleConns = 100

	// defaultIdleConnTimeout represents the maximum amount of time an idle
	// (keep-alive) connection will remain idle before closing itself.
	defaultIdleConnTimeout = 100

	// defaultTLSHandshakeTimeout represents the default timeout (in seconds)
	// for TLS handshake.
	defaultTLSHandshakeTimeout = 60

	// defaultExpectContinueTimeout represents the default amount of time to
	// wait for a server's first response headers.
	defaultExpectContinueTimeout = 1
)

// const spAPIEndpoint = "records"

// Client stores details that are needed to work with ServicePipe API.
type Client struct {
	// HTTPClient represents an initialized HTTP client that will be used to do requests.
	HTTPClient *http.Client

	// Token is a client authentication token.
	Token string

	// Endpoint represents an endpoint that will be used in all requests.
	Endpoint string

	// UserAgent contains user agent that will be used in all requests.
	UserAgent string
}

// NewClientV1 initializes a new client for the ServicePipe API V1.
func NewClientV1(token, endpoint string) *Client {
	return &Client{
		HTTPClient: newHTTPClient(),
		Token:      token,
		Endpoint:   endpoint,
		UserAgent:  userAgent,
	}
}

// NewClientV1WithDefaultEndpoint initializes a new client for the  API V1
// with default endpoint.
func NewClientV1WithDefaultEndpoint(token string) *Client {
	return &Client{
		HTTPClient: newHTTPClient(),
		Token:      token,
		Endpoint:   defaultEndpoint,
		UserAgent:  userAgent,
	}
}

// NewClientV1WithCustomHTTP initializes a new client for the  API V1
// using custom HTTP client.
// If customHTTPClient is nil - default HTTP client will be used.
func NewClientV1WithCustomHTTP(customHTTPClient *http.Client, token, endpoint string) *Client {
	if customHTTPClient == nil {
		customHTTPClient = newHTTPClient()
	}
	return &Client{
		HTTPClient: customHTTPClient,
		Token:      token,
		Endpoint:   endpoint,
		UserAgent:  userAgent,
	}
}

// newHTTPClient returns a reference to an initialized and configured HTTP client.
func newHTTPClient() *http.Client {
	return &http.Client{
		Timeout:   defaultHTTPTimeout * time.Second,
		Transport: newHTTPTransport(),
	}
}

// newHTTPTransport returns a reference to an initialized and configured HTTP transport.
func newHTTPTransport() *http.Transport {
	return &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   defaultDialTimeout * time.Second,
			KeepAlive: defaultKeepaliveTimeout * time.Second,
		}).DialContext,
		MaxIdleConns:          defaultMaxIdleConns,
		IdleConnTimeout:       defaultIdleConnTimeout * time.Second,
		TLSHandshakeTimeout:   defaultTLSHandshakeTimeout * time.Second,
		ExpectContinueTimeout: defaultExpectContinueTimeout * time.Second,
	}
}

// DoRequest performs the HTTP request with the current ServiceClient's HTTPClient.
// Authentication and optional headers will be added automatically.
func (client *Client) DoRequest(ctx context.Context, method, path string, body io.Reader) (*ResponseResult, error) {
	// Prepare an HTTP request with the provided context.
	request, err := http.NewRequest(method, path, body)
	if err != nil {
		return nil, err
	}

	request.Header.Set("User-Agent", client.UserAgent)
	tok := strings.Join([]string{"Bearer", client.Token}, " ")
	request.Header.Set("Authorization", tok)

	if body != nil {
		request.Header.Set("Content-Type", "application/json")
	}
	request = request.WithContext(ctx)

	// Send the HTTP request and populate the ResponseResult.
	response, err := client.HTTPClient.Do(request)
	if err != nil {
		return nil, err
	}

	responseResult := &ResponseResult{
		Response:    response,
		RequestBody: body,
	}

	// Check status code and populate custom error body with extended error message if it's possible.
	if response.StatusCode >= http.StatusBadRequest {
		err = responseResult.extractErr()
	}
	if err != nil {
		return nil, err
	}

	return responseResult, nil
}

func (client *Client) Echo(ctx context.Context) (bool, *ResponseResult, error) {

	l7ResourcePath := "l7/resource"
	url := strings.Join([]string{client.Endpoint, l7ResourcePath}, "/")
	responseResult, err := client.DoRequest(ctx, http.MethodGet, url, nil)
	if err != nil {
		return false, nil, err
	}
	if responseResult.Err != nil {
		return false, responseResult, responseResult.Err
	}

	return true, responseResult, responseResult.Err
}

// ResponseResult represents a result of an HTTP request.
// It embeds standard http.Response and adds custom API error representations.
type ResponseResult struct {
	*http.Response

	*ErrNotFound

	*ErrGeneric

	// Err contains an error that can be provided to a caller.
	Err error

	RequestBody io.Reader
}

// ErrNotFound represents 404 status code error of an HTTP response.
type ErrNotFound struct {
	Error string `json:"error"`
}

// ErrGeneric represents a generic error of an HTTP response.
type ErrGeneric struct {
	Error string `json:"error"`
}

// ExtractResult allows to provide an object into which ResponseResult body will be extracted.
func (result *ResponseResult) ExtractResult(to interface{}) error {
	body, err := io.ReadAll(result.Body)
	if err != nil {
		return err
	}
	defer result.Body.Close()

	return json.Unmarshal(body, to)
}

// ExtractRaw extracts ResponseResult body into the slice of bytes without unmarshalling.
func (result *ResponseResult) ExtractRaw() ([]byte, error) {
	bytes, err := io.ReadAll(result.Body)
	if err != nil {
		return nil, err
	}
	defer result.Body.Close()

	return bytes, nil
}

// extractErr populates an error message and error structure in the ResponseResult body.
func (result *ResponseResult) extractErr() error {
	body, err := io.ReadAll(result.Body)
	if err != nil {
		return err
	}
	defer result.Body.Close()

	if len(body) == 0 {
		result.Err = fmt.Errorf("sp-go: got the %d status code from the server", result.StatusCode)
		return nil
	}

	switch result.StatusCode {
	case http.StatusNotFound:
		err = json.Unmarshal(body, &result.ErrNotFound)
	default:
		err = json.Unmarshal(body, &result.ErrGeneric)
	}
	if err != nil {
		result.Err = fmt.Errorf("sp-go: got invalid response from the server, status code %d",
			result.StatusCode)
		return err
	}

	result.Err = fmt.Errorf("sp-go: got the %d status code from the server: %s", result.StatusCode, string(body))

	return nil
}
